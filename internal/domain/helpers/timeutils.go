package helpers

import (
	"time"

	"github.com/critma/tgsheduler/internal/store"
)

func TimeToUserTZ(user *store.User, t time.Time) time.Time {
	userTZ := time.FixedZone("User_loc", int(time.Hour.Seconds())*int(user.UTC))
	result := t.In(userTZ)
	return result
}
