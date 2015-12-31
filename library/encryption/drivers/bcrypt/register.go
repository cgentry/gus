package bcrypt

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
)

const (
	DRIVER_NAME = "bcrypt"
	STORAGE_IDENTITY	= "BCrypt"
	SHORT_HELP			= "Standard high-quality encryption using BCRYPT methods"
	HELP_TEMPLATE		= `
  The bcrypt function is the default password hash algorithm for BSD and many other systems.
  Besides incorporating a salt to protect against rainbow table attacks, bcrypt is an adaptive
  function: over time, the iteration count can be increased to make it slower, so it remains
  resistant to brute-force search attacks even with increasing computation power.

  Options: There are two options that should be passed by JSON strings. They are:
      "Cost" and "Salt". Cost is the number of iterations you want for the function, making
      it more costly to encrypt (which is a good thing). Salt is an additional bit of
      encryption you want added when it is encrypting the password. The salt should
      be a long, random string of any characters. Do not include quotes.

      The cost defaults to '7' and the salt has a long, random string built in. You must
      not change the salt after you have set it or passwords will never match again.

  Option format: {"Cost" : 7, "Salt": "abcd...........xyz" }

`
)

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
		return STORAGE_IDENTITY
	case gdriver.IDENT_SHORT:
		return 	SHORT_HELP
	case gdriver.IDENT_LONG:
		return HELP_TEMPLATE
	}
	return gdriver.IDENT_UNKNOWN
}
