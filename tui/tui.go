package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"time"
)

type Model struct {
	Summary         string
	ReportGenerated bool
	AwaitLength     bool
	Token           string
	Period          string
	Since           time.Time
	Done            bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if m.AwaitLength {
			switch key {
			case "1", "2", "3", "4", "5":
				prefix := ""
				var since time.Time
				now := time.Now()
				switch key {
				case "1":
					prefix = "daily"
					since = now.AddDate(0, 0, -1)
				case "2":
					prefix = "weekly"
					since = now.AddDate(0, 0, -7)
				case "3":
					prefix = "monthly"
					since = now.AddDate(0, -1, 0)
				case "4":
					prefix = "6-month"
					since = now.AddDate(0, -6, 0)
				case "5":
					prefix = "yearly"
					since = now.AddDate(-1, 0, 0)
				}
				m.Period = prefix
				m.Since = since
				m.Done = true
				return m, tea.Quit
			default:
			}
			return m, nil
		}

		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "g":
			m.AwaitLength = true
			m.ReportGenerated = false
			return m, nil
		case "b":
			m.Period = "badges"
			m.Since = time.Now()
			m.Done = true
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m Model) View() string {
	reportMessage := ""
	if m.ReportGenerated {
		reportMessage = "\n📁 Report exported successfully"
	}

	lengthPrompt := ""
	if m.AwaitLength {
		lengthPrompt = `
Choose report period:
1) Daily
2) Weekly
3) Monthly
4) Last 6 months
5) Yearly
`
	}

	return fmt.Sprintf(`
📊 GitHub Developer Metrics
----------------------------
%s%s%s

Press g to generate a report.
Press b to view badges.
Press q to quit.
`, m.Summary, lengthPrompt, reportMessage)
}

func Run(summary string) {
	p := tea.NewProgram(Model{Summary: summary})
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

// RunWithTokenWithSummary initializes the TUI model with both the token and prebuilt summary.
func RunWithTokenWithSummary(token string, summary string) (string, time.Time, bool) {
	model := Model{
		Token:   token,
		Summary: summary,
	}
	p := tea.NewProgram(model)
	finalModel, err := p.StartReturningModel()
	if err != nil {
		log.Fatal(err)
	}
	m := finalModel.(Model)
	return m.Period, m.Since, m.Done
}
