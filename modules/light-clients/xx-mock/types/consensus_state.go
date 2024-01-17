package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ConsensusState = &ConsensusState{}

// ClientType returns mock-client type.
func (ConsensusState) ClientType() string {
	return Mock
}

// GetTimestamp returns zero.
func (cs ConsensusState) GetTimestamp() uint64 {
	return cs.Timestamp
}

// GetRoot returns nil since mock-client do not have roots.
func (cs ConsensusState) GetRoot() exported.Root {
	return nil
}

// ValidateBasic defines basic validation for the mock-client consensus state.
func (cs ConsensusState) ValidateBasic() error {
	if cs.Timestamp == 0 {
		return cosmossdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "timestamp cannot be 0")
	}
	return nil
}
