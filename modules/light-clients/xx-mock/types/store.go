package types

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var (
	// keyProcessedTime is appended to consensus state key to store the processed time
	keyProcessedTime = []byte("/processedTime")
	// keyProcessedHeight is appended to consensus state key to store the processed height
	keyProcessedHeight = []byte("/processedHeight")
)

// setClientState stores the client state
func setClientState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, clientState *ClientState) {
	key := host.ClientStateKey()
	val := clienttypes.MustMarshalClientState(cdc, clientState)
	clientStore.Set(key, val)
}

// setConsensusState stores the consensus state at the given height.
func setConsensusState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, consensusState *ConsensusState, height exported.Height) {
	key := host.ConsensusStateKey(height)
	val := clienttypes.MustMarshalConsensusState(cdc, consensusState)
	clientStore.Set(key, val)
}

// getConsensusState retrieves the consensus state from the client prefixed store.
// If the ConsensusState does not exist in state for the provided height a nil value and false boolean flag is returned
func getConsensusState(store storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height) (*ConsensusState, bool) {
	bz := store.Get(host.ConsensusStateKey(height))
	if len(bz) == 0 {
		return nil, false
	}

	consensusStateI := clienttypes.MustUnmarshalConsensusState(cdc, bz)
	return consensusStateI.(*ConsensusState), true
}

// processedTimeKey returns the key under which the processed time will be stored in the client store.
func processedTimeKey(height exported.Height) []byte {
	return append(host.ConsensusStateKey(height), keyProcessedTime...)
}

// setProcessedTime stores the time at which a header was processed and the corresponding consensus state was created.
// This is useful when validating whether a packet has reached the time specified delay period in the tendermint client's
// verification functions
func setProcessedTime(clientStore storetypes.KVStore, height exported.Height, timeNs uint64) {
	key := processedTimeKey(height)
	val := sdk.Uint64ToBigEndian(timeNs)
	clientStore.Set(key, val)
}

// getProcessedTime gets the time (in nanoseconds) at which this chain received and processed a tendermint header.
// This is used to validate that a received packet has passed the time delay period.
func getProcessedTime(clientStore storetypes.KVStore, height exported.Height) (uint64, bool) {
	key := processedTimeKey(height)
	bz := clientStore.Get(key)
	if len(bz) == 0 {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

// processedHeightKey returns the key under which the processed height will be stored in the client store.
func processedHeightKey(height exported.Height) []byte {
	return append(host.ConsensusStateKey(height), keyProcessedHeight...)
}

// setProcessedHeight stores the height at which a header was processed and the corresponding consensus state was created.
// This is useful when validating whether a packet has reached the specified block delay period in the tendermint client's
// verification functions
func setProcessedHeight(clientStore storetypes.KVStore, consHeight, processedHeight exported.Height) {
	key := processedHeightKey(consHeight)
	val := []byte(processedHeight.String())
	clientStore.Set(key, val)
}

// getProcessedHeight gets the height at which this chain received and processed a tendermint header.
// This is used to validate that a received packet has passed the block delay period.
func getProcessedHeight(clientStore storetypes.KVStore, height exported.Height) (exported.Height, bool) {
	key := processedHeightKey(height)
	bz := clientStore.Get(key)
	if len(bz) == 0 {
		return nil, false
	}
	processedHeight, err := clienttypes.ParseHeight(string(bz))
	if err != nil {
		return nil, false
	}
	return processedHeight, true
}

// setConsensusMetadata sets context time as processed time and set context height as processed height.
// This is same logic as tendermint LC.
func setConsensusMetadata(ctx sdk.Context, clientStore storetypes.KVStore, height exported.Height) {
	setConsensusMetadataWithValues(clientStore, height, clienttypes.GetSelfHeight(ctx), uint64(ctx.BlockTime().UnixNano()))
}

// setConsensusMetadataWithValues sets the consensus metadata with the provided values
func setConsensusMetadataWithValues(
	clientStore storetypes.KVStore, height,
	processedHeight exported.Height,
	processedTime uint64,
) {
	setProcessedTime(clientStore, height, processedTime)
	setProcessedHeight(clientStore, height, processedHeight)
}
