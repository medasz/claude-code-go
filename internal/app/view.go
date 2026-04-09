package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	backgroundColor = lipgloss.Color("#111315")
	panelColor      = lipgloss.Color("#171A1F")
	panelColorSoft  = lipgloss.Color("#1D222A")
	panelColorMute  = lipgloss.Color("#14181D")
	panelBorder     = lipgloss.Color("#2B313C")
	mutedColor      = lipgloss.Color("#7D8594")
	softTextColor   = lipgloss.Color("#D7DBE2")
	brandColor      = lipgloss.Color("#F3EBD2")
	noticeAccent    = lipgloss.Color("#E6C17A")
	userAccent      = lipgloss.Color("#89D3C7")
	toolAccent      = lipgloss.Color("#BFC7D5")
	errorAccent     = lipgloss.Color("#FF8B8B")
)

const messageGutterWidth = 14

func renderView(m model) string {
	transcriptWidth := max(40, min(96, m.viewport.Width))
	if m.startup {
		startupWidth := max(84, min(132, m.width-4))
		return renderStartupView(m, startupWidth)
	}

	header := lipgloss.PlaceHorizontal(
		transcriptWidth,
		lipgloss.Center,
		renderHeader(m, transcriptWidth),
	)

	transcript := lipgloss.PlaceHorizontal(
		transcriptWidth,
		lipgloss.Center,
		m.viewport.View(),
	)

	composer := lipgloss.PlaceHorizontal(
		transcriptWidth,
		lipgloss.Center,
		renderComposer(m, transcriptWidth),
	)

	status := lipgloss.PlaceHorizontal(
		transcriptWidth,
		lipgloss.Center,
		statusLineStyle(transcriptWidth).Render(buildStatusLine(m)),
	)

	body := lipgloss.JoinVertical(lipgloss.Left, header, transcript, composer, status)
	return lipgloss.NewStyle().Background(backgroundColor).Render(appFrameStyle.Render(body))
}

func renderStartupView(m model, width int) string {
	card := lipgloss.PlaceHorizontal(width, lipgloss.Left, renderStartupCard(width))
	notice := lipgloss.PlaceHorizontal(width, lipgloss.Left, renderStartupNotice(width))
	composer := lipgloss.PlaceHorizontal(width, lipgloss.Left, renderStartupComposer(m, width))
	footer := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Left,
		lipgloss.NewStyle().
			Width(width).
			Foreground(mutedColor).
			Render("? for shortcuts"),
	)

	body := lipgloss.JoinVertical(lipgloss.Left, card, "", notice, "", composer, footer)
	return lipgloss.NewStyle().Background(backgroundColor).Render(appFrameStyle.Render(body))
}

func renderHeader(m model, width int) string {
	logo := lipgloss.NewStyle().
		Foreground(brandColor).
		Bold(true).
		Render("Claude Code")

	subtitle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Terminal agent interface replica")

	notice := "Idle"
	if m.running {
		notice = fmt.Sprintf("%s Running %s", spinnerFrames[m.spinnerFrame], m.currentCmd)
	}

	versionChip := lipgloss.NewStyle().
		Foreground(mutedColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Background(panelColor).
		Padding(0, 1).
		Render("v0.1")

	pathChip := lipgloss.NewStyle().
		Foreground(softTextColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Background(panelColor).
		Padding(0, 1).
		Render("D:\\code\\claude-code-go")

	noticeBox := lipgloss.NewStyle().
		Foreground(noticeAccent).
		Background(panelColorSoft).
		Padding(0, 1).
		Render(notice)

	sessionBox := lipgloss.NewStyle().
		Foreground(mutedColor).
		Background(panelColorMute).
		Padding(0, 1).
		Render(fmt.Sprintf("%d messages", len(m.messages)))

	notices := renderHeaderNotices(m, width)

	headerCard := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, logo, " ", versionChip),
		subtitle,
		lipgloss.JoinHorizontal(lipgloss.Left, noticeBox, " ", sessionBox),
		lipgloss.NewStyle().Foreground(mutedColor).Render("workspace"),
		pathChip,
		notices,
		lipgloss.NewStyle().
			Width(width).
			Foreground(panelBorder).
			Render(strings.Repeat("-", max(0, width-2))),
	)

	return lipgloss.NewStyle().
		Width(width).
		MarginBottom(1).
		Render(headerCard)
}

