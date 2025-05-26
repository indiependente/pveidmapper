# PVE ID Mapper

A command-line tool for managing UID/GID mappings in Proxmox VE LXC containers. This tool helps generate the necessary configuration for both the container and the host system to properly map user and group IDs.

## Installation

### Using Go Install

```bash
go install github.com/indiependente/pveidmapper@latest
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/indiependente/pveidmapper.git
cd pveidmapper
```

2. Build and install:
```bash
make install
```

## Usage

The tool provides a `generate` command to create ID mappings:

```bash
pveidmapper generate -i <mapping>
```

### Mapping Format

The mapping format is:
```
containeruid[:containergid][=hostuid[:hostgid]]
```

### Examples

1. Map a single UID:
```bash
pveidmapper generate -i 1000=1000
```

2. Map both UID and GID:
```bash
pveidmapper generate -i 1000:1000=1000:1000
```

3. Map multiple IDs:
```bash
pveidmapper generate -i 1000=1000 -i 1001=1001
```

### Output

The tool generates three sections of configuration:

1. Container configuration to add to `/etc/pve/lxc/<container_id>.conf`
2. User ID mappings to add to `/etc/subuid`
3. Group ID mappings to add to `/etc/subgid`

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using the Makefile)

### Common Tasks

- Build the binary:
```bash
make build
```

- Install the tool:
```bash
make install
```

- Run tests:
```bash
make test
```

- Update dependencies:
```bash
make update
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
