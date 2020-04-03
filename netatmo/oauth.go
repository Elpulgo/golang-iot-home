package netatmo

import "time"

type NetatmoOAuth struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	ExpiresIn    int      `json:"expires_in"`
}

func (OAuth *NetatmoOAuth) Expires() time.Time {

	return time.Now().Add(time.Duration(transformToHours(OAuth.ExpiresIn)))
}

func transformToHours(expiresIn int) time.Duration {
	return time.Duration(expiresIn * 1000 * 1000 * 1000)
}
