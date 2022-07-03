package main

import (
	"fmt"
	"testing"
)

type checksumTest struct {
	arg      string
	expected uint
}

var checksumTests = []checksumTest{
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

func TestChecksum(t *testing.T) {
	for _, test := range checksumTests {
		testName := fmt.Sprintf("checksum(%q)", test.arg)
		t.Run(testName, func(t *testing.T) {
			if output := checksum(test.arg); output != test.expected {
				t.Errorf("got %v, expected %v", output, test.expected)
			}
		})
	}
}

type intToKeyTest struct {
	arg      uint
	expected string
}

var intToKeyTests = []intToKeyTest{
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

func TestIntToKey(t *testing.T) {
	for _, test := range intToKeyTests {
		testName := fmt.Sprintf("intToKey(%v)", test.arg)
		t.Run(testName, func(t *testing.T) {
			if output := intToKey(test.arg); output != test.expected {
				t.Errorf("got %q, expected %q", output, test.expected)
			}
		})
	}
}
