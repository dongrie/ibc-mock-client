package types

import (
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/modules/core/exported"
)

var _ exported.Header = &Header{}

// ClientType defines that the Header is a Multisig.
func (Header) ClientType() string {
	return Mock
}

// GetHeight returns the current sequence number as the height.
// Return clientexported.Height to satisfy interface
// Revision number is always 0 for a solo-machine
func (h Header) GetHeight() exported.Height {
	return clienttypes.NewHeight(0, h.Height)
}

// ValidateBasic ensures that the sequence, signature and public key have all
// been initialized.
func (h Header) ValidateBasic() error {
	return nil
}
