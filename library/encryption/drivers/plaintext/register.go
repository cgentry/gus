package plaintext

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
)

// DriverName Specifies the specific identity of this driver within a group
const DriverName = "plaintext"

type registerDriver struct{}

// Register is a simple wrapper to make sure registration occurs properly
func Register() {
	gdriver.Register(encryption.DriverGroup, &registerDriver{})
}

// SetDefault will set THIS driver as the default encryption driver.
func SetDefault() {
	gdriver.Default(encryption.DriverGroup, DriverName)
}

// New() will return a plaintext driver. the registerDriver function, being generic to all encryption drivers
// instantiates a specific driver
func (r *registerDriver) New() interface{} {
	return New()
}

// Identity will return the defined values for help, depending on whether they want a short or longg identifier.
func (r *registerDriver) Identity(id int) string {
	switch id {

	case gdriver.IdentityShort:
		return "For testing only! Do not use in production"
	case gdriver.IdentityLong:
		return constPlainTextHelpTempate
	}
	return DriverName
}
