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

type ESClient struct {
	server string
}

func NewClient(server string) *ESClient {
	c := ESClient{
		server: server,
	}
	return &c
}

func encodeJson(record LogRecord) []byte {
	str, err := json.Marshal(record)
	if err == nil {
		return str
	}
	return nil
}

func (c *ESClient) Test() (string, error) {
	tmp := "first_form?corpname=omezeni%2Fsyn2010&format=json"
	return c.jsonRequest(tmp, LogRecord{})
}

func (c *ESClient) jsonRequest(path string, record LogRecord) (string, error) {
	body := bytes.NewBuffer(encodeJson(record))
	req, err := http.NewRequest("GET", c.server+"/"+path, body)
	client := http.Client{}
	if err == nil {
		resp, err := client.Do(req)
		ans, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			fmt.Println("RESP ALL: ", resp)
			var respMap map[string]string
			json.Unmarshal(ans, respMap)
			fmt.Println("RESPONSE: ", respMap)
			return respMap["id"], nil
		}
	}
	return "", err
}
