package plaintext

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/drivers/encryption"
)

const DRIVER_NAME = "plaintext"

type Register struct{}

// init will register this driver for use.
func init() {
	gdriver.Register(encryption.DRIVER_GROUP, &Register)
}

func (r *Register) New() EncryptDriver {
	return New()
}

func (r *Register) Identity(id int) string {
	switch id {
	case IDENT_NAME:
		return "plaintext"
	case IDENT_SHORT:
		return "For testing only! Do not use in production"
	case IDENT_LONG:
		return const_plain_help_template
	}
	return "unknown"
}
