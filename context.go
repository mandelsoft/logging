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
	lock  sync.RWMutex
	base  Context
	level int
	sink  logr.LogSink
	rules []Rule
}

var _ Context = (*context)(nil)
var _ ContextProvider = (*context)(nil)

func NewDefault() Context {
	return NewWithBase(nil)
}

func New(logger logr.Logger) Context {
	sink := shifted(logger)
	ctx := &context{
		sink:  sink,
		level: -1,
	}
	ctx.level = InfoLevel
	return ctx
}

func NewWithBase(base Context, baselogger ...logr.Logger) Context {
	ctx := &context{
		base:  base,
		level: -1,
	}
	if len(baselogger) > 0 {
		ctx.SetBaseLogger(baselogger[0])
	}
	if base == nil && len(baselogger) == 0 {
		l := logrus.New()
		l.SetLevel(9)
		ctx.level = InfoLevel
		ctx.SetBaseLogger(logrusr.New(l))
	}
	return ctx
}

func (c *context) LoggingContext() Context {
	return c
}

func (c *context) GetSink() logr.LogSink {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.getSink()
}

func (c *context) getSink() logr.LogSink {
	if c.sink != nil {
		return c.sink
	}
	return c.base.GetSink()
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
	return NewLogger(DynSink(c.GetDefaultLevel, 0, c.GetSink))
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
}

func (c *context) SetBaseLogger(logger logr.Logger, plain ...bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(plain) == 0 || !plain[0] {
		c.sink = shifted(logger)
	} else {
		c.sink = logger.GetSink()
	}
}

func (c *context) AddRule(rules ...Rule) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, rule := range rules {
		if rule != nil {
			c.rules = append(append(c.rules[:0:0], rule), c.rules...)
		}
	}
}

func (c *context) AddRulesTo(ctx Context) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	rules := make([]Rule, len(c.rules))
	for i := range rules {
		rules[i] = c.rules[len(c.rules)-i-1]
	}
	ctx.AddRule(rules...)
}

func (c *context) ResetRules() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.rules = nil
}

func (c *context) Logger(messageContext ...MessageContext) Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	l := c.evaluate(c.GetSink, messageContext...)
	if l != nil {
		return l
	}
	return c.getDefaultLogger()
}

func (c *context) V(level int, messageContext ...MessageContext) logr.Logger {
	return c.Logger(messageContext...).V(level)
}

func (c *context) Evaluate(base SinkFunc, messageContext ...MessageContext) Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.evaluate(base, messageContext...)
}

func (c *context) evaluate(base SinkFunc, messageContext ...MessageContext) Logger {
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

type levelCatcher struct {
	sink
}

func (h *levelCatcher) Enabled(l int) bool {
	h.level = l
	return false
}

func getLogrLevel(l logr.Logger) int {
	var s levelCatcher
	l.WithSink(&s).Enabled()
	return s.level
}

func shifted(logger logr.Logger) logr.LogSink {
	level := getLogrLevel(logger)
	return WrapSink(level+9, level, logger.GetSink())
}
