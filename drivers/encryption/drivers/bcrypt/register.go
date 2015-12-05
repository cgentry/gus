package bcrypt

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/drivers/encryption"
)

const DRIVER_NAME = "bcrypt"

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
		return "BCrypt"
	case gdriver.IDENT_SHORT:
		return "Standard high-quality encryption using BCRYPT methods"
	case gdriver.IDENT_LONG:
		return const_bcrypt_help_template
	}
	return gdriver.IDENT_UNKNOWN
}
