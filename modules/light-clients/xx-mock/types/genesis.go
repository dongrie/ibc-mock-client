package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// ExportMetadata always panics. This function is used only for testing.
func (cs ClientState) ExportMetadata(_ sdk.KVStore) []exported.GenesisMetadata {
	panic(fmt.Errorf("not implemented"))
}
