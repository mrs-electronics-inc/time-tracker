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

### Design Philosophy
- **Decision**: TUI design inspired by helix editor and zellij - intuitive with keybindings prominently displayed in the interface
- **Rationale**: Reduces learning curve and helps users discover features without external documentation

## Task List

### Setup Bubble Tea Framework
- [x] Add charmbracelet/bubbles and charmbracelet/bubbletea dependencies
- [x] Create `cmd/tui/` directory structure
- [x] Create main TUI model struct implementing bubbletea.Model interface
- [x] Implement basic TUI initialization and lifecycle (Init, Update, View)
- [x] Add TUI entry point to main CLI (show when no args provided)

### Build Entry List View
- [ ] Create list view showing recent time entries
- [ ] Display columns: Start, End, Project, Title, Duration (matching CLI list command output)
- [ ] Implement scrolling through entry list
- [ ] Add navigation to start/stop options from list view
- [ ] Display keybindings at bottom of screen (e.g., "s: start, e: stop, q: quit")

### Build Autocomplete Component
- [ ] Extract unique projects and tasks from historical entries, weighted by recency
- [ ] Implement filtering logic as user types
- [ ] Display top suggestions dynamically
- [ ] Allow keyboard navigation (arrow keys) and selection (Enter)
- [ ] Allow free-text entry not limited to historical values

### Build Start Entry Flow
- [ ] Create screen showing 10 most recent unique (project, task) combinations
- [ ] Set "something else" as default selection
- [ ] When "something else" selected, show manual entry input screens:
  - Project name input with autocomplete from historical data
  - Task/title input with autocomplete from historical data
  - Submit to start new entry
- [ ] When recent task selected, immediately start entry with that project/task
- [ ] Display keybindings for navigation and actions (e.g., "↑↓: navigate, enter: select, esc: back")

### Build Stop Entry Flow
- [ ] Create "Stop Entry" screen
- [ ] Show currently running entry (if any)
- [ ] Display project, task, and elapsed time
- [ ] Add confirmation prompt before stopping
- [ ] Handle case when no entry is running
- [ ] Display keybindings for confirmation and navigation (e.g., "y: confirm, n: cancel, esc: back")

### Integration with Data Store
- [ ] Connect TUI start action to existing data store start logic
- [ ] Connect TUI stop action to existing data store stop logic
- [ ] Ensure data consistency between CLI and TUI operations
- [ ] Load historical entries for recent task list and autocomplete suggestions
- [ ] Test concurrent operations (CLI + TUI)

### Polish & UX
- [ ] Implement consistent keybinding display across all screens (footer/status bar)
- [ ] Improve visual styling with colors and spacing (inspired by helix/zellij)
- [ ] Add loading states for data operations
- [ ] Handle edge cases (empty data, no entries yet, errors)
- [ ] Add ability to go back to previous screen (Escape key)
- [ ] Smooth transitions between screens
- [ ] Ensure all interactive elements show available keybindings in the interface

### README Updates
- [ ] Add "Using the TUI" section with:
  - How to launch TUI (run without args)
  - Keyboard shortcuts reference
  - Autocomplete behavior explanation
  - How to duplicate entries
- [ ] Update "Development" section with:
  - `just` recipes for running dev commands
  - Remove out-of-date build/run instructions
- [ ] Clean up any deprecated information

### Testing
- [ ] Add unit tests for model state transitions (sending messages and verifying state changes)
- [ ] Add integration tests for TUI data operations:
  - Start entry via TUI and verify data file
  - Stop entry via TUI and verify data file
  - Load recent entries and verify they match CLI list output
- [ ] Test autocomplete filtering and ranking logic
- [ ] Verify data consistency between CLI and TUI operations
- [ ] Test edge cases (no data, invalid input, concurrent operations)
