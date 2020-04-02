package netatmo

import "time"

type NetatmoOAuth struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	ExpiresIn    int      `json:"expires_in"`
}

func (OAuth *NetatmoOAuth) Expires() time.Time {
	return time.Now().Add(time.Duration(OAuth.ExpiresIn))
}
