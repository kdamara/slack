package slack

import (
	"errors"
	"net/url"
)

// UserProfile contains all the information details of a given user
type UserProfile struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	RealName           string `json:"real_name"`
	RealNameNormalized string `json:"real_name_normalized"`
	Email              string `json:"email"`
	Skype              string `json:"skype"`
	Phone              string `json:"phone"`
	Image24            string `json:"image_24"`
	Image32            string `json:"image_32"`
	Image48            string `json:"image_48"`
	Image72            string `json:"image_72"`
	Image192           string `json:"image_192"`
	ImageOriginal      string `json:"image_original"`
	Title              string `json:"title"`
	BotID              string `json:"bot_id,omitempty"`
	ApiAppID           string `json:"api_app_id,omitempty"`
}

// User contains all the information of a user
type User struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Deleted           bool        `json:"deleted"`
	Color             string      `json:"color"`
	RealName          string      `json:"real_name"`
	TZ                string      `json:"tz,omitempty"`
	TZLabel           string      `json:"tz_label"`
	TZOffset          int         `json:"tz_offset"`
	Profile           UserProfile `json:"profile"`
	IsBot             bool        `json:"is_bot"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	Has2FA            bool        `json:"has_2fa"`
	HasFiles          bool        `json:"has_files"`
	Presence          string      `json:"presence"`
}

// UserPresence contains details about a user online status
type UserPresence struct {
	Presence        string   `json:"presence,omitempty"`
	Online          bool     `json:"online,omitempty"`
	AutoAway        bool     `json:"auto_away,omitempty"`
	ManualAway      bool     `json:"manual_away,omitempty"`
	ConnectionCount int      `json:"connection_count,omitempty"`
	LastActivity    JSONTime `json:"last_activity,omitempty"`
}

type UserIdentityResponse struct {
	User UserIdentity `json:"user"`
	Team TeamIdentity `json:"team"`
	SlackResponse
}

type UserIdentity struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Image24  string `json:"image_24"`
	Image32  string `json:"image_32"`
	Image48  string `json:"image_48"`
	Image72  string `json:"image_72"`
	Image192 string `json:"image_192"`
	Image512 string `json:"image_512"`
}

type TeamIdentity struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	Image34       string `json:"image_34"`
	Image44       string `json:"image_44"`
	Image68       string `json:"image_68"`
	Image88       string `json:"image_88"`
	Image102      string `json:"image_102"`
	Image132      string `json:"image_132"`
	Image230      string `json:"image_230"`
	ImageDefault  bool   `json:"image_default"`
	ImageOriginal string `json:"image_original"`
}

type userResponseFull struct {
	Members      []User                  `json:"members,omitempty"` // ListUsers
	User         `json:"user,omitempty"` // GetUserInfo
	UserPresence                         // GetUserPresence
	SlackResponse
	Offset string `json:"offset"`
}

func userRequest(path string, values url.Values, debug bool) (*userResponseFull, error) {
	response := &userResponseFull{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// GetUserPresence will retrieve the current presence status of given user.
func (api *Client) GetUserPresence(user string) (*UserPresence, error) {
	values := url.Values{
		"token": {api.config.token},
		"user":  {user},
	}
	response, err := userRequest("users.getPresence", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.UserPresence, nil
}

// GetUserInfo will retrieve the complete user information
func (api *Client) GetUserInfo(user string) (*User, error) {
	values := url.Values{
		"token": {api.config.token},
		"user":  {user},
	}
	response, err := userRequest("users.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.User, nil
}

// GetUsers returns the list of users (with their detailed information)
func (api *Client) GetUsers(offset string) ([]User, string, error) {
	values := url.Values{
		"token": {api.config.token},
		"limit": {"1000"},
    "offset": {offset},
	}
	response, err := userRequest("users.list", values, api.debug)
	if err != nil {
		return nil, ``, err
	}
	return response.Members, response.Offset, nil
}

// SetUserAsActive marks the currently authenticated user as active
func (api *Client) SetUserAsActive() error {
	values := url.Values{
		"token": {api.config.token},
	}
	_, err := userRequest("users.setActive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// SetUserPresence changes the currently authenticated user presence
func (api *Client) SetUserPresence(presence string) error {
	values := url.Values{
		"token":    {api.config.token},
		"presence": {presence},
	}
	_, err := userRequest("users.setPresence", values, api.debug)
	if err != nil {
		return err
	}
	return nil

}

// GetUserIdentity will retrieve user info available per identity scopes
func (api *Client) GetUserIdentity() (*UserIdentityResponse, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response := &UserIdentityResponse{}
	err := post("users.identity", values, response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}
