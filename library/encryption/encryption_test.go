package encryption

import (
	"github.com/cgentry/gdriver"
	"testing"
)

type mockDriver struct{}
type mockCrypt struct{}

func (m *mockDriver) Init() string { return "init" }

type tDriver1 struct{}

func (t *tDriver1) New() interface{} { return &mockCrypt{} }
func (t *tDriver1) Identity(id int) string {
	switch id {
	case gdriver.IdentityName:
		return "name"
	case gdriver.IdentityShort:
		return "short"
	case gdriver. IdentityLong:
		return "long"
	}
	return "unknown"
}

func (m *mockCrypt ) EncryptPassword(password string, salt string) string {
	return password +"/" + salt
}
func (m *mockCrypt) ComparePasswords( a string, b string, c string ) bool {
	return a==b
}

func (m *mockCrypt) Setup( a string ) EncryptDriver {
	return m
}
func (t *mockCrypt) Id() string {
	return gdriver.Help(DriverGroup, "name", gdriver.IdentityName)
}
func (t *mockCrypt) ShortHelp() string {
	return gdriver.Help(DriverGroup, "name", gdriver.IdentityShort)
}
func (t *mockCrypt) LongHelp() string {
	return gdriver.Help(DriverGroup, "name", gdriver.IdentityLong)
}


func TestRegister(t *testing.T) {

	gdriver.Register( DriverGroup, &tDriver1{})

	drv := SetDefault( "name")
	if "name" != drv.Id() {
		t.Error("Name returned was not 'name': " + drv.Id() )
	}

	tstString := drv.EncryptPassword( "password","salt" )
	if "password/salt" != tstString {
		t.Error("Invalid return from encrypt: " + tstString)
	}else {

		if ! drv.ComparePasswords(tstString, "password/salt", "") {
			t.Error("Invalid match from ComparePasswords")
		}

		if drv.ComparePasswords(tstString, "password?salt", "") {
			t.Error("Invalid match from ComparePasswords")
		}
	}
}
