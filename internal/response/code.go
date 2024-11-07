package response

const (
	Success           Code = "CS0000"
	ServerError       Code = "CS0001"
	BadRequest        Code = "CS0002"
	InvalidRequest    Code = "CS0004"
	Failed            Code = "CS0073"
	Pending           Code = "CS0050"
	InvalidInputParam Code = "CS0032"
	DuplicateUser     Code = "CS0033"
	NotFound          Code = "CS0034"

	Unauthorized   Code = "CS0502"
	Forbidden      Code = "CS0503"
	GatewayTimeout Code = "CS0048"
)

type Code string

var codeMap = map[Code]string{
	Success:           "success",
	Failed:            "failed",
	Pending:           "pending",
	BadRequest:        "bad or invalid request",
	Unauthorized:      "Unauthorized Token",
	GatewayTimeout:    "Gateway Timeout",
	ServerError:       "Internal Server Error",
	InvalidInputParam: "Other invalid argument",
	DuplicateUser:     "duplicate user",
	NotFound:          "Not found",
}

func (c Code) AsString() string {
	return string(c)
}

func (c Code) GetStatus() string {
	switch c {
	case Success:
		return "SUCCESS"

	default:
		return "FAILED"
	}
}

func (c Code) GetMessage() string {
	return codeMap[c]
}

func (c Code) GetVersion() string {
	return "1"
}
