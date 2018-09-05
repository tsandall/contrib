// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/olivere/elastic"
	"github.com/open-policy-agent/contrib/data_filter_elasticsearch/internal/es"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

const policyFileName = "example.rego"

type Result struct {
	Defined bool
	Query   elastic.Query
}

func Compile(ctx context.Context, input map[string]interface{}) (Result, error) {
	unknowns := []string{"data.posts"}

	inputBytes, oerr := json.Marshal(input)
	if oerr != nil {
		return Result{}, fmt.Errorf("JSON Encoding error %v", oerr)
	}

	inputTerm, err := ast.ParseTerm(string(inputBytes))
	if err != nil {
		return Result{}, err
	}

	module, err := ioutil.ReadFile(policyFileName)
	if err != nil {
		return Result{}, fmt.Errorf("failed to read policy: %v", err)
	}

	r := rego.New(
		rego.Query("data.example.allow == true"),
		rego.Module(policyFileName, string(module)),
		rego.ParsedInput(inputTerm.Value),
		rego.Unknowns(unknowns),
	)

	pq, err := r.Partial(ctx)
	if err != nil {
		return Result{}, err
	}

	if len(pq.Queries) == 0 {
		// always deny
		return Result{Defined: false}, nil
	} else {
		for _, query := range pq.Queries {
			if len(query) == 0 {
				// always allow
				return Result{Defined: true}, nil
			}
		}
	}
	return processQuery(pq)
}

func processQuery(pq *rego.PartialQueries) (Result, error) {

	queries := []elastic.Query{}
	for i := range pq.Queries {
		fmt.Printf("Query #%d: %v\n\n", i+1, pq.Queries[i])

		exprQueries := []elastic.Query{}
		for _, expr := range pq.Queries[i] {
			if len(expr.Operands()) != 2 {
				continue
			}

			var value string
			var processedTerm []string
			for _, term := range expr.Operands() {
				if term.IsGround() {
					value = strings.Trim(term.String(), "\"\"")
				} else {
					processedTerm = processTerm(term.String())
				}
			}

			var esQuery elastic.Query

			if isEqualityOperator(expr.Operator().String()) {
				// generate ES Term query
				esQuery = es.GenerateTermQuery(processedTerm[1], value)
			} else if isRangeOperator(expr.Operator().String()) {
				// generate ES Range query
				if expr.Operator().String() == "lt" {
					esQuery = es.GenerateRangeQueryLt(processedTerm[1], value)
				} else if expr.Operator().String() == "gt" {
					esQuery = es.GenerateRangeQueryGt(processedTerm[1], value)
				} else if expr.Operator().String() == "lte" {
					esQuery = es.GenerateRangeQueryLte(processedTerm[1], value)
				} else {
					esQuery = es.GenerateRangeQueryGte(processedTerm[1], value)
				}
			} else if expr.Operator().String() == "neq" {
				// generate ES Must Not query
				esQuery = es.GenerateBoolMustNotQuery(processedTerm[1], value)
			} else if isContainsOperator(expr.Operator().String()) {
				// generate ES Query String query
				esQuery = es.GenerateQueryStringQuery(processedTerm[1], value)
			} else if isRegexpMatchOperator(expr.Operator().String()) {
				// generate ES Regexp query
				esQuery = es.GenerateRegexpQuery(processedTerm[1], value)
			} else {
				return Result{}, fmt.Errorf("Unsupported Operator: %v", expr.Operator().String())
			}
			exprQueries = append(exprQueries, esQuery)
		}

		if len(exprQueries) == 1 {
			queries = append(queries, exprQueries[0])
		} else {
			// ES queries generated within a rule are And'ed
			boolQuery := es.GenerateBoolFilterQuery(exprQueries)
			queries = append(queries, boolQuery)
		}
	}

	// ES queries generated from partial eval queries
	// are Or'ed
	combinedQuery := es.GenerateBoolShouldQuery(queries)
	return Result{Defined: true, Query: combinedQuery}, nil

}

func processTerm(query string) []string {
	splitQ := strings.Split(query, ".")

	indexName := strings.Split(splitQ[1], "[")[0]
	fieldName := splitQ[2]
	return []string{indexName, fieldName}
}

func isEqualityOperator(op string) bool {
	return op == "eq" || op == "equal"
}

func isContainsOperator(op string) bool {
	return op == "contains"
}

func isRegexpMatchOperator(op string) bool {
	return op == "re_match"
}

func isRangeOperator(op string) bool {
	return op == "lt" || op == "gt" || op == "lte" || op == "gte"
}
