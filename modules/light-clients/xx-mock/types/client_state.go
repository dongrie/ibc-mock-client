package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	cosmossdkerrors "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

const (
	Mock string = "mock-client"
)

var _ exported.ClientState = (*ClientState)(nil)

// NewClientState creates a new ClientState instance.
func NewClientState(latestHeight clienttypes.Height) *ClientState {
	return &ClientState{
		LatestHeight: latestHeight,
	}
}

// ClientType returns a type of the client.
func (cs ClientState) ClientType() string {
	return Mock
}

// GetLatestHeight returns the latest height.
// Return exported.Height to satisfy ClientState interface
func (cs ClientState) GetLatestHeight() exported.Height {
	return cs.LatestHeight
}

// GetTimestampAtHeight returns the timestamp in nanoseconds of the consensus state at the given height.
func (cs ClientState) GetTimestampAtHeight(
	ctx sdk.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
) (uint64, error) {
	// get consensus state at height from clientStore to check for expiry
	consState, found := getConsensusState(clientStore, cdc, height)
	if !found {
		return 0, cosmossdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound, "height (%s)", height)
	}
	return consState.GetTimestamp(), nil
}

// Status returns the status of the mock client.
// It always returns active.
func (cs ClientState) Status(_ sdk.Context, _ storetypes.KVStore, _ codec.BinaryCodec) exported.Status {
	return exported.Active
}

// Validate performs a basic validation of the client state fields.
func (cs ClientState) Validate() error {
	return nil
}

// ZeroCustomFields returns a ClientState that is a copy of the current ClientState
// with all client customizable fields zeroed out (but Mock ClientState has no such field)
func (cs ClientState) ZeroCustomFields() exported.ClientState {
	return &ClientState{
		LatestHeight: cs.LatestHeight,
	}
}

// Initialize will check that initial consensus state is equal to the latest consensus state of the initial client.
func (cs ClientState) Initialize(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, consState exported.ConsensusState) error {
	consensusState, ok := consState.(*ConsensusState)
	if !ok {
		return cosmossdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "invalid initial consensus state. expected type: %T, got: %T",
			&ConsensusState{}, consState)
	}

	setClientState(clientStore, cdc, &cs)
	setConsensusState(clientStore, cdc, consensusState, cs.GetLatestHeight())
	setConsensusMetadata(ctx, clientStore, cs.GetLatestHeight())

	return nil

}

// VerifyMembership is a generic proof verification method which verifies a proof of the existence of a value at a given CommitmentPath at the specified height.
// The caller is expected to construct the full CommitmentPath from a CommitmentPrefix and a standardized path (as defined in ICS 24).
func (cs ClientState) VerifyMembership(
	ctx sdk.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
	value []byte,
) error {
	if cs.GetLatestHeight().LT(height) {
		return cosmossdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height,
		)
	}

	if err := verifyDelayPeriodPassed(ctx, clientStore, height, delayTimePeriod, delayBlockPeriod); err != nil {
		return err
	}

	if _, found := getConsensusState(clientStore, cdc, height); !found {
		return cosmossdkerrors.Wrap(clienttypes.ErrConsensusStateNotFound, "please ensure the proof was constructed against a height that exists on the client")
	}

	// sha256(abi.encodePacked(height.toUint128(), sha256(prefix), sha256(path), sha256(value)))
	revisionNumber := height.GetRevisionNumber()
	revisionHeight := height.GetRevisionHeight()

	heightBuf := make([]byte, 16)
	binary.BigEndian.PutUint64(heightBuf[:8], revisionNumber)
	binary.BigEndian.PutUint64(heightBuf[8:], revisionHeight)

	merklePath := path.(commitmenttypes.MerklePath)
	mPrefix, err := merklePath.GetKey(0)
	if err != nil {
		return cosmossdkerrors.Wrapf(err, "invalid merkle path key at index 0")
	}
	mPath, err := merklePath.GetKey(1)
	if err != nil {
		return cosmossdkerrors.Wrapf(err, "invalid merkle path key at index 1")
	}

	hashPrefix := sha256.Sum256([]byte(mPrefix))
	hashPath := sha256.Sum256([]byte(mPath))
	hashValue := sha256.Sum256([]byte(value))

	var combined []byte
	combined = append(combined, heightBuf...)
	combined = append(combined, hashPrefix[:]...)
	combined = append(combined, hashPath[:]...)
	combined = append(combined, hashValue[:]...)
	h := sha256.Sum256(combined)

	if !bytes.Equal(proof, h[:]) {
		return cosmossdkerrors.Wrapf(ErrInvalidProof, "expected the proof '%X', actually got '%X'", h, proof)
	}

	return nil
}

