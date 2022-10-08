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
)

var _ = Describe("slice test", func() {

	It("provides append and copied slice", func() {

		orig := make([]int, 0, 10)

		init := append(orig, 5, 6, 7)
		tmp := append(orig, 0, 1, 2, 3, 4)

		Expect(init).To(Equal([]int{0, 1, 2}))
		Expect(tmp).To(Equal([]int{0, 1, 2, 3, 4}))

		// copy and append
		appended := append(tmp[:len(tmp):len(tmp)], 5)

		tmp = append(tmp, 9)
		tmp[0] = 10
		Expect(tmp).To(Equal([]int{10, 1, 2, 3, 4, 9}))
		Expect(appended).To(Equal([]int{0, 1, 2, 3, 4, 5}))
	})
})
