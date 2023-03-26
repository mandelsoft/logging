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

package logging

import (
	"reflect"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type nonReplacableRule struct {
	Rule
}

func (r *nonReplacableRule) MatchRule(o Rule) (result bool) {
	return false
}

func NewNonReplacableRule(level int, tag Tag) Rule {
	return &nonReplacableRule{NewConditionRule(level, tag)}
}

var _ = ginkgo.Describe("slice test", func() {

	ginkgo.It("be sure comparing functions does not panic", func() {
		f := func() {}
		gomega.Expect(reflect.DeepEqual(f, f)).To(gomega.BeFalse())
	})

	ginkgo.It("appends new rules", func() {
		ctx := NewWithBase(nil)
		ctx.AddRule(NewConditionRule(DebugLevel, NewTag("test")))
		ctx.AddRule(NewConditionRule(TraceLevel, NewTag("other")))
		gomega.Expect(len(ctx.(*context).rules)).To(gomega.Equal(2))
		l := ctx.Logger(NewTag("test"))
		gomega.Expect(l.Enabled(DebugLevel)).To(gomega.BeTrue())
		gomega.Expect(l.Enabled(TraceLevel)).To(gomega.BeFalse())
	})

	ginkgo.It("gets rid of matching old", func() {
		ctx := NewWithBase(nil)
		ctx.AddRule(NewConditionRule(DebugLevel, NewTag("test")))
		ctx.AddRule(NewConditionRule(TraceLevel, NewTag("other")))
		ctx.AddRule(NewConditionRule(TraceLevel, NewTag("test")))
		gomega.Expect(len(ctx.(*context).rules)).To(gomega.Equal(2))
		l := ctx.Logger(NewTag("test"))
		gomega.Expect(l.Enabled(DebugLevel)).To(gomega.BeTrue())
		gomega.Expect(l.Enabled(TraceLevel)).To(gomega.BeTrue())
	})

	ginkgo.It("override non-replacable old rules", func() {
		ctx := NewWithBase(nil)
		ctx.AddRule(NewNonReplacableRule(DebugLevel, NewTag("test")))
		ctx.AddRule(NewConditionRule(TraceLevel, NewTag("other")))
		ctx.AddRule(NewConditionRule(TraceLevel, NewTag("test")))
		gomega.Expect(len(ctx.(*context).rules)).To(gomega.Equal(3))
		l := ctx.Logger(NewTag("test"))
		gomega.Expect(l.Enabled(DebugLevel)).To(gomega.BeTrue())
		gomega.Expect(l.Enabled(TraceLevel)).To(gomega.BeTrue())
	})
})
