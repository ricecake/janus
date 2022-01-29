package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"janus/util"
)

var _ = Describe("String", func() {
	It("Handles simple list", func() {
		list := []string{"a", "b", "c"}
		Expect(util.UniquifyStringSlice(list)).To(Equal(list))
	})
	It("Handles simple list", func() {
		Expect(util.UniquifyStringSlice([]string{"a", "a", "b", "c"})).To(Equal([]string{"a", "b", "c"}))
	})
	It("Ignores not-consecutive duplicates", func() {
		Expect(util.UniquifyStringSlice([]string{"a", "b", "a", "c"})).To(Equal([]string{"a", "b", "a", "c"}))
	})
})
