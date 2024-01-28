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

package logrusfmt_test

import (
	"bytes"

	"github.com/mandelsoft/logging/logrusl"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tonglil/buflogr"

	"github.com/mandelsoft/logging"
)

var pkg = logging.Package()

var _ = Describe("text formatter test", func() {
	var buf bytes.Buffer
	var ctx logging.Context

	buflogr.NameSeparator = ":"

	Context("human provider", func() {
		BeforeEach(func() {
			buf.Reset()
			ctx = logrusl.Human().WithWriter(&buf).New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})

		It("formats fields", func() {
			ctx.Logger().Info("ages", "alice", 25, "bob", 26)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    ages alice=25 bob=26\n`))
		})

		It("formats fields with realm", func() {
			ctx.Logger(logging.NewRealm("realm")).Info("ages", "alice", 25, "bob", 26)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    \[realm\] ages alice=25 bob=26\n`))
		})

		It("formats fields with name", func() {
			ctx.Logger().WithName("name").Info("ages", "alice", 25, "bob", 26)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    name ages alice=25 bob=26\n`))
		})

		It("formats fields with realm and name", func() {
			ctx.Logger(logging.NewRealm("realm")).WithName("name").Info("ages", "alice", 25, "bob", 26)
			Expect(buf.String()).To(MatchRegexp(`.{25} info    name \[realm\] ages alice=25 bob=26\n`))
		})

	})

	Context("padded human provider", func() {
		BeforeEach(func() {
			buf.Reset()
			ctx = logrusl.Human(true).WithWriter(&buf).New()
			ctx.SetDefaultLevel(logging.InfoLevel)
		})

		It("formats fields", func() {
			log := ctx.Logger()

			log.Info("ages", "alice", 25, "bob", 26)

			logr1 := ctx.Logger(logging.NewRealm("realm"))
			logr1.Info("message", "alice", 25, "bob", 26)

			log.Info("other", "alice", 25, "bob", 26)

			log.WithName("test").Info("name", "alice", 25, "bob", 26)

			logr1.Info("message", "alice", 25, "bob", 26)
			logr1.WithName("bla").Info("message", "alice", 25, "bob", 26)

			ctx.Logger(logging.NewRealm("test")).WithName("bla").Info("message", "alice", 25, "bob", 26)
			Expect("\n" + buf.String()).To(MatchRegexp(`
.{25} info    ages alice=25 bob=26
.{25} info    \[realm\] message alice=25 bob=26
.{25} info            other alice=25 bob=26
.{25} info    test         name alice=25 bob=26
.{25} info         \[realm\] message alice=25 bob=26
.{25} info    bla  \[realm\] message alice=25 bob=26
.{25} info    bla  \[test\]  message alice=25 bob=26
`))
		})
	})
})
