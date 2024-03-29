/*
 * Copyright 2024 Mandelsoft. All rights reserved.
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

	"github.com/mandelsoft/logging"
	"github.com/mandelsoft/logging/logrusl"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Attribution Test", func() {
	Context("", func() {
		var ctx logging.Context
		var buf bytes.Buffer

		BeforeEach(func() {
			buf.Reset()
			ctx = logrusl.Human().WithWriter(&buf).New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})

		It("attribution with realm", func() {
			ctx.Logger().Info("plain")
			attr := ctx.AttributionContext().WithContext(logging.NewRealm("test"))

			attr.Logger().Info("info")
			Expect(buf.String()).To(MatchRegexp(`^.{25} info    plain
.{25} info    \[test\] info
$`))
		})

		It("attribution with values", func() {
			ctx.Logger().Info("plain")
			attr := ctx.AttributionContext().WithValues("foo", "bar")

			attr.Logger().Info("info")
			Expect(buf.String()).To(MatchRegexp(`^.{25} info    plain
.{25} info    info foo=bar
$`))
		})

		It("attribution with name", func() {
			ctx.Logger().Info("plain")
			attr := ctx.AttributionContext().WithName("foobar")
			attr2 := ctx.AttributionContext().WithName("other")

			attr.Logger().Info("info")
			attr2.Logger().Info("info")
			Expect(buf.String()).To(MatchRegexp(`^.{25} info    plain
.{25} info    foobar info
.{25} info    other info
$`))
		})

		It("attribution with string context", func() {
			ctx.Logger().Info("plain")
			attr := ctx.AttributionContext().WithContext("foobar")

			attr.Logger().Info("info")
			Expect(buf.String()).To(MatchRegexp(`^.{25} info    plain
.{25} info    info message=foobar
$`))
		})
	})
})
