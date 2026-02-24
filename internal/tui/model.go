package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stonecharioteer/goforgo/internal/exercise"
	"github.com/stonecharioteer/goforgo/internal/runner"
	"github.com/stonecharioteer/goforgo/internal/watcher"
)

// ViewMode represents the current view state of the TUI.
type ViewMode int

const (
	ViewSplash  ViewMode = iota
	ViewWelcome
	ViewMain
	ViewList
	ViewHint
	ViewOutput
)

// Model represents the TUI application state
type Model struct {
	// Exercise management
	exerciseManager *exercise.ExerciseManager
	currentExercise *exercise.Exercise
	currentIndex    int
	exercises       []*exercise.Exercise

	// Execution and validation
	runner        *runner.Runner
	lastResult    *runner.Result
	isRunning     bool
	currentHintLevel int  // Track current hint level (0=none, 1=level1, 2=level1+2, 3=all)

	// File watching
	watcher    *watcher.Watcher
	watcherErr error

	// UI state
	viewMode ViewMode
	width    int
	height   int
	ready    bool
	
	// List view state for scrollable exercise list
	listSelectedIndex int // Currently selected item in list
	listScrollOffset  int // Scroll offset for list view
	listViewHeight    int // Available height for list items
	
	// Filter state
	filterMode bool   // Whether we're in filter mode
	filterText string // Current filter text
	
	// Output view state
	outputScrollPos  int  // Current scroll position in output view
	outputViewHeight int  // Available height for output content
	
	// Progress and statistics
	// Counts are now calculated dynamically via exerciseManager methods
	
	// Messages and status
	statusMessage string
	splashFrame   int
}

// NewModel creates a new TUI model
func NewModel(exerciseManager *exercise.ExerciseManager, runner *runner.Runner) *Model {
	exercises := exerciseManager.GetExercises()
	currentEx := exerciseManager.GetNextExercise()
	currentIndex := 0
	
	// Find the index of the current exercise
	for i, ex := range exercises {
		if currentEx != nil && ex.Info.Name == currentEx.Info.Name {
			currentIndex = i
			break
		}
	}

	// Exercise counts are now handled dynamically by ExerciseManager

	return &Model{
		exerciseManager: exerciseManager,
		currentExercise: currentEx,
		currentIndex:    currentIndex,
		exercises:       exercises,
		runner:          runner,
		viewMode:        ViewSplash,
		splashFrame:     0,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.runCurrentExercise(),
		m.startFileWatcher(),
		m.splashTick(), // Start splash animation
	)
}

