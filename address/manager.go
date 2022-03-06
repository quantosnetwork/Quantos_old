package address

import (
	"github.com/quantosnetwork/Quantosuint512"
	"context"
	"github.com/holiman/uint256"
)

type addressManager interface {
	New(args ...interface{}) *Address
	SetContext(ctx context.Context, aCtx AddressContext)
	ToUint512() *uint512.Int
	FromUint512() *uint512.Address
	ValidateAddress() bool
	Sign() error
	VerifySignature() error
	NewMaster() *Address
	Derive(aCtx AddressContext) *Address
	FromPublicKey() *Address
	AuthorizeUsage(secret string) bool
	GetBalance() uint256.Int
}

type contextManager interface {
	New(ctx context.Context, args ...interface{}) *AddressContext
	DecodeFromAddress(addr *Address) *AddressContext
	EncodeWithSecret(addr *Address, secret string)
}
