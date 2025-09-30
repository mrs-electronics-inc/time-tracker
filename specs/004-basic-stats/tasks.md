# Tasks: Basic stats

**Input**: Design documents from `/specs/004-basic-stats/`
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
- Paths assume single project structure

## Phase 3.1: Setup
- [ ] T001 Configure Go linting and formatting tools (gofmt, go vet)
- [ ] T002 Verify existing project structure matches plan.md

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T003 [P] Contract test for stats command output format in tests/contract/stats_contract_test.go
- [ ] T004 [P] Integration test for daily totals display in tests/integration/stats_daily_test.go
- [ ] T005 [P] Integration test for weekly totals display in tests/integration/stats_weekly_test.go
- [ ] T006 [P] Integration test for project totals display in tests/integration/stats_projects_test.go
- [ ] T007 [P] Integration test for no data scenario in tests/integration/stats_no_data_test.go

## Phase 3.3: Core Implementation

- [x] T008 Create stats command structure in src/cmd/stats.go
- [x] T009 Implement time calculation utilities in src/utils/stats_calculations.go
- [x] T010 Add flag parsing for --daily, --weekly, --projects in src/cmd/stats.go
- [x] T011 Implement data aggregation logic in src/utils/stats_calculations.go
- [x] T012 Add output formatting for table display in src/cmd/stats.go

## Phase 3.4: Integration

- [x] T013 Connect stats command to file storage in src/cmd/stats.go
- [x] T014 Handle time zone and date calculations in src/utils/stats_calculations.go

## Phase 3.5: Polish

- [x] T015 [P] Unit tests for calculation functions in tests/unit/stats_calculations_test.go
- [x] T016 [P] Update README.md with stats command documentation
- [x] T017 Run quickstart.md validation steps
- [x] T018 Performance verification (<100ms startup)

## Dependencies

- Tests (T003-T007) before implementation (T008-T014)
- T008 blocks T010, T012, T013
- T009 blocks T011, T014
- Implementation before polish (T015-T018)

## Parallel Example

```
# Launch T003-T007 together:
Task: "Contract test for stats command output format in tests/contract/stats_contract_test.go"
Task: "Integration test for daily totals display in tests/integration/stats_daily_test.go"
Task: "Integration test for weekly totals display in tests/integration/stats_weekly_test.go"
Task: "Integration test for project totals display in tests/integration/stats_projects_test.go"
Task: "Integration test for no data scenario in tests/integration/stats_no_data_test.go"
```

## Notes

- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task following TDD
- Avoid: vague tasks, same file conflicts

## Task Generation Rules

_Applied during main() execution_

1. **From Contracts**:
   - stats.json → contract test task [P]

2. **From Data Model**:
   - TimeEntry entity exists, no new model task
   - Project entity exists, no new model task

3. **From User Stories**:
   - Daily totals story → integration test [P]
   - Weekly totals story → integration test [P]
   - Project totals story → integration test [P]
   - No data scenario → integration test [P]

4. **Ordering**:
   - Setup → Tests → Core → Integration → Polish
   - Dependencies block parallel execution

## Validation Checklist

_GATE: Checked by main() before returning_

- [x] All contracts have corresponding tests
- [x] All entities have model tasks (existing)
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task</content>
      </xai:function_call">## Validation Checklist
      _GATE: Checked by main() before returning_

- [x] All contracts have corresponding tests
- [x] All entities have model tasks (existing)
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
