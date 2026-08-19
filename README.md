# nomad-selection-firewall

Executable **non-interference boundary** for Nomad research.

The Selection Firewall is not a packet filter. It is an architectural rule: local reading/search/reconstruction state must have no code path into externally observable traffic scheduling.

This repository makes that rule explicit:

- `NetworkConfig` is the complete input to emission planning.
- `SelectionState` is a separate type with no conversion to network scheduling.
- `Plan` is a pure function of public network state and epoch.
- million-epoch tests compare idle and active local worlds.

A real browser integration would additionally require OS/process sandboxing and review of telemetry, storage, retransmission and cache side channels. This repository proves only the application-level API boundary it implements.

```bash
go test ./...
go test -race ./...
go vet ./...
go run ./cmd/firewall-check
```
