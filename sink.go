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

type sink struct {
	level int
	sink  logr.LogSink
}

var _ logr.LogSink = (*sink)(nil)

func WrapSink(level int, orig logr.LogSink) logr.LogSink {
	return &sink{
		level: level,
		sink:  orig,
	}
}

func (s *sink) Init(info logr.RuntimeInfo) {
	s.sink.Init(info)
}

func (s *sink) Enabled(level int) bool {
	return s.level >= level
}

func (s *sink) Info(level int, msg string, keysAndValues ...interface{}) {
	if !s.Enabled(level) {
		return
	}
	s.sink.Info(level, msg, keysAndValues...)
}

func (s sink) Error(err error, msg string, keysAndValues ...interface{}) {
	s.sink.Error(err, msg, keysAndValues...)
}

func (s *sink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &sink{
		level: s.level,
		sink:  s.sink.WithValues(keysAndValues...),
	}
}

func (s sink) WithName(name string) logr.LogSink {
	return &sink{
		level: s.level,
		sink:  s.sink.WithName(name),
	}
}
