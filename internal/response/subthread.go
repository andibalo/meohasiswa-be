package response

import "github.com/andibalo/meowhasiswa-be/internal/model"

type GetSubThreadListResponse struct {
	Data []model.SubThread `json:"subthreads"`
	Meta PaginationMeta    `json:"meta"`
}

type GetSubThreadByIDResponse struct {
	Data *model.SubThread `json:"subthread"`
}
