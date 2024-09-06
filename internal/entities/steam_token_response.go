package entities

type SteamTokenResponse struct {
	TokenId string `json:"token_id,omitempty"`
	Message string `json:"message"`
}
