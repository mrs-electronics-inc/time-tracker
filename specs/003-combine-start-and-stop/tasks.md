# Tasks: combine-start-and-stop

**Input**: Design documents from `/specs/003-combine-start-and-stop/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [X] T001 Ensure Go project is properly initialized with Cobra dependencies

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [X] T012 Create new combined command in src/cmd/track.go
- [X] T013 Update src/cmd/root.go to register track command with aliases
- [X] T014 Remove old src/cmd/start.go and src/cmd/stop.go files

## Phase 3.4: Integration
- [X] T015 Update task manager to work with new command structure

## Phase 3.5: Polish
- [X] T016 [P] Update existing unit tests in tests/unit/ to use new command
- [X] T017 [P] Update existing integration tests in tests/integration/ to use new command
- [X] T018 Run quickstart.md scenarios to validate functionality
- [X] T019 Update README.md with new command usage

## Dependencies
- T012 blocks T013, T014
- T013 blocks T015
- Implementation before polish (T016-T019)

## Parallel Example
```
# No parallel tests in this phase
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Each contract invocation → contract test task [P]
   - Each command → implementation task

2. **From Data Model**:
   - No new entities → no model tasks

3. **From User Stories**:
   - Each acceptance scenario → integration test [P]
   - Quickstart scenarios → validation tasks

4. **Ordering**:
   - Setup → Tests → Core → Integration → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task