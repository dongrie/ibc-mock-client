package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "mock-client"
)

var (
	ErrInvalidHeaderHeight     = sdkerrors.Register(ModuleName, 5, "invalid header height")
	ErrInvalidProof            = sdkerrors.Register(ModuleName, 6, "invalid Mock proof")
	ErrProcessedTimeNotFound   = sdkerrors.Register(ModuleName, 8, "processed time not found")
	ErrProcessedHeightNotFound = sdkerrors.Register(ModuleName, 9, "processed height not found")
	ErrDelayPeriodNotPassed    = sdkerrors.Register(ModuleName, 10, "packet-specified delay period has not been reached")
)
