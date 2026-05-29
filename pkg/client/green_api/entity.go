package green_api

type Toggle string

const (
	ToggleYes Toggle = "yes"
	ToggleNo  Toggle = "no"
)

func (t Toggle) Bool() bool { return t == ToggleYes }

type InstanceState string

const (
	StateNotAuthorized InstanceState = "notAuthorized"
	StateAuthorized    InstanceState = "authorized"
	StateBlocked       InstanceState = "blocked"
	StateSleepMode     InstanceState = "sleepMode"
	StateStarting      InstanceState = "starting"
	StateYellowCard    InstanceState = "yellowCard" // partial/full send block imposed by WhatsApp
)

type PreviewSize string

const (
	PreviewLarge PreviewSize = "large"
	PreviewSmall PreviewSize = "small"
)

type TypingType string

const (
	TypingTypeRecording TypingType = "recording"
)

type Settings struct {
	WID                               string `json:"wid"`
	WebhookURL                        string `json:"webhookUrl"`
	WebhookURLToken                   string `json:"webhookUrlToken"`
	DelaySendMessagesMilliseconds     int    `json:"delaySendMessagesMilliseconds"`
	LinkPreview                       bool   `json:"linkPreview"`
	AutoTyping                        int    `json:"autoTyping"`
	OutgoingWebhook                   Toggle `json:"outgoingWebhook"`
	OutgoingMessageWebhook            Toggle `json:"outgoingMessageWebhook"`
	OutgoingAPIMessageWebhook         Toggle `json:"outgoingAPIMessageWebhook"`
	IncomingWebhook                   Toggle `json:"incomingWebhook"`
	StateWebhook                      Toggle `json:"stateWebhook"`
	StatusInstanceWebhook             Toggle `json:"statusInstanceWebhook"`
	DeviceWebhook                     Toggle `json:"deviceWebhook"`
	PollMessageWebhook                Toggle `json:"pollMessageWebhook"`
	IncomingCallWebhook               Toggle `json:"incomingCallWebhook"`
	IncomingBlockWebhook              Toggle `json:"incomingBlockWebhook"`
	EditedMessageWebhook              Toggle `json:"editedMessageWebhook"`
	DeletedMessageWebhook             Toggle `json:"deletedMessageWebhook"`
	MarkIncomingMessagesReaded        Toggle `json:"markIncomingMessagesReaded"`
	MarkIncomingMessagesReadedOnReply Toggle `json:"markIncomingMessagesReadedOnReply"`
	EnableMessagesHistory             Toggle `json:"enableMessagesHistory"`
	KeepOnlineStatus                  Toggle `json:"keepOnlineStatus"`
	SharedSession                     Toggle `json:"sharedSession"`
}

type StateInstanceResponse struct {
	State InstanceState `json:"stateInstance"`
}

func (s *StateInstanceResponse) IsAuthorized() bool {
	return s.State == StateAuthorized
}

type CustomPreview struct {
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	Link          string `json:"link,omitempty"`
	URLFile       string `json:"urlFile,omitempty"`
	JPEGThumbnail string `json:"jpegThumbnail,omitempty"`
}

type SendMessageRequest struct {
	ChatID          string         `json:"chatId"`
	Message         string         `json:"message"`
	QuotedMessageID string         `json:"quotedMessageId,omitempty"`
	LinkPreview     *bool          `json:"linkPreview,omitempty"`
	TypePreview     PreviewSize    `json:"typePreview,omitempty"`
	CustomPreview   *CustomPreview `json:"customPreview,omitempty"`
	TypingTime      int            `json:"typingTime,omitempty"`
}

type SendFileByURLRequest struct {
	ChatID          string     `json:"chatId"`
	URLFile         string     `json:"urlFile"`
	FileName        string     `json:"fileName"`
	Caption         string     `json:"caption,omitempty"`
	QuotedMessageID string     `json:"quotedMessageId,omitempty"`
	TypingTime      int        `json:"typingTime,omitempty"`
	TypingType      TypingType `json:"typingType,omitempty"`
}

type SendResponse struct {
	IDMessage string `json:"idMessage"`
}
