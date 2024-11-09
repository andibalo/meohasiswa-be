package pkg

import (
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"net/http"
	"strconv"
)

type Pagination struct {
	CurrentCursor string `json:"current_cursor"`
	NextCursor    string `json:"next_cursor"`
}

func GetIntQueryParams(c *gin.Context, defValue int, key string) (int, error) {
	if c.Query(key) == "" {
		return defValue, nil
	}

	val, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return 0, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("%s should be integer, got error: %v", key, err)
	}

	if val > 0 {
		return val, nil
	}

	return defValue, nil
}

func GetBoolQueryParams(c *gin.Context, key string) (bool, error) {
	if c.Query(key) == "" {
		return false, nil
	}

	boolValue, err := strconv.ParseBool(c.Query(key))

	if err != nil {
		return false, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("%s should be boolean, got error: %v", key, err)
	}

	return boolValue, nil
}
