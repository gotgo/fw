package zoo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestZoo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zoo Suite")
}
