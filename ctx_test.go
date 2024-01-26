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

package logging_test

import (
	"bytes"
	"fmt"

	"github.com/mandelsoft/logging/logrusl"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/tonglil/buflogr"

	"github.com/mandelsoft/logging"
)

var pkg = logging.Package()

var _ = Describe("context test", func() {
	var buf bytes.Buffer
	var ctx logging.Context

	buflogr.NameSeparator = ":"

	BeforeEach(func() {
		buf.Reset()
		def := buflogr.NewWithBuffer(&buf)

		ctx = logging.New(def)
	})

	It("just default logs", func() {
		ctx.Logger().Trace("trace")
		ctx.Logger().Debug("debug")
		ctx.Logger().Info("info")
		ctx.Logger().Warn("warn")
		ctx.Logger().Error("error")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[3] info
V[2] warn
ERROR <nil> error
`))
	})

	It("just default logs with package", func() {
		ctx.Logger(pkg).Trace("trace")
		ctx.Logger(pkg).Debug("debug")
		ctx.Logger(pkg).Info("info")
		ctx.Logger(pkg).Warn("warn")
		ctx.Logger(pkg).Error("error")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[3] info realm github.com/mandelsoft/logging_test
V[2] warn realm github.com/mandelsoft/logging_test
ERROR <nil> error realm github.com/mandelsoft/logging_test
`))
	})

	It("just logs", func() {
		ctx.SetDefaultLevel(9)
		ctx.Logger().Trace("trace")
		ctx.Logger().Debug("debug")
		ctx.Logger().Info("info")
		ctx.Logger().Warn("warn")
		ctx.Logger().Error("error")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[5] trace
V[4] debug
V[3] info
V[2] warn
ERROR <nil> error
`))
	})

	It("just logs with level rule", func() {
		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))

		ctx.Logger().Trace("trace")
		ctx.Logger().Debug("debug")
		ctx.Logger().Info("info")
		ctx.Logger().Warn("warn")
		ctx.Logger().Error("error")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[2] warn
ERROR <nil> error
`))
	})

	It("just logs with context realm rule", func() {
		realm := logging.NewRealm("realm")

		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
		ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

		ctx.Logger().Trace("trace")
		ctx.Logger().Debug("debug")
		ctx.Logger().Info("info")
		ctx.Logger().Warn("warn")
		ctx.Logger().Error("error")
		ctx.Logger(realm).Debug("test debug")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[2] warn
ERROR <nil> error
V[4] test debug realm realm
`))
	})

	It("just logs with context realm prefix rule", func() {
		realm := logging.NewRealm("prefix/realm")
		prefix := logging.NewRealmPrefix("prefix")

		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
		ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, prefix))

		ctx.Logger().Trace("trace")
		ctx.Logger().Debug("debug")
		ctx.Logger().Info("info")
		ctx.Logger().Warn("warn")
		ctx.Logger().Error("error")
		ctx.Logger(realm).Debug("test prefix debug")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[2] warn
ERROR <nil> error
V[4] test prefix debug realm prefix/realm
`))
	})

	It("just logs with context value rule", func() {
		realm := logging.NewRealm("realm")
		attr := logging.NewAttribute("attr", "test")

		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
		ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm, attr))

		ctx.Logger(realm).Debug("test debug")
		ctx.Logger(realm, attr).Debug("test debug value")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[4] test debug value realm realm attr test
`))
	})

	It("just logs with context NOT value rule", func() {
		realm := logging.NewRealm("realm")
		attr := logging.NewAttribute("attr", "test")

		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
		ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm, logging.Not(attr)))

		ctx.Logger(realm).Debug("test debug")
		ctx.Logger(realm, attr).Debug("test debug value")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[4] test debug realm realm
`))
	})

	It("just logs with context OR value rule", func() {
		realm := logging.NewRealm("realm")
		attr := logging.NewAttribute("attr", "test")
		attr2 := logging.NewAttribute("attr", "other")

		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
		ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm, logging.Or(attr2, attr)))

		ctx.Logger(realm).Debug("test debug")
		ctx.Logger(realm, attr).Debug("test debug value")
		ctx.Logger(realm, attr2).Debug("test debug value")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[4] test debug value realm realm attr test
