package notifsvc

type CreateNotifTemplateReq struct {
	TemplateName string `json:"template_name"`
}

type SendPushNotificationReq struct {
	NotificationTokens []string          `json:"notification_tokens"`
	Title              string            `json:"title"`
	Content            string            `json:"content"`
	Data               map[string]string `json:"data"`
}
