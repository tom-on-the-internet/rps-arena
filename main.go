package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	game     *game
	width    int
	height   int
	speed    string
	showHelp bool
}

const (
	turtle    = "turtle"
	slow      = "slow"
	normal    = "normal"
	fast      = "fast"
	lightning = "lightning"
)

const (
	rock     = "rock"
	paper    = "paper"
	scissors = "scissors"
)

type TickMsg time.Time

func doTick(speed string) tea.Cmd {
	var ms time.Duration

	switch speed {
	case lightning:
		ms = 20
	case fast:
		ms = 75
	case normal:
		ms = 125
	case slow:
		ms = 250
	case turtle:
		ms = 1000
	}

	return tea.Tick(time.Millisecond*ms, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return doTick(m.speed)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "left":
			switch m.speed {
			case slow:
				m.speed = turtle
			case normal:
				m.speed = slow
			case fast:
				m.speed = normal
			case lightning:
				m.speed = fast
			}

			return m, nil
		case "right":
			switch m.speed {
			case fast:
				m.speed = lightning
			case normal:
				m.speed = fast
			case slow:
				m.speed = normal
			case turtle:
				m.speed = slow
			}

			return m, nil
		case "N", "n":
			m.game.initialize(40)

			return m, nil
		case "H", "h":
			m.showHelp = !m.showHelp
			return m, nil
		case "p", " ":
			m.game.paused = !m.game.paused
			if !m.game.paused {
				return m, doTick(m.speed)
			}

			return m, nil
		}
	case TickMsg:
		if m.game.paused {
			return m, nil
		}

		m.game.takeTurn()

		return m, doTick(m.speed)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, tea.ClearScreen
	}

	return m, nil
}

func (m model) View() string {
	if m.showHelp {
		return showHelp()
	}

	if m.width == 0 {
		// initial width has not yet been set
		return ""
	}

	titleView := generateTitleView()
	widthForScoreboard := m.width - lipgloss.Width(titleView)
	scoreboardView := generateScoreboardView(m.game, widthForScoreboard)

	top := lipgloss.JoinHorizontal(lipgloss.Left, titleView, scoreboardView)

	heightForArena := m.height - lipgloss.Height(top) - 6

	m.game.maxX = (m.width / 2) - 2 // handle padding and double wide emojis
	m.game.maxY = heightForArena

	if !m.game.initialized {
		m.game.initialize(40)
	}

	if m.game.maxX < 5 || m.game.maxY < 5 {
		return "view port too small."
	}

	m.game.removeOutOfBoundsPlayers()

	arenaView := generateArenaView(m.game)

	footerView := generateFooterView(m.game, m.speed)

	x := lipgloss.JoinVertical(lipgloss.Left, top, arenaView, footerView)

	return x
}

func main() {
	rand.Seed(time.Now().UnixNano())

	m := model{
		game:  newGame(),
		speed: "normal",
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
