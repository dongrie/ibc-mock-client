package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	ics23 "github.com/confio/ics23/go"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/modules/core/exported"
)

const (
	Mock string = "mock-client"
)

var _ exported.ClientState = (*ClientState)(nil)

// NewClientState creates a new ClientState instance.
func NewClientState(latestHeight clienttypes.Height, allowUpdateAfterProposal bool) *ClientState {
	return &ClientState{
		LatestHeight: latestHeight.RevisionHeight,
	}
}

// ClientType returns a type of the client.
func (cs ClientState) ClientType() string {
	return Mock
}

// GetLatestHeight returns the latest height.
// Return exported.Height to satisfy ClientState interface
func (cs ClientState) GetLatestHeight() exported.Height {
	return clienttypes.NewHeight(0, cs.LatestHeight)
}

// Status returns the status of the mock client.
// It always returns active.
func (cs ClientState) Status(_ sdk.Context, _ sdk.KVStore, _ codec.BinaryCodec) exported.Status {
	return exported.Active
}

// GetProofSpecs returns nil proof specs since client state verification uses signatures.
func (cs ClientState) GetProofSpecs() []*ics23.ProofSpec {
	return nil
}

// Validate performs basic validation of the client state fields.
func (cs ClientState) Validate() error {
	return nil
}

// ZeroCustomFields returns Mock client state with client-specific fields FrozenSequence,
// and AllowUpdateAfterProposal zeroed out
func (cs ClientState) ZeroCustomFields() exported.ClientState {
	return &ClientState{}
}

// Initialize will check that initial consensus state is equal to the latest consensus state of the initial client.
func (cs ClientState) Initialize(_ sdk.Context, _ codec.BinaryCodec, _ sdk.KVStore, consState exported.ConsensusState) error {
	return nil
}

// ExportMetadata is a no-op since Mock does not store any metadata in client store
func (cs ClientState) ExportMetadata(_ sdk.KVStore) []exported.GenesisMetadata {
	return nil
}

// VerifyUpgradeAndUpdateState returns an error since Mock client does not support upgrades
func (cs ClientState) VerifyUpgradeAndUpdateState(
	_ sdk.Context, _ codec.BinaryCodec, _ sdk.KVStore,
	_ exported.ClientState, _ exported.ConsensusState, _, _ []byte,
) (exported.ClientState, exported.ConsensusState, error) {
	return nil, nil, sdkerrors.Wrap(clienttypes.ErrInvalidUpgradeClient, "cannot upgrade Mock client")
}

// VerifyClientState verifies a proof of the client state of the running chain
// stored on the Mock.
func (cs ClientState) VerifyClientState(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	prefix exported.Prefix,
	counterpartyClientIdentifier string,
	proof []byte,
	clientState exported.ClientState,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return err
	}

	bz, err := cdc.Marshal(anyClientState)
	if err != nil {
		return err
	}

	h := sha256.Sum256(bz)
	if !bytes.Equal(proof, h[:]) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, h)
	}
	return nil
}

// VerifyClientConsensusState verifies a proof of the consensus state of the
// running chain stored on the Mock.
func (cs ClientState) VerifyClientConsensusState(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	counterpartyClientIdentifier string,
	consensusHeight exported.Height,
	prefix exported.Prefix,
	proof []byte,
	consensusState exported.ConsensusState,
) error {
	// NOTE In cosmos/ibc-go, it cannot give a consensus state of an external prover(e.g. mock-client) to the client, so we skip this verification for now.
	return nil
}

// VerifyConnectionState verifies a proof of the connection state of the
// specified connection end stored on the target machine.
func (cs ClientState) VerifyConnectionState(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	connectionID string,
	connectionEnd exported.ConnectionI,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	connection, ok := connectionEnd.(connectiontypes.ConnectionEnd)
	if !ok {
		return sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnection,
			"expected type %T, got %T", connectiontypes.ConnectionEnd{}, connectionEnd,
		)
	}

	bz, err := cdc.Marshal(&connection)
	if err != nil {
		return err
	}

	h := sha256.Sum256(bz)
	if !bytes.Equal(proof, h[:]) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, h)
	}
	return nil
}

