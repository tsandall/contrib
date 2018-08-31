# OPA-Elasticsearch Data Filtering Example

This directory contains an example of how to perform data filtering in
Elasticsearch using the queries provided by OPA's Compile API.

The example server is written in Go and when it receives API requests it asks
OPA for a set of conditions to apply to the Elasticsearch query that serves
the request. OPA is integrated as a library.

## Building

Build the example by running `make build`

## Running the example

1. Run Elasticsearch. See Elasticsearch's [Installation](https://www.elastic.co/guide/en/elasticsearch/reference/current/_installation.html) to get started.

2. Open a new window and start the example server:
   ```bash
   ./opa-es-filtering
   ```

   The server listens on `:8080` and loads sample data. The server exposes
   two endpoints `/posts` and `/posts/{post_id}`.
   OPA is loaded with an example policy which has rules related to both these
   endpoints.

3. Open a new window and make a request:
   ```bash
   curl  -H "Authorization: bob" localhost:8080/posts |  jq .
   ```


## Supported OPA Built-in Functions

### Comparison

- [x] ==
- [x] !=
- [x] <
- [x] <=
- [x] >
- [x] >=

### Strings

- [x] contains

### Regex

- [x] re_match
