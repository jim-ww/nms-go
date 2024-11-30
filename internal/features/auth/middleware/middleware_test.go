package middleware_test

import (
	"testing"

	mw "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	"github.com/jim-ww/nms-go/internal/features/user"
)

func TestIsAllowed(t *testing.T) {
	testsIsAllowed := []struct {
		role     user.Role
		access   mw.AccessLevel
		expected bool
	}{
		{user.ROLE_USER, mw.AllowUnauthorized, true},
		{user.ROLE_USER, mw.OnlyUnauthorized, false},
		{user.ROLE_USER, mw.OnlyAuthorized, true},
		{user.ROLE_USER, mw.OnlyAdmins, false},
		{user.ROLE_ADMIN, mw.AllowUnauthorized, true},
		{user.ROLE_ADMIN, mw.OnlyUnauthorized, false},
		{user.ROLE_ADMIN, mw.OnlyAuthorized, true},
		{user.ROLE_ADMIN, mw.OnlyAdmins, true},
	}

	// testsCheckDynamicPaths := []struct {
	// 	path              string
	// 	expectedAccessLvl mw.AccessLevel
	// 	expectedOk        bool
	// }{
	// 	{"/static/somedir/somefile.js", mw.AllowUnauthorized, true},
	// 	{"/static/css/style.css", mw.AllowUnauthorized, true},
	// 	{"/api/notes/12345", mw.OnlyAuthorized, true},
	// 	{"/api/admin/notes/12345", mw.OnlyAdmins, true},
	// 	{"/api/admin/users/54321", mw.OnlyAdmins, true},
	// 	{"/unknown/path", 0, false},
	// }

	for _, test := range testsIsAllowed {
		result := mw.IsAllowed(test.role, test.access)
		if result != test.expected {
			t.Errorf("IsAllowed(%v, %v) = %v; want %v", test.role, test.access, result, test.expected)
		}
	}

	// for _, test := range testsCheckDynamicPaths {
	// 	result, ok := mw.CheckDynamicPath(test.path)
	// 	if result != test.expectedAccessLvl || ok != test.expectedOk {
	// 		t.Errorf("CheckDynamicPath(%v) = %v, %v; want %v, %v", test.path, result, ok, test.expectedAccessLvl, test.expectedOk)
	// 	}
	// }

}
