package plaintext

import (
	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/library/encryption"
	"github.com/cgentry/gus/record/tenant"
	"testing"
)

func TestSetup(t *testing.T ){
	Register()
	if ! gdriver.IsRegistered( encryption.DRIVER_GROUP, DRIVER_NAME ) {
		t.Errorf("%s is not registered", DRIVER_NAME)
	}
}
func TestGenerate(t *testing.T) {

	user := tenant.NewTestUser()
	expected := "hello;" + user.Salt + ";SALT;Plaintext"
	pwd := encryption.GetDefaultDriver().EncryptPassword("hello", user.Salt)
	if pwd != expected {
		t.Errorf("Passwords don't match encrypted (%s != %s)", pwd, expected)
	}
}

func TestRepeatable(t *testing.T) {
	user := tenant.NewTestUser()
	pwd := encryption.GetDefaultDriver().EncryptPassword("123456", user.Salt)
	pwd2 := encryption.GetDefaultDriver().EncryptPassword("123456", user.Salt)
	if pwd != pwd2 {
		t.Errorf("Passwords didn't match: '%s' and '%s'", pwd, pwd2)
	}

}

func TestIsLongEnough(t *testing.T) {
	user := tenant.NewTestUser()
	pwd := encryption.GetDefaultDriver().EncryptPassword("hello", user.Salt)
	pwdLen := len(pwd)
	sbLen := len("hello;" + user.Salt + ";SALT;Plaintext")
	if pwdLen != sbLen {
		t.Errorf("PWD isn't long enough %d", pwdLen)
	}
}

func TestSimilarUserDifferntPwd(t *testing.T) {
	user := tenant.NewTestUser()
	pwd := encryption.GetDefaultDriver().EncryptPassword("123456", user.Salt)
	user2 := tenant.NewTestUser()
	pwd2 := encryption.GetDefaultDriver().EncryptPassword("123456", user2.Salt)
	if pwd == pwd2 {
		t.Errorf("Passwords for different users should not match: '%s' and '%s'", pwd, pwd2)
	}
}

func TestAfterChangingSalt(t *testing.T) {
	user := tenant.NewTestUser()
	drv := encryption.GetDefaultDriver()
	if drv.Id( ) != DRIVER_NAME {
		t.Errorf("Driver identity is wrong: %s != %s", DRIVER_NAME, drv.Id())
	}
	pwd := drv.EncryptPassword("123456", user.Salt)
	drv.Setup(`{ "Salt": "hello - this should screw up password" }`)
	pwd2 := drv.EncryptPassword("123456", user.Salt)

	if pwd == pwd2 {
		t.Errorf("Passwords with different salts should not match: '%s' and '%s'", pwd, pwd2)
	}
}
