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
	"sync"

	"github.com/go-logr/logr"
	"github.com/mandelsoft/logging/logrusr"
	"github.com/sirupsen/logrus"
)

type context struct {
	lock          sync.RWMutex
	level         int
	baseLogger    logr.Logger
	defaultLogger Logger
	rules         []Rule
}

var _ Context = (*context)(nil)

func Default() Context {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	return New(logrusr.New(l))
}

func New(base logr.Logger) Context {
	ctx := &context{
		baseLogger: base,
	}
	ctx.setDefaultLevel(InfoLevel)
	return ctx
}

func (c *context) SetDefaultLevel(level int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.setDefaultLevel(level)
}

func (c *context) setDefaultLevel(level int) {
	c.level = level
	c.defaultLogger = NewLogger(c.baseLogger.WithSink(WrapSink(level, c.baseLogger.GetSink())))
}

func (c *context) SetBaseLogger(logger logr.Logger) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.baseLogger = logger
	c.setDefaultLevel(c.level)
}

func (c *context) AddRule(rules ...Rule) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, rule := range rules {
		if rule != nil {
			c.rules = append(append([]Rule{}, rule), c.rules...)
		}
	}
}

func (c *context) ResetRules() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.rules = nil
}

func (c *context) Logger(messageContext ...MessageContext) Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, rule := range c.rules {
		l := rule.Match(c.baseLogger, messageContext...)
		if l != nil {
			for _, c := range messageContext {
				if a, ok := c.(Attacher); ok {
					l = a.Attach(l)
				}
			}
			return l
		}
	}
	return c.defaultLogger
}
