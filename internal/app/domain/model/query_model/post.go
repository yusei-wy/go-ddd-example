package query_model

import (
	"time"

	"ddd_go_example/internal/app/domain/model/value_object/value_post"
	"ddd_go_example/internal/app/domain/model/value_object/value_user"
)

type Post struct {
	PostId    value_post.PostId
	CreatedBy value_user.UserId
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
