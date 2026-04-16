package amcrest

import (
	"context"
	"fmt"
	"io"
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

// ModifyUser modifies attributes of an existing user. The params map should
// contain user fields such as "user.Memo", "user.Group", etc.
// CGI: userManager.cgi?action=modifyUser&name=X&user.Memo=Y&...
func (s *UserService) ModifyUser(ctx context.Context, name string, params map[string]string) error {
	qv := url.Values{
		"name": {name},
	}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "userManager.cgi", "modifyUser", qv)
}

// ModifyPasswordByManager changes a user's password using manager credentials.
// CGI: userManager.cgi?action=modifyPasswordByManager&name=X&pwd=Y&managerName=M&managerPwd=P
func (s *UserService) ModifyPasswordByManager(ctx context.Context, managerName, managerPwd, userName, newPwd string) error {
	return s.client.cgiAction(ctx, "userManager.cgi", "modifyPasswordByManager", url.Values{
		"name":       {userName},
		"pwd":        {newPwd},
		"managerName": {managerName},
		"managerPwd":  {managerPwd},
	})
}

// SetLoginAuthPolicy sets the login authentication policy.
// Policy values: 0 = off, 1 = on, etc. as defined by the device.
func (s *UserService) SetLoginAuthPolicy(ctx context.Context, policy int) error {
	return s.client.setConfig(ctx, map[string]string{
		"LoginAuthCtrl.PriSvrPolicy": fmt.Sprintf("%d", policy),
	})
}

// ExportAccounts exports user account data as binary. FileType specifies the
// export format (e.g., 0 for default).
// POST /cgi-bin/api/userManager/accountFileExport with JSON {FileType: N}.
func (s *UserService) ExportAccounts(ctx context.Context, fileType int) ([]byte, error) {
	reqBody := struct {
		FileType int `json:"FileType"`
	}{FileType: fileType}

	resp, err := s.client.postRawResponse(ctx, "/cgi-bin/api/userManager/accountFileExport", reqBody)
	if err != nil {
		return nil, fmt.Errorf("amcrest: export accounts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading export body: %w", err)
	}
	return data, nil
}
