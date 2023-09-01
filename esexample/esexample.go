package esexample

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var es_client *elasticsearch.Client

// 连接es
func InitEs() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.214.133:31200",
		},
		Username: "elastic",
		Password: "ellischen",
	}

	es_client, _ = elasticsearch.NewClient(cfg)
	boolquery = &esquery.BoolQuery{}
}

// 列出索引
func ListIndex() {
	res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), es_client)
	if err != nil {
		return
	}
	defer res.Body.Close()

	fmt.Println(res.String())
}

// 查询索引文档个数
func CalculateIndexDocCount() float64 {
	// res, err := esapi.CountRequest{Index: []string{"ellis"}}.Do(context.Background(), es_client)
	// if err != nil {
	// 	return 0
	// }
	// defer res.Body.Close()
	// var resMap map[string]interface{}
	// json.NewDecoder(res.Body).Decode(&resMap)
	// fmt.Printf("resMap: %v\n", resMap["count"])
	// // fmt.Printf("res.Header: %v\n", res.Header)
	// // fmt.Println(res.String())
	// return resMap["count"].(float64)

	r, err := es_client.Count(es_client.Count.WithIndex("ellis"))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var value map[string]interface{}
	json.NewDecoder(r.Body).Decode(&value)
	return value["count"].(float64)
}

// 插入一个文档
func IndexOneDocument() {
	//method 1
	// document := Ellis{Name: "haha"}
	// data, _ := json.Marshal(document)
	// req := esapi.IndexRequest{
	// 	Index:      "ellis",
	// 	DocumentID: "3",
	// 	Body:       strings.NewReader(string(data)),
	// 	Refresh:    "true",
	// }
	// res, err := req.Do(context.TODO(), es_client)
	// if err != nil {
	// 	log.Fatalf("IndexRequest ERROR: %s", err)
	// }
	// defer res.Body.Close()

	// if res.IsError() {
	// 	log.Printf("%s ERROR indexing document ID=%d", res.Status(), 3)
	// } else {

	// 	// Deserialize the response into a map.
	// 	var resMap map[string]interface{}
	// 	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
	// 		log.Printf("Error parsing the response body: %s", err)
	// 	} else {
	// 		log.Printf("\nIndexRequest() RESPONSE:")
	// 		// Print the response status and indexed document version.
	// 		fmt.Println("Status:", res.Status())
	// 		fmt.Println("Result:", resMap["result"])
	// 		fmt.Println("Version:", int(resMap["_version"].(float64)))
	// 		fmt.Println("resMap:", resMap)
	// 	}
	// }

	//method 2
	document := Ellis{Name: "haha"}
	data, _ := json.Marshal(document)
	r, err2 := es_client.Index("ellis", strings.NewReader(string(data)), es_client.Index.WithDocumentID("4"))
	if err2 != nil {
		fmt.Printf("err2: %v\n", err2)
	} else {
		defer r.Body.Close()
		var value map[string]interface{}
		json.NewDecoder(r.Body).Decode(&value)

		vv, _ := json.Marshal(value)
		fmt.Printf("string(vv): %v\n", string(vv))
	}
}

// 通过ID查询
func GetByID(id string) (value any) {
	// method 1
	// r, err := esapi.GetRequest{Index: "ellis", DocumentID: id}.Do(context.Background(), es_client)
	// if err != nil {
	// 	return nil
	// } else {
	// 	defer r.Body.Close()
	// 	fmt.Printf("r.String(): %v\n", r.String())
	// 	var value interface{}
	// 	json.NewDecoder(r.Body).Decode(&value)
	// 	fmt.Printf("value: %v\n", value)
	// 	return value
	// }

	// method 2
	r, err := es_client.Get("ellis", id, es_client.Get.WithRefresh(true))
	if err != nil {
		return nil
	} else {
		defer r.Body.Close()
		var value map[string]interface{}
		json.NewDecoder(r.Body).Decode(&value)
		fmt.Printf("value: %v\n", value)
		// 将value转换成JSON
		vv, _ := json.Marshal(value)
		fmt.Printf("string(vv): %v\n", string(vv))
		return value
	}
}

// 通过DSL查询
func SearchByDSL() {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"name": "haha",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es_client.Search(
		es_client.Search.WithContext(context.Background()),
		es_client.Search.WithIndex("ellis"),
		es_client.Search.WithBody(&buf),
		es_client.Search.WithTrackTotalHits(true),
		es_client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
}

// update by query
func UpdateByQuery() {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"name": "haha",
			},
		},
		"script": map[string]interface{}{
			"source": "ctx._source[\"name\"]=params.name",
			"params": map[string]string{
				"name": "ellis",
			},
			"lang": "painless",
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	es_client.UpdateByQuery([]string{"ellis"}, es_client.UpdateByQuery.WithBody(&buf))
}

