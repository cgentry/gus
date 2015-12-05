// The encryption drivers are used to encrypt and decrypt passwords stored for a user.
// Standard encryption drivers are sha512, bcrypt and plaintext (not suitable for
// production systems)
//
// All drivers need to call Register in order to be usable by the system. Failure to select
// an encryption driver will cause a panic during runtime.
//
// Drivers are selected by:
//    crypt := encryption.Select( driverName ).Setup( options )
// The call to Setup() is optional and driver specific.

package encryption

import (
	"encoding/json"
	"github.com/cgentry/gdriver"
	"strings"
)

const DRIVER_GROUP = "encryption"

// The interface gives the set of methods that an encryption driver must implement.
type EncryptDriver interface {
	EncryptPassword(password string, salt string) string
	ComparePasswords(string, string, string) bool
	Setup(string) EncryptDriver

	//  The following are wrappers for the gdriver routines.
	Id() string
	ShortHelp() string
	LongHelp() string
}

// These are common parameters used by many drivers. Each driver may use structures that are specific to
// that driver.
type CryptOptions struct {
	StaticSaltIndex int			`json:"StaticSaltIndex"`
	Cost       int			`json:"Cost"`
	Salt       string		`json:"Salt"`
}

// Unmarshal a json string containing the common options defined in CryptOptions and return
// the option structure
func UnmarshalOptions(jsonOption string) (opt *CryptOptions, err error) {
	opt = &CryptOptions{}
	jsonOption = strings.TrimSpace(jsonOption)
	if jsonOption != "" {
		err = json.Unmarshal([]byte(jsonOption), opt)
	}
	return
}

// Pick a registered driver for use in the system. Only one driver can be selected at a time.
// This will panic if no drivers have been registered
func Select(name string) EncryptDriver {
	gdriver.Default(DRIVER_GROUP,name)
	return gdriver.MustNew(DRIVER_GROUP, name).(EncryptDriver)
}


// This will panic if no drivers have been registered
func GetDriver() EncryptDriver {
	return gdriver.MustNewDefault(DRIVER_GROUP).(EncryptDriver)
}

func GetStaticSalt(offset int) string {
	modIndex := offset % len(encryption_salts)
	return encryption_salts[modIndex]
}
