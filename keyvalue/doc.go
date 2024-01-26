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

// Package keyvalue provides some standard key/value pair values usable for
// the key/value list of logging messages. They bundle key and value and
// provide standardizes keys for various use cases, like errors, ids or
// names.
//
// Own standard key/value pairs can be defined in own packages by using
// the logging.KeyValue function.
//
// Those values can be used as single argument representing a key/value pair
// together with a sequence of key and value
// arguments in the argument list of the logging methods.
package keyvalue
