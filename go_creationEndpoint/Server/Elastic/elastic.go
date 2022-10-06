package MessageQueue

import (
	logger "Server/Logger"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

func New(elasticAddr string) *Elastic {
	cfg := elasticsearch.Config{
		Addresses: []string{elasticAddr},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		logger.FailOnError(logger.ELASTIC, logger.ESSENTIAL, "Unable to setup elastic search")
	}

	toRet := &Elastic{
		elastic: es,
		ctx:     context.Background(),
		timeout: time.Second * 10,
	}

	res, err := es.Indices.Exists([]string{MESSAGES_INDEX})
	if err != nil {
		logger.FailOnError(logger.ELASTIC, logger.ESSENTIAL, "Unable to check if elastic index exists")
	}
	if res.StatusCode == 200 {
		return toRet
	}

	//index doesn't exist, so create it
	_, err = es.Indices.Create(MESSAGES_INDEX)
	if err != nil {
		logger.FailOnError(logger.ELASTIC, logger.ESSENTIAL, "Unable to create elastic index")
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

	res, err := req.Do(ctx, &e.elastic.BaseClient)
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
