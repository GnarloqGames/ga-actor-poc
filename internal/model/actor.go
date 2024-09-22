package model

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Address struct {
	Kind string
	ID   uuid.UUID
}

func (a Address) Hash() uuid.UUID {
	return uuid.NewHash(
		sha256.New(),
		uuid.NameSpaceOID,
		[]byte(fmt.Sprintf("%s:%s", a.Kind, a.ID.String())),
		5,
	)
}

type Actor interface {
	GetID() string
	GetKind() string
	Start(ctx context.Context)
	Destroy(ctx context.Context)
	Receive(ctx context.Context, msg proto.Message, res proto.Message) error
}
