// Copyright 2017 Tomas Machalek <tomas.machalek@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ESClient is a simple ElasticSearch client
type ESClient struct {
	server string
	index  string
}

// NewClient returns an instance of ESClient
func NewClient(server string, index string) *ESClient {
	c := ESClient{
		server: server,
		index:  index,
	}
	return &c
}

// UpdateSetAPIFlag sets a (new) attribute "isApi" to "true" for all
// the matching documents
func (c *ESClient) UpdateSetAPIFlag(conf ConfUpdate) (string, error) {
	if !conf.Disabled {
		ans, err := CreateQuery(conf.FromDate, conf.ToDate, conf.IPAddress, conf.UserAgent)
		if err == nil {
			return c.Do("GET", "/"+c.index+"/_search", ans)
		}
		return "", err
	}
	return "", nil
}

/*
TODO: cheatsheet

partial update;
POST /website/blog/1/_update
{
   "doc" : {
      "tags" : [ "testing" ],
      "views": 0
   }
}
*/

// Do sends a general request to ElasticSearch server
func (c *ESClient) Do(method string, path string, query []byte) (string, error) {
	body := bytes.NewBuffer(query)
	client := http.Client{}
	req, err := http.NewRequest(method, c.server+path, body)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	fmt.Println("STATUS CODE: ", resp.StatusCode)
	ans, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("RESP ALL: ", string(ans))
	var srchResult Result
	err2 := json.Unmarshal(ans, &srchResult)
	fmt.Println("ERR: ", err2)
	fmt.Println("RESPONSE: ", srchResult)
	return "foo", err2
}
