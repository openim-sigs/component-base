package uuid

import (
	"github.com/google/uuid"

	"github.com/openim-sigs/component-base/pkg/types"
)

func NewUUID() types.UID {
	return types.UID(uuid.New().String())
}
