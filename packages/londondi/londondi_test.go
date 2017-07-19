package londondi

import (
	"bytes"
	"testing"
)

//decoding cases
var CASE_0 = []byte{0x8d, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8c}
var RESULT_0 = []byte{0x8d, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8c}

var CASE_1 = []byte{0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x1b, 0x82, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec}
var RESULT_1 = []byte{0x8d, 0x00, 0x01, 0x03, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec}

var CASE_2 = []byte{0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82}
var RESULT_2 = []byte{0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82}

var CASE_3 = []byte{0x02, 0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x1b, 0x82, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82, 0x03}
var RESULT_3 = []byte{0x02, 0x8d, 0x00, 0x01, 0x03, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82, 0x03}

func TestChecksum(t *testing.T) {

}

func TestSubstitutions(t *testing.T) {

	//test decoding a slice where there are no subsitutions necessary when decoding, e.g. no instances of escape byte
	result, err := MakeSubstitutions(CASE_0, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_0) {
		t.Error("Error on case 0: decode a slice with no necessary substitutions")
	}

	//slice with substitutions 0x1b82 -> 0x02, 0x1b83 -> 0x03
	result, err = MakeSubstitutions(CASE_1, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_1) {
		t.Error("Error on case 1: expected two substitutions")
	}

	//slice with unescaped subsitution (no substitutions)
	result, err = MakeSubstitutions(CASE_2, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_2) {
		t.Error("Error on case 2: unescaped substitution (expected no substitutions)")
	}

	//the kitchen sink
	result, err = MakeSubstitutions(CASE_3, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_3) {
		t.Error("Houston, we have a problem")
	}

}
