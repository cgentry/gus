package bcrypt

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/drivers/encryption"
)

const DRIVER_NAME = "bcrypt"

type Register struct {}

// init will register this driver for use.
func init() {
	gdriver.Register(encryption.DRIVER_GROUP,&RegisterBCrypt  )
}


func ( r *Register ) New() EncryptDriver {
	return New()
}


func ( r *Register ) Identity( id int ) string {
	switch id {
	case IDENT_NAME:
		return "BCrypt"
	case IDENT_SHORT:
		return "Standard high-quality encryption using BCRYPT methods"
	case IDENT_LONG:
		return const_bcrypt_help_template
	}
	return "unknown"
}