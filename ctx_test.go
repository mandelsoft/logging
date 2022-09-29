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

	"github.com/mandelsoft/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tonglil/buflogr"
)

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
V[1] error
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
V[1] error
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
V[1] error
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
V[1] error
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
V[1] error
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

})
