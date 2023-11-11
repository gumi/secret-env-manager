package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

type checkboxListModel struct {
	Secrets  model.Secrets
	Selected []bool
	Index    int
}

func CheckBoxList(secrets model.Secrets) checkboxListModel {
	initialModel := checkboxListModel{
		Secrets:  secrets,
		Selected: make([]bool, len(secrets.Secrets)),
		Index:    0,
	}

	p := tea.NewProgram(initialModel)
	model, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return model.(checkboxListModel)
}

func (m checkboxListModel) Init() tea.Cmd {
	return nil
}

func (m checkboxListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			return m, tea.Quit
		}
		if msg.String() == "j" || msg.String() == "down" {
			m.Index = (m.Index + 1) % len(m.Secrets.Secrets)
			return m, nil
		}
		if msg.String() == "k" || msg.String() == "up" {
			m.Index = (m.Index - 1 + len(m.Secrets.Secrets)) % len(m.Secrets.Secrets)
			return m, nil
		}
		if msg.String() == " " {
			m.Selected[m.Index] = !m.Selected[m.Index]
			return m, nil
		}
	}

	return m, nil
}

func (m checkboxListModel) View() string {
	var items string

	for i, secret := range m.Secrets.Secrets {
		text := ""
		if i == m.Index {
			text += "▶ "
		} else {
			text += "  "
		}

		if m.Selected[i] {
			text += "[x] "
		} else {
			text += "[ ] "
		}

		items += fmt.Sprintf("%s%s (%s)\n", text, secret.Name, secret.CreatedAt)
	}
	return fmt.Sprintf("Select secrets:\nsecret name (created at)\n%s", items)
}
