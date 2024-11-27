package notifsvc

type CreateNotifTemplateReq struct {
	TemplateName string `json:"template_name"`
}

type SendPushNotificationReq struct {
	NotificationToken string `json:"notification_token"`
	Title             string `json:"title"`
	Content           string `json:"content"`
}
