# SchemaGit

> A Git-inspired schema management and versioning system for structured data.

SchemaGit is an open-source project designed to bring **Git-like versioning concepts to schemas and structured data**.

Instead of treating schema changes as isolated files or undocumented database migrations, SchemaGit provides a structured way to track, inspect, compare, and manage schema evolution over time.

The project is written in **Go** and is designed around a lightweight local architecture with SQLite persistence and a dedicated user interface.

---

## ✨ Why SchemaGit?

Schemas evolve.

Fields are added, removed, renamed, or modified. Relationships change. Constraints evolve. In real-world projects, these changes can become difficult to understand and reproduce.

Traditional version control works extremely well for source code, but schema evolution often requires a different workflow.

SchemaGit aims to provide a workflow where schema changes can be treated more like source-code changes:

```text
Schema
  │
  ├── Snapshot
  │
  ├── Change
  │
  ├── Diff
  │
  ├── History
  │
  └── Version
```

The goal is to make schema evolution:

- Observable
- Reproducible
- Versioned
- Comparable
- Auditable
- Easier to understand

---

## 🚀 Core Idea

The central idea behind SchemaGit is simple:

```text
Git manages changes to code.
SchemaGit manages changes to schemas.
```

A schema can be captured as a versioned state, allowing different states to be inspected and compared.

For example:

```text
v1
 └── users
      ├── id
      ├── username
      └── email

        │

        ▼

v2
 └── users
      ├── id
      ├── username
      ├── email
      └── created_at
```

Instead of manually figuring out what changed, SchemaGit can represent the evolution between these states.

---

## 🧩 Features

### Schema Versioning

Track different versions of a schema and preserve its history.

### Schema Diffing

Compare schema states and identify structural changes.

Conceptually:

```diff
 users
   id          INTEGER
   username    TEXT
   email       TEXT
+  created_at  DATETIME
```

### History

Keep a historical view of schema evolution instead of treating the current schema as the only source of truth.

### Local Persistence

SchemaGit uses SQLite for lightweight local persistence.

### CLI

The project provides a command-line interface through:

```text
cmd/schemagit
```

### UI

SchemaGit also contains a dedicated UI layer:

```text
ui/
```

This allows the project to evolve beyond a purely command-line workflow.

---

## 🏗️ Architecture

The repository is organized into several major components:

```text
SchemaGit/
│
├── cmd/
│   └── schemagit/
│       └── CLI entrypoint
│
├── internal/
│   └── Core application logic
│
├── ui/
│   └── User interface
│
├── docker/
│   └── Container-related resources
│
├── go.mod
└── README.md
```

The project follows a Go-oriented structure that keeps the executable entrypoint separate from internal implementation details.

---

## 🛠️ Requirements

- Go 1.22+
- Git
- SQLite support through `go-sqlite3`

Check your Go installation:

```bash
go version
```

---

## 📦 Installation

Clone the repository:

```bash
git clone https://github.com/sepantartd/SchemaGit.git
cd SchemaGit
```

Download dependencies:

```bash
go mod download
```

Build the project:

```bash
go build ./...
```

---

## ▶️ Running

Run the CLI directly:

```bash
go run ./cmd/schemagit
```

Or build it first:

```bash
go build -o schemagit ./cmd/schemagit
```

Then:

```bash
./schemagit
```

> Available commands and flags may evolve as SchemaGit develops.

---

## 🐳 Docker

SchemaGit also contains Docker-related resources under:

```text
docker/
```

If you are using Docker-based development or deployment, inspect the files in that directory for the currently supported configuration.

---

## 🧪 Development

Clone the repository:

```bash
git clone https://github.com/sepantartd/SchemaGit.git
cd SchemaGit
```

Run the test suite:

```bash
go test ./...
```

Run static checks:

```bash
go vet ./...
```

Build everything:

```bash
go build ./...
```

---

## 🔍 Project Philosophy

SchemaGit is built around a few principles:

### 1. History should be first-class

A schema is not just its current state.

Its history explains how it got there.

### 2. Changes should be inspectable

A developer should be able to understand what changed without manually comparing unrelated schema definitions.

### 3. Local-first should remain practical

Schema management should not require a large external infrastructure stack just to inspect or version a schema.

### 4. Automation should be possible

A versioned schema model creates opportunities for CI/CD, validation, migration generation, compatibility checks, and other automated workflows.

---

## 🔮 Potential Future Directions

SchemaGit is an evolving project.

Possible future capabilities include:

- Schema branching
- Schema merging
- Migration generation
- Rollback support
- Compatibility analysis
- Breaking-change detection
- Schema validation
- CI/CD integration
- Remote repositories
- Schema registries
- Database introspection
- Multiple database backends
- Machine-readable schema formats
- API support
- Team collaboration
- Visual schema history
- Import/export tooling

These ideas are part of the broader direction of the project and may change as development continues.

---

## 🤝 Contributing

Contributions are welcome.

Before opening a pull request:

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Add or update tests where appropriate.
5. Run:

```bash
go test ./...
go vet ./...
go build ./...
```

6. Open a pull request with a clear description of the change.

For larger changes, opening an issue first can help discuss the design before implementation.

---

## 📄 License

See the repository for the current license information.

---

## 🌐 Repository

GitHub:

https://github.com/sepantartd/SchemaGit

---

## ⭐ Project Status

SchemaGit is an actively evolving open-source project.

The architecture and feature set may change as the project grows.

If you find the project useful, consider giving it a star and contributing ideas, issues, or pull requests.

---

## 💡 Name

**SchemaGit** combines:

- **Schema** — the structure and definition of data
- **Git** — version control and change history

The name represents the project's main concept:

> **Version control for schemas.**
