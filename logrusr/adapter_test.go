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

package logrusr_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/sirupsen/logrus"

	"github.com/mandelsoft/logging"
	"github.com/mandelsoft/logging/logrusr"
)

var _ = Describe("mapping test", func() {
	It("maps Info to Info", func() {
		buf := &bytes.Buffer{}
		log := logrus.New()
		log.SetLevel(9)
		log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
		log.SetOutput(buf)
		ctx := logging.New(logrusr.New(log))
		ctx.SetDefaultLevel(logging.InfoLevel)
		ctx.Logger().Info("test")
		Expect(buf.String()).To(Equal("{\"level\":\"info\",\"msg\":\"test\"}\n"))
	})

	It("maps Error to Error", func() {
		buf := &bytes.Buffer{}
		log := logrus.New()
		log.SetLevel(9)
		log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
		log.SetOutput(buf)
		ctx := logging.New(logrusr.New(log))
		ctx.Logger().Error("test")
		Expect(buf.String()).To(Equal("{\"level\":\"error\",\"msg\":\"test\"}\n"))
	})

	It("maps Error with err to Error", func() {
		buf := &bytes.Buffer{}
		log := logrus.New()
		log.SetLevel(9)
		log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
		log.SetOutput(buf)
		ctx := logging.New(logrusr.New(log))
		ctx.Logger().LogError(fmt.Errorf("errmsg"), "test")
		Expect(buf.String()).To(Equal("{\"error\":\"errmsg\",\"level\":\"error\",\"msg\":\"test\"}\n"))
	})

})
