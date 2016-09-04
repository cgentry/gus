package plaintext

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
)

// PwdPlaintext defines the variables required and initialised
type PwdPlaintext struct {
	Salt string
}

// New will return an address pointing to a new PwdPlaintext structure
func New() *PwdPlaintext {
	c := &PwdPlaintext{
		Salt: "SALT",
	}
	return c
}

// Id returns the string identifier for this driver
func (t *PwdPlaintext) Id() string {
	return gdriver.Help(encryption.DriverGroup, DriverName, gdriver.IdentityName)
}

// ShortHelp returns a short string identifier for the identity.
func (t *PwdPlaintext) ShortHelp() string {
	return gdriver.Help(encryption.DriverGroup, DriverName, gdriver.IdentityShort)
}

// LongHelp returns a longer descriptive text for the help
func (t *PwdPlaintext) LongHelp() string {
	return gdriver.Help(encryption.DriverGroup, DriverName, gdriver.IdentityLong)
}

// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdPlaintext) EncryptPassword(clearPassword, userSalt string) string {

	return clearPassword + ";" + userSalt + ";" + t.Salt + ";Plaintext"
}

// Setup should be called only when the driver has been selected for use.
func (t *PwdPlaintext) Setup(json string) encryption.EncryptDriver {
	opt, err := encryption.UnmarshalOptions(json)
	if err != nil {
		panic(err.Error())
	}

	t.setSalt(opt.Salt)

	return t
}

func (t *PwdPlaintext) setSalt(newEncryptionSalt string) {
	if len(newEncryptionSalt) > 0 {
		t.Salt = newEncryptionSalt
	}
}

// ComparePasswords must check to see if the passed password is equal to the stored password
func (t *PwdPlaintext) ComparePasswords(hashedPassword, password, salt string) bool {
	return hashedPassword == t.EncryptPassword(password, salt)
}

const constPlainTextHelpTempate = `
  This does not encrypt passwords and should never be selected for production use. It
  is only to be used by developers and for testing purposes. The format of the password
  output is:
           [user password];[user salt];[driver's salt];Plaintext
  If a user has a salt of 'kjldoeuifnfl203294fkf' and the password is 'BadPassword', with
  defaults it would become:
           BadPassword;kjldoeuifnfl203294fkf;SALT;Plaintext

  Options: There is one option that can be passed in JSON format: "Salt". The default is "SALT".

  Option format: {"Salt": "Salty" }
`
