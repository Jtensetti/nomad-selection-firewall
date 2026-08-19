package firewall

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

// NetworkConfig contains every input that the emission planner is allowed to
// use. Private reader state intentionally has no representation in this package.
type NetworkConfig struct {
	CellsPerEpoch uint32
	CellSize      uint32
	PeerSlots     uint16
	PublicSeed    [32]byte
}

func (c NetworkConfig) Validate() error {
	if c.CellsPerEpoch == 0 {
		return errors.New("cells per epoch must be positive")
	}
	if c.CellSize == 0 {
		return errors.New("cell size must be positive")
	}
	if c.PeerSlots == 0 {
		return errors.New("peer slots must be positive")
	}
	return nil
}

type Emission struct {
	Epoch    uint64
	Slot     uint32
	PeerSlot uint16
	Size     uint32
}

// Plan is a pure function of public network configuration and epoch.
func Plan(cfg NetworkConfig, epoch uint64) ([]Emission, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	out := make([]Emission, cfg.CellsPerEpoch)
	for i := uint32(0); i < cfg.CellsPerEpoch; i++ {
		h := sha256.New()
		_, _ = h.Write([]byte("nomad-selection-firewall-plan-v1"))
		_, _ = h.Write(cfg.PublicSeed[:])
		var b [12]byte
		binary.BigEndian.PutUint64(b[:8], epoch)
		binary.BigEndian.PutUint32(b[8:], i)
		_, _ = h.Write(b[:])
		sum := h.Sum(nil)
		peer := binary.BigEndian.Uint16(sum[:2]) % cfg.PeerSlots
		out[i] = Emission{Epoch: epoch, Slot: i, PeerSlot: peer, Size: cfg.CellSize}
	}
	return out, nil
}
