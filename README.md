# 🚀 sdk-go: Powerful Go Libraries Toolkit

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/downsized-devs/sdk-go)
![Build Status](https://img.shields.io/github/actions/workflow/status/downsized-devs/sdk-go/go.yml)
![Coverage](https://img.shields.io/codecov/c/github/downsized-devs/sdk-go)
![Version](https://img.shields.io/github/v/release/downsized-devs/sdk-go)

## 📖 Overview
`sdk-go` is a monorepo of Go libraries maintained by Downsized Devs.  Each
top-level directory exposes a focused package that can be imported on its own or
pulled in as part of the complete toolkit.  The project aims to streamline Go
development with well-tested utilities ranging from logging and scheduled jobs
to data storage helpers.

## ✨ Features
- 🔧 Modular library architecture
- 🚄 High-performance implementations
- 🛡️ Robust error handling
- 📦 Easy integration
- 🧪 Thoroughly tested components

## 🛠️ Installation

### Quick Start
```bash
go get -u github.com/downsized-devs/sdk-go
```

### Individual Library Installation
```bash
# Install specific libraries as needed
go get -u github.com/downsized-devs/sdk-go/<package-name>
```

## 💻 Usage Example

```go
import "github.com/downsized-devs/sdk-go/<package-name>"
```

## 📂 Repository Structure
Each top-level directory houses a standalone Go package. A few notable examples
include:

- `appcontext` – request-scoped context helpers
- `logger` – structured logging based on Zerolog
- `scheduler` – wrappers around `gocron` for background jobs
- `redis` – thin client with distributed locking support
- `translator` – i18n solution using `universal-translator`

Packages can be imported individually or as part of the entire toolkit.

## 🛠 Code Generator
The `generator/` folder contains a scaffolding tool for creating boilerplate in
other projects. Run it with:

```bash
go run ./generator --name <EntityName> --path <output-path> --api
```

## 🔧 Testing & Tooling
Use the provided `Makefile` for common tasks:

- `make build` – compile all packages
- `make run-tests` – execute the unit test suite
- `make mock-all` – generate GoMock stubs

## 🌱 Explore Further
- Dive into packages such as `auth`, `storage`, or `messaging` to see available
  APIs.
- Review `errors/` and `codes/` to understand custom error handling.
- Check `instrument/` for Prometheus instrumentation examples.

## 📜 License
Distributed under the MIT License. See `LICENSE` for more information.

## 🌟 Support
If you encounter any problems or have suggestions, please [open an issue](https://github.com/downsized-devs/sdk-go/issues).


## 🏆 Quality Metrics
![GitHub Issues](https://img.shields.io/github/issues/downsized-devs/sdk-go)
![GitHub Pull Requests](https://img.shields.io/github/issues-pr/downsized-devs/sdk-go)
![GitHub License](https://img.shields.io/github/license/downsized-devs/sdk-go)
![Code Quality](https://img.shields.io/codefactor/grade/github/downsized-devs/sdk-go)

## 📊 Repository Stats
![GitHub Contributors](https://img.shields.io/github/contributors/downsized-devs/sdk-go)
![GitHub Last Commit](https://img.shields.io/github/last-commit/downsized-devs/sdk-go)
![Repo Size](https://img.shields.io/github/repo-size/downsized-devs/sdk-go)
![Language](https://img.shields.io/github/languages/top/downsized-devs/sdk-go)

## 🛡️ Code Health
![Go Report Card](https://goreportcard.com/badge/github.com/downsized-devs/sdk-go)

## 🌐 Community & Engagement
![GitHub Stars](https://img.shields.io/github/stars/downsized-devs/sdk-go?style=social)
![GitHub Forks](https://img.shields.io/github/forks/downsized-devs/sdk-go?style=social)

---

**Happy Coding! 👨‍💻👩‍💻**
