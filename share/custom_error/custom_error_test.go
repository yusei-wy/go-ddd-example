package customerror

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type UserId string

func Test_NotFoundError(t *testing.T) {
	userId := UserId("user_0")

	err := NotFoundError("User", userId, nil)

	if diff := cmp.Diff("User not found UserId(\"user_0\")", err.Error()); diff != "" {
		t.Errorf("Error: (-got +want)\n%s", diff)
	}

	t.Parallel()
}
