package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type textFieldModel struct {
	Value string
}

type msgUpdate string

func (m textFieldModel) Init() tea.Cmd {
	return nil
}

func (m textFieldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			return m, tea.Quit
		}
		if msg.String() == "backspace" {
			if len(m.Value) > 0 {
				m.Value = m.Value[:len(m.Value)-1]
			}
			return m, nil
		}
		if len(msg.String()) == 1 {
			m.Value += msg.String()
			return m, nil
		}
	}

	return m, nil
}

func (m textFieldModel) View() string {
	return fmt.Sprintf("Please Select GCP Project : %s ", m.Value)
}

func TextField() string {
	initialModel := textFieldModel{}

	p := tea.NewProgram(initialModel)
	model, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return model.(textFieldModel).Value
}
