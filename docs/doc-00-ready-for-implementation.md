# Ready for Implementation

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## What Has Been Decided

Forge will be built in Go using bubbletea v2 for TUI, cobra for CLI, chi + templ + HTMX for web, and Wails v2 (deferred) for desktop. The core layer (`core/tools/`) contains pure functions implementing 6 Tier-1 tools — base64, jwt, json, hash, url, uuid — with a Tool interface for metadata/discovery and a Registry for search and smart detection. The architecture enforces a hard boundary: `core/` never imports from `ui/` or `cmd/`. Each UI surface adapts core logic to its rendering model independently. The project structure comprises ~115 files across `cmd/`, `core/`, `ui/`, and `internal/` packages. All research documents, mockups, CLI reference, ADRs, and risk analysis are complete and ready to guide implementation.

## First Implementation Task

**Start with `core/tools/`.** Create `go.mod`, then `core/tools/tool.go` (the Tool interface and Result types), then implement all 6 Tier-1 tools as pure functions with comprehensive unit tests. This is the foundation everything else depends on, requires no UI libraries, and is the fastest way to become productive in Go before tackling bubbletea. Target: all 6 tools implemented and tested within the first week.

## Open Questions Before Coding Begins

1. **Bubbletea v2 import path:** Confirm whether to use `charm.land/bubbletea/v2` (vanity URL) or `github.com/charmbracelet/bubbletea/v2` (GitHub path). Both work, but the project should be consistent. Recommendation: use `charm.land/*` as the Charm team is standardizing on it.

2. **UUID v7 implementation:** Go's `google/uuid` library supports v4 natively. Verify that it supports v7 (time-ordered UUIDs) in its current version, or whether a different library is needed.

3. **YAML library choice:** The Go stdlib has no YAML package. The two main options are `gopkg.in/yaml.v3` and `github.com/goccy/go-yaml`. Decide before implementing the YAML tool.

4. **Templ version pinning:** Pin templ to a specific v0.3.x release in `go.mod` and document the pinned version, since templ is pre-v1.

---

## Document Index

| # | Document | File |
|---|----------|------|
| 1 | Tech Stack Validation Report | `doc-01-tech-stack-validation.md` |
| 2 | Competitive Analysis | `doc-02-competitive-analysis.md` |
| 3 | Project Structure Proposal | `doc-03-project-structure-proposal.md` |
| 4 | Tool Inventory | `doc-04-tool-inventory.md` |
| 5 | TUI Mockups | `doc-05-tui-mockups.md` |
| 6 | Web UI Mockups | `doc-06-web-ui-mockups.md` |
| 7 | CLI Command Reference | `doc-07-cli-command-reference.md` |
| 8 | Architecture Decision Records | `doc-08-architecture-decision-records.md` |
| 9 | Risk Register | `doc-09-risk-register.md` |
