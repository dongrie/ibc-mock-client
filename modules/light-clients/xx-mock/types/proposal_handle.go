package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// CheckSubstituteAndUpdateState always returns an error because Mock client doesn't support substitute.
func (cs ClientState) CheckSubstituteAndUpdateState(
	_ sdk.Context, _ codec.BinaryCodec, _, _ storetypes.KVStore, _ exported.ClientState,
) error {
	return cosmossdkerrors.Wrapf(clienttypes.ErrInvalidSubstitute, "cannot substribute Mock client")
}
