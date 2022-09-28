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
	logger logr.Logger
}

var _ Logger = (*logger)(nil)

func NewLogger(l logr.Logger) Logger {
	return &logger{l}
}

func (l logger) LogError(err error, msg string, keypairs ...interface{}) {
	l.logger.V(1).Info(msg, append(append(keypairs[:0:0], "error", err), keypairs...))
}

func (l logger) Error(msg string, keypairs ...interface{}) {
	l.logger.V(1).Info(msg, keypairs...)
}

func (l logger) Warn(msg string, keypairs ...interface{}) {
	l.logger.V(2).Info(msg, keypairs...)
}

func (l logger) Info(msg string, keypairs ...interface{}) {
	l.logger.V(3).Info(msg, keypairs...)
}

func (l logger) Debug(msg string, keypairs ...interface{}) {
	l.logger.V(4).Info(msg, keypairs...)
}

func (l logger) Trace(msg string, keypairs ...interface{}) {
	l.logger.V(5).Info(msg, keypairs...)
}

func (l logger) WithName(name string) Logger {
	return &logger{l.logger.WithName(name)}
}

func (l logger) WithValues(keypairs ...interface{}) Logger {
	return &logger{l.logger.WithValues(keypairs...)}
}

func (l logger) Enabled(level int) bool {
	return l.logger.GetSink().Enabled(level)
}
