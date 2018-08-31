// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package es

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
)

// Posts is a structure used for serializing/deserializing data in Elasticsearch.
type Post struct {
	Id         string   `json:"id"`
	Author     string   `json:"author"`
	Message    string   `json:"message"`
	Department string   `json:"department"`
	Email      string   `json:"email"`
	Clearance  int      `json:"clearance"`
	Likes      []string `json: likes`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"_doc":{
			"properties":{
				"id":{
					"type":"keyword"
				},
				"author":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"fields": {
						"raw": {
							"type": "keyword"
						}
					}
				},
				"department":{
					"type":"keyword"
				},
				"email":{
					"type":"keyword"
				},
				"clearance":{
					"type":"integer"
				},
				"likes":{
					"type":"nested"
				}
			}
		}
	}
}`

func NewPost(id, author, message, department, email string, clearance int, likes []string) *Post {
	return &Post{Id: id, Author: author, Message: message, Department: department, Email: email, Clearance: clearance, Likes: likes}
}

func NewESClient() (*elastic.Client, error) {
	return elastic.NewClient()
}

func GetIndexMapping() string {
	return mapping
}

// Elasticsearch queries

func GenerateTermQuery(fieldName, fieldValue string) *elastic.TermQuery {
	return elastic.NewTermQuery(fieldName, fieldValue)
}

func GenerateBoolFilterQuery(filters []elastic.Query) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	for _, filter := range filters {
		q = q.Filter(filter)
	}
	return q

}

func GenerateBoolShouldQuery(queries []elastic.Query) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	for _, query := range queries {
		q = q.Should(query)
	}
	return q
}

func GenerateBoolMustNotQuery(fieldName, fieldValue string) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	q = q.MustNot(elastic.NewTermQuery(fieldName, fieldValue))
	return q
}

func GenerateMatchAllQuery() *elastic.MatchAllQuery {
	return elastic.NewMatchAllQuery()
}

func GenerateMatchQuery(fieldName, fieldValue string) *elastic.MatchQuery {
	return elastic.NewMatchQuery(fieldName, fieldValue)
}

func GenerateQueryStringQuery(fieldName, fieldValue string) *elastic.QueryStringQuery {
	queryString := fmt.Sprintf("*%s*", fieldValue)
	q := elastic.NewQueryStringQuery(queryString)
	q = q.DefaultField(fieldName)
	return q
}

func GenerateRegexpQuery(fieldName, fieldValue string) *elastic.RegexpQuery {
	return elastic.NewRegexpQuery(fieldName, fieldValue)
}

func GenerateRangeQueryLt(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Lt(val)
}

func GenerateRangeQueryLte(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Lte(val)
}

func GenerateRangeQueryGt(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Gt(val)
}

func GenerateRangeQueryGte(fieldName string, val interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName).Gte(val)
}

func ExecuteEsSearch(ctx context.Context, client *elastic.Client, indexName string, query elastic.Query) (*elastic.SearchResult, error) {
	searchResult, err := client.Search().
		Index(indexName).
		Query(query). // specify the query
		Pretty(true). // pretty print request and response JSON
		Do(ctx)       // execute
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func AnalyzeSearchResult(searchResult *elastic.SearchResult) {

	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d posts\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// Deserialize hit
			var t Post
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}

			// Print with post
			fmt.Printf("\nPost ID: %s\nAuthor: %s\nMessage: %s\nDepartment: %s\nClearance: %d\n", t.Id, t.Author, t.Message, t.Department, t.Clearance)
		}
	} else {
		// No hits
		fmt.Print("Found no posts\n")
	}
}

func GetPrettyResult(searchResult *elastic.SearchResult) []Post {

	result := []Post{}
	if searchResult.Hits.TotalHits > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// Deserialize hit
			var t Post
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}
			result = append(result, t)
		}
	}
	return result
}
