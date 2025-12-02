# TUI Implementation

## Description

The current CLI-only interface requires users to remember or look up project and task names when starting/stopping time entries. A Text User Interface (TUI) will improve usability by providing autocomplete suggestions, visual feedback, and the ability to duplicate recent entries. This makes time tracking more accessible and efficient for everyday use.

## Design Decisions

### Framework Choice: Bubble Tea
- **Decision**: Use charmbracelet/bubbletea and charmbracelet/bubbles for TUI framework
- **Rationale**: Referenced in issue #2, purpose-built for Go TUIs with good component library and documentation

### Autocomplete Suggestions
- **Decision**: Project and task suggestions come from existing time entries, weighted toward most recent ones
- **Rationale**: Reduces typing for frequently-used projects/tasks and reflects actual user patterns

### Start Entry Flow
- **Decision**: When user selects "start" option, show a list of the 10 most recent unique (project, task) combinations. Default selection is "something else" which allows entering a new project/task manually
- **Rationale**: Quick access to duplicating recent tasks with option to create new ones, reducing friction for both common and new entries

### Main Interface & Navigation
- **Decision**: Default TUI interface is the list of entries. Menu options are: start, stop, exit. Stats will be added in a future spec
- **Rationale**: Keep initial scope focused on core time tracking actions with list view as the hub

## Task List

### Phase 1: Setup & Core Infrastructure
- [ ] Add charmbracelet/bubbles and charmbracelet/bubbletea dependencies
- [ ] Create `cmd/tui/` directory structure
- [ ] Create main TUI model struct implementing bubbletea.Model interface
- [ ] Implement basic TUI initialization and lifecycle (Init, Update, View)
- [ ] Add TUI entry point to main CLI (show when no args provided)

### Phase 2: Main Menu
- [ ] Design main menu with options:
  - Start new time entry
  - Stop current entry
  - List entries
  - View statistics
  - Exit
- [ ] Implement menu navigation and selection
- [ ] Create navigation stack system for moving between screens

### Phase 3: Start/Stop Entry Screens
- [ ] Create "Start Entry" screen with:
  - Project name input with autocomplete
  - Task/title input with autocomplete
  - Submit button
- [ ] Create "Stop Entry" screen
  - Show currently running entry (if any)
  - Confirmation prompt
- [ ] Autocomplete component:
  - Filter suggestions as user types
  - Display top 5-10 matches
  - Allow arrow keys to select suggestion
  - Allow Enter to confirm selection

### Phase 4: Recent Entries & Duplicate
- [ ] Create "Recent Entries" screen showing last 10 entries
- [ ] Display: project, title, duration, start time
- [ ] Implement duplicate entry functionality:
  - Select entry from list
  - Press 'd' to duplicate (start new entry with same project/title)
- [ ] Gracefully handle duplicate of currently running entry

### Phase 5: Integration with Data Store
- [ ] Connect TUI start/stop actions to existing data store logic
- [ ] Ensure data consistency between CLI and TUI operations
- [ ] Test concurrent operations (CLI + TUI)

### Phase 6: Polish & UX
- [ ] Add help text and keybinding hints
- [ ] Improve visual styling with colors and spacing
- [ ] Add loading states for data operations
- [ ] Handle edge cases (empty data, no projects yet, errors)
- [ ] Add ability to go back to previous screen (Escape key)

### Phase 7: README Updates
- [ ] Add "Using the TUI" section with:
  - How to launch TUI (run without args)
  - Keyboard shortcuts reference
  - Autocomplete behavior explanation
  - How to duplicate entries
- [ ] Update "Development" section with:
  - `just` recipes for running dev commands
  - Remove out-of-date build/run instructions
- [ ] Clean up any deprecated information

### Phase 8: Testing
- [ ] Add integration tests for TUI data operations
- [ ] Test autocomplete filtering logic
- [ ] Test navigation between screens
- [ ] Verify data consistency between CLI and TUI
