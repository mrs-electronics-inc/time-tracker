# TUI Implementation

Add a TUI to improve the user experience.

## Design Decisions

### Framework Choice: Bubble Tea

- **Decision**: Use charmbracelet/bubbletea and charmbracelet/bubbles for TUI framework
- **Rationale**: Referenced in issue #2, purpose-built for Go TUIs with good component library and documentation

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

- [x] Create list view showing recent time entries
- [x] Display columns: Start, End, Project, Title, Duration (matching CLI list command output)
- [x] Implement scrolling through entry list
- [x] Add navigation to start/stop options from list view

### Build Start Entry Flow

- [x] Add dialog state to Model (dialogMode, inputs []textinput.Model, focusIndex)
- [x] Initialize textinput components with proper styling (focused/blurred styles)
- [x] Implement dialog mode in Update handler:
  - [x] Trigger dialog on 's' key in list mode, pre-populate with selected entry values
  - [x] Handle Tab/Up/Down/Shift+Tab to cycle focus between project and title inputs
  - [x] Handle Enter to submit and call taskManager.StartEntry()
  - [x] Handle Esc to cancel and return to list
  - [x] Route character input to focused textinput
- [x] Implement dialog rendering in View:
  - [x] Show modal dialog with project and title inputs
  - [x] Display keybinding hints (Enter: Submit | Esc: Cancel)
  - [x] Handle special case: if selected entry is running, show "stop and start blank" option
- [x] Reload entries after successful start and return to list view

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
  - Basic keyboard shortcuts reference
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
- [ ] Verify data consistency between CLI and TUI operations
- [ ] Test edge cases (no data, invalid input, concurrent operations)
