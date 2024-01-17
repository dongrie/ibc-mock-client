package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

const (
	ModuleName = "mock-client"
)

var (
	ErrInvalidHeaderHeight     = cosmossdkerrors.Register(ModuleName, 5, "invalid header height")
	ErrInvalidProof            = cosmossdkerrors.Register(ModuleName, 6, "invalid Mock proof")
	ErrProcessedTimeNotFound   = cosmossdkerrors.Register(ModuleName, 8, "processed time not found")
	ErrProcessedHeightNotFound = cosmossdkerrors.Register(ModuleName, 9, "processed height not found")
	ErrDelayPeriodNotPassed    = cosmossdkerrors.Register(ModuleName, 10, "packet-specified delay period has not been reached")
)
