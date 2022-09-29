/*
 * Copyright 2022 Mandelsoft. All rights reserved.
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

package config_test

import (
	"bytes"

	"github.com/mandelsoft/logging"
	"github.com/mandelsoft/logging/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tonglil/buflogr"
)

type Test struct {
	A string
	B string
}

type TestType struct {
	A string
	B string
}

func (r *TestType) Create(_ config.Registry) (interface{}, error) {
	return &Test{
		A: r.A,
		B: r.B,
	}, nil
}

func init() {
	config.RegisterValueType("logging.test", &TestType{})
}

var _ = Describe("externaliized data", func() {
	reg := config.DefaultRegistry()

	Context("values", func() {
		It("deserializes generic value", func() {
			data := `
value:
  x: 1
`
			v, err := reg.CreateValue([]byte(data))
			Expect(err).To(Succeed())
			m, ok := v.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(m)).To(Equal(1))
		})

		It("deserializes generic", func() {
			data := `
logging.test:
  A: a
  B: b
`
			v, err := reg.CreateValue([]byte(data))
			Expect(err).To(Succeed())
			m, ok := v.(*Test)
			Expect(ok).To(BeTrue())
			Expect(m).To(Equal(&Test{
				A: "a",
				B: "b",
			}))
		})
	})

	Context("conditions", func() {

		It("deserializes tag", func() {
			data := `
tag: test
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(logging.Tag)
			Expect(ok).To(BeTrue())
			Expect(c.Name()).To(Equal("test"))
		})

		It("deserializes realm", func() {
			data := `
realm:
  test
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(logging.Realm)
			Expect(ok).To(BeTrue())
			Expect(c.Name()).To(Equal("test"))
		})

		It("deserializes realm prefix", func() {
			data := `
realmprefix:
  test
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(logging.RealmPrefix)
			Expect(ok).To(BeTrue())
			Expect(c.Name()).To(Equal("test"))
		})

		It("deserializes and", func() {
			data := `
and:
- tag: test
- realm: mine
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(*logging.AndExpr)
			Expect(ok).To(BeTrue())
			Expect(c.Conditions()).To(Equal([]logging.Condition{logging.NewTag("test"), logging.NewRealm("mine")}))
		})

		It("deserializes or", func() {
			data := `
or:
- tag: test
- realm: mine
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(*logging.OrExpr)
			Expect(ok).To(BeTrue())
			Expect(c.Conditions()).To(Equal([]logging.Condition{logging.NewTag("test"), logging.NewRealm("mine")}))
		})

		It("deserializes not", func() {
			data := `
not:
  tag: test
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(*logging.NotExpr)
			Expect(ok).To(BeTrue())
			Expect(c.Condition()).To(Equal(logging.NewTag("test")))
		})

		It("deserializes attribute", func() {
			data := `
attribute:
  name: test
  value:
    logging.test:
      A: a
      B: b
`
			cond, err := reg.CreateCondition([]byte(data))
			Expect(err).To(Succeed())
			c, ok := cond.(logging.Attribute)
			Expect(ok).To(BeTrue())
			Expect(c.Name()).To(Equal("test"))
			Expect(c.Value()).To(Equal(&Test{
				A: "a",
				B: "b",
			}))
		})
	})

	Context("rules", func() {
		It("deserializes rule", func() {
			data := `
rule:
  level: Warn
  conditions:
    - realm: test
`
			rule, err := reg.CreateRule([]byte(data))
			Expect(err).To(Succeed())
			r, ok := rule.(*logging.ConditionRule)
			Expect(ok).To(BeTrue())
			Expect(r.Level()).To(Equal(logging.WarnLevel))
			Expect(r.Conditions()).To(Equal([]logging.Condition{logging.NewRealm("test")}))
		})
	})

	Context("configure", func() {
		It("configures context", func() {
			var buf bytes.Buffer

			buf.Reset()
			def := buflogr.NewWithBuffer(&buf)

			ctx := logging.New(def)
			data := `
defaultLevel: Warn
rules:
- rule:
    level: Debug
    conditions:
    - realm: test
`
			err := reg.Configure(ctx, []byte(data))
			Expect(err).To(Succeed())

			ctx.Logger().Debug("debug")
			ctx.Logger(logging.NewRealm("test")).Debug("debug")

			Expect("\n" + buf.String()).To(Equal(`
V[4] test debug
`))
		})

		It("configures context", func() {
			var buf bytes.Buffer

			buf.Reset()
			def := buflogr.NewWithBuffer(&buf)

			ctx := logging.New(def)
			data := `
defaultLevel: Warn
rules:
  - rule:
      level: Trace
      conditions:
        - attribute:
            name: test
            value:
               value: testvalue
`
			err := reg.Configure(ctx, []byte(data))
			Expect(err).To(Succeed())

			ctx.Logger().Trace("debug")
			ctx.Logger(logging.NewAttribute("test", "testvalue")).Trace("debug")

			Expect("\n" + buf.String()).To(Equal(`
V[5] debug test testvalue
`))
		})
	})
})