// Update handles messages and state changes
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case exerciseResultMsg:
		m.lastResult = msg.result
		m.isRunning = false
		m.statusMessage = ""
		
		// Mark exercise as completed if successful and not already completed
		if msg.result.Success && m.currentExercise != nil && !m.currentExercise.Completed {
			if err := m.exerciseManager.MarkExerciseCompleted(m.currentExercise.Info.Name); err == nil {
				// Update local completion tracking
				m.currentExercise.Completed = true
				
				// Update exercises list with fresh completion status
				m.exercises = m.exerciseManager.GetExercises()
			}
		}
		
		return m, m.waitForFileChange(m.watcher)

	case exerciseRunningMsg:
		m.isRunning = true
		m.statusMessage = "Running exercise..."
		return m, nil

	case fileChangedMsg:
		if !m.isRunning {
			return m, m.runCurrentExercise()
		}
		return m, nil

	case watcherErrorMsg:
		m.watcherErr = msg.err
		return m, nil

	case continueWatchingMsg:
		// Continue listening for file changes
		if m.watcher != nil {
			return m, m.waitForFileChange(m.watcher)
		}
		return m, nil

	case splashTickMsg:
		if m.viewMode == ViewSplash {
			m.splashFrame++
			if m.splashFrame >= splashFrameCount {
				m.viewMode = ViewWelcome
				return m, nil
			}
			return m, m.splashTick()
		}
		return m, nil

	case statusMsg:
		m.statusMessage = msg.message
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "n":
		// Next exercise
		if m.viewMode == ViewHint || m.viewMode == ViewList {
			m.viewMode = ViewMain
			return m, nil
		}
		return m, m.nextExercise()

	case "p":
		// Previous exercise
		if m.viewMode == ViewHint || m.viewMode == ViewList {
			m.viewMode = ViewMain
			return m, nil
		}
		return m, m.previousExercise()

	case "h":
		// Show next hint level or hide if at max
		if m.viewMode != ViewHint {
			m.currentHintLevel = 1
			m.viewMode = ViewHint
		} else {
			maxLevel := m.getMaxHintLevel()
			if m.currentHintLevel < maxLevel {
				m.currentHintLevel++
			} else {
				m.viewMode = ViewMain
				m.currentHintLevel = 0
			}
		}
		return m, nil

	case "l":
		// Toggle exercise list
		if m.viewMode != ViewList {
			m.viewMode = ViewList
			m.listSelectedIndex = m.currentIndex
			m.listScrollOffset = 0
			m.listViewHeight = max(m.height-listReservedHeight, minListHeight)
			m.ensureSelectedVisible()
		} else {
			m.viewMode = ViewMain
		}
		return m, nil

	case "r":
		// Manually run exercise
		if !m.isRunning {
			return m, m.runCurrentExercise()
		}
		return m, nil

	case "s":
		// Show output view
		if m.viewMode == ViewMain && m.lastResult != nil {
			m.viewMode = ViewOutput
			m.outputScrollPos = 0
			m.outputViewHeight = max(m.height-listReservedHeight, minListHeight)
			return m, nil
		}
		return m, nil

	case "up", "k":
		if m.viewMode == ViewList && !m.filterMode {
			return m, m.moveListSelection(-1)
		}
		if m.viewMode == ViewOutput {
			return m, m.scrollOutput(-1)
		}
		return m, nil

	case "down", "j":
		if m.viewMode == ViewList && !m.filterMode {
			return m, m.moveListSelection(1)
		}
		if m.viewMode == ViewOutput {
			return m, m.scrollOutput(1)
		}
		return m, nil

	case "page_up":
		if m.viewMode == ViewList && !m.filterMode {
			return m, m.moveListSelection(-m.listViewHeight)
		}
		if m.viewMode == ViewOutput {
			return m, m.scrollOutput(-m.outputViewHeight)
		}
		return m, nil

	case "page_down":
		if m.viewMode == ViewList && !m.filterMode {
			return m, m.moveListSelection(m.listViewHeight)
		}
		if m.viewMode == ViewOutput {
			return m, m.scrollOutput(m.outputViewHeight)
		}
		return m, nil

	case "home":
		if m.viewMode == ViewList && !m.filterMode {
			m.listSelectedIndex = 0
			m.ensureSelectedVisible()
			return m, nil
		}
		if m.viewMode == ViewOutput {
			m.outputScrollPos = 0
			return m, nil
		}
		return m, nil

	case "end":
		if m.viewMode == ViewList && !m.filterMode {
			filteredExercises := m.getFilteredExercises()
			m.listSelectedIndex = len(filteredExercises) - 1
			m.ensureSelectedVisible()
			return m, nil
		}
		if m.viewMode == ViewOutput {
			return m, m.scrollToBottom()
		}
		return m, nil

	case "backspace":
		// Handle backspace in filter mode
		if m.filterMode && len(m.filterText) > 0 {
			m.filterText = m.filterText[:len(m.filterText)-1]
			return m, nil
		}
		return m, nil

	default:
		// Handle text input in filter mode
		if m.filterMode && len(msg.String()) == 1 {
			char := msg.String()
			// Only allow alphanumeric characters, underscore, and space
			if (char >= "a" && char <= "z") || (char >= "A" && char <= "Z") || 
			   (char >= "0" && char <= "9") || char == "_" || char == " " {
				m.filterText += char
				return m, nil
			}
		}
		return m, nil

	case "/":
		// Enter filter mode when in list view
		if m.viewMode == ViewList && !m.filterMode {
			m.filterMode = true
			m.filterText = ""
			return m, nil
		}
		return m, nil

	case "enter", "esc":
		if m.viewMode == ViewSplash {
			m.viewMode = ViewWelcome
			return m, nil
		}
		if m.viewMode == ViewWelcome {
			m.viewMode = ViewMain
			return m, nil
		}
		if m.viewMode == ViewList && msg.String() == "enter" {
			if m.filterMode {
				m.filterMode = false
				m.listSelectedIndex = 0
				m.listScrollOffset = 0
				return m, nil
			}
			// Select the highlighted exercise
			filteredExercises := m.getFilteredExercises()
			if m.listSelectedIndex >= 0 && m.listSelectedIndex < len(filteredExercises) {
				selectedExercise := filteredExercises[m.listSelectedIndex]
				for i, ex := range m.exercises {
					if ex == selectedExercise {
						m.currentIndex = i
						m.currentExercise = ex
						break
					}
				}
				m.currentHintLevel = 0
				m.viewMode = ViewMain
				m.filterMode = false
				m.filterText = ""
				return m, m.runCurrentExercise()
			}
		}
		if msg.String() == "esc" {
			if m.filterMode {
				m.filterMode = false
				m.filterText = ""
				return m, nil
			} else if m.viewMode == ViewOutput {
				m.viewMode = ViewMain
				m.outputScrollPos = 0
				return m, nil
			} else if m.viewMode == ViewList && m.filterText != "" {
				m.filterText = ""
				m.listSelectedIndex = 0
				m.listScrollOffset = 0
				return m, nil
			}
		}
		// Dismiss any overlay view
		m.viewMode = ViewMain
		m.filterMode = false
		m.filterText = ""
		m.currentHintLevel = 0
		return m, nil
	}

	return m, nil
}

