package utilities

/* ErrorReponse */
type ErrorReponse struct {
	Status        string `json:"status"`
	ErrorType     string `json:"error_type"`
	StatusCode    int    `json:"status_code"`
	ErrorMsg      string `json:"error_msg"`
	AdditionalMsg string `json:"additional_msg"`
}
type SuccessReponse struct {
	Status     string                 `json:"status"`
	StatusCode int                    `json:"status_code"`
	Msg        string                 `json:"msg"`
	Data       map[string]interface{} `json:"data"`
}

var SUCCESS = map[string]SuccessReponse{
	"SUCCESS": {"success", 200, "", nil},
}

var ERRORS = map[string]ErrorReponse{
	"INTERNAL_ERROR": {"error", "INTERNAL_ERROR", 400, "Internal error", ""},
}

var SUCCESS_LOGIN = map[string]interface{}{
	"status":      "success",
	"csrf":        "",
	"sessionid":   "",
	"status_code": 100,
}
