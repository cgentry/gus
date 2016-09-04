// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//

package jsonfile

import (
	"fmt"
	"github.com/cgentry/gus/library/storage"
	"github.com/cgentry/gus/record/tenant"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func getRidOfFile(fname string) {
	if os.Remove(fname) == nil {
		fmt.Println("Deleted temp file ", fname)
	}
}
func TestSimpleRegisterCycle(t *testing.T) {

	// need a new test file
	fp, err := ioutil.TempFile("", "jsonstore_")
	if err != nil {
		t.Errorf("Could not create temporary file. '%s'", err.Error())
	}
	fname := fp.Name()
	fp.Close()
	defer getRidOfFile(fname)

	fmt.Printf("File is %s\n", fname)
	dbGeneralCon, err := NewJsonFileDriver().Open(fname, ``)

	Convey("Create User", t, func() {
		So(err, ShouldBeNil)
		dbConn, ok := dbGeneralCon.(*JsonFileConn) // To force getting at the raw calls...
		So(ok, ShouldBeTrue)

		user := tenant.NewTestUser()
		user.SetDomain("Register")
		user.SetToken("TestToken")
		user.SetName("Just a test name")
		user.SetEmail("et@home.com")
		user.SetLoginName("justlogin")

		serr := dbConn.UserInsert(user) // Register new user
		So(serr, ShouldBeNil)

		// FETCH BY EMAIL
		user2, err := dbConn.UserFetch(storage.MatchAnyDomain, storage.FieldGUID, user.Guid)
		So(err, ShouldBeNil)
		So(user2.Domain, ShouldEqual, user.Domain)
		So(user2.Token, ShouldEqual, user.Token)
		So(user2.FullName, ShouldEqual, user.FullName)

		// Fetch by TOKEN
		user3, err := dbConn.UserFetch(storage.MatchAnyDomain, storage.FieldToken, user.Token)
		So(err, ShouldBeNil)
		So(user3.Domain, ShouldEqual, user.Domain)
		So(user3.Token, ShouldEqual, user.Token)
		So(user3.FullName, ShouldEqual, user.FullName)

		// FETCH BY EMAIL
		user4, err := dbConn.UserFetch(user.Domain, storage.FieldEmail, user.Email)
		So(err, ShouldBeNil)
		So(user4.Domain, ShouldEqual, user.Domain)
		So(user4.Token, ShouldEqual, user.Token)
		So(user4.FullName, ShouldEqual, user.FullName)

		// FETCH BY Login name
		user5, err := dbConn.UserFetch(user.Domain, storage.FieldLogin, user.LoginName)
		So(err, ShouldBeNil)
		So(user5.Domain, ShouldEqual, user.Domain)
		So(user5.Token, ShouldEqual, user.Token)
		So(user5.FullName, ShouldEqual, user.FullName)

	})
	err = os.Remove(fname)
	if err == nil {
		fmt.Println("Deleted temp file ", fname)
	} else {
		fmt.Println("ERROR! ", err.Error())
	}
}
