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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tonglil/buflogr"

	"github.com/mandelsoft/logging"
)

var _ = Describe("context test", func() {
	var buf bytes.Buffer
	var ctx logging.Context

	buflogr.NameSeparator = ":"

	BeforeEach(func() {
		buf.Reset()
		def := buflogr.NewWithBuffer(&buf)

		ctx = logging.New(def)
	})

	Context("context update handling", func() {
		It("root context", func() {
			ctx.SetDefaultLevel(logging.DebugLevel)
			ctx.SetBaseLogger(buflogr.NewWithBuffer(&buf))
			u := ctx.Tree().Updater()
			Expect(u.Watermark()).To(Equal(int64(2)))
			Expect(u.SeenWatermark()).To(Equal(u.Watermark()))
		})

		It("nested context", func() {
			nctx := logging.NewWithBase(ctx)
			ctx.SetDefaultLevel(logging.DebugLevel)
			ctx.SetBaseLogger(buflogr.NewWithBuffer(&buf))
			ur := ctx.Tree().Updater()
			Expect(ur.Watermark()).To(Equal(int64(2)))
			Expect(ur.SeenWatermark()).To(Equal(ur.Watermark()))

			// modify leaf
			nctx.Logger().Info("test")
			ul := nctx.Tree().Updater()
			Expect(nctx.GetDefaultLevel()).To(Equal(logging.DebugLevel))
			Expect(ul.Watermark()).To(Equal(int64(2)))
			Expect(ul.SeenWatermark()).To(Equal(ul.Watermark()))

			// modify root
			nctx.SetDefaultLevel(logging.TraceLevel)
			nctx.Logger().Info("test")
			Expect(ur.Watermark()).To(Equal(int64(2)))
			Expect(ur.SeenWatermark()).To(Equal(ur.Watermark()))
			Expect(nctx.GetDefaultLevel()).To(Equal(logging.TraceLevel))
			Expect(ul.Watermark()).To(Equal(int64(3)))
			Expect(ul.SeenWatermark()).To(Equal(ur.Watermark()))

			// modify leaf - root
			nctx.SetDefaultLevel(logging.TraceLevel)
			Expect(ul.Watermark()).To(Equal(int64(4)))
			Expect(ur.Watermark()).To(Equal(int64(2)))
			ctx.SetDefaultLevel(logging.InfoLevel)
			Expect(ul.Watermark()).To(Equal(int64(5)))
			Expect(ul.SeenWatermark()).To(Equal(int64(2))) // not yet updated
			nctx.Logger().Info("test")
			Expect(nctx.GetDefaultLevel()).To(Equal(logging.TraceLevel))
			Expect(ul.Watermark()).To(Equal(int64(5)))
			Expect(ul.SeenWatermark()).To(Equal(ur.Watermark()))
			Expect(ur.Watermark()).To(Equal(int64(5)))
			Expect(ur.SeenWatermark()).To(Equal(ur.Watermark()))

			// modify root - leaf
			ctx.SetDefaultLevel(logging.InfoLevel)
			Expect(ul.Watermark()).To(Equal(int64(6)))
			Expect(ul.SeenWatermark()).To(Equal(int64(5)))
			Expect(ur.Watermark()).To(Equal(int64(6)))
			Expect(ur.SeenWatermark()).To(Equal(int64(6)))
			nctx.SetDefaultLevel(logging.TraceLevel)
			Expect(ul.Watermark()).To(Equal(int64(7)))
			Expect(ul.SeenWatermark()).To(Equal(int64(5))) // not yet updated
			nctx.Logger().Info("test")
			Expect(nctx.GetDefaultLevel()).To(Equal(logging.TraceLevel))
			Expect(ul.Watermark()).To(Equal(int64(7)))
			Expect(ur.Watermark()).To(Equal(int64(6)))
			Expect(ul.SeenWatermark()).To(Equal(ur.Watermark()))
			Expect(ur.SeenWatermark()).To(Equal(ur.Watermark()))
		})
	})

	Context("bound loggers", func() {
		It("level update", func() {
			logger := ctx.Logger()

			logger.Debug("debug")
			Expect("\n" + buf.String()).To(Equal(`
`))

			ctx.SetDefaultLevel(logging.DebugLevel)
			logger.Debug("debug")
			fmt.Printf("%s\n", buf.String())
			Expect("\n" + buf.String()).To(Equal(`
V[4] debug
`))
		})

		It("no level update with matching rule", func() {
			realm := logging.NewRealm("realm")
			ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			logger := ctx.Logger(realm)
			logger.Debug("debug")
			logger.Trace("trace")
			fmt.Printf("%s\n", buf.String())
			Expect("\n" + buf.String()).To(Equal(`
V[4] debug realm realm
`))
			// do not use default level, but fixed rule level
			buf.Reset()
			ctx.SetDefaultLevel(logging.TraceLevel)
			logger.Trace("trace")
			Expect("\n" + buf.String()).To(Equal(`
`))

			// change rule level, but instantiated logger should not reflect this
			ctx.AddRule(logging.NewConditionRule(logging.TraceLevel, realm))
			buf.Reset()
			logger.Trace("trace")
			Expect("\n" + buf.String()).To(Equal(`
`))
		})
	})

	Context("unbound loggers", func() {
		It("level update", func() {
			logger := logging.DynamicLogger(ctx)

			logger.Debug("debug")
			Expect("\n" + buf.String()).To(Equal(`
`))

			ctx.SetDefaultLevel(logging.DebugLevel)
			logger.Debug("debug")
			fmt.Printf("%s\n", buf.String())
			Expect("\n" + buf.String()).To(Equal(`
V[4] debug
`))
		})

		DescribeTable("no level update with matching rule", func(names string, values []interface{}) {
			realm := logging.NewRealm("realm")
			ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, realm))

			var logger logging.Logger = logging.DynamicLogger(ctx, realm)
			if len(names) > 0 {
				for _, n := range strings.Split(names, ":") {
					logger = logger.WithName(n)
				}
				names = " " + names
			}
			logger = logger.WithValues(values...)
			vstring := ""
			if len(values) > 0 {
				vstring = " " + strings.Trim(fmt.Sprintf("%v", values), "[]")
			}

			logger.Debug("debug")
			logger.Trace("trace")
			fmt.Printf("%s\n", buf.String())
			Expect("\n" + buf.String()).To(Equal(fmt.Sprintf(`
V[4]%s debug realm realm%s
`, names, vstring)))
			// do not use default level, but fixed rule level
			buf.Reset()
			ctx.SetDefaultLevel(logging.TraceLevel)
			logger.Trace("trace")
			Expect("\n" + buf.String()).To(Equal(`
`))

			// change rule level, but instantiated logger should not reflect this
			ctx.AddRule(logging.NewConditionRule(logging.TraceLevel, realm))
			buf.Reset()
			logger.Trace("trace")
			Expect("\n" + buf.String()).To(Equal(fmt.Sprintf(`
V[5]%s trace realm realm%s
`, names, vstring)))
		},
			Entry("without name", nil, nil),
			Entry("with single name", "name", nil),
			Entry("with multiple names", "name1:name2", nil),
			Entry("with values", nil, []interface{}{"arg", "value"}),
			Entry("with name and values", "name", []interface{}{"arg", "value"}),
		)
	})

})
