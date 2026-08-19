package firewall

import (
	"crypto/sha256"
	"reflect"
	"testing"
)

func testConfig() NetworkConfig {
	return NetworkConfig{
		CellsPerEpoch: 32,
		CellSize:      1200,
		PeerSlots:     8,
		PublicSeed:    sha256.Sum256([]byte("selection-firewall-test-seed")),
	}
}

func TestPlanShape(t *testing.T) {
	cfg := testConfig()
	plan, err := Plan(cfg, 42)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != int(cfg.CellsPerEpoch) {
		t.Fatalf("got %d emissions, want %d", len(plan), cfg.CellsPerEpoch)
	}
	for i, emission := range plan {
		if emission.Epoch != 42 || emission.Slot != uint32(i) || emission.Size != cfg.CellSize {
			t.Fatalf("invalid emission %d: %#v", i, emission)
		}
		if emission.PeerSlot >= cfg.PeerSlots {
			t.Fatalf("peer slot %d out of range", emission.PeerSlot)
		}
	}
}

func TestPlanDeterministicForSamePublicInputs(t *testing.T) {
	cfg := testConfig()
	a, err := Plan(cfg, 42)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Plan(cfg, 42)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("same public inputs produced different plans")
	}
}

func TestPeerSequenceDependsOnEpoch(t *testing.T) {
	cfg := testConfig()
	a, err := Plan(cfg, 1)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Plan(cfg, 2)
	if err != nil {
		t.Fatal(err)
	}
	same := true
	for i := range a {
		if a[i].PeerSlot != b[i].PeerSlot {
			same = false
			break
		}
	}
	if same {
		t.Fatal("peer-slot sequence did not change across epochs")
	}
}

func TestInvalidConfig(t *testing.T) {
	cases := []NetworkConfig{
		{},
		{CellsPerEpoch: 1, CellSize: 1200},
		{CellsPerEpoch: 1, PeerSlots: 1},
		{CellSize: 1200, PeerSlots: 1},
	}
	for _, cfg := range cases {
		if _, err := Plan(cfg, 0); err == nil {
			t.Fatalf("expected validation error for %#v", cfg)
		}
	}
}
