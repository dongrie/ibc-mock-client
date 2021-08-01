# IBC Mock Client

This client is intended to be used for testing purpose. Therefore, it is not generally available in a production, except in a fully trusted environment.

## Spec

NOTE: A full spec is WIP

The client verifies that the serialization of Connection, Channel and other commitments are compatible on several different IBC implementations.

Each verification function is given a sha256 hash of the value to be verified as proof. If the value and an expected value (e.g. ConnectionEnd in the counterparty) match, the verification succeeds.

## Implementations

- [Go](./modules/light-clients/xx-mock)
- [Solidity](https://github.com/hyperledger-labs/yui-ibc-solidity/blob/main/contracts/core/MockClient.sol)
