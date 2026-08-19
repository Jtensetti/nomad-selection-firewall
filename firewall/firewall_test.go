package firewall

import (
	"crypto/rand"
	"testing"
)

func TestSelectionCannotChangeTrace(t *testing.T) {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		t.Fatal(err)
	}
	cfg := NetworkConfig{CellsPerEpoch: 32, CellSize: 1200, PeerSlots: 8, PublicSeed: seed}
	idle := SelectionState{}
	active := SelectionState{PrivateQuery: "vapensystem i irans militär", SelectedBasins: []uint64{1, 2, 3, 4}}
	ok, err := SameObservableTrace(cfg, 10000, idle, active)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("selection influenced observable trace")
	}
}

func TestPlanDeterministic(t *testing.T) {
	cfg := NetworkConfig{CellsPerEpoch: 4, CellSize: 1200, PeerSlots: 3}
	a, _ := Plan(cfg, 42)
	b, _ := Plan(cfg, 42)
	for i := range a {
		if a[i] != b[i] {
			t.Fatal("plan is not deterministic")
		}
	}
}
