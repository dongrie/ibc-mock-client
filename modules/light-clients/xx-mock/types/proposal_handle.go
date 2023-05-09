package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// CheckSubstituteAndUpdateState always returns an error because Mock client doesn't support substitute.
func (cs ClientState) CheckSubstituteAndUpdateState(
	_ sdk.Context, _ codec.BinaryCodec, _, _ sdk.KVStore, _ exported.ClientState,
) error {
	return sdkerrors.Wrapf(clienttypes.ErrInvalidSubstitute, "cannot substribute Mock client")
}
