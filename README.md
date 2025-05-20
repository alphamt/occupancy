# GPU Device HTTP Service

A simple Go HTTP service to collect NVIDIA-SMI dumps and expose in-memory device information.

## Features

- **`POST /dump`**: Accepts JSON dumps from `nvidia-smi` output (host → devices) and stores in a thread-safe in-memory map.
- **`GET  /devices`**: Returns the full map of hosts and their device lists.
- **`GET  /summary`**: Returns per-host summaries:
  - `total_gpus`
  - `total_used_mem` (MiB)
  - `total_free_mem` (MiB)
  - `avg_free_mem_per_gpu` (MiB)

## Prerequisites

- Go 1.24.2+
- `nvidia-smi` available on agents sending dumps

## Installation

```bash
go run github.com/alphamt/occupancy/...@latest
```

