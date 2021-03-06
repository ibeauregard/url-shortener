package handling

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	generateKeyTests := []struct {
		arg1     string
		arg2     uint
		expected string
	}{
		{"http://fgr.ca", 42, "2Ms"},
		{"http://foobar.com", 235, "Pw69"},
		{"http://google.ca", 568, "=yC="},
		{"http://twitter.com", 735298, "p=5=M="},
		{"http://stackoverflow.com", 0, "Vf2"},
		{"http://qwasar.io", 235, "mP69"},
		{"http://linkedin.com", 9128, "Rd4xB"},
		{"http://facebook.com", 19, "bwP"},
	}
	for _, test := range generateKeyTests {
		testName := fmt.Sprintf("generateKey(%q, %v)", test.arg1, test.arg2)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, generateKey(test.arg1, test.arg2))
		})
	}
}

func TestChecksum(t *testing.T) {
	checksumTests := []struct {
		arg      string
		expected uint
	}{
		{"0", 4108050209},
		{"1", 2212294583},
		{"a", 3904355907},
		{"b", 1908338681},
		{"A", 3554254475},
		{"B", 1255198513},
		{"The quick brown fox jumps over the lazy dog", 1095738169},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, " +
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", 1196127599},
		{"https://www.google.com", 857627499},
		{"https://www.twitter.com", 1116836626},
	}

	for _, test := range checksumTests {
		testName := fmt.Sprintf("checksum(%q)", test.arg)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, checksum(test.arg))
		})
	}
}

func TestIntToKey(t *testing.T) {
	intToKeyTests := []struct {
		arg      uint
		expected string
	}{
		{0, "2"},
		{1, "3"},
		{2, "4"},
		{3, "5"},
		{4, "6"},
		{5, "7"},
		{10, "D"},
		{42, "s"},
		{275, "6y"},
		{5868, "3w&"},
		{92_840, "Zfv"},
		{644_539, "5YRq"},
		{5_063_859, "YPgk"},
		{49_812_303, "6q=g5"},
		{851_646_681, "3Tnpvz"},
		{3_071_006_908, "77!s4D"},
	}

	for _, test := range intToKeyTests {
		testName := fmt.Sprintf("intToKey(%v)", test.arg)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, intToKey(test.arg))
		})
	}
}
