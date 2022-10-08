package MessageQueue

import (
	logger "Worker/Logger"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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
