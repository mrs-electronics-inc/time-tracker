# Tasks: Rework the data model

**Input**: Design documents from `/specs/001-rework-the-data/`
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
- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize Go project with Cobra dependencies
- [ ] T003 [P] Configure linting and formatting tools

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test for list command in tests/contract/test_list.go
- [ ] T005 [P] Contract test for start command in tests/contract/test_start.go
- [ ] T006 [P] Contract test for stop command in tests/contract/test_stop.go
- [ ] T007 [P] Integration test for start and stop scenario in tests/integration/test_start_stop.go
- [ ] T008 [P] Integration test for auto-stop scenario in tests/integration/test_auto_stop.go
- [ ] T009 [P] Integration test for stop when no active in tests/integration/test_stop_no_active.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T010 [P] TimeEntry model in src/models/time_entry.go
- [ ] T011 TaskManager service in src/utils/task_manager.go
- [ ] T012 [P] list command implementation in src/cmd/list.go
- [ ] T013 [P] start command implementation in src/cmd/start.go
- [ ] T014 [P] stop command implementation in src/cmd/stop.go

## Phase 3.4: Integration
- [ ] T015 Connect TaskManager to JSON file storage
- [ ] T016 Error handling and logging

## Phase 3.5: Polish
- [ ] T017 [P] Unit tests for TimeEntry model in tests/unit/test_time_entry.go
- [ ] T018 Performance tests (<100ms startup)
- [ ] T019 [P] Update README.md with usage

## Dependencies
- Tests (T004-T009) before implementation (T010-T014)
- T010 blocks T011
- T011 blocks T012-T014, T015
- T015 blocks T016
- Implementation before polish (T017-T019)

## Parallel Example
```
# Launch T004-T009 together:
Task: "Contract test for list command in tests/contract/test_list.go"
Task: "Contract test for start command in tests/contract/test_start.go"
Task: "Contract test for stop command in tests/contract/test_stop.go"
Task: "Integration test for start and stop scenario in tests/integration/test_start_stop.go"
Task: "Integration test for auto-stop scenario in tests/integration/test_auto_stop.go"
Task: "Integration test for stop when no active in tests/integration/test_stop_no_active.go"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Each contract file → contract test task [P]
   - Each endpoint → implementation task
   
2. **From Data Model**:
   - Each entity → model creation task [P]
   - Relationships → service layer tasks
   
3. **From User Stories**:
   - Each story → integration test [P]
   - Quickstart scenarios → validation tasks

4. **Ordering**:
   - Setup → Tests → Models → Services → Endpoints → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task