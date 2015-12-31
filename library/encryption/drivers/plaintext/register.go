package plaintext

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
)

const DRIVER_NAME = "plaintext"

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(encryption.DRIVER_GROUP, &registerDriver{})
}

func (r *registerDriver) New() interface{} {
	return New()
}

func (r *registerDriver) Identity(id int) string {
	switch id {
	case gdriver.IDENT_NAME:
		return "plaintext"
	case gdriver.IDENT_SHORT:
		return "For testing only! Do not use in production"
	case gdriver.IDENT_LONG:
		return const_plain_help_template
	}
	return gdriver.IDENT_UNKNOWN
}