func renderStartupCard(width int) string {
	leftWidth := max(31, width/3)
	rightWidth := max(38, width-leftWidth-5)

	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(noticeAccent).Bold(true).Render(" Claude Code "),
		lipgloss.NewStyle().Foreground(mutedColor).Bold(true).Render("v999.0.0-restored"),
	)

	left := lipgloss.NewStyle().
		Width(leftWidth).
		Padding(2, 2, 1, 1).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Width(leftWidth-3).Align(lipgloss.Center).Bold(true).Foreground(softTextColor).Render("Welcome back!"),
			"",
			lipgloss.NewStyle().Width(leftWidth-3).Align(lipgloss.Center).Foreground(noticeAccent).Render(startupMascot()),
			"",
			lipgloss.NewStyle().Width(leftWidth-3).Align(lipgloss.Center).Foreground(mutedColor).Render("Sonnet 4.6 | API Usage Billing"),
			lipgloss.NewStyle().Width(leftWidth-3).Align(lipgloss.Center).Foreground(mutedColor).Render("D:\\code\\claude-code"),
		))

	right := lipgloss.NewStyle().
		Width(rightWidth).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(noticeAccent).
		Padding(2, 1, 1, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(noticeAccent).Bold(true).Render("Tips for getting started"),
			lipgloss.NewStyle().Foreground(softTextColor).Render("Run /init to create a CLAUDE.md file with instructions for Claude"),
			lipgloss.NewStyle().Foreground(noticeAccent).Render(strings.Repeat("-", max(0, rightWidth-8))),
			lipgloss.NewStyle().Foreground(noticeAccent).Bold(true).Render("Recent activity"),
			lipgloss.NewStyle().Foreground(mutedColor).Render("No recent activity"),
		))

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	return renderStartupFrame(content, title, width)
}

func renderStartupNotice(width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Foreground(mutedColor).
		Padding(1, 1, 0, 1).
		Render("^ Opus now defaults to 1M context | 5x more room, same pricing")
}

func renderStartupComposer(m model, width int) string {
	line := lipgloss.NewStyle().
		Width(width).
		Foreground(mutedColor).
		Render(strings.Repeat("-", max(0, width-2)))

	promptText := strings.TrimSpace(m.input.Value())
	promptColor := mutedColor
	if promptText == "" {
		promptText = `Try "fix typecheck errors"`
	} else {
		promptColor = softTextColor
	}

	cursor := lipgloss.NewStyle().
		Background(softTextColor).
		Foreground(backgroundColor).
		Render(" ")

	prompt := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(softTextColor).Bold(true).Render(">"),
		" ",
		cursor,
		lipgloss.NewStyle().Foreground(promptColor).Render(promptText),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		line,
		lipgloss.NewStyle().Width(width).Padding(0, 1).Render(prompt),
		line,
	)
}

func renderTranscript(m model, width int) string {
	lines := make([]string, 0, len(m.messages))
	for _, item := range m.messages {
		if item.kind == messageKindNotice {
			continue
		}
		lines = append(lines, renderMessage(item, width))
	}
	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 0, 1, 0).
		Render(strings.Join(lines, "\n\n"))
}

func renderMessage(item message, width int) string {
	bodyWidth := max(18, width-messageGutterWidth-2)

	gutterLabel := strings.ToLower(item.title)
	gutterColor := mutedColor
	if item.pending {
		gutterLabel = fmt.Sprintf("%s %s", spinnerFrames[currentSpinnerFrame(item)], gutterLabel)
		gutterColor = noticeAccent
	}

	gutterTitle := lipgloss.NewStyle().
		Foreground(gutterColor).
		Render(gutterLabel)

	gutterTime := lipgloss.NewStyle().
		Foreground(panelBorder).
		Render(item.at.Format("15:04"))

	gutter := lipgloss.NewStyle().
		Width(messageGutterWidth).
		PaddingRight(1).
		Render(lipgloss.JoinVertical(lipgloss.Right, gutterTitle, gutterTime))

	content := lipgloss.NewStyle().
		Width(bodyWidth).
		Foreground(softTextColor).
		Render(indentMultiline(item.content))

	var bodyStyle lipgloss.Style
	switch item.kind {
	case messageKindUser:
		bodyStyle = lipgloss.NewStyle().
			Width(bodyWidth).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(userAccent).
			Background(panelColorSoft).
			Padding(0, 1)
	case messageKindTool:
		bodyStyle = lipgloss.NewStyle().
			Width(bodyWidth).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(toolAccent).
			Background(panelColorMute).
			Padding(0, 1)
	case messageKindError:
		bodyStyle = lipgloss.NewStyle().
			Width(bodyWidth).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(errorAccent).
			Background(panelColorMute).
			Padding(0, 1)
	default:
		bodyStyle = lipgloss.NewStyle().
			Width(bodyWidth).
			Padding(0, 1).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(panelBorder)
	}

	if item.pending {
		bodyStyle = bodyStyle.
			Foreground(noticeAccent).
			Background(panelColorSoft)
	}

	body := bodyStyle.Render(content)
	row := lipgloss.JoinHorizontal(lipgloss.Top, gutter, body)
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, row)
}

