package stringutil

import (
	"github.com/ghetzel/uuid"
)

type Uuid struct {
	uuid.UUID
}

func UuidFromBytes(b []byte) (*Uuid, error) {
	if uuid, err := uuid.FromBytes(b); err == nil {
		return &Uuid{uuid}, nil
	} else {
		return nil, err
	}
}

func UUID() *Uuid {
	return &Uuid{
		UUID: uuid.New(),
	}
}

func (self *Uuid) Bytes() []byte {
	return []byte(self.UUID[:])
}
