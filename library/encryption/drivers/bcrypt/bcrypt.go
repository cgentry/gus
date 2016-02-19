// Package bcrypt will encrypt passwords using a strong algorithim, bcrypt,
// and is the recommended driver. The setup option passed in can be:
// { Cost: n, Salt: "string"}
// Where cost is how many iterations bcrypt should perform. The higher the
// number the longer it takes to generate and break. The default is 7.
// Salt is an additional string you want to add in make the hash harder
// to guess. If you don't include one, the internal salt will be used.

// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package

package bcrypt

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
	"log"
)

type PwdBcrypt struct {
	Salt string
	Cost int
}

func (t *PwdBcrypt) Id() string {
	return gdriver.Help(encryption.DRIVER_GROUP, DRIVER_NAME, gdriver.IDENT_NAME)
}
func (t *PwdBcrypt) ShortHelp() string {
	return gdriver.Help(encryption.DRIVER_GROUP, DRIVER_NAME, gdriver.IDENT_SHORT)
}
func (t *PwdBcrypt) LongHelp() string {
	return gdriver.Help(encryption.DRIVER_GROUP, DRIVER_NAME, gdriver.IDENT_LONG)
}

// New will create a BCRYPT strucutre. The salt is given a static string but
// can be set up on selection from the driver. This must be the same with every
// load or you won't be able to login anymore.
func New() *PwdBcrypt {
	c := &PwdBcrypt{
		Cost: 7,
		Salt: "vniiO5UD0w5GpJkPijwQCT63MuMjyWnyi5TtUWBGInCq84zaFFsSwGm9DK8UyUeQp{2h&gV,KoQi9ysC",
	}
	return c
}

// EncryptPassword will encrypt the password using the magic number within the record.
// This should be sufficient to protect it but still allow us to re-create later on.
// (The magic number will never alter for the life of the record
func (t *PwdBcrypt) EncryptPassword(clearPassword, userSalt string) string {
	saltyPassword := []byte(clearPassword + t.Salt + userSalt + encryption.GetStaticSalt(0))
	pass1, _ := bcrypt.GenerateFromPassword(saltyPassword, t.Cost)
	return string(pass1)
}

// Setup should be called only when the driver has been selected for use.
func (t *PwdBcrypt) Setup(jsonOptions string) encryption.EncryptDriver {
	if jsonOptions != "" {
		opt, err := encryption.UnmarshalOptions(jsonOptions)
		if err != nil {
			log.Printf("Bcrypt: Could not unmarshal '%s' options: ignored.", jsonOptions)
			return t
		}
		t.setCost( opt.Cost )
		t.setSalt( opt.Salt )
	}
	return t
}

func ( t *PwdBcrypt )setCost( newCostValue int ){
	if newCostValue > 0 {
		t.Cost = newCostValue
	}
}

func ( t *PwdBcrypt ) setSalt( newEncryptionSalt string ){
	if len(newEncryptionSalt) > 0 {
		t.Salt = newEncryptionSalt
	}
}


// ComparePasswords must be called with a bcrypt password.
func (t *PwdBcrypt) ComparePasswords(hashedPassword, clearPassword, userSalt string) bool {
	saltyPassword := []byte(clearPassword + t.Salt + userSalt + encryption.GetStaticSalt(0))
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), saltyPassword)
	return err == nil
}