func currentSpinnerFrame(item message) int {
	if len(spinnerFrames) == 0 {
		return 0
	}
	tick := int(item.at.UnixNano()/int64(120000000)) % len(spinnerFrames)
	if tick < 0 {
		return 0
	}
	return tick
}

func renderHeaderNotices(m model, width int) string {
	noticeItems := make([]string, 0, 3)
	for _, item := range m.messages {
		if item.kind != messageKindNotice {
			continue
		}
		noticeItems = append(noticeItems, lipgloss.NewStyle().
			Width(max(0, width-4)).
			Foreground(softTextColor).
			Background(panelColorMute).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(noticeAccent).
			Padding(0, 1).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(noticeAccent).Bold(true).Render(item.title),
					lipgloss.NewStyle().Foreground(softTextColor).Render(item.content),
				),
			))
		if len(noticeItems) == 3 {
			break
		}
	}

	if len(noticeItems) == 0 {
		return ""
	}

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(strings.Join(noticeItems, "\n"))
}

func renderComposer(m model, width int) string {
	title := lipgloss.NewStyle().
		Foreground(softTextColor).
		Bold(true).
		Render("Message")

	hint := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Enter to send  Ctrl+C to quit")

	shortcuts := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Esc clears focus  Shell backend enabled")

	toolbar := lipgloss.JoinHorizontal(
		lipgloss.Left,
		title,
		"  ",
		lipgloss.NewStyle().Foreground(mutedColor).Render("mode"),
		" ",
		lipgloss.NewStyle().
			Foreground(softTextColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(panelBorder).
			Background(panelColorSoft).
			Padding(0, 1).
			Render("shell"),
		" ",
		lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(panelColorMute).
			Padding(0, 1).
			Render(currentShellLabel()),
	)

	inputWrap := lipgloss.NewStyle().
		Width(max(0, width-4)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Background(panelColorSoft).
		Padding(0, 1).
		Render(m.input.View())

	inputBody := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Background(panelColor).
		Padding(1, 1, 0, 1).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			toolbar,
			inputWrap,
			lipgloss.JoinVertical(lipgloss.Left, hint, shortcuts),
		))

	return lipgloss.NewStyle().MarginTop(1).Render(inputBody)
}

func startupMascot() string {
	return strings.Join([]string{
		" " + "\u2590" + "\u259b\u2588\u2588\u2588\u259c" + "\u258c" + " ",
		"\u259d\u259c" + "\u2588\u2588\u2588\u2588\u2588" + "\u259b\u2598",
		"  " + "\u2598\u2598" + " " + "\u259d\u259d" + "  ",
	}, "\n")
}

func renderStartupFrame(content, title string, width int) string {
	innerWidth := max(10, width-2)
	lines := strings.Split(content, "\n")
	framed := make([]string, 0, len(lines)+2)

	topWidth := lipgloss.Width(title)
	topLine := lipgloss.NewStyle().Foreground(noticeAccent).Render("\u256d")
	topLine += title
	topLine += lipgloss.NewStyle().Foreground(noticeAccent).Render(strings.Repeat("\u2500", max(0, innerWidth-topWidth)))
	topLine += lipgloss.NewStyle().Foreground(noticeAccent).Render("\u256e")
	framed = append(framed, topLine)

	leftBorder := lipgloss.NewStyle().Foreground(noticeAccent).Render("\u2502")
	rightBorder := lipgloss.NewStyle().Foreground(noticeAccent).Render("\u2502")
	for _, line := range lines {
		padding := max(0, innerWidth-lipgloss.Width(line))
		framed = append(framed, leftBorder+line+strings.Repeat(" ", padding)+rightBorder)
	}

	bottom := lipgloss.NewStyle().Foreground(noticeAccent).Render("\u2570" + strings.Repeat("\u2500", innerWidth) + "\u256f")
	framed = append(framed, bottom)
	return strings.Join(framed, "\n")
}

func buildStatusLine(m model) string {
	state := "idle"
	if m.running {
		state = "running"
	}

	errorText := "none"
	if m.lastError != "" {
		errorText = m.lastError
	}

	return fmt.Sprintf("model: shell | state: %s | messages: %d | runtime: %s | last-error: %s", state, len(m.messages), currentShellLabel(), errorText)
}

func indentMultiline(s string) string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return s
	}
	return strings.ReplaceAll(s, "\n", "\n")
}

func statusLineStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Foreground(mutedColor).
		Background(backgroundColor).
		MarginTop(1)
}