// View renders the TUI
func (m *Model) View() string {
	if !m.ready {
		return "Initializing GoForGo..."
	}

	switch m.viewMode {
	case ViewSplash:
		return m.renderSplash()
	case ViewWelcome:
		return m.renderWelcome()
	case ViewList:
		return m.renderExerciseList()
	case ViewHint:
		return m.renderHint()
	case ViewOutput:
		return m.renderOutput()
	default:
		return m.renderMain()
	}
}

// Custom messages for the tea program
type exerciseResultMsg struct {
	result *runner.Result
}

type exerciseRunningMsg struct{}

type fileChangedMsg struct {
	path string
}

type watcherErrorMsg struct {
	err error
}

type continueWatchingMsg struct{}

type statusMsg struct {
	message string
}

type splashTickMsg struct{}

// Commands
func (m *Model) runCurrentExercise() tea.Cmd {
	if m.currentExercise == nil {
		return nil
	}

	return tea.Batch(
		func() tea.Msg { return exerciseRunningMsg{} },
		func() tea.Msg {
			result, _ := m.runner.RunExercise(m.currentExercise)
			return exerciseResultMsg{result: result}
		},
	)
}

func (m *Model) nextExercise() tea.Cmd {
	if m.currentIndex < len(m.exercises)-1 {
		m.currentIndex++
		m.currentExercise = m.exercises[m.currentIndex]
		m.currentHintLevel = 0  // Reset hint level for new exercise
		return m.runCurrentExercise()
	}
	return func() tea.Msg {
		return statusMsg{message: "You've reached the last exercise!"}
	}
}

func (m *Model) previousExercise() tea.Cmd {
	if m.currentIndex > 0 {
		m.currentIndex--
		m.currentExercise = m.exercises[m.currentIndex]
		m.currentHintLevel = 0  // Reset hint level for new exercise
		return m.runCurrentExercise()
	}
	return func() tea.Msg {
		return statusMsg{message: "You're at the first exercise!"}
	}
}

func (m *Model) startFileWatcher() tea.Cmd {
	w, err := watcher.NewWatcher()
	if err != nil {
		return func() tea.Msg {
			return watcherErrorMsg{err: err}
		}
	}

	m.watcher = w

	// Watch the exercises directory recursively
	exercisesDir := m.exerciseManager.ExercisesPath
	if err := w.WatchRecursive(exercisesDir); err != nil {
		return func() tea.Msg {
			return watcherErrorMsg{err: err}
		}
	}

	// Start watching for file changes
	return m.waitForFileChange(w)
}

func (m *Model) waitForFileChange(w *watcher.Watcher) tea.Cmd {
	return func() tea.Msg {
		select {
		case event := <-w.Events():
			if m.shouldProcessFileEvent(event) {
				return fileChangedMsg{path: event.Name}
			}
			// Event not relevant, continue listening
			return continueWatchingMsg{}
		case err := <-w.Errors():
			return watcherErrorMsg{err: err}
		}
	}
}

func (m *Model) shouldProcessFileEvent(event watcher.Event) bool {
	// Many editors use atomic writes (create, rename), so we watch for more than just Write events.
	isModification := event.IsWrite() || event.IsCreate() || event.IsRename()
	if !isModification {
		return false
	}

	if !strings.HasSuffix(event.Name, ".go") {
		return false
	}

	if m.currentExercise == nil {
		return false
	}

	// Check if it's the current exercise file
	return strings.Contains(event.Name, m.currentExercise.Info.Name)
}

// Styles
var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("#7C3AED"))

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1F2937"))

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)

	hintStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Italic(true)

	codeStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#F3F4F6")).
		Foreground(lipgloss.Color("#1F2937")).
		Padding(0, 1)

	progressBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C3AED"))

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
)

