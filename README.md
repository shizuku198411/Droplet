# Droplet

Droplet is the **low-level container runtime** of the `Raind` container runtime stack.

It is responsible for executing OCI-style containers based on `config.json`
by setting up the filesystem, applying resource limits, and spawning the container
process on Linux.

Droplet plays a similar role to tools such as `runc`, but is being developed as
part of the Raind stack with a focus on extensibility and experimental features.

---

## 🧩 Project Architecture

`Raind` stack consists of multiple layers:

- **Condenser** – high-level runtime
- **Droplet** – low-level runtime (this repository)
- **Raind CLI** – developer-facing control interface

In the `Raind` stack, the higher-level runtime is responsible for managing
container metadata (image layers, command, networking, etc.).

`Droplet` receives these parameters as execution inputs, **builds the OCI
`config.json` from them**, and then:

- prepares the root filesystem (including overlay mounts)
- configures namespaces and cgroups
- starts the container process

---

## ✨ Current Features

- Parse OCI-style `config.json`
- Mount required filesystems and user-specified mounts
- Apply basic CPU / memory resource limits
- Execute process inside namespace-isolated environment

> ⚠ This project is currently experimental and under active development.

Planned / WIP:

- Improved networking support
- Additional isolation features
- Error handling improvements

---

## 🔧 Build & Run

### Requirements

- Go 1.21+
- Linux environment

### Build

```bash
git clone https://github.com/pyxgun/Droplet
cd droplet
go build ./cmd/droplet
```
