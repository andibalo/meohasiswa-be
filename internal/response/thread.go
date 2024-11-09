package response

import "github.com/andibalo/meowhasiswa-be/internal/model"

type GetThreadListResponse struct {
	Data []model.Thread `json:"threads"`
	Meta PaginationMeta `json:"meta"`
}
