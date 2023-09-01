/*
 * Copyright 2023 Mandelsoft. All rights reserved.
 *  This file is licensed under the Apache Software License, v. 2 except as noted
 *  otherwise in the LICENSE file
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package logrusfmt_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	me "github.com/mandelsoft/logging/logrusfmt"
)

var _ = Describe("json formatter", func() {

	It("renders fields", func() {
		formatter := me.JSONFormatter{}

		entry := &me.Entry{
			Logger: nil,
			Data: map[string]interface{}{
				"a": map[string]interface{}{
					"field": "value",
				},
				"z": "value2",
			},
			Time:    time.Unix(0, 0),
			Level:   me.InfoLevel,
			Caller:  nil,
			Message: "test message",
			Buffer:  nil,
			Context: nil,
		}

		data, err := formatter.Format(entry)
		Expect(err).To(Succeed())
		Expect(string(data)).To(Equal(`{"time":"1970-01-01T01:00:00+01:00","level":"info","msg":"test message","a":{"field":"value"},"z":"value2"}`))
		m := map[string]interface{}{}
		Expect(json.Unmarshal(data, &m)).To(Succeed())

		formatter.PrettyPrint = true
		data, err = formatter.Format(entry)
		Expect(err).To(Succeed())
		Expect(string(data)).To(Equal(`{
  "time": "1970-01-01T01:00:00+01:00",
  "level": "info",
  "msg": "test message",
  "a": {
    "field": "value"
  },
  "z": "value2"
}
`))
		m = map[string]interface{}{}
		Expect(json.Unmarshal(data, &m)).To(Succeed())
	})

	It("renders custom fixed fields", func() {
		formatter := me.JSONFormatter{
			FixedFields: []string{
				me.FieldKeyTime,
				me.FieldKeyLevel,
				"z",
				me.FieldKeyMsg,
			},
		}

		entry := &me.Entry{
			Logger: nil,
			Data: map[string]interface{}{
				"a": map[string]interface{}{
					"field": "value",
				},
				"z": "value2",
			},
			Time:    time.Unix(0, 0),
			Level:   me.InfoLevel,
			Caller:  nil,
			Message: "test message",
			Buffer:  nil,
			Context: nil,
		}

		data, err := formatter.Format(entry)
		Expect(err).To(Succeed())
		Expect(string(data)).To(Equal(`{"time":"1970-01-01T01:00:00+01:00","level":"info","z":"value2","msg":"test message","a":{"field":"value"}}`))
		m := map[string]interface{}{}
		Expect(json.Unmarshal(data, &m)).To(Succeed())

		formatter.PrettyPrint = true
		data, err = formatter.Format(entry)
		Expect(err).To(Succeed())
		Expect(string(data)).To(Equal(`{
  "time": "1970-01-01T01:00:00+01:00",
  "level": "info",
  "z": "value2",
  "msg": "test message",
  "a": {
    "field": "value"
  }
}
`))
		m = map[string]interface{}{}
		Expect(json.Unmarshal(data, &m)).To(Succeed())
	})
})
