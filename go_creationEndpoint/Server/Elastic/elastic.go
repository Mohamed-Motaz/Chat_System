package MessageQueue

import (
	logger "Server/Logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mitchellh/mapstructure"
)

func New(elasticAddr string) *Elastic {
	cfg := elasticsearch.Config{
		Addresses: []string{elasticAddr},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		logger.FailOnError(logger.ELASTIC, logger.ESSENTIAL, "Unable to setup elastic search with err %+v", err)
	}

	toRet := &Elastic{
		elastic: es,
		ctx:     context.Background(),
		timeout: time.Second * 10,
	}

	ctr := 0
	res, err := es.Indices.Exists([]string{MESSAGES_INDEX})
	for ctr < 5 {
		ctr++
		if err == nil {
			continue
		}
		res, err = es.Indices.Exists([]string{MESSAGES_INDEX})
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		logger.FailOnError(logger.ELASTIC, logger.ESSENTIAL, "Unable to check if elastic index exists with err %+v", err)
	}
	if res.StatusCode == 200 {
		return toRet
	}

	//index doesn't exist, so create it
	_, err = es.Indices.Create(MESSAGES_INDEX)
	if err != nil {
		logger.FailOnError(logger.ELASTIC, logger.ESSENTIAL, "Unable to create elastic index with err %+v", err)
	}

	logger.LogInfo(logger.ELASTIC, logger.ESSENTIAL, "Elastic search setup complete")

	return toRet
}

//updates and inserts
func (e Elastic) Index(id string, body []byte) error {
	req := esapi.IndexRequest{
		Index:      MESSAGES_INDEX,
		Body:       bytes.NewReader(body),
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(e.ctx, e.timeout)
	defer cancel()

	res, err := req.Do(ctx, e.elastic)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return fmt.Errorf("conflict while inserting")
	}

	if res.IsError() {
		return fmt.Errorf(res.String())
	}

	return nil
}

func (e Elastic) Get(chatId int, partialMatch string) ([]ElasticObj, error) {
	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match_phrase": map[string]interface{}{
						"body": fmt.Sprintf("(.|\n)*%s(.|\n)*", partialMatch),
					},
					},
					{"match": map[string]interface{}{
						"chat_id": chatId,
					}},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// Perform the search request.
	res, err := e.elastic.Search(
		e.elastic.Search.WithContext(context.Background()),
		e.elastic.Search.WithIndex(MESSAGES_INDEX),
		e.elastic.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	arr := []ElasticObj{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		obj := &ElasticObj{}
		err := mapstructure.Decode(hit.(map[string]interface{})["_source"], obj)
		if err != nil {
			return nil, err
		}
		arr = append(arr, *obj)
	}

	fmt.Printf("final %+v\n", arr)

	return arr, nil
}
