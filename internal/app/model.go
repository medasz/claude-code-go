package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type messageKind string

const (
	messageKindNotice    messageKind = "notice"
	messageKindUser      messageKind = "user"
	messageKindAssistant messageKind = "assistant"
	messageKindTool      messageKind = "tool"
	messageKindError     messageKind = "error"
)

type message struct {
	kind    messageKind
	title   string
	content string
	at      time.Time
	pending bool
}

type commandFinishedMsg struct {
	command string
	output  string
	errText string
	err     error
}

type spinnerTickMsg time.Time

type model struct {
	viewport      viewport.Model
	input         textarea.Model
	messages      []message
	width         int
	height        int
	startup       bool
	running       bool
	currentCmd    string
	lastError     string
	spinnerFrame  int
	maxTranscript int
}

func New() model {
	input := textarea.New()
	input.Placeholder = "Message Claude Code or enter a shell command"
	input.Prompt = ""
	input.CharLimit = 0
	input.SetWidth(80)
	input.SetHeight(4)
	input.ShowLineNumbers = false
	input.KeyMap.InsertNewline.SetEnabled(false)
	input.Focus()

	vp := viewport.New(0, 0)

	m := model{
		viewport:      vp,
		input:         input,
		startup:       true,
		maxTranscript: 96,
		messages: []message{
			newMessage(messageKindNotice, "Welcome", "This phase focuses on the Claude Code shell layout: header, transcript, composer, and status line."),
			newMessage(messageKindNotice, "Current Scope", "The backend still runs shell commands for now, while the outer UI is being rebuilt toward a Claude Code style REPL."),
			newMessage(messageKindAssistant, "Claude", "You can enter a command directly, or use this as a chat-style prompt."),
		},
	}
	m.refreshViewport()
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, spinnerTick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		return m, nil
	case spinnerTickMsg:
		if m.running {
			m.spinnerFrame = (m.spinnerFrame + 1) % len(spinnerFrames)
			m.refreshViewport()
		}
		return m, spinnerTick()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.running {
				return m, nil
			}

			command := strings.TrimSpace(m.input.Value())
			if command == "" {
				return m, nil
			}

			m.startup = false
			m.running = true
			m.currentCmd = command
			m.lastError = ""
			m.spinnerFrame = 0
			m.messages = append(m.messages,
				newMessage(messageKindUser, "You", command),
				newPendingMessage(messageKindAssistant, "Claude", fmt.Sprintf("Running %s", command)),
			)
			m.trimMessages()
			m.input.Reset()
			m.refreshViewport()
			return m, runCommand(command)
		}
	case commandFinishedMsg:
		m.running = false
		m.currentCmd = ""

		output := strings.TrimSpace(msg.output)
		errOutput := strings.TrimSpace(msg.errText)

		if output != "" {
			m.messages = append(m.messages, newMessage(messageKindTool, "Command Output", output))
		}

		if errOutput != "" {
			m.messages = append(m.messages, newMessage(messageKindError, "stderr", errOutput))
		}

		m.resolvePendingAssistant(msg, output, errOutput)

		if msg.err != nil {
			m.lastError = msg.err.Error()
		}

		m.trimMessages()
		m.refreshViewport()
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return renderView(m)
}

func (m *model) resize() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	horizontalFrame, verticalFrame := appFrameStyle.GetFrameSize()
	headerHeight := 11
	statusHeight := 1
	inputHeight := 8
	if m.startup {
		headerHeight = 16
		inputHeight = 4
	}

	m.viewport.Width = max(0, min(110, m.width-horizontalFrame))
	m.viewport.Height = max(0, m.height-verticalFrame-headerHeight-statusHeight-inputHeight)
	m.input.SetWidth(max(0, m.viewport.Width-4))
	m.refreshViewport()
}

func (m *model) refreshViewport() {
	m.viewport.SetContent(renderTranscript(*m, max(40, m.viewport.Width)))
	m.viewport.GotoBottom()
}

func (m *model) trimMessages() {
	if len(m.messages) <= m.maxTranscript {
		return
	}
	m.messages = m.messages[len(m.messages)-m.maxTranscript:]
}

func newMessage(kind messageKind, title, content string) message {
	return message{
		kind:    kind,
		title:   title,
		content: content,
		at:      time.Now(),
	}
}

func newPendingMessage(kind messageKind, title, content string) message {
	msg := newMessage(kind, title, content)
	msg.pending = true
	return msg
}

func (m *model) resolvePendingAssistant(result commandFinishedMsg, output, errOutput string) {
	for i := len(m.messages) - 1; i >= 0; i-- {
		msg := m.messages[i]
		if msg.kind != messageKindAssistant || !msg.pending {
			continue
		}

		msg.pending = false
		msg.at = time.Now()

		switch {
		case result.err != nil:
			msg.content = fmt.Sprintf("Command failed: %v", result.err)
		case output == "" && errOutput == "":
			msg.content = "Command finished with no output."
		default:
			msg.content = "Command completed."
		}

		m.messages[i] = msg
		return
	}
}

func spinnerTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	appFrameStyle = lipgloss.NewStyle().
			Padding(1, 2)
	spinnerFrames = []string{"-", "\\", "|", "/"}
)
