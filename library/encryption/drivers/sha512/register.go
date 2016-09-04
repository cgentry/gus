package sha512

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
)

// DriverName Specifies the specific identity of this driver within a group
const DriverName = "SHA512"

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(encryption.DriverGroup, &registerDriver{})
}

// SetDefault will set THIS driver as the default encryption driver.
func SetDefault() {
	gdriver.Default(encryption.DriverGroup, DriverName)
}

// New() will return a sha512 driver. the registerDriver function, being generic to all encryption drivers
// instantiates a specific driver
func (r *registerDriver) New() interface{} {
	return New()
}

// Identity will return the identifier string for this driver. The string can be either short or long
func (r *registerDriver) Identity(id int) string {
	switch id {

	case gdriver.IdentityShort:
		return "Standard quality encryption using SHA512 methods"
	case gdriver.IdentityLong:
		return constSha512HelpTemplate
	}
	return DriverName
}
