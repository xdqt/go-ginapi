package esexample

import (
	"github.com/elastic/go-elasticsearch/v7"
)

var es_client *elasticsearch.Client

func InitEs() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.214.133:31200",
		},
		Username: "elastic",
		Password: "ellischen",
	}

	es_client, _ = elasticsearch.NewClient(cfg)
}

func QueryByDsl(value map[string]interface{}) {

}
