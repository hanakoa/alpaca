package services

type PickOptionRequest struct {
	// Account can be an email address, phone number, or username
	Account string `json:"account"`
}