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

package logging

type tag struct {
	name string
}

var _ Tag = (*tag)(nil)

func NewTag(name string) Tag {
	return &tag{name}
}

func (r *tag) Match(messageContext ...MessageContext) bool {
	for _, c := range messageContext {
		if e, ok := c.(Tag); ok && e.Name() == r.name {
			return true
		}
	}
	return false
}

func (r *tag) Name() string {
	return r.name
}
