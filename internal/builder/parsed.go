package builder

import "haisite/internal/model"

// ParsedData is the parsed data.
type ParsedData struct {
	Posts []*model.Post
	Pages []*model.Page
}
