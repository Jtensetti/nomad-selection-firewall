# nomad-selection-firewall

A narrow network-planning package for one Nomad architectural constraint: private reader state must not be an input to the public emission plan.

The useful property here is **API shape**, not a statistical result. The package contains only public scheduling inputs (`NetworkConfig`) and planned emissions (`Emission`). It has no query, selected-basin, reading-state or reconstruction-state type.

## Implemented

- validation of public scheduling dimensions,
- deterministic per-epoch emission plans,
- fixed cell count/size in the plan,
- deterministic peer-slot selection from public epoch state.

## What this does not enforce

This package is not a firewall process, sandbox or network stack. A browser, transport, cache, OS scheduler or caller can still leak private state through other channels. The package only makes one dependency harder to introduce accidentally: there is no private-selection input to `Plan`.

Cross-process and packet-level non-interference belongs in browser integration and `nomad-testnet`.

```bash
go test -race ./...
go vet ./...
```
