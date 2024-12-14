package middleware_test

import (
	"testing"

	mw "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
)

func TestIsAllowed(t *testing.T) {
	testsIsAllowed := []struct {
		role     role.Role
		access   mw.AccessLevel
		expected bool
	}{
		{role.ROLE_USER, mw.AllowUnauthorized, true},
		{role.ROLE_USER, mw.OnlyUnauthorized, false},
		{role.ROLE_USER, mw.OnlyAuthorized, true},
		{role.ROLE_USER, mw.OnlyAdmins, false},
		{role.ROLE_ADMIN, mw.AllowUnauthorized, true},
		{role.ROLE_ADMIN, mw.OnlyUnauthorized, false},
		{role.ROLE_ADMIN, mw.OnlyAuthorized, true},
		{role.ROLE_ADMIN, mw.OnlyAdmins, true},
	}
	_ = testsIsAllowed // TODO remove

	// testsCheckDynamicPaths := []struct {
	// 	path              string
	// 	expectedAccessLvl mw.AccessLevel
	// 	expectedOk        bool
	// }{
	// 	{"/web/js/somefile.js", mw.AllowUnauthorized, true},
	// 	{"/web/css/style.css", mw.AllowUnauthorized, true},
	// 	{"/api/notes/12345", mw.OnlyAuthorized, true},
	// 	{"/api/admin/notes/12345", mw.OnlyAdmins, true},
	// 	{"/api/admin/users/54321", mw.OnlyAdmins, true},
	// 	{"/unknown/path", 0, false},
	// }

	// for _, test := range testsIsAllowed {
	// 	result := mw.IsAllowed(test.role, test.access)
	// 	if result != test.expected {
	// 		t.Errorf("IsAllowed(%v, %v) = %v; want %v", test.role, test.access, result, test.expected)
	// 	}
	// }

	// for _, test := range testsCheckDynamicPaths {
	// 	result, ok := mw.CheckDynamicPath(test.path)
	// 	if result != test.expectedAccessLvl || ok != test.expectedOk {
	// 		t.Errorf("CheckDynamicPath(%v) = %v, %v; want %v, %v", test.path, result, ok, test.expectedAccessLvl, test.expectedOk)
	// 	}
	// }

}
