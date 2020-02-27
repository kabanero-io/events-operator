/*
Copyright 2020 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils_test

import (
	"fmt"
	"github.com/kabanero-io/events-operator/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

/* constants*/
const (
	KABANEROINDEXURLTEST = "https://github.com/kabanero-io/collections/releases/download/v0.1.2/kabanero-index.yaml"

	GZIPTAR0 = "../../test_data/gZipTarDir0.tar.gz"
)

func TestReadUrl(t *testing.T) {
	url := KABANEROINDEXURLTEST
	bytes, err := utils.ReadHTTPURL(url)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("URL %s has content: %s", url, string(bytes))
}

var kabaneroIndex = `triggers:
 - description: sample trigger
   url: https://host/trigger-path
   sha256: 270fe2e576132a6fb247755762b4556d6ff215d3507d66643304ea97cb2ec58e
stacks:
- default-image: java-microprofile
  default-pipeline: default
  default-template: default
  description: Eclipse MicroProfile on Open Liberty & OpenJ9 using Maven
  id: java-microprofile
  images:
  - id: java-microprofile
    image: kabanero/java-microprofile:0.2
  language: java
  license: Apache-2.0
  maintainers:
  - email: emijiang6@googlemail.com
    github-id: Emily-Jiang
    name: Emily Jiang
  - email: neeraj.laad@gmail.com
    github-id: neeraj-laad
    name: Neeraj Laad
  name: Eclipse MicroProfile®
  pipelines:
  - id: default
    sha256: a59a779825c543e829d7a51e383f26c2089b4399cf39a89c10b563597d286991
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.common.pipeline.default.tar.gz
  templates:
  - id: default
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.java-microprofile.v0.2.11.templates.default.tar.gz
  version: 0.2.11
- default-image: java-spring-boot2
  default-pipeline: default
  default-template: default
  description: Spring Boot using OpenJ9 and Maven
  id: java-spring-boot2
  images:
  - id: java-spring-boot2
    image: kabanero/java-spring-boot2:0.3
  language: java
  license: Apache-2.0
  maintainers:
  - email: schnabel@us.ibm.com
    github-id: ebullient
    name: Erin Schnabel
  name: Spring Boot®
  pipelines:
  - id: default
    sha256: a59a779825c543e829d7a51e383f26c2089b4399cf39a89c10b563597d286991
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.common.pipeline.default.tar.gz
  templates:
  - id: default
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.java-spring-boot2.v0.3.9.templates.default.tar.gz
  - id: kotlin
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.java-spring-boot2.v0.3.9.templates.kotlin.tar.gz
  version: 0.3.9
- default-image: nodejs-express
  default-pipeline: default
  default-template: simple
  description: Express web framework for Node.js
  id: nodejs-express
  images:
  - id: nodejs-express
    image: kabanero/nodejs-express:0.2
  language: nodejs
  license: Apache-2.0
  maintainers:
  - email: cnbailey@gmail.com
    github-id: seabaylea
    name: Chris Bailey
  - email: neeraj.laad@gmail.com
    github-id: neeraj-laad
    name: Neeraj Laad
  name: Node.js Express
  pipelines:
  - id: default
    sha256: a59a779825c543e829d7a51e383f26c2089b4399cf39a89c10b563597d286991
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.common.pipeline.default.tar.gz
  templates:
  - id: simple
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.nodejs-express.v0.2.5.templates.simple.tar.gz
  - id: skaffold
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.nodejs-express.v0.2.5.templates.skaffold.tar.gz
  version: 0.2.5
- default-image: nodejs-loopback
  default-pipeline: default
  default-template: scaffold
  description: LoopBack 4 API Framework for Node.js
  id: nodejs-loopback
  images:
  - id: nodejs-loopback
    image: kabanero/nodejs-loopback:0.1
  language: nodejs
  license: Apache-2.0
  maintainers:
  - email: enjoyjava@gmail.com
    github-id: raymondfeng
    name: Raymond Feng
  name: LoopBack 4
  pipelines:
  - id: default
    sha256: a59a779825c543e829d7a51e383f26c2089b4399cf39a89c10b563597d286991
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.common.pipeline.default.tar.gz
  templates:
  - id: scaffold
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.nodejs-loopback.v0.1.4.templates.scaffold.tar.gz
  version: 0.1.4
- default-image: nodejs
  default-pipeline: default
  default-template: simple
  description: Runtime for Node.js applications
  id: nodejs
  images:
  - id: nodejs
    image: kabanero/nodejs:0.2
  language: nodejs
  license: Apache-2.0
  maintainers:
  - email: cnbailey@gmail.com
    github-id: seabaylea
    name: Chris Bailey
  - email: neeraj.laad@gmail.com
    github-id: neeraj-laad
    name: Neeraj Laad
  name: Node.js
  pipelines:
  - id: default
    sha256: a59a779825c543e829d7a51e383f26c2089b4399cf39a89c10b563597d286991
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.common.pipeline.default.tar.gz
  templates:
  - id: simple
    url: https://github.com/kabanero-io/collections/releases/download/v0.1.2/incubator.nodejs.v0.2.5.templates.simple.tar.gz
  version: 0.2.5
`

type mergePathData struct {
	dir     string
	toMerge string
	succeed bool
}

var mergedPathTestData = []mergePathData{
	{".", "abc", true},
	{".", "abc/def", true},
	{".", "abc/def/ghi", true},
	{".", "abc/../def", true},
	{".", "..", false},
	{".", "../..", false},
	{".", "/abc", true},  // having '/' will still append to the paath
	{".", "\\abc", true}, // having '\\' will still append to the path
}

func TestMergePathWithErrorCheck(t *testing.T) {

	for _, testData := range mergedPathTestData {
		mergedPath, err := utils.MergePathWithErrorCheck(testData.dir, testData.toMerge)
		succeeded := err == nil
		if testData.succeed != succeeded {
			t.Fatal(fmt.Errorf("unexpected error when merging %s with %s, error: %s, mergedPath: %s", testData.dir, testData.toMerge, err, mergedPath))
		}
	}

}

var gZipTar0Files = []string{
	"eventTriggers.yaml",
	"subdir0/file0",
	"subdir0/file1",
	"subdir1/file0",
	"subdir1/file1",
	"subdir1/file2",
}

func TestGUnzipUnTar(t *testing.T) {
	dir, err := ioutil.TempDir("", "webhook-unittest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.Open(GZIPTAR0)
	err = utils.DecompressGzipTar(file, dir)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the files have been expanded
	for _, fileName := range gZipTar0Files {
		tempFile := filepath.Join(dir, fileName)
		fileInfo, err := os.Stat(tempFile)
		if err != nil {
			t.Fatal(err)
		}
		if fileInfo.IsDir() {
			t.Fatal(fmt.Errorf("file %s is directory", fileName))
		}
	}
}
