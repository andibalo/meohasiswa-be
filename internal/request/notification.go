package request

type SendPushNotificationReq struct {
	NotificationTokens []string `json:"notification_tokens"`
	Title              string   `json:"title"`
	Content            string   `json:"content"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}
