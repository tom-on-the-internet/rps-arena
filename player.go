package main

import "math/rand"

type player struct {
	kind      string
	turnTaken bool
}

// defeats checks if the current player is of a stronger
// type than the other player. ex: paper beats rock
// If there is no other player (nil), then we default
// to true
func (p *player) defeats(otherPlayer *player) bool {
	if otherPlayer == nil {
		return true
	}

	if p.kind == rock && otherPlayer.kind == scissors {
		return true
	}

	if p.kind == scissors && otherPlayer.kind == paper {
		return true
	}

	if p.kind == paper && otherPlayer.kind == rock {
		return true
	}

	return false
}

// newPlayer returns a new player with a random kind
func newPlayer() *player {
	player := player{}

	num := rand.Intn(3)

	switch num {
	case 0:
		player.kind = rock
	case 1:
		player.kind = paper
	case 2:
		player.kind = scissors
	}

	return &player
}
