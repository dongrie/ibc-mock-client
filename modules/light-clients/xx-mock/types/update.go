package types

import (
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// VerifyClientMessage checks if the clientMessage is of type Header
func (cs *ClientState) VerifyClientMessage(
	ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore,
	clientMsg exported.ClientMessage,
) error {
	switch msg := clientMsg.(type) {
	case *Header:
		return cs.verifyHeader(ctx, clientStore, cdc, msg)
	default:
		return clienttypes.ErrInvalidClientType
	}
}

// verifyHeader returns an error if:
// - header revision is not equal to latest header revision
func (cs *ClientState) verifyHeader(
	ctx sdk.Context, clientStore storetypes.KVStore, cdc codec.BinaryCodec,
	header *Header,
) error {
	if header.GetHeight().GetRevisionNumber() != cs.LatestHeight.RevisionNumber {
		return cosmossdkerrors.Wrapf(
			ErrInvalidHeaderHeight,
			"header height revision %d does not match latest header revision %d",
			header.GetHeight().GetRevisionNumber(), cs.LatestHeight.RevisionNumber,
		)

	}
	return nil
}

// UpdateState may be used to either create a consensus state for:
// - a future height greater than the latest client state height
// - a past height that was skipped during bisection
// If we are updating to a past height, a consensus state is created for that height to be persisted in client store
// If we are updating to a future height, the consensus state is created and the client state is updated to reflect
// the new latest height
// A list containing the updated consensus height is returned.
// UpdateState must only be used to update within a single revision, thus header revision number and trusted height's revision
// number must be the same. To update to a new revision, use a separate upgrade path
func (cs ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	header, ok := clientMsg.(*Header)
	if !ok {
		panic(fmt.Errorf("expected type %T, got %T", &Header{}, clientMsg))
	}

	// check for duplicate update
	if _, found := getConsensusState(clientStore, cdc, header.GetHeight()); found {
		// perform no-op
		return []exported.Height{header.GetHeight()}
	}

	height := header.GetHeight().(clienttypes.Height)
	if height.GT(cs.LatestHeight) {
		cs.LatestHeight = height
	}

	consensusState := &ConsensusState{
		Timestamp: header.Timestamp,
	}

	// set client state, consensus state and asssociated metadata
	setClientState(clientStore, cdc, &cs)
	setConsensusState(clientStore, cdc, consensusState, header.GetHeight())
	setConsensusMetadata(ctx, clientStore, header.GetHeight())

	return []exported.Height{height}
}

// UpdateStateOnMisbehaviour updates state upon misbehaviour, freezing the ClientState.
// For Mock, misbehaviour isn't defined and so this function never be called.
func (cs ClientState) UpdateStateOnMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, _ exported.ClientMessage) {
	panic(fmt.Errorf("misbehaviour is unexpected"))
}
