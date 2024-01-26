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

package keyvalue

import (
	"github.com/mandelsoft/logging"
)

const ERROR = "error"

func Error(v interface{}) interface{} {
	return logging.KeyValue(ERROR, v)
}

const ID = "id"

func Id(v interface{}) interface{} {
	return logging.KeyValue(ERROR, v)
}

const NAME = "name"

func Name(v interface{}) interface{} {
	return logging.KeyValue(NAME, v)
}

const NAMESPACE = "namespace"

func Namespace(v interface{}) interface{} {
	return logging.KeyValue(NAMESPACE, v)
}

const ELEMENT = "element"

func Element(v interface{}) interface{} {
	return logging.KeyValue(ELEMENT, v)
}
