package sha512

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
)

const DRIVER_NAME = "SHA512"

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(encryption.DRIVER_GROUP, &registerDriver{})
}

func SetDefault(){
	gdriver.Default( encryption.DRIVER_GROUP, DRIVER_NAME )
}

func (r *registerDriver) New() interface{} {
	return New()
}

func (r *registerDriver) Identity(id int) string {
	switch id {

	case gdriver.IDENT_SHORT:
		return "Standard quality encryption using SHA512 methods"
	case gdriver.IDENT_LONG:
		return const_sha512_help_template
	}
	return DRIVER_NAME
}
