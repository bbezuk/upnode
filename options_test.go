package upnode

import (
	"fmt"
	"testing"
)

func TestParseFromString(t *testing.T) {
	samples := make(map[string]Options)

	samples["p:=80#kwd:=55#kwd:=56"] = Options{80, "default", []uint64{55, 56}, "", 3, false, false, true}
	samples["p:=80#ste:=200 OK#flw:=0"] = Options{80, "200 OK", []uint64{}, "", 3, false, false, false}
	samples["p:=123#wb:=1"] = Options{123, "default", []uint64{}, "", 3, false, true, true}
	samples["p:=33#fwww:=1"] = Options{33, "default", []uint64{}, "", 3, true, false, true}
	samples["kwd:=43"] = Options{0, "default", []uint64{43}, "", 3, false, false, true}

	for key, val := range samples {
		generated := new(Options)
		generated.ParseFromString(key)

		//t.Logf("Generated: %v", *generated)

		optionsExpectedString := fmt.Sprintf("%v", val)
		optionsRecievedString := fmt.Sprintf("%v", *generated)

		if optionsExpectedString != optionsRecievedString {
			t.Logf("Wrong parsing of options, expected %s, and recieved %s\n", optionsExpectedString, optionsRecievedString)
			t.Fail()
		}
	}
}
