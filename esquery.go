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

import "encoding/json"

// ----------------- query --------------------------

type datetimeRangeExpr struct {
	From string `json:"gte"`
	To   string `json:"lt"`
}

type datetimeRangeQuery struct {
	Datetime datetimeRangeExpr `json:"datetime"`
}

type rangeObj struct {
	Range datetimeRangeQuery `json:"range"`
}

type iPAddressExpr struct {
	IPAddress string `json:"ipAddress"`
}

type iPAddressTermObj struct {
	Term iPAddressExpr `json:"term"`
}

type userAgentExpr struct {
	UserAgent string `json:"userAgent"`
}

type userAgentMatchObj struct {
	Match userAgentExpr `json:"match"`
}

type boolObj struct {
	Must []interface{} `json:"must"`
}

type query struct {
	Bool boolObj `json:"bool"`
}

type srchQuery struct {
	Query query `json:"query"`
}

// ----------------- result -------------------------

// ResultHit represents an individual result
type ResultHit struct {
	Index  string      `json:"_index"`
	Type   string      `json:"_type"`
	ID     string      `json:"_id"`
	Score  string      `json:"_score"`
	Source interface{} `json:"_source"`
}

// Hits is a "hits" part of the ElasticSearch query result object
type Hits struct {
	Total    int         `json:"total"`
	MaxScore float32     `json:"max_score"`
	Hits     []ResultHit `json:"hits"`
}

// Result represents an ElasticSearch query result object
type Result struct {
	Took     int         `json:"took"`
	TimedOut bool        `json:"timed_out"`
	Shards   interface{} `json:"_shards"`
	Hits     Hits        `json:"hits"`
}

// CreateQuery generate JSON-encoded query for ElastiSearch to
// find documents matching specified datetime range, optional IP
// address and optional userAgent substring/pattern
func CreateQuery(fromDate string, toDate string, ipAddress string, userAgent string) ([]byte, error) {
	m := boolObj{Must: make([]interface{}, 1)}
	dateInterval := datetimeRangeExpr{From: fromDate, To: toDate}
	m.Must[0] = &rangeObj{Range: datetimeRangeQuery{Datetime: dateInterval}}
	if ipAddress != "" {
		ipAddrObj := iPAddressTermObj{Term: iPAddressExpr{IPAddress: ipAddress}}
		m.Must = append(m.Must, ipAddrObj)
	}
	if userAgent != "" {
		userAgentObj := userAgentMatchObj{userAgentExpr{UserAgent: userAgent}}
		m.Must = append(m.Must, userAgentObj)
	}
	q := srchQuery{Query: query{Bool: m}}
	return json.Marshal(q)
}
