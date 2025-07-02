package dto

type WebhookPayload struct {
	IPAddress string `json:"ip_address"`
	Event     string `json:"event"`
}
