package main

import "errors"

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
