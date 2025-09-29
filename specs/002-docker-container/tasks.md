# Tasks: 002-docker-container

**Input**: Design documents from `/specs/002-docker-container/`
**Prerequisites**: plan.md (required), research.md, data-model.md, quickstart.md

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
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [x] T001 Create Dockerfile in repository root
- [x] T002 Create .github/workflows/ci.yml directory and file

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [x] T003 [P] Integration test: Docker build succeeds in tests/integration/docker_build_test.go
- [x] T004 [P] Integration test: Container runs and executes help command in tests/integration/docker_run_test.go
- [x] T005 [P] Integration test: Container can start/stop time tracking in tests/integration/docker_commands_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [x] T006 [P] Implement Dockerfile with multi-stage Go build in Dockerfile
- [x] T007 [P] Implement CI workflow to build and push Docker image on main branch in .github/workflows/ci.yml
- [x] T008 [P] Create docker-compose.yml for easy testing in repository root

## Phase 3.4: Integration
- [x] T009 Test CI workflow by pushing to main branch (workflow implemented and validated locally)

## Phase 3.5: Polish
- [x] T010 [P] Update README.md with Docker build and run instructions
- [x] T011 Update AGENTS.md with instructions on using the application through docker compose ONLY
- [x] T012 Run quickstart.md scenarios manually

## Dependencies
- Tests (T003-T005) before implementation (T006-T008)
- Implementation before integration (T009)
- Everything before polish (T010-T012)

## Parallel Example
```
# Launch T003-T005 together:
Task: "Integration test: Docker build succeeds in tests/integration/docker_build_test.go"
Task: "Integration test: Container runs and executes help command in tests/integration/docker_run_test.go"
Task: "Integration test: Container can start/stop time tracking in tests/integration/docker_commands_test.go"

# Launch T006-T008 together:
Task: "Implement Dockerfile with multi-stage Go build in Dockerfile"
Task: "Implement CI workflow to build and push Docker image on main branch in .github/workflows/ci.yml"
Task: "Create docker-compose.yml for easy testing in repository root"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - No contracts directory - no contract test tasks
   
2. **From Data Model**:
   - No new entities - no model creation tasks
   
3. **From User Stories**:
   - Quickstart scenarios → integration tests [P]
   - Building Docker image → T003
   - Running container → T004
   - Testing in container → T005

4. **Ordering**:
   - Setup → Tests → Core → Integration → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding tests (N/A - no contracts)
- [ ] All entities have model tasks (N/A - no new entities)
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task
- [ ] Docker compose included for agent testing