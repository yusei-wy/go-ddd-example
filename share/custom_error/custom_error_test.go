package customerror

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type UserID string

func Test_NotFoundError(t *testing.T) {
	userID := UserID("user_0")

	err := NotFoundError("User", userID, nil)

	if diff := cmp.Diff("User not found UserId(\"user_0\")", err.Error()); diff != "" {
		t.Errorf("Error: (-got +want)\n%s", diff)
	}

	t.Parallel()
}