// VerifyNonMembership is a generic proof verification method which verifies the absence of a given CommitmentPath at a specified height.
// The caller is expected to construct the full CommitmentPath from a CommitmentPrefix and a standardized path (as defined in ICS 24).
func (cs ClientState) VerifyNonMembership(
	ctx sdk.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
) error {
	if cs.GetLatestHeight().LT(height) {
		return cosmossdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height,
		)
	}

	if err := verifyDelayPeriodPassed(ctx, clientStore, height, delayTimePeriod, delayBlockPeriod); err != nil {
		return err
	}

	if _, found := getConsensusState(clientStore, cdc, height); !found {
		return cosmossdkerrors.Wrap(clienttypes.ErrConsensusStateNotFound, "please ensure the proof was constructed against a height that exists on the client")
	}

	if len(proof) != 0 {
		return cosmossdkerrors.Wrapf(ErrInvalidProof, "expected the empty proof, actually got '%X'", proof)
	}

	return nil
}

// VerifyUpgradeAndUpdateState returns an error since Mock client does not support upgrades
func (cs ClientState) VerifyUpgradeAndUpdateState(
	_ sdk.Context, _ codec.BinaryCodec, _ storetypes.KVStore,
	_ exported.ClientState, _ exported.ConsensusState, _, _ []byte,
) error {
	return cosmossdkerrors.Wrap(clienttypes.ErrInvalidUpgradeClient, "cannot upgrade Mock client")
}

// verifyDelayPeriodPassed will ensure that at least delayTimePeriod amount of time and delayBlockPeriod number of blocks have passed
// since consensus state was submitted before allowing verification to continue.
func verifyDelayPeriodPassed(ctx sdk.Context, store storetypes.KVStore, proofHeight exported.Height, delayTimePeriod, delayBlockPeriod uint64) error {
	if delayTimePeriod != 0 {
		// check that executing chain's timestamp has passed consensusState's processed time + delay time period
		processedTime, ok := getProcessedTime(store, proofHeight)
		if !ok {
			return cosmossdkerrors.Wrapf(ErrProcessedTimeNotFound, "processed time not found for height: %s", proofHeight)
		}

		currentTimestamp := uint64(ctx.BlockTime().UnixNano())
		validTime := processedTime + delayTimePeriod

		// NOTE: delay time period is inclusive, so if currentTimestamp is validTime, then we return no error
		if currentTimestamp < validTime {
			return cosmossdkerrors.Wrapf(ErrDelayPeriodNotPassed, "cannot verify packet until time: %d, current time: %d",
				validTime, currentTimestamp)
		}

	}

	if delayBlockPeriod != 0 {
		// check that executing chain's height has passed consensusState's processed height + delay block period
		processedHeight, ok := getProcessedHeight(store, proofHeight)
		if !ok {
			return cosmossdkerrors.Wrapf(ErrProcessedHeightNotFound, "processed height not found for height: %s", proofHeight)
		}

		currentHeight := clienttypes.GetSelfHeight(ctx)
		validHeight := clienttypes.NewHeight(processedHeight.GetRevisionNumber(), processedHeight.GetRevisionHeight()+delayBlockPeriod)

		// NOTE: delay block period is inclusive, so if currentHeight is validHeight, then we return no error
		if currentHeight.LT(validHeight) {
			return cosmossdkerrors.Wrapf(ErrDelayPeriodNotPassed, "cannot verify packet until height: %s, current height: %s",
				validHeight, currentHeight)
		}
	}

	return nil
}
