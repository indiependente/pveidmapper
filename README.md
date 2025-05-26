![Go Report Card](https://goreportcard.com/badge/github.com/indiependente/pveidmapper)](https://goreportcard.com/report/github.com/indiependente/pveidmapper)
[![Workflow Status](https://github.com/indiependente/pveidmapper/workflows/lint-test/badge.svg)](https://github.com/indiependente/pveidmapper/actions)

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

The tool provides a `map` command to create ID mappings:

```bash
pveidmapper map -i <mapping>
```

### Mapping Format

The mapping format is:
```
containeruid[:containergid][=hostuid[:hostgid]]
```

### Examples

1. Map a single UID:
```bash
pveidmapper map -i 1000=1000
```

2. Map both UID and GID:
```bash
pveidmapper map -i 1000:1000
```
Output:
```
# Add to /etc/pve/lxc/<container_id>.conf:
lxc.idmap: u 0 100000 1000
lxc.idmap: u 1000 1000 1
lxc.idmap: u 1001 101001 64535
lxc.idmap: g 0 100000 1000
lxc.idmap: g 1000 1000 1
lxc.idmap: g 1001 101001 64535

# Add to /etc/subuid:
root:1000:1

# Add to /etc/subgid:
root:1000:1
```

3. Map multiple IDs:
```bash
pveidmapper map -i 1000=1000 -i 1001=1001
```

### Configuration Files

The tool generates configuration that needs to be added to three different files:

1. Container configuration file:
   - Path: `/etc/pve/lxc/<container_id>.conf`
   - Replace `<container_id>` with your actual container ID
   - Add the `lxc.idmap` lines to this file

2. User ID mappings:
   - Path: `/etc/subuid`
   - Add the `root:UID:1` lines to this file

3. Group ID mappings:
   - Path: `/etc/subgid`
   - Add the `root:GID:1` lines to this file

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using the Makefile)
- golangci-lint (for code quality checks)

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

- Run linter:
```bash
make lint
```

- Run linter with auto-fix:
```bash
make lint-fix
```

- Update dependencies:
```bash
make update
```

### Code Quality

This project uses golangci-lint for code quality checks. The configuration can be found in `.golangci.yml`. The linter checks for:

- Code style and formatting
- Common programming mistakes
- Potential bugs
- Security issues
- Performance issues
- Code complexity
- Duplicate code

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
