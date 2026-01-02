# Architecture: goprod

## 🏗 Philosophy

We use **Hexagonal Architecture** to keep business logic pure and decoupled from technical details like databases or web frameworks.

## 📁 Directory Structure

```text
goprod/
├── cmd/                # App entry points (main.go)
├── internal/
│   ├── domain/         # Entities & Interfaces (No external imports)
│   ├── service/        # Business logic & Orchestration (Uses domain)
│   └── infra/          # Adapters (Postgres, HTTP, External APIs)
├── pkg/                # Shared helper utilities
└── Architecture.md     # This document

```

## 📜 Core Rules

1. **Dependency Direction:** Always point inwards. `Infra -> Service -> Domain`.
2. **Domain Purity:** The `internal/domain` folder must not import any other internal packages.
3. **Context:** Every function in Service and Infra must take `context.Context` as its first argument.
4. **Error Handling:** Wrap errors with context: `fmt.Errorf("service.method: %w", err)`.
5. **ID Generation:** Generate UUIDs in the **Service** layer before saving to the database.
