<!-- Sync Impact Report
Version change: none → 1.0.0
List of modified principles: N/A (new constitution)
Added sections: All sections added
Removed sections: None
Templates requiring updates: None (templates align with new principles)
Follow-up TODOs: None
-->

# Time Tracker Constitution

## Core Principles

### I. CLI Simplicity

The Time Tracker tool MUST provide a simple, intuitive CLI interface. Commands SHOULD be straightforward with clear help messages and consistent argument patterns. Output MUST be human-readable by default, with optional JSON for scripting.

### II. Data Integrity

Task data MUST be accurately tracked and persisted without loss. The system MUST handle concurrent access safely and provide data recovery mechanisms.

### III. Test-First Development

All features MUST be developed following Test-Driven Development principles. Tests MUST be written before implementation, ensuring red-green-refactor cycle is strictly enforced. ALL Go test files MUST use the `_test.go` suffix (e.g., `example_test.go`) instead of a `test_` prefix, as the Go test runner only recognizes the `_test.go` suffix.

### IV. Performance Efficiency

The tool MUST be fast and lightweight, with minimal resource usage. Startup time SHOULD be under 100ms, and memory usage SHOULD remain low.

### V. Code Modularity

The codebase MUST be modular, with clear separation of concerns. Components SHOULD be independently testable and reusable.

## Technical Constraints

Primary programming language: Go. CLI framework: Cobra. Data storage: JSON file format. Dependencies: Minimal external libraries to maintain simplicity.

## Development Workflow

Use AI-driven workflow with spec-kit for specifications, opencode for implementation. Follow TDD, commit after each task, and ensure constitution compliance.

## Governance

Constitution supersedes all other practices. Amendments require consensus among maintainers. Version updates follow semantic versioning. Compliance MUST be verified in all PRs.

**Version**: 1.0.0 | **Ratified**: 2025-09-24 | **Last Amended**: 2025-09-24
