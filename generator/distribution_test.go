package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestNormalCompatibleDomain(t *testing.T) {
	norm := &NormalDistribution{}
	Assert(t, norm.isCompatibleDomain("float"), "floats should be a compatible domain for normal distributions")
	Assert(t, !norm.isCompatibleDomain("integer"), "integers should not be a compatible domain for normal distributions")
}

func TestUniformCompatibleDomain(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, uni.isCompatibleDomain("float"), "floats should be a compatible domain for uniform distributions")
	Assert(t, uni.isCompatibleDomain("integer"), "integers should be a compatible domain for uniform distributions")
	Assert(t, !uni.isCompatibleDomain("string"), "strings should not be a compatible domain for uniform distributions")
}

func TestNormalShouldntSupportMultipleDomains(t *testing.T) {
	norm := &NormalDistribution{}
	Assert(t, !norm.supportsMultipleDomains(), "normal distributions don't support multiple domains")
}

func TestUniformShouldntSupportMultipleDomains(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, !uni.supportsMultipleDomains(), "uniform distributions don't support multiple domains")
}

func TestWeightedShouldSupportMultipleDomains(t *testing.T) {
	w := &WeightedDistribution{}
	Assert(t, w.supportsMultipleDomains(), "weighted distributions should support multiple domains")
}
