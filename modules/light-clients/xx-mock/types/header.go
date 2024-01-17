package types

import (
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ClientMessage = &Header{}

// ClientType return the client identifier of mock-client.
func (Header) ClientType() string {
	return Mock
}

// GetHeight returns the current sequence number as the height.
// Return clientexported.Height to satisfy interface
// Revision number is always 0 for a solo-machine
func (h Header) GetHeight() exported.Height {
	return h.Height
}

// ValidateBasic ensures that the sequence, signature and public key have all
// been initialized.
func (h Header) ValidateBasic() error {
	return nil
}
