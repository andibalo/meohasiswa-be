package response

type PaginationMeta struct {
	CurrentCursor string `json:"current_cursor"`
	NextCursor    string `json:"next_cursor"`
}
