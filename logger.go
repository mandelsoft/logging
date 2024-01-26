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

import (
	"github.com/go-logr/logr"
)

type logger struct {
	sink logr.LogSink
}

var _ Logger = (*logger)(nil)

func NewLogger(s logr.LogSink) Logger {
	return &logger{s}
}

func (l *logger) V(delta int) logr.Logger {
	return logr.New(l.sink).V(delta)
}

func (l *logger) LogError(err error, msg string, keypairs ...interface{}) {
	if l.Enabled(ErrorLevel) {
		l.sink.Error(err, msg, prepare(keypairs)...)
	}
}

func (l *logger) Error(msg string, keypairs ...interface{}) {
	if l.Enabled(ErrorLevel) {
		l.sink.Error(nil, msg, prepare(keypairs)...)
	}
}

func (l *logger) Warn(msg string, keypairs ...interface{}) {
	l.sink.Info(WarnLevel, msg, prepare(keypairs)...)
}

func (l *logger) Info(msg string, keypairs ...interface{}) {
	l.sink.Info(InfoLevel, msg, prepare(keypairs)...)
}

func (l *logger) Debug(msg string, keypairs ...interface{}) {
	l.sink.Info(DebugLevel, msg, prepare(keypairs)...)
}

func (l *logger) Trace(msg string, keypairs ...interface{}) {
	l.sink.Info(TraceLevel, msg, prepare(keypairs)...)
}

func (l logger) WithName(name string) Logger {
	return &logger{l.sink.WithName(name)}
}

func (l logger) WithValues(keypairs ...interface{}) Logger {
	return &logger{l.sink.WithValues(prepare(keypairs)...)}
}

func (l logger) Enabled(level int) bool {
	return l.sink.Enabled(level)
}

type keyValue struct {
	name  string
	value interface{}
}

var _ keyvalue = (*keyValue)(nil)

type keyvalue interface {
	Name() string
	Value() interface{}
}

func (kv *keyValue) Name() string {
	return kv.name
}

func (kv *keyValue) Value() interface{} {
	return kv.value
}

// KeyValue provide a key/value pair for the argument list of logging methods.
//
// Those values can be used as single argument representing a key/value pair
// together with a sequence of key and value
// arguments in the argument list of the logging methods.
//
// This function can be used to define standardized keys for various
// use cases(see package [github.com/mandelsoft/logging/keyvalue]).
func KeyValue(key string, value interface{}) *keyValue {
	return &keyValue{
		name:  key,
		value: value,
	}
}

// prepare
func prepare(keypairs []interface{}) []interface{} {
	for i, e := range keypairs {
		if i%2 == 0 {
			if _, ok := e.(keyvalue); ok {
				var r []interface{}
				for i := 0; i < len(keypairs); i += 2 {
					if v, ok := keypairs[i].(keyvalue); ok {
						r = append(r, v.Name(), v.Value())
						i--
					} else {
						if i+1 < len(keypairs) {
							r = append(r, keypairs[i], keypairs[i+1])
						} else {
							r = append(r, keypairs[i]) // preserve the erroneous value
						}
					}
				}
				return r
			}
		}
	}
	return keypairs
}
