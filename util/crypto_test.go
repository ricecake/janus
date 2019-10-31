package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ricecake/janus/util"
)

var _ = Describe("Crypto", func() {
	Describe("CompactUUID", func() {
		uuid := util.CompactUUID()
		It("returns a well-formed string", func() {
			Expect(uuid).NotTo(Equal(""))
			Expect(uuid).To(MatchRegexp("^[A-Za-z0-9+_-]+$"))
		})
	})
})