// VerifyChannelState verifies a proof of the channel state of the specified
// channel end, under the specified port, stored on the target machine.
func (cs ClientState) VerifyChannelState(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	channelEnd exported.ChannelI,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	channel, ok := channelEnd.(channeltypes.Channel)
	if !ok {
		return sdkerrors.Wrapf(
			channeltypes.ErrInvalidChannel,
			"expected channel type %T, got %T", channeltypes.Channel{}, channelEnd)
	}

	bz, err := cdc.Marshal(&channel)
	if err != nil {
		return err
	}

	h := sha256.Sum256(bz)
	if !bytes.Equal(proof, h[:]) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, h)
	}
	return nil
}

// VerifyPacketCommitment verifies a proof of an outgoing packet commitment at
// the specified port, specified channel, and specified sequence.
func (cs ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	_ uint64,
	_ uint64,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	packetSequence uint64,
	commitmentBytes []byte,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	if !bytes.Equal(proof, commitmentBytes) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, commitmentBytes)
	}
	return nil
}

// VerifyPacketAcknowledgement verifies a proof of an incoming packet
// acknowledgement at the specified port, specified channel, and specified sequence.
func (cs ClientState) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	_ uint64,
	_ uint64,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	packetSequence uint64,
	acknowledgement []byte,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	commitmentBytes := channeltypes.CommitAcknowledgement(acknowledgement)
	if !bytes.Equal(proof, commitmentBytes) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, commitmentBytes)
	}
	return nil
}

// VerifyPacketReceiptAbsence verifies a proof of the absence of an
// incoming packet receipt at the specified port, specified channel, and
// specified sequence.
func (cs ClientState) VerifyPacketReceiptAbsence(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	_ uint64,
	_ uint64,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	packetSequence uint64,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	commitmentBytes := sha256.Sum256([]byte(fmt.Sprintf("%v/%v/%v", portID, channelID, packetSequence)))
	if !bytes.Equal(proof, commitmentBytes[:]) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, commitmentBytes)
	}
	return nil
}

// VerifyNextSequenceRecv verifies a proof of the next sequence number to be
// received of the specified channel at the specified port.
func (cs ClientState) VerifyNextSequenceRecv(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	_ uint64,
	_ uint64,
	prefix exported.Prefix,
	proof []byte,
	portID,
	channelID string,
	nextSequenceRecv uint64,
) error {
	_, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	commitmentBytes := sha256.Sum256([]byte(fmt.Sprintf("%v/%v/%v", portID, channelID, nextSequenceRecv)))
	if !bytes.Equal(proof, commitmentBytes[:]) {
		return fmt.Errorf("expected the proof '%X', actually got '%X'", proof, commitmentBytes)
	}
	return nil
}

// produceVerificationArgs perfoms the basic checks on the arguments that are
// shared between the verification functions and returns the public key of the
// consensus state, the unmarshalled proof representing the signature and timestamp
// along with the solo-machine sequence encoded in the proofHeight.
func produceVerificationArgs(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	cs ClientState,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
) (*ConsensusState, error) {
	if revision := height.GetRevisionNumber(); revision != 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidHeight, "revision must be 0 for Mock, got revision-number: %d", revision)
	}

	if prefix == nil {
		return nil, sdkerrors.Wrap(commitmenttypes.ErrInvalidPrefix, "prefix cannot be empty")
	}

	_, ok := prefix.(*commitmenttypes.MerklePrefix)
	if !ok {
		return nil, sdkerrors.Wrapf(commitmenttypes.ErrInvalidPrefix, "invalid prefix type %T, expected MerklePrefix", prefix)
	}

	if proof == nil {
		return nil, sdkerrors.Wrap(ErrInvalidProof, "proof cannot be empty")
	}

	cons, err := getConsensusState(store, cdc, height)
	if err != nil {
		return nil, err
	}
	return cons, nil
}
