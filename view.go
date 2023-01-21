package main

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	rockChar     = "🪨"
	paperChar    = "📄"
	scissorsChar = string('\u2702') + " " // solves emoji single/double width issue.
	deadChar     = "😵"
)

func generateArenaView(g *game) string {
	board := make([][]string, g.maxY+1)

	for idx := range board {
		row := make([]string, g.maxX+1)
		for idx2 := range row {
			row[idx2] = "  "
		}

		board[idx] = row
	}

	for loc, player := range g.players {
		if loc.x > g.maxX || loc.y > g.maxY {
			continue
		}

		var char string

		switch player.kind {
		case rock:
			char = rockChar
		case paper:
			char = paperChar
		case scissors:
			char = scissorsChar
		}

		board[loc.y][loc.x] = char
	}

	arena := ""

	for _, l := range board {
		line := ""
		for _, q := range l {
			line += q
		}

		arena += line + "\n"
	}

	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).
		Render(arena)
}

func generateScoreboardView(g *game, availableWidth int) string {
	title := lipgloss.NewStyle().
		Underline(true).
		Background(lipgloss.Color("#CC33BB")).
		Render("SCOREBOARD")

	rocks := strings.Repeat(rockChar, g.playerCountOfKind(rock))
	if len(rocks) == 0 {
		rocks = deadChar
	}

	rockLine := "ROCK:     " + rocks

	papers := strings.Repeat(paperChar, g.playerCountOfKind(paper))
	if len(papers) == 0 {
		papers = deadChar
	}

	paperLine := "PAPER:    " + papers

	scissors := strings.Repeat(scissorsChar, g.playerCountOfKind(scissors))

	if len(scissors) == 0 {
		scissors = deadChar
	}

	scissorsLine := "SCISSORS: " + scissors

	scores := []string{rockLine, paperLine, scissorsLine}
	sort.Slice(scores, func(i, j int) bool {
		return len(scores[i]) > len(scores[j])
	})

	board := title + "\n"
	for _, v := range scores {
		board += "\n" + v
	}

	scoreBoard := lipgloss.NewStyle().
		PaddingTop(0).
		Width(availableWidth-2).
		PaddingLeft(0).Border(lipgloss.RoundedBorder(), true, true, true, true).Render(board)

	return scoreBoard
}

func generateTitleView() string {
	title := `
 ___ ___  ___     _   ___ ___ _  _   _   
| _ \ _ \/ __|   /_\ | _ \ __| \| | /_\  
|   /  _/\__ \  / _ \|   / _|| .  |/ _ \ 
|_|_\_|  |___/ /_/ \_\_|_\___|_|\_/_/ \_\
    `
	text := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#CC33BB")).
		Border(lipgloss.RoundedBorder()).
		Render(title)

	return text
}

func showHelp() string {
	title := makePink("HELP / ABOUT")

	help := `
RPS Arena is a zero player game. The game plays itself. You watch.

[keys]
n         start a new game
p/space   pause or unpause
h         toggle help
q         quit
left      slow down game
right     speed up game

[looks broken]
Might be your terminal emulator. Or your font. Or both. Sorry.

[why did you bring this into existence?]
I don't know why I made this. It was not the best use of my time.


                                                              Tom
`

	help = lipgloss.NewStyle().
		Render(help)

	joined := lipgloss.JoinVertical(lipgloss.Left, title, help)

	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).Padding(2).Render(joined)
}

func generateFooterView(g *game, speed string) string {
	if g.isOver() {
		player := getSomePlayer(g)
		winner := makePink(strings.Title(player.kind))

		return " " + winner + " wins. Press \"n\" to play again."
	}

	email := makePink("tom@tomontheinternet.com")

	return " [ speed = " + speed + " ]   RPS Arena by Tom. Press \"h\" for help. Contact: " + email
}

func makePink(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#CC33BB")).Render(s)
}
