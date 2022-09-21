package code

import (
	"context"

	"github.com/gofrs/uuid"
)

type (
	RecoveryCodePersister interface {
		CreateRecoveryCode(ctx context.Context, code *RecoveryCode) error
		UseRecoveryCode(ctx context.Context, fID uuid.UUID, code string) (*RecoveryCode, error)
	}

	RecoveryCodePersistenceProvider interface {
		RecoveryCodePersister() RecoveryCodePersister
	}
)
