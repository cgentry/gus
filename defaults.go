package main

// This is the default configuration file for the GUS system. You should select any
// database packages here, alter the constants that are used in the code and put any
// system-wide configurations here.
import (
	/*
	 *  DATABASE SUPPORT:
	 *		Include what you want to use here, then perform the registration below
	 */
	"github.com/cgentry/gus/library/storage/drivers/jsonfile"
	"github.com/cgentry/gus/library/storage/drivers/mock"
	"github.com/cgentry/gus/library/storage/drivers/sqlite"

	/*
	*  ENCRYPTION SUPPORT:
	*		Include what you want to use here, then perform the registration below
	 */
	"github.com/cgentry/gus/library/encryption/drivers/bcrypt"
	"github.com/cgentry/gus/library/encryption/drivers/sha512"
	/* REMOVE WHEN IN PRODUCTION */
	"github.com/cgentry/gus/library/encryption/drivers/plaintext"
)

// DefaultConfigFilename is where you will find the configuration file for GUS
// DefaultConfigPermissions is the octal UNIX permissions for the config file
const (
	DefaultConfigFilename    = "/etc/gus/config.json"
	DefaultConfigPermissions = 0600
)

func init() {
	/* DATABASE SUPPORT */
	jsonfile.Register()
	mock.Register()
	sqlite.Register()

	/* ENCRYPTION SUPPORT */
	bcrypt.Register()
	sha512.Register()
	plaintext.Register()
}
