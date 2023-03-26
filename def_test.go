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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/logging"
)

var _ = Describe("definitions", func() {

	It("defines tags and realms", func() {
		logging.DefineTag("tag1", "")
		logging.DefineTag("tag1", "tag1 desc 2")
		logging.DefineTag("tag1", "tag1 desc 1")
		logging.DefineTag("tag2", "tag2 desc 1")
		logging.DefineTag("tag2", "")
		logging.DefineTag("tag1", "")

		logging.DefineRealm("realm1", "")
		logging.DefineRealm("realm1", "realm1 desc 2")
		logging.DefineRealm("realm1", "realm1 desc 1")
		logging.DefineRealm("realm2", "realm2 desc 1")
		logging.DefineRealm("realm2", "")
		logging.DefineRealm("realm1", "")

		Expect(logging.GetTagDefinitions()).To(Equal(logging.Definitions{
			"tag1": []string{"tag1 desc 1", "tag1 desc 2"},
			"tag2": []string{"tag2 desc 1"},
		}))
		Expect(logging.GetRealmDefinitions()).To(Equal(logging.Definitions{
			"realm1": []string{"realm1 desc 1", "realm1 desc 2"},
			"realm2": []string{"realm2 desc 1"},
		}))
	})
})
