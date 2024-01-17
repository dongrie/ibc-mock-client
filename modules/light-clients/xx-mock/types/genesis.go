package types

import (
	fmt "fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// ExportMetadata always panics. This function is used only for testing.
func (cs ClientState) ExportMetadata(_ storetypes.KVStore) []exported.GenesisMetadata {
	panic(fmt.Errorf("not implemented"))
}
