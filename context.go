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
	base          Context
	level         int
	baseLogger    *logr.Logger
	defaultLogger Logger
	rules         []Rule
}

var _ Context = (*context)(nil)

func NewDefault() Context {
	return NewWithBase(nil)
}

func New(base logr.Logger) Context {
	ctx := &context{
		baseLogger: &base,
		level:      -1,
	}
	ctx.setDefaultLevel(InfoLevel)
	return ctx
}

func NewWithBase(base Context, baselogger ...logr.Logger) Context {
	ctx := &context{
		base: base,
	}
	if len(baselogger) > 0 {
		ctx.SetBaseLogger(baselogger[0])
	}
	if base == nil && len(baselogger) == 0 {
		l := logrus.New()
		l.SetLevel(9)
		ctx.SetBaseLogger(logrusr.New(l))
	}
	return ctx
}

func (c *context) GetBaseLogger() logr.Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.getBaseLogger()
}

func (c *context) getBaseLogger() logr.Logger {
	if c.baseLogger != nil {
		return *c.baseLogger
	}
	return c.base.GetBaseLogger()
}

func (c *context) GetBaseContext() Context {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.base
}

func (c *context) GetDefaultLogger() Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.getDefaultLogger()
}

func (c *context) getDefaultLogger() Logger {
	if c.defaultLogger != nil {
		return c.defaultLogger
	}
	return c.base.GetDefaultLogger()
}

func (c *context) GetDefaultLevel() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.getDefaultLevel()
}

func (c *context) getDefaultLevel() int {
	if c.level < 0 {
		return c.base.GetDefaultLevel()
	}
	return c.level
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

	c.baseLogger = &logger
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

	l := c.evaluate(c.getBaseLogger(), messageContext...)
	if l != nil {
		return l
	}
	return c.getDefaultLogger()
}

func (c *context) V(level int, messageContext ...MessageContext) logr.Logger {
	return c.Logger(messageContext...).V(level)
}

func (c *context) Evaluate(base logr.Logger, messageContext ...MessageContext) Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.evaluate(base, messageContext...)
}

func (c *context) evaluate(base logr.Logger, messageContext ...MessageContext) Logger {
	for _, rule := range c.rules {
		l := rule.Match(base, messageContext...)
		if l != nil {
			for _, c := range messageContext {
				if a, ok := c.(Attacher); ok {
					l = a.Attach(l)
				}
			}
			return l
		}
	}
	if c.base != nil {
		return c.base.Evaluate(base, messageContext...)
	}
	return nil
}
