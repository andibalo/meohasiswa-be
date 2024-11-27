package request

type SendPushNotificationReq struct {
	NotificationToken string `json:"notification_token"`
	Title             string `json:"title"`
	Content           string `json:"content"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}
