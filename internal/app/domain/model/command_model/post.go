package command_model

import (
	"ddd_go_example/internal/app/domain/model/value_object/value_post"
	"ddd_go_example/internal/app/domain/model/value_object/value_user"
	"time"
)

type Post struct {
	PostId    value_post.PostId
	CreatedBy value_user.UserId
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