// 删除
func Delete() {
	r, err := es_client.Delete("ellis", "1", es_client.Delete.WithRefresh("true"))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		defer r.Body.Close()
		var value map[string]interface{}
		json.NewDecoder(r.Body).Decode(&value)
		vv, _ := json.Marshal(value)
		fmt.Printf("string(vv): %v\n", string(vv))
	}
}

// search after
func SearchAfter() {
	query := `{
		"query": {
		  "match_all": {}
		},
		"sort": [
		  {
			"_id": {
			  "order": "desc"
			}
		  }
		],
		"size": 1,
		"search_after":["3"]
	  }`
	res, err := es_client.Search(
		es_client.Search.WithIndex("ellis"),
		es_client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
}

// scroll 查询
func Scroll() {
	query := `{
			"query": {
			"match_all": {}
			}
		}`
	log.Println("Scrolling the index...")
	log.Println(strings.Repeat("-", 80))
	res, err := es_client.Search(
		es_client.Search.WithBody(strings.NewReader(query)),
		es_client.Search.WithIndex("ellis"),
		// es_client.Search.WithSort("_doc"),
		es_client.Search.WithSize(1),
		es_client.Search.WithScroll(time.Minute),
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var mapvalue map[string]interface{}
	json.NewDecoder(res.Body).Decode(&mapvalue)
	value1, _ := json.Marshal(mapvalue)
	jsonvalue := string(value1)

	defer res.Body.Close()

	scrollID := gjson.Get(jsonvalue, "_scroll_id").String()

	log.Println("ScrollID", scrollID)
	log.Println("IDs     ", gjson.Get(jsonvalue, "hits.hits.#._id"))
	log.Println(strings.Repeat("-", 80))
	for _, hit := range mapvalue["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	for {
		// Perform the scroll request and pass the scrollID and scroll duration
		//
		res, err := es_client.Scroll(es_client.Scroll.WithScrollID(scrollID), es_client.Scroll.WithScroll(time.Minute))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		if res.IsError() {
			log.Fatalf("Error response: %s", res)
		}

		defer res.Body.Close()

		// Extract the scrollID from response
		var mapvalue map[string]interface{}
		json.NewDecoder(res.Body).Decode(&mapvalue)
		value1, _ := json.Marshal(mapvalue)
		jsonvalue := string(value1)
		scrollID = gjson.Get(jsonvalue, "_scroll_id").String()

		hits := gjson.Get(jsonvalue, "hits.hits")

		if len(hits.Array()) < 1 {
			log.Println("Finished scrolling")
			break
		} else {
			log.Println("ScrollID", scrollID)
			log.Println("IDs     ", gjson.Get(hits.Raw, "#._id"))
			log.Println(strings.Repeat("-", 80))
			for _, v := range hits.Array() {
				fmt.Printf("v.Raw: %v\n", v.Raw)
			}
		}
	}

}

var boolquery *esquery.BoolQuery

func DynamicDSL(c *gin.Context) {

	var body []Dynamic
	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		if len(body) > 1 {
			for _, v := range body {
				if v.Operation == "=" {
					boolquery.Should(esquery.Term(v.Field, v.Value))
				}
			}
			boolquery.MinimumShouldMatch(1)
		} else {
			for _, v := range body {
				if v.Operation == "=" {
					boolquery.Must(esquery.Term(v.Field, v.Value))
				}
			}
		}

	}

	res, err := esquery.Search().
		Query(boolquery).Sort("_id", esquery.OrderDesc).Run(es_client, es_client.Search.WithIndex("ellis"))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var value map[string]interface{}
	json.NewDecoder(res.Body).Decode(&value)

	b, _ := json.Marshal(value)
	stringjson := string(b)
	log.Println("gjson.Get(stringjson, \"hits.hits.#._id\"):\n", gjson.Get(stringjson, "hits.hits.#._id"))
	defer res.Body.Close()
}

func SearchAfterSecond() {

	res, err := esquery.Search().Query(esquery.MatchAll()).Sort("_id", esquery.OrderDesc).SearchAfter("3").Size(1).Run(es_client, es_client.Search.WithIndex("ellis"))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var value map[string]interface{}
	json.NewDecoder(res.Body).Decode(&value)

	b, _ := json.Marshal(value)
	stringjson := string(b)
	log.Println("gjson.Get(stringjson, \"hits.hits.#._id\"):\n", gjson.Get(stringjson, "hits.hits.#._id"))
	defer res.Body.Close()
}