V[4] test debug value realm realm attr other
`))
	})

	It("just logs with context AND value rule", func() {
		realm := logging.NewRealm("realm")
		attr := logging.NewAttribute("attr", "test")
		attr2 := logging.NewAttribute("attr", "other")

		ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
		ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm, logging.And(attr2, attr)))

		ctx.Logger(realm).Debug("test debug")
		ctx.Logger(realm, attr).Debug("test debug value")
		ctx.Logger(realm, attr2).Debug("test debug value")
		ctx.Logger(realm, attr, attr2).Debug("test debug value")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[4] test debug value realm realm attr test attr other
`))
	})

	It("handles regular error level", func() {
		ctx.Logger().V(logging.ErrorLevel).Info("error")

		fmt.Printf("%s\n", buf.String())

		Expect("\n" + buf.String()).To(Equal(`
V[1] error
`))
	})

	Context("package realm", func() {
		local := "github.com/mandelsoft/logging_test"

		It("provides package", func() {
			Expect(pkg.Name()).To(Equal(local))
		})

		It("provides package in function expression", func() {
			Expect(logging.Package().Name()).To(Equal("github.com/mandelsoft/logging_test"))
		})

		It("provides package in struct nested function", func() {
			Expect((&X{}).Package().Name()).To(Equal("github.com/mandelsoft/logging_test"))
		})
	})

	Context("with realm and name", func() {
		It("adds realm and name", func() {
			realm := logging.NewRealm("test")
			ctx = ctx.WithContext(logging.NewName("name"))
			enriched := ctx.WithContext(realm, logging.NewName("sub"))
			ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			enriched.Logger().Debug("debug pkg")
			enriched.Logger().WithName("cmd").Info("info pkg")
			ctx.Logger().Debug("debug plain")
			ctx.Logger().Info("info plain")

			fmt.Printf("\n" + buf.String())
			Expect("\n" + buf.String()).To(Equal(`
V[4] name:sub debug pkg realm test
V[3] name:sub:cmd info pkg realm test
V[3] name info plain
`))
		})
	})

	Context("with standard message context", func() {
		It("adds standard message context", func() {
			realm := logging.NewRealm("test")
			enriched := ctx.WithContext(realm)
			ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			enriched.Logger().Debug("debug pkg")
			enriched.Logger().Info("info pkg")
			ctx.Logger().Debug("debug plain")
			ctx.Logger().Info("info plain")

			fmt.Printf("\n" + buf.String())
			Expect("\n" + buf.String()).To(Equal(`
V[4] debug pkg realm test
V[3] info pkg realm test
V[3] info plain
`))
		})

		It("adds nested standard message context", func() {
			realm := logging.NewRealm("test")
			enriched := ctx.WithContext(realm)
			enriched2 := enriched.WithContext(logging.NewRealm("detail"))
			ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			enriched2.Logger().Debug("debug detail")
			enriched2.Logger().Info("info detail")
			enriched.Logger().Debug("debug pkg")
			enriched.Logger().Info("info pkg")
			ctx.Logger().Debug("debug plain")
			ctx.Logger().Info("info plain")

			fmt.Printf("\n" + buf.String())
			Expect("\n" + buf.String()).To(Equal(`
V[3] info detail realm test realm detail
V[4] debug pkg realm test
V[3] info pkg realm test
V[3] info plain
`))
		})
	})

	Context("nested contexts", func() {
		var nested logging.Context

		BeforeEach(func() {
			nested = logging.NewWithBase(ctx)
		})

		It("just logs with context realm rule", func() {
			realm := logging.NewRealm("realm")

			ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
			ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			nested.Logger().Trace("trace")
			nested.Logger().Debug("debug")
			nested.Logger().Info("info")
			nested.Logger().Warn("warn")
			nested.Logger().Error("error")
			nested.Logger(realm).Debug("test debug")

			fmt.Printf("%s\n", buf.String())

			Expect("\n" + buf.String()).To(Equal(`
V[2] warn
ERROR <nil> error
V[4] test debug realm realm
`))
		})

		It("just logs with nested context realm rule", func() {
			realm := logging.NewRealm("realm")

			ctx.AddRule(logging.NewConditionRule(logging.WarnLevel))
			nested.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			nested.Logger().Trace("trace")
			nested.Logger().Debug("debug")
			nested.Logger().Info("info")
			nested.Logger().Warn("warn")
			nested.Logger().Error("error")
			nested.Logger(realm).Debug("test debug")
			ctx.Logger(realm).Debug("ctx test debug")

			fmt.Printf("%s\n", buf.String())

			Expect("\n" + buf.String()).To(Equal(`
V[2] warn
ERROR <nil> error
V[4] test debug realm realm
`))
		})

		Context("dynamic changes", func() {
			It("adapts provided default logger", func() {
				logger := nested.Logger()
				logger.Info("info before")
				logger.Debug("debug before")

				ctx.SetDefaultLevel(logging.DebugLevel)
				logger.Info("info after")
				logger.Debug("debug after")

				Expect("\n" + buf.String()).To(Equal(`
V[3] info before
V[3] info after
V[4] debug after
`))
			})
			It("adapts provided sink", func() {
				var buf bytes.Buffer
				def := buflogr.NewWithBuffer(&buf)
				ctx.SetBaseLogger(def)

				logger := nested.Logger()
				logger.Info("info before")
				logger.Debug("debug before")

				ctx.SetDefaultLevel(logging.DebugLevel)
				logger.Info("info after")
				logger.Debug("debug after")

				Expect("\n" + buf.String()).To(Equal(`
V[3] info before
V[3] info after
V[4] debug after
`))
			})
			It("adapts provided ruke based logger", func() {
				tag := logging.NewTag("tag")

				nested.AddRule(logging.NewConditionRule(logging.DebugLevel, tag))
				logger := nested.Logger(tag)

				var buf bytes.Buffer
				def := buflogr.NewWithBuffer(&buf)
				ctx.SetBaseLogger(def)

				logger.Info("info before")
				logger.Debug("debug before")

				ctx.SetDefaultLevel(logging.DebugLevel)
				logger.Info("info after")
				logger.Debug("debug after")

				Expect("\n" + buf.String()).To(Equal(`
V[3] info before
V[4] debug before
V[3] info after
V[4] debug after
`))
			})
		})

		Context("update nested contexts", func() {
			It("plain", func() {
				Expect(nested.Tree().Updater().SeenWatermark()).To((Equal(int64(0))))
				Expect(nested.Tree().Updater().SeenWatermark()).To(Equal(nested.Tree().Updater().Watermark()))
			})

			It("plain", func() {
				ctx.SetDefaultLevel(9)
				Expect(nested.GetDefaultLevel()).To((Equal(9)))
				Expect(nested.Tree().Updater().SeenWatermark()).To((Equal(int64(1))))
				Expect(nested.Tree().Updater().SeenWatermark()).To(Equal(nested.Tree().Updater().Watermark()))
			})

		})
	})

	Context("shifted", func() {
		var shifted logging.Context

		BeforeEach(func() {
			shifted = logging.New(buflogr.NewWithBuffer(&buf).V(2))
		})

		It("shift initial level", func() {
			shifted.Logger().Info("info")
			Expect("\n" + buf.String()).To(Equal(`
V[5] info
`))
		})
		It("keeps error", func() {
			shifted.Logger().Error("error")
			Expect("\n" + buf.String()).To(Equal(`
ERROR <nil> error
`))
		})

		It("disables error", func() {
			shifted.SetDefaultLevel(0)
			shifted.Logger().Error("error")
			Expect(buf.String()).To(Equal(""))
		})
	})

	Context("deriving standard message context", func() {
		It("inherits context", func() {
			name := logging.NewName("test")
			enriched := ctx.WithContext(name)

			eff := logging.NewWithBase(enriched)

			eff.Logger().Error("with realm")
			Expect(buf.String()).To(Equal("ERROR <nil> test with realm\n"))
		})

		It("extended message context", func() {
			name := logging.NewName("test")
			enriched := ctx.WithContext(name)

			eff := logging.NewWithBase(enriched)

			ext := logging.NewName("extended")

			eff.Logger(ext).Error("with realm")
			Expect(buf.String()).To(Equal("ERROR <nil> test:extended with realm\n"))
		})
	})

	Context("logrus human", func() {

		BeforeEach(func() {
			ctx = logrusl.WithWriter(&buf).Human().New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})

		It("regular", func() {
			ctx.Logger().Info("test", logrus.FieldKeyLevel, 1000)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    test fields.level=1000\n`))
		})

		It("substitution", func() {
			ctx.Logger().Info("test {{value}}", "value", "test")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "test test"`))

			buf.Reset()
			ctx.Logger().Info("test {{value}}", "value", 4711)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "test 4711"`))

			buf.Reset()
			ctx.Logger().Info("test {{value}}", "value", true)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "test true"`))

			buf.Reset()
			ctx.Logger().Info("test {{value}}", "value", struct {
				Field1 string `json:"field1"`
				Field2 string `json:"field2"`
			}{
				Field1: "value1",
				Field2: "value2",
			})
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "test {\\"field1\\":\\"value1\\",\\"field2\\":\\"value2\\"}"`))
		})

	})

	Context("name and realm handling", func() {
		BeforeEach(func() {
			ctx = logrusl.Human().WithWriter(&buf).New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})

		It("adds name to logger name", func() {
			enriched := ctx.WithContext(logging.NewName("test")).WithContext(logging.NewName("nested"))
			enriched.Logger().Error("with name")
			Expect(buf.String()).To(MatchRegexp(`.{25} error   test.nested "with name" error=<nil>\n`))

			buf.Reset()
			enriched.Logger(logging.NewName("sub")).Error("with name")
			Expect(buf.String()).To(MatchRegexp(`.{25} error   test.nested.sub "with name" error=<nil>\n`))

			buf.Reset()
			enriched.Logger(logging.Realm("absolute")).Error("with name")
			Expect(buf.String()).To(MatchRegexp(`.{25} error   test.nested \[absolute\] "with name" error=<nil>\n`))
		})

		It("adds realm as key/value pair", func() {
			enriched := ctx.WithContext(logging.NewRealm("test")).WithContext(logging.NewRealm("other"), logging.NewName("test"))
			enriched.Logger().Error("with realm")
			Expect(buf.String()).To(MatchRegexp(`.{25} error   test \[other\] "with realm" error=<nil>\n`))

			buf.Reset()
			enriched.Logger(logging.NewRealm("sub")).Error("with realm")
			Expect(buf.String()).To(MatchRegexp(`.{25} error   test \[sub\] "with realm" error=<nil>\n`))

			buf.Reset()
			enriched.Logger(logging.NewRealm("sub"), logging.NewName("sub")).Error("with realm")
			Expect(buf.String()).To(MatchRegexp(`.{25} error   test.sub \[sub\] "with realm" error=<nil>\n`))
		})
	})

	Context("fields", func() {
		BeforeEach(func() {
			ctx = logrusl.Human().WithWriter(&buf).New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})

		It("logs simple attribute values", func() {
			enriched := ctx.WithContext(logging.NewAttribute("attr", "value"))
			enriched.Logger().Info("with attribute")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with attribute" attr=value\n`))
		})
		It("logs value attribute pointer stringer", func() {
			enriched := ctx.WithContext(logging.NewAttribute("attr", &Pointer{"value"}))
			enriched.Logger().Info("with attribute")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with attribute" attr=value\n`))
		})
		It("logs value attribute values", func() {
			enriched := ctx.WithContext(logging.NewAttribute("attr", Value{"value"}))
			enriched.Logger().Info("with attribute")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with attribute" attr=value\n`))
		})
		It("logs value simple underlying type", func() {
			enriched := ctx.WithContext(logging.NewAttribute("attr", Underlying("value")))
			enriched.Logger().Info("with attribute")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with attribute" attr=value\n`))
		})
		It("logs value simple underlying type with substitution", func() {
			enriched := ctx.WithContext(logging.NewAttribute("attr", Underlying("value")))
			enriched.Logger().Info("with attribute value {{attr}}")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with attribute value value"\n`))
		})
		It("logs value simple underlying type with substitution", func() {
			enriched := ctx.WithContext(logging.NewAttribute("attr", fmt.Errorf("value")))
			enriched.Logger().Info("with attribute value {{attr}}")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with attribute value value"\n`))
		})
	})

	Context("fields values", func() {
		BeforeEach(func() {
			ctx = logrusl.Human().WithWriter(&buf).New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})
		It("key value", func() {
			log := ctx.Logger()
			log.Info("with key value", logging.KeyValue("name", "value"))
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with key value" name=value\n`))
		})
		It("key value sequence", func() {
			log := ctx.Logger()
			log.Info("with key value", logging.KeyValue("name", "value"), logging.KeyValue("other", "othervalue"))
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with key value" name=value other=othervalue\n`))
		})
		It("mixed key value sequence", func() {
			log := ctx.Logger()
			log.Info("with key value", logging.KeyValue("name", "value"), "other", "othervalue")
			log.Info("with key value", "first", "firstvalue", logging.KeyValue("name", "value"), "other", "othervalue")
			Expect(buf.String()).To(MatchRegexp(`.{25} info    "with key value" name=value other=othervalue\n.{25} info    "with key value" first=firstvalue name=value other=othervalue\n`))
		})
	})
})

type Value struct {
	msg string
}

type Pointer struct {
	msg string
}

type Underlying string

func (v Value) String() string {
	return v.msg
}

func (v *Pointer) String() string {
	return v.msg
}

type X struct{}

func (x *X) Package() logging.Realm {
	return func() logging.Realm {
		return logging.Package()
	}()
}
