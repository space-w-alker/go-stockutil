package stringutil

import (
	"fmt"

	"github.com/satori/go.uuid"
)

type Uuid struct {
	uuid.UUID
}

func UUID() *Uuid {
	if u, err := uuid.NewV4(); err == nil {
		return &Uuid{u}
	} else {
		panic(fmt.Sprintf("uuid error: %v", err))
	}
}
