package code

import (
	"context"
	"time"

	"github.com/ory/kratos/selfservice/flow"

	"github.com/gofrs/uuid"
	errors "github.com/pkg/errors"

	"github.com/ory/x/randx"

	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/selfservice/flow/recovery"
	"github.com/ory/kratos/x"
)

type RecoveryCodeType int

const (
	RecoveryCodeTypeAdmin RecoveryCodeType = iota + 1
	RecoveryCodeTypeSelfService
)

type RecoveryCode struct {
	// ID represents the code's unique ID.
	//
	// required: true
	// type: string
	// format: uuid
	ID uuid.UUID `json:"id" db:"id" faker:"-"`

	// Code represents the recovery code
	Code string `json:"-" db:"code"`

	// RecoveryAddress links this code to a recovery address.
	// required: true
	RecoveryAddress *identity.RecoveryAddress `json:"recovery_address" belongs_to:"identity_recovery_addresses" fk_id:"RecoveryAddressID"`

	// CodeType is the type of the code - either "admin" or "selfservice"
	CodeType RecoveryCodeType `json:"-" faker:"-" db:"code_type"`

	// ExpiresAt is the time (UTC) when the code expires.
	// required: true
	ExpiresAt time.Time `json:"expires_at" faker:"time_type" db:"expires_at"`

	// IssuedAt is the time (UTC) when the code was issued.
	// required: true
	IssuedAt time.Time `json:"issued_at" faker:"time_type" db:"issued_at"`

	// CreatedAt is a helper struct field for gobuffalo.pop.
	CreatedAt time.Time `json:"-" faker:"-" db:"created_at"`
	// UpdatedAt is a helper struct field for gobuffalo.pop.
	UpdatedAt time.Time `json:"-" faker:"-" db:"updated_at"`
	// RecoveryAddressID is a helper struct field for gobuffalo.pop.
	RecoveryAddressID *uuid.UUID `json:"-" faker:"-" db:"identity_recovery_address_id"`
	// FlowID is a helper struct field for gobuffalo.pop.
	FlowID     uuid.UUID `json:"-" faker:"-" db:"selfservice_recovery_flow_id"`
	NID        uuid.UUID `json:"-"  faker:"-" db:"nid"`
	IdentityID uuid.UUID `json:"identity_id"  faker:"-" db:"identity_id"`
}

func (RecoveryCode) TableName(ctx context.Context) string {
	return "identity_recovery_codes"
}

func NewSelfServiceRecoveryCode(address *identity.RecoveryAddress, f *recovery.Flow, expiresIn time.Duration) *RecoveryCode {
	now := time.Now().UTC()
	var identityID = uuid.UUID{}
	var recoveryAddressID = uuid.UUID{}
	if address != nil {
		identityID = address.IdentityID
		recoveryAddressID = address.ID
	}
	return &RecoveryCode{
		ID:                x.NewUUID(),
		Code:              randx.MustString(8, randx.Numeric),
		RecoveryAddress:   address,
		ExpiresAt:         now.Add(expiresIn),
		IssuedAt:          now,
		IdentityID:        identityID,
		FlowID:            f.ID,
		RecoveryAddressID: &recoveryAddressID,
		CodeType:          RecoveryCodeTypeSelfService,
	}
}

func NewAdminRecoveryCode(identityID uuid.UUID, fID uuid.UUID, expiresIn time.Duration) *RecoveryCode {
	now := time.Now().UTC()
	return &RecoveryCode{
		ID:         x.NewUUID(),
		Code:       randx.MustString(8, randx.Numeric),
		ExpiresAt:  now.Add(expiresIn),
		IssuedAt:   now,
		IdentityID: identityID,
		FlowID:     fID,
		CodeType:   RecoveryCodeTypeAdmin,
	}
}

func (f RecoveryCode) Valid() error {
	if f.ExpiresAt.Before(time.Now()) {
		return errors.WithStack(flow.NewFlowExpiredError(f.ExpiresAt))
	}
	return nil
}
