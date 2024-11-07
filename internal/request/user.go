package request

type TestLogWithBodyReq struct {
	Msg string `json:"msg"`
}

type TestCallNotifServiceReq struct {
	TemplateName string `json:"template_name"`
}
