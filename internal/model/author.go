package model

import (
	"net/url"
	"pugo/internal/utils"
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

// NewDemoAuthor return a new author with demo fulfilled information.
func NewDemoAuthor(name string) *Author {
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