// getTotalCount returns the total number of exercises (dynamic)
func (m *Model) getTotalCount() int {
	return m.exerciseManager.GetTotalExerciseCount()
}

// getCompletedCount returns the number of completed exercises (dynamic)
func (m *Model) getCompletedCount() int {
	return m.exerciseManager.GetCompletedExerciseCount()
}

// getMaxHintLevel returns the maximum hint level available for the current exercise
func (m *Model) getMaxHintLevel() int {
	if m.currentExercise == nil {
		return 0
	}
	
	maxLevel := 0
	if m.currentExercise.Hints.Level1 != "" {
		maxLevel = 1
	}
	if m.currentExercise.Hints.Level2 != "" {
		maxLevel = 2
	}
	if m.currentExercise.Hints.Level3 != "" {
		maxLevel = 3
	}
	
	return maxLevel
}

// splashTick creates a command for splash screen animation
func (m *Model) splashTick() tea.Cmd {
	return tea.Tick(time.Millisecond*splashTickMs, func(time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

// moveListSelection moves the selection in the list view
func (m *Model) moveListSelection(delta int) tea.Cmd {
	filteredExercises := m.getFilteredExercises()
	newIndex := m.listSelectedIndex + delta
	
	// Clamp to valid range (no wrapping)
	if newIndex < 0 {
		newIndex = 0
	} else if newIndex >= len(filteredExercises) {
		newIndex = len(filteredExercises) - 1
	}
	
	m.listSelectedIndex = newIndex
	m.ensureSelectedVisible()
	
	return nil
}

// ensureSelectedVisible adjusts scroll offset to keep selected item visible
func (m *Model) ensureSelectedVisible() {
	filteredExercises := m.getFilteredExercises()
	if m.listSelectedIndex < m.listScrollOffset {
		// Selected item is above visible area
		m.listScrollOffset = m.listSelectedIndex
	} else if m.listSelectedIndex >= m.listScrollOffset+m.listViewHeight {
		// Selected item is below visible area
		m.listScrollOffset = m.listSelectedIndex - m.listViewHeight + 1
	}
	
	// Ensure scroll offset doesn't go negative or exceed filtered exercise count
	if m.listScrollOffset < 0 {
		m.listScrollOffset = 0
	}
	if len(filteredExercises) > 0 && m.listScrollOffset >= len(filteredExercises) {
		m.listScrollOffset = len(filteredExercises) - 1
	}
}

// getFilteredExercises returns exercises filtered by the current filter text
func (m *Model) getFilteredExercises() []*exercise.Exercise {
	if m.filterText == "" {
		return m.exercises
	}
	
	var filtered []*exercise.Exercise
	filterLower := strings.ToLower(m.filterText)
	
	for _, ex := range m.exercises {
		// Check exercise name
		if strings.Contains(strings.ToLower(ex.Info.Name), filterLower) {
			filtered = append(filtered, ex)
			continue
		}
		
		// Check exercise title
		if strings.Contains(strings.ToLower(ex.Description.Title), filterLower) {
			filtered = append(filtered, ex)
			continue
		}
		
		// Check category
		if strings.Contains(strings.ToLower(ex.Info.Category), filterLower) {
			filtered = append(filtered, ex)
			continue
		}
		
		// Check difficulty
		difficulty := ex.GetDifficultyString()
		if strings.Contains(strings.ToLower(difficulty), filterLower) {
			filtered = append(filtered, ex)
			continue
		}
	}
	
	return filtered
}

// scrollOutput scrolls the output view by the given delta
func (m *Model) scrollOutput(delta int) tea.Cmd {
	if m.lastResult == nil {
		return nil
	}
	
	// Split output into lines for scrolling
	outputLines := strings.Split(m.lastResult.Output, "\n")
	maxScroll := max(0, len(outputLines)-m.outputViewHeight)
	
	m.outputScrollPos += delta
	
	// Clamp scroll position
	if m.outputScrollPos < 0 {
		m.outputScrollPos = 0
	} else if m.outputScrollPos > maxScroll {
		m.outputScrollPos = maxScroll
	}
	
	return nil
}

// scrollToBottom scrolls the output view to the bottom
func (m *Model) scrollToBottom() tea.Cmd {
	if m.lastResult == nil {
		return nil
	}
	
	outputLines := strings.Split(m.lastResult.Output, "\n")
	maxScroll := max(0, len(outputLines)-m.outputViewHeight)
	m.outputScrollPos = maxScroll
	
	return nil
}