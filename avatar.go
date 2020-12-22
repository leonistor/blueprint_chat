package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL is the error retrurned when
// an Avatar instance cannot provide an avatar URL
var ErrNoAvatarURL = errors.New("chat: unable to get avatar URL")

// Avatar represents types capable
// of representing user profile pictures
type Avatar interface {
	// GetAvatarURL gets the avatar url for the specified client
	// or ErrNoAvatarURL
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar gets avatar from an OAuth2 provider
type AuthAvatar struct{}

// UseAuthAvatar selects auth as avatar source
var UseAuthAvatar AuthAvatar

// GetAvatarURL for AuthAvatar
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	if urlStr, ok := url.(string); ok {
		return urlStr, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar gets avatar from gravatar.com
type GravatarAvatar struct{}

// UseGravatarAvatar selects gravatar as avatar source
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL for GravatarAvatar
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	email, ok := c.userData["email"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	if emailStr, ok := email.(string); ok {
		m := md5.New()
		io.WriteString(m, strings.ToLower(emailStr))
		return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
	}
	return "", ErrNoAvatarURL
}
