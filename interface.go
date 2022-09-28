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

import "github.com/go-logr/logr"

// These are the different logging levels. You can set the logging level to log
// on your instance of logger.
const (
	// None level. No logging,
	None = iota
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type Logger interface {
	Error(msg string, keypairs ...interface{})
	Warn(msg string, keypairs ...interface{})
	Info(msg string, keypairs ...interface{})
	Debug(msg string, keypairs ...interface{})
	Trace(msg string, keypairs ...interface{})

	WithName(name string) Logger
	WithValues(keypairs ...interface{}) Logger
	Enabled(level int) bool
}

type MessageContext interface {
}

type Condition interface {
	Match(...MessageContext) bool
}

type Rule interface {
	Match(logr.Logger, ...MessageContext) Logger
}

type Context interface {
	SetBaseLogger(logger logr.Logger)

	AddRule(...Rule)
	ResetRules()

	Logger(...MessageContext) Logger
}

type Attacher interface {
	Attach(l Logger) Logger
}

type Realm interface {
	Condition
	Attacher

	Name() string
}

type RealmPrefix interface {
	Condition
	Attacher

	Name() string
	IsPrefix() bool
}

type Attribute interface {
	Condition
	Attacher

	Name() string
	Value() interface{}
}

type Tag interface {
	Condition

	Name() string
}
