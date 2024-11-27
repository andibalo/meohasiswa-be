package notifsvc

type CreateNotifTemplateResp struct {
	Success bool `json:"success"`
}

type SendPushNotificationResp struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
