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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type Conf struct {
	WorklogPath   string `json:"worklogPath"`
	LogDir        string `json:"logDir"`
	ElasticServer string `json:"elasticServer"`
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic("Invalid arguments")
	}
	rawData, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	var conf Conf
	json.Unmarshal(rawData, &conf)

	worklog, err := LoadWorklog(conf.WorklogPath)
	if err != nil {
		panic(err)
	}
	last := worklog.FindLastRecord()
	fmt.Println(worklog, last)
	files := getFilesInDir(conf.LogDir)
	fmt.Println("FILES: ", files)
	for _, file := range files {
		p := NewParser(file)
		p.Parse(last)
	}

	client := NewClient(conf.ElasticServer)
	ans, err := client.Test()
	fmt.Println("CLIENT ANS: ", ans, err)
}
