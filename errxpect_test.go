// Copyright 2020 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package errxpect

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModelMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "geps module")
}

var _ = Describe("Errpect", func() {

	It("Succeed()s with nil error", func() {
		Errxpect(func() (string, error) { return "", nil }()).To(Succeed())
		Errxpect(func() (string, error) { return "", nil }()).To(Succeed(), "DOH!")
	})

	It("HaveOccured()s with non-nil error", func() {
		Errxpect(func() (string, error) { return "", errors.New("42") }()).To(HaveOccurred())
	})

	It("", func() {
		s := InterceptGomegaFailures(
			func() { Errxpect(func() (string, error) { return "", nil }()).To(HaveOccurred(), "DOH!") })
		Expect(s).To(ConsistOf(HavePrefix("DOH!\nExpected an error to have occurred.")))

		//
		Errxpect(func() (string, error) { return "true", nil }()).To(Succeed(), "DOH!")

		// When final err return value is non-nil, all preceeding return values
		// must be zero.
		s = InterceptGomegaFailures(
			func() {
				Errxpect(func() (string, error) { return "foobar", errors.New("dOH! ") }()).To(Succeed(), "DOH!")
			})
		Expect(s).To(ConsistOf(HavePrefix("DOH!\nUnexpected non-nil/non-zero actual non-error argument at index 1:")))

	})

})
