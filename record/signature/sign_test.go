package signature_test

import (
	"testing"

	"github.com/cgentry/gus/record/signature"
)

func TestSignature(t *testing.T) {
	s := signature.New()

	if s.IsSignatureSet() == true {
		t.Error("Signature is set but shouldn't be (initialization error)")
	}

	if _, err := s.GetSignature(); err == nil {
		t.Error("Signature returned value but it should be enpty")
	}

	s.SetSignature([]byte(`abcdefg`))
	if s.IsSignatureSet() == false {
		t.Error("Signature isnt set but should be (SetSignature error)")
	}
	data, err := s.GetSignature()
	if err != nil {
		t.Errorf("Error is set but should be nil '%s'\n", err.Error())
	}
	if string(data) != `abcdefg` {
		t.Errorf("Signature back is invalid: %s\n", string(data))
	}

	s.ClearSignature()
	if s.IsSignatureSet() == true {
		t.Error("Signature is set but shouldn't be (ClearSignature error)")
	}
}
