package command_model

import (
	"time"

	"ddd_go_example/internal/app/domain/model/value_object/post"
	"ddd_go_example/internal/app/domain/model/value_object/user"
)

type Post struct {
	PostId    post.PostId
	CreatedBy user.UserId
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
