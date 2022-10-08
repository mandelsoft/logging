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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/logging"
	"github.com/tonglil/buflogr"
)

var pkg = logging.Package()

var _ = Describe("context test", func() {
	var buf bytes.Buffer
	var ctx logging.Context

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
V[4] realm test debug
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
V[4] prefix/realm test prefix debug
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
V[4] realm test debug value attr test
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
V[4] realm test debug
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
V[4] realm test debug value attr test
V[4] realm test debug value attr other
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
V[4] realm test debug value attr test attr other
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
V[4] realm test debug
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
V[4] realm test debug
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
				Expect(nested.Tree().GetWatermark()).To((Equal(int64(0))))
			})

			It("plain", func() {
				ctx.SetDefaultLevel(9)
				Expect(nested.GetDefaultLevel()).To((Equal(9)))
				Expect(nested.Tree().GetWatermark()).To((Equal(int64(1))))
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
})

type X struct{}

func (x *X) Package() logging.Realm {
	return func() logging.Realm {
		return logging.Package()
	}()
}
