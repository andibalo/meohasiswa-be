package request

type TestLogWithBodyReq struct {
	Msg string `json:"msg"`
}

type TestCallNotifServiceReq struct {
	TemplateName string `json:"template_name"`
}

type GetUserProfileReq struct {
	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type CreateUserDeviceReq struct {
	Brand                string `json:"brand"`
	Type                 string `json:"type"`
	Model                string `json:"model"`
	NotificationToken    string `json:"notification_token" binding:"required"`
	IsNotificationActive bool   `json:"is_notification_active"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type GetUserDevicesReq struct {
	NotificationToken string `json:"notification_token"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type BanUserReq struct {
	BanUserID string `json:"-"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type UnBanUserReq struct {
	UnBanUserID string `json:"-"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}
