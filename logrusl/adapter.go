/*
 * Copyright 2023 Mandelsoft. All rights reserved.
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

package logrusl

import (
	"bytes"
	"io"
	"os"

	"github.com/mandelsoft/logging"
	"github.com/mandelsoft/logging/logrusfmt"
	"github.com/mandelsoft/logging/logrusr"
	"github.com/sirupsen/logrus"
)

const FieldKeyRealm = logging.FieldKeyRealm

var defaultFixedKeys = []string{
	logrusfmt.FieldKeyTime,
	logrusfmt.FieldKeyLevel,
	logrusr.FieldKeyLogger,
	FieldKeyRealm,
	logrusfmt.FieldKeyMsg,
}

var defaultFieldFormatters = logrusfmt.FieldFormatters{
	logrusfmt.FieldKeyTime:  logrusfmt.PlainValue,
	logrusfmt.FieldKeyLevel: logrusfmt.LevelValue,
	logrusr.FieldKeyLogger:  logrusfmt.PlainValue,
	FieldKeyRealm:           logrusfmt.BracketValue,
	logrusfmt.FieldKeyMsg:   logrusfmt.PlainValue,
}

////////////////////////////////////////////////////////////////////////////////

// Settings is an composition environment to configure a
// logrus.Logger or a logging.Context.
type Settings struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

func (s Settings) WithWriter(w io.Writer) Settings {
	s.Writer = w
	return s
}

func (s Settings) WithFormatter(f logrus.Formatter) Settings {
	s.Formatter = f
	return s
}

func (s Settings) Human() Settings {
	s.Formatter = NewTextFmtFormatter()
	return s
}

func (s Settings) JSON() Settings {
	s.Formatter = NewJSONFormatter()
	return s
}

func (s Settings) NewLogrus() *logrus.Logger {
	logger := newLogger()
	logger.Out = s.Writer
	if logger.Out == nil {
		logger.Out = os.Stderr
	}
	logger.Formatter = s.Formatter
	if logger.Formatter == nil {
		logger.Formatter = NewTextFormatter()
	}
	return logger
}

func (s Settings) New() logging.Context {
	return logging.New(logrusr.New(s.NewLogrus()))
}

////////////////////////////////////////////////////////////////////////////////

func New() logging.Context {
	return Settings{}.New()
}

func WithWriter(w io.Writer) Settings {
	return Settings{}.WithWriter(w)
}

func WithFormatter(f logrus.Formatter) Settings {
	return Settings{}.WithFormatter(f)
}

func Human() Settings {
	return Settings{}.Human()
}

func JSON() Settings {
	return Settings{}.JSON()
}

func newLogger(buf ...*bytes.Buffer) *logrus.Logger {
	logger := logrus.New()
	if len(buf) > 0 && buf[0] != nil {
		logger.Out = buf[0]
	}
	logger.Level = 9
	return logger
}

func NewTextFormatter() *logrusfmt.TextFormatter {
	return &logrusfmt.TextFormatter{
		FixedFields:     defaultFixedKeys,
		FieldFormatters: defaultFieldFormatters,
	}
}

func NewTextFmtFormatter() *logrusfmt.TextFmtFormatter {
	return &logrusfmt.TextFmtFormatter{*NewTextFormatter()}
}

func NewJSONFormatter() *logrusfmt.JSONFormatter {
	return &logrusfmt.JSONFormatter{
		FixedFields: defaultFixedKeys,
	}
}
