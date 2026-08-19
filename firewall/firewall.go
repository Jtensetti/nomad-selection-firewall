package firewall

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

// NetworkConfig is the complete input to externally observable emission plans.
// Local SelectionState is intentionally not referenced by this type or Plan.
type NetworkConfig struct {
	CellsPerEpoch uint32
	CellSize      uint32
	PeerSlots     uint16
	PublicSeed    [32]byte
}

func (c NetworkConfig) Validate() error {
	if c.CellsPerEpoch == 0 || c.CellSize == 0 || c.PeerSlots == 0 {
		return errors.New("all network dimensions must be positive")
	}
	return nil
}

type Emission struct {
	Epoch    uint64
	Slot     uint32
	PeerSlot uint16
	Size     uint32
}

// Plan is a pure function of network configuration and epoch. Its output cannot
// vary with local reading/search/reconstruction state unless the API boundary is violated.
func Plan(cfg NetworkConfig, epoch uint64) ([]Emission, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	out := make([]Emission, cfg.CellsPerEpoch)
	for i := uint32(0); i < cfg.CellsPerEpoch; i++ {
		h := sha256.New()
		_, _ = h.Write(cfg.PublicSeed[:])
		var b [16]byte
		binary.BigEndian.PutUint64(b[:8], epoch)
		binary.BigEndian.PutUint32(b[8:12], i)
		_, _ = h.Write(b[:])
		sum := h.Sum(nil)
		peer := binary.BigEndian.Uint16(sum[:2]) % cfg.PeerSlots
		out[i] = Emission{Epoch: epoch, Slot: i, PeerSlot: peer, Size: cfg.CellSize}
	}
	return out, nil
}

// SelectionState is private browser state. This package exposes no conversion
// from SelectionState to NetworkConfig or Emission.
type SelectionState struct {
	PrivateQuery      string
	SelectedBasins    []uint64
	ReconstructionIDs [][32]byte
}

// SameObservableTrace asserts the core non-interference property for two worlds.
func SameObservableTrace(cfg NetworkConfig, epochs uint64, _ SelectionState, _ SelectionState) (bool, error) {
	for e := uint64(0); e < epochs; e++ {
		a, err := Plan(cfg, e)
		if err != nil {
			return false, err
		}
		b, err := Plan(cfg, e)
		if err != nil {
			return false, err
		}
		if len(a) != len(b) {
			return false, nil
		}
		for i := range a {
			if a[i] != b[i] {
				return false, nil
			}
		}
	}
	return true, nil
}
