// Copyright 2021 The utg Authors
// This file is part of the utg library.
//
// The utg library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The utg library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the utg library. If not, see <http://www.gnu.org/licenses/>.

package params

import "github.com/UltronGlow/UltronGlow-Origin/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main utg network.
var MainnetBootnodes = []string{
	// utg Foundation Go Bootnodes
	"enode://617ff1e65455373593e08f08bb44e1c3e5dc3736e824870e5e2ab42b501610570b7bbaeb9621af0df04919cd0d4782924465ca39c97a811be0f3600dfe1d0fff@65.19.174.250:30313",
	"enode://9591316492a851ec7631f46f56633172804d7d5f81a56b98b3e378cb48fc01844b690d923fd55a6b5d5bbed0e96fc65da967f1752d1d5f4b30b76aa118af4a1f@65.19.174.250:30314",
	"enode://42255d9d3bd5b4eb9c19ae0f6a6f03dfe66eec89c69fa7936fcb340ad9b2eca7ad51b2ceee67373e7d82538e0f506df6d9deccc4edfdf7d7899896af42b237ca@65.19.174.250:30315",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Testnet test network.
var TestnetBootnodes = []string{
	"enode://2ffed1bb6b475259c1448dc93b639569886999e51ade144451877a706d2a9b71eff8eb067d289fde48ba4807370034d851553746fac8816af27f5a922703e2e4@127.0.0.1:30311",
}

var V5Bootnodes = []string{
}

const dnsPrefix = "enrtree://AKA3AM6LPBYEUDMVNU3BSVQJ5AD45Y7YPOHJLEF6W26QOE4VTUDPE@"

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	var net string
	switch genesis {
	case MainnetGenesisHash:
		net = "mainnet"
	case TestnetGenesisHash:
		net = "testnet"
	default:
		return ""
	}
	return dnsPrefix + protocol + "." + net + ".ethdisco.net"
}
