// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/api"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
)

func main() {
	ctx := context.Background()

	// Create an ES client.
	client, err := es.NewESClient()
	if err != nil {
		panic(err)
	}

	indexName := "posts"

	// Check if a specified index exists.
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(indexName).BodyString(es.GetIndexMapping()).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic("Index creation not acknowledged")
		}
	}

	// Index posts.
	createTestPosts(ctx, client, indexName)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	// Start server.
	if err := api.New(client, indexName).Run(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Shutting down.")
}

func createTestPosts(ctx context.Context, client *elastic.Client, indexName string) {

	post1 := es.NewPost("post1", "bob", "My first post", "dev", "bob@abc.com", 2, []string{})
	_, err := client.Index().
		Index(indexName).
		Type("_doc").
		Id("1").
		BodyJson(post1).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post2 := es.NewPost("post2", "bob", "My second post", "dev", "bob@abc.com", 2, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("2").
		BodyJson(post2).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post3 := es.NewPost("post3", "charlie", "Hello world", "it", "charlie@xyz.com", 1, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("3").
		BodyJson(post3).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post4 := es.NewPost("post4", "alice", "Hii world", "hr", "alice@xyz.com", 3, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("4").
		BodyJson(post4).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post5 := es.NewPost("post5", "ben", "Hii from Ben", "ceo", "ben@opa.com", 10, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("5").
		BodyJson(post5).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post6 := es.NewPost("post6", "ken", "Hii form Ken", "ceo", "ken@opa.com", 5, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("6").
		BodyJson(post6).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post7 := es.NewPost("post7", "john", "OPA Good", "dev", "john@blah.com", 6, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("7").
		BodyJson(post7).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post8 := es.NewPost("post8", "ben", "This is OPA's time", "ceo", "ben@opa.com", 10, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("8").
		BodyJson(post8).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post9 := es.NewPost("post9", "jane", "Hello from Jane", "it", "jane@opa.org", 7, []string{})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("9").
		BodyJson(post9).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	post10 := es.NewPost("post10", "ross", "Hello from Ross", "it", "ross@opal.eu", 9, []string{"bob", "alice"})
	_, err = client.Index().
		Index(indexName).
		Type("_doc").
		Id("10").
		BodyJson(post10).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}
