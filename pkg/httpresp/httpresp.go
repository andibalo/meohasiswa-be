package httpresp

import (
	"fmt"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"net/http"
	"reflect"
	"time"
)

const (
	StatusCodeCtxKey = "httpStatusCode"
)

type Meta struct {
	Path       string `json:"path"`
	Code       string `json:"code"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty" swaggerignore:"true"`
	Timestamp  string `json:"timestamp"`
}

type Response struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Success    string      `json:"success"`
}

type Pagination struct {
	CurrentPage     int64   `json:"current_page"`
	CurrentElements int64   `json:"current_elements"`
	TotalPages      int64   `json:"total_pages"`
	TotalElements   int64   `json:"total_elements"`
	SortBy          string  `json:"sort_by"`
	CursorStart     *string `json:"cursor_start,omitempty"`
	CursorEnd       *string `json:"cursor_end,omitempty"`
}

// HTTPErrResp http error response
type HTTPErrResp struct {
	Meta Meta `json:"metadata"`
}

func HttpRespError(c *gin.Context, err error) {

	statusCode := http.StatusInternalServerError
	respCode := response.ServerError.AsString()

	oopsErr, ok := oops.AsOops(err)

	if ok {
		respCode = oopsErr.Code()
		errCtx := oopsErr.Context()

		sc, exists := errCtx[StatusCodeCtxKey]
		if exists {
			statusCode = sc.(int)
		}
	}

	jsonErrResp := &HTTPErrResp{
		Meta: Meta{
			Path:       c.Request.URL.Path,
			Code:       respCode,
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
			Message:    fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)),
			Error:      err.Error(),
			Timestamp:  time.Now().Format(time.RFC3339),
		},
	}

	c.Set("status_code", statusCode)
	c.Set("status", http.StatusText(statusCode))
	c.Set("error", fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)))

	c.JSON(statusCode, jsonErrResp)
}

func HttpRespSuccess(c *gin.Context, data interface{}, pagination *Pagination) {

	//check typenya slice / array , soalnya kalo bukan slice / array ga perlu di dikosongin datanya soalnya udah kena error di handler
	kind := reflect.ValueOf(data).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		//check kalo data nya nil / kosong
		if data == nil || reflect.ValueOf(data).IsNil() {
			//kalo data arraynya kosong returnnya "data": []
			data = []interface{}{}
		}
	}

	c.JSON(http.StatusOK, Response{
		Data:       data,
		Pagination: pagination,
		Success:    "success",
	})
}

func ResetPagination() *Pagination {

	return &Pagination{
		CurrentPage:     1,
		CurrentElements: 0,
		TotalPages:      0,
		TotalElements:   0,
	}
}
