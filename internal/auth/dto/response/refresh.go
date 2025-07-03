package response

type Refresh struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}