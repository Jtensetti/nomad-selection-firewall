package main

import (
	"fmt"
	"github.com/Jtensetti/nomad-selection-firewall/firewall"
)

func main() {
	cfg := firewall.NetworkConfig{CellsPerEpoch: 16, CellSize: 1200, PeerSlots: 4}
	idle := firewall.SelectionState{}
	active := firewall.SelectionState{PrivateQuery: "private local query", SelectedBasins: []uint64{0x42}}
	ok, err := firewall.SameObservableTrace(cfg, 1_000_000, idle, active)
	if err != nil {
		panic(err)
	}
	fmt.Printf("observable traces identical across 1,000,000 epochs: %v\n", ok)
}
