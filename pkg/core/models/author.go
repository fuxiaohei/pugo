package models

import (
	"net/url"
	"pugo/pkg/utils"
	"strings"
)

type Author struct {
	Name        string            `toml:"name"`
	Email       string            `toml:"email"`
	Website     string            `toml:"website"`
	Bio         string            `toml:"bio"`
	Avatar      string            `toml:"avatar"`
	UseGravatar bool              `toml:"use_gravatar"`
	Slug        string            `toml:"slug"`
	Social      map[string]string `toml:"social"`
}

// NewAuthor return a new author with demo fulfilled information.
func NewAuthor(name string) *Author {
	return &Author{
		Name:        name,
		Email:       name + "@example.com",
		Bio:         "user bio",
		UseGravatar: true,
		Slug:        "/author/" + url.PathEscape(name) + "/",
		Social: map[string]string{
			"github": "https://github.com/" + name,
		},
	}
}

// AvatarLink returns the avatar link
func (a *Author) AvatarLink() string {
	if !a.UseGravatar {
		return a.Avatar
	}
	return "https://www.gravatar.com/avatar/" + utils.MD5String(strings.ToLower(a.Email))
}

// HasSocials returns true if the author has social information
func (a *Author) HasSocials() bool {
	return len(a.Social) > 0
}

// HasSocial returns true if the author has the specified social information
func (a *Author) HasSocial(key string) bool {
	return a.Social[key] != ""
}

// GetSocial returns the specified social information
func (a *Author) GetSocial(key string) string {
	return a.Social[key]
}

func (a *Author) Valid() bool {
	return a.Name != ""
}
