package domain

type GreenAPICredentials struct {
	IDInstance       string
	APITokenInstance string
}

type GreenAPIGetInput struct {
	Credentials GreenAPICredentials
}

type GreenAPISendMessageInput struct {
	Credentials     GreenAPICredentials
	ChatID          string
	Message         string
	QuotedMessageID string
	LinkPreview     *bool
}

type GreenAPISendFileByURLInput struct {
	Credentials     GreenAPICredentials
	ChatID          string
	URLFile         string
	FileName        string
	Caption         string
	QuotedMessageID string
}
