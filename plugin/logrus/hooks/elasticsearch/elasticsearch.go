package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lsjhtang/kid/logger"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

type EsHook struct {
	opts   Options
	client *elastic.Client
}

func (es *EsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

func (es *EsHook) Fire(entry *logrus.Entry) error {
	doc := make(map[string]interface{})
	for k, v := range entry.Data {
		doc[k] = v
	}
	doc["time"] = time.Now().Local()
	doc["lvl"] = entry.Level
	doc["message"] = entry.Message
	doc["caller"] = fmt.Sprintf("%s:%d %#v", entry.Caller.File, entry.Caller.Line, entry.Caller.Func)
	go func() {
		_, err := es.client.Index().Index(es.opts.indexName()).Type("_doc").BodyJson(doc).Do(context.Background())
		fmt.Println(err)
		if err != nil {
			fmt.Fprintf(os.Stderr, "HandleError to fire hook: %v\n", err)
		}
	}()
	return nil
}

type Option func(*Options)

func NewEsHook(options ...Option) *EsHook {
	opts := Options{
		logLevel:   logger.InfoLevel.String(),
		esAddress:  []string{},
		esUser:     "",
		esPassword: "",
		cmd:        "",
		indexName: func() string {
			return fmt.Sprintf("%s_doc", time.Now().Format("2006-01-02-15:04:05"))
		},
		health: time.Second * 15,
	}
	for _, o := range options {
		o(&opts)
	}

	if len(opts.esAddress) <= 0 {
		log.Fatal("logger.EsHook: elasticSearch address is required")
	}

	h := &EsHook{
		opts: opts,
	}
	es, err := elastic.NewClient(
		elastic.SetURL(h.opts.esAddress...),
		elastic.SetBasicAuth(h.opts.esUser, h.opts.esPassword),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeout(h.opts.health),
		elastic.SetErrorLog(log.New(os.Stderr, "ES:", log.LstdFlags)),
	)
	if err != nil {
		log.Fatal("failed to create ElasticSearch v6 client: ", err.Error())
	}
	h.client = es
	return h
}
