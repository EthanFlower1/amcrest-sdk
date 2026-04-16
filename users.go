package amcrest

import (
	"context"
	"net/url"
)

// GetUserInfo returns information about a single user as key-value pairs.
// CGI: userManager.cgi?action=getUserInfo&name=<name>
func (s *UserService) GetUserInfo(ctx context.Context, name string) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "userManager.cgi", "getUserInfo", url.Values{
		"name": {name},
	})
	if err != nil {
		return nil, err
	}
	return parseKVWithPrefix(body, "users."), nil
}

// GetAllUsers returns raw body of all user information.
// The response uses an indexed format (e.g., users[0].Name=admin) that is best
// consumed as raw text by the caller.
// CGI: userManager.cgi?action=getUserInfoAll
func (s *UserService) GetAllUsers(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "userManager.cgi", "getUserInfoAll", nil)
}

// GetActiveUsers returns raw body of all currently active/logged-in users.
// CGI: userManager.cgi?action=getActiveUserInfoAll
func (s *UserService) GetActiveUsers(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "userManager.cgi", "getActiveUserInfoAll", nil)
}

// GetGroupInfo returns information about a single user group as key-value pairs.
// CGI: userManager.cgi?action=getGroupInfo&name=<name>
func (s *UserService) GetGroupInfo(ctx context.Context, name string) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "userManager.cgi", "getGroupInfo", url.Values{
		"name": {name},
	})
	if err != nil {
		return nil, err
	}
	return parseKVWithPrefix(body, "group."), nil
}

// GetAllGroups returns raw body of all group information.
// The response uses an indexed format that is best consumed as raw text.
// CGI: userManager.cgi?action=getGroupInfoAll
func (s *UserService) GetAllGroups(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "userManager.cgi", "getGroupInfoAll", nil)
}

// AddUser creates a new user with the given name, password, and group.
// CGI: userManager.cgi?action=addUser&user.Name=X&user.Password=Y&user.Group=Z&user.Sharable=true
func (s *UserService) AddUser(ctx context.Context, name, password, group string) error {
	return s.client.cgiAction(ctx, "userManager.cgi", "addUser", url.Values{
		"user.Name":     {name},
		"user.Password": {password},
		"user.Group":    {group},
		"user.Sharable": {"true"},
	})
}

// DeleteUser removes a user by name.
// CGI: userManager.cgi?action=deleteUser&name=<name>
func (s *UserService) DeleteUser(ctx context.Context, name string) error {
	return s.client.cgiAction(ctx, "userManager.cgi", "deleteUser", url.Values{
		"name": {name},
	})
}

// ModifyPassword changes a user's password.
// CGI: userManager.cgi?action=modifyPassword&name=X&pwd=NEW&pwdOld=OLD
func (s *UserService) ModifyPassword(ctx context.Context, name, oldPwd, newPwd string) error {
	return s.client.cgiAction(ctx, "userManager.cgi", "modifyPassword", url.Values{
		"name":   {name},
		"pwd":    {newPwd},
		"pwdOld": {oldPwd},
	})
}
