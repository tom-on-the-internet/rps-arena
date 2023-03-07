package main

import (
	"math/rand"
)

type game struct {
	maxX        int
	maxY        int
	players     map[location]*player
	paused      bool
	initialized bool
}

func (g *game) playerCountOfKind(kind string) int {
	count := 0

	for _, p := range g.players {
		if p.kind == kind {
			count++
		}
	}

	return count
}

func (g *game) takeTurn() {
	g.convertPlayers()
	g.movePlayers()

	// now that the turn is over
	// reset the turns
	for _, p := range g.players {
		p.turnTaken = false
	}
}

func (g *game) convertPlayers() {
	// if a player is beside a weaker enemy player
	// force the enemy to become their kind
	for location, player := range g.players {
		if player.turnTaken {
			continue
		}

		weakerPlayers := g.getWeakerPlayersNearby(location)

		if len(weakerPlayers) == 0 {
			continue
		}

		// There is at least one weaker player
		// so we take our turn.
		player.turnTaken = true

		for _, weakerPlayer := range weakerPlayers {
			weakerPlayer.kind = player.kind
			weakerPlayer.turnTaken = true
		}
	}
}

func (g *game) movePlayers() {
	for locationOfPlayerWithTurnNotTaken(g) != nil {
		playerLocation := *locationOfPlayerWithTurnNotTaken(g)
		player := g.players[playerLocation]
		player.turnTaken = true

		goalLocation := getGoalLocation(g, playerLocation)

		// can move to goal location
		existingPlayer := g.players[goalLocation]
		if existingPlayer == nil {
			delete(g.players, playerLocation)
			g.players[goalLocation] = player

			continue
		}

		// try horizontal
		if playerLocation.x != goalLocation.x {
			secondaryLocation := location{x: goalLocation.x, y: playerLocation.y}

			// can move to goal location
			existingPlayer := g.players[secondaryLocation]
			if existingPlayer == nil {
				delete(g.players, playerLocation)
				g.players[secondaryLocation] = player

				continue
			}
		}

		// try vertical
		if playerLocation.y != goalLocation.y {
			secondaryLocation := location{x: playerLocation.x, y: goalLocation.y}

			// can move to goal location
			existingPlayer := g.players[secondaryLocation]
			if existingPlayer == nil {
				delete(g.players, playerLocation)
				g.players[secondaryLocation] = player

				continue
			}
		}
	}
}

func (g *game) getClosestWeakerEnemyPlayerLocation(loc location) *location {
	var (
		closestLocation location
		minDistance     int
		firstEnemyFound bool
	)

	thisPlayer := g.players[loc]

	for enemyLocation, p := range g.players {
		if p.kind == thisPlayer.kind {
			continue
		}

		if !thisPlayer.defeats(p) {
			continue
		}

		if !firstEnemyFound {
			closestLocation = enemyLocation
			minDistance = loc.relativeDistance(closestLocation)
			firstEnemyFound = true

			continue
		}

		relativeDistance := loc.relativeDistance(enemyLocation)

		if relativeDistance < minDistance {
			minDistance = relativeDistance
			closestLocation = enemyLocation
		}
	}

	if !firstEnemyFound {
		return nil
	}

	return &closestLocation
}

func (g *game) getClosestEnemyPlayerLocation(loc location) *location {
	var (
		closestLocation location
		minDistance     int
		firstEnemyFound bool
	)

	thisPlayer := g.players[loc]

	for enemyLocation, p := range g.players {
		if p.kind == thisPlayer.kind {
			continue
		}

		if !firstEnemyFound {
			closestLocation = enemyLocation
			minDistance = loc.relativeDistance(closestLocation)
			firstEnemyFound = true

			continue
		}

		relativeDistance := loc.relativeDistance(enemyLocation)

		if relativeDistance < minDistance {
			minDistance = relativeDistance
			closestLocation = enemyLocation
		}
	}

	if !firstEnemyFound {
		return nil
	}

	return &closestLocation
}

func (g *game) getWeakerPlayersNearby(playerLocation location) []*player {
	locations := g.surroundingLocations(playerLocation)

	players := []*player{}
	thisPlayer := g.players[playerLocation]

	for _, l := range locations {
		p, ok := g.players[l]
		if ok && thisPlayer.defeats(p) {
			players = append(players, p)
		}
	}

	return players
}

// surroundingLocations return locations surrounding a location
// that are within the game.
func (g *game) surroundingLocations(l location) []location {
	locations := []location{}

	// top left
	if l.x > 0 && l.y > 0 {
		locations = append(locations, location{x: l.x - 1, y: l.y - 1})
	}
	// top middle
	if l.y > 0 {
		locations = append(locations, location{x: l.x, y: l.y - 1})
	}
	// top right
	if l.x < g.maxX && l.y > 0 {
		locations = append(locations, location{x: l.x + 1, y: l.y - 1})
	}
	// left
	if l.x > 0 {
		locations = append(locations, location{x: l.x - 1, y: l.y})
	}
	// right
	if l.x < g.maxX {
		locations = append(locations, location{x: l.x + 1, y: l.y})
	}
	// bottom left
	if l.x > 0 && l.y < g.maxY {
		locations = append(locations, location{x: l.x - 1, y: l.y + 1})
	}
	// bottom middle
	if l.y < g.maxY {
		locations = append(locations, location{x: l.x, y: l.y + 1})
	}
	// bottom right
	if l.x < g.maxX && l.y < g.maxY {
		locations = append(locations, location{x: l.x + 1, y: l.y + 1})
	}

	return locations
}

func (g *game) isOver() bool {
	// no players somehow, so game is over
	if len(g.players) == 0 {
		return true
	}

	firstKind := getSomePlayer(g).kind

	// if any player is of a different kind
	// than the first kind, the game is not over
	for _, player := range g.players {
		if player.kind != firstKind {
			return false
		}
	}

	return true
}

func (g *game) randomLocation() location {
	x := rand.Intn(g.maxX + 1)
	y := rand.Intn(g.maxY + 1)

	return location{x: x, y: y}
}

func (g *game) randomEmptyLocation() location {
	for {
		l := g.randomLocation()

		_, isSet := g.players[l]
		if !isSet {
			return l
		}
	}
}

func (g *game) removeOutOfBoundsPlayers() {
	for l := range g.players {
		if l.x > g.maxX || l.y > g.maxY {
			delete(g.players, l)
		}
	}
}

func getSomePlayer(g *game) *player {
	for _, p := range g.players {
		return p
	}

	return nil
}

func locationOfPlayerWithTurnNotTaken(g *game) *location {
	for l, p := range g.players {
		if !p.turnTaken {
			return &l
		}
	}

	return nil
}

func newGame() *game {
	return &game{}
}

func (g *game) initialize(playerCount int) {
	g.players = map[location]*player{}

	for i := 0; i < playerCount; i++ {
		player := newPlayer()
		location := g.randomEmptyLocation()
		g.players[location] = player
	}

	g.initialized = true
}

// getGoalLocation returns a location adjacent to the player
// that I believe is a pretty good place to move to.
// This function is gross.
func getGoalLocation(g *game, playerLocation location) location {
	player := g.players[playerLocation]
	enemyLocation := g.getClosestEnemyPlayerLocation(playerLocation)
	closestWeakerEnemyLocation := g.getClosestWeakerEnemyPlayerLocation(playerLocation)
	goalLocation := playerLocation

	// party dance at end
	if enemyLocation == nil {
		return getRandomGoalLocation(g, playerLocation)
	}

	// give up if no weaker enemies
	if closestWeakerEnemyLocation == nil {
		return playerLocation
	}

	// if enemyLocation far enough, enemyLocation is weakest location
	// basically, don't be a coward. if the enemy is far away, act like
	// he's not here. get after it.
	if playerLocation.relativeDistance(*enemyLocation) > 10 {
		enemyLocation = closestWeakerEnemyLocation
	}

	enemy := g.players[*enemyLocation]

	if player.defeats(enemy) {
		// move toward location or as close as possible
		if playerLocation.x < enemyLocation.x {
			goalLocation.x++
		} else if playerLocation.x > enemyLocation.x {
			goalLocation.x--
		}

		if playerLocation.y < enemyLocation.y {
			goalLocation.y++
		} else if playerLocation.y > enemyLocation.y {
			goalLocation.y--
		}

		return goalLocation
	}

	// move away from location
	// but be careful of going out of bounds
	if playerLocation.x < enemyLocation.x && playerLocation.x != 0 {
		goalLocation.x--
	} else if playerLocation.x > enemyLocation.x && playerLocation.x != g.maxX {
		goalLocation.x++
	}

	if playerLocation.y < enemyLocation.y && playerLocation.y != 0 {
		goalLocation.y--
	} else if playerLocation.y > enemyLocation.y && playerLocation.y != g.maxY {
		goalLocation.y++
	}

	// because we are running away, there's a chance we are being
	// chased in a straight line. if that's the case, _maybe_ move off the line.
	if goalLocation.x == enemyLocation.x {
		rndNum := rand.Intn(9)
		if rndNum == 0 && goalLocation.x != 0 {
			goalLocation.x--
		}

		if rndNum == 1 && goalLocation.x != g.maxX {
			goalLocation.x++
		}
	}

	if goalLocation.y == enemyLocation.y {
		rndNum := rand.Intn(9)
		if rndNum == 0 && goalLocation.y != 0 {
			goalLocation.y--
		}

		if rndNum == 1 && goalLocation.y != g.maxY {
			goalLocation.y++
		}
	}

	if goalLocation != playerLocation {
		return goalLocation
	}

	// if the goal location is the same as the player location at this point
	// it means the player is stuck and being run down.
	// shuck and give.
	// but only accept this location if it's not actually closer to the enemy
	randomLocation := getRandomGoalLocation(g, playerLocation)

	if randomLocation.relativeDistance(
		*enemyLocation,
	) >= playerLocation.relativeDistance(
		*enemyLocation,
	) {
		return randomLocation
	}

	return playerLocation
}

func getRandomGoalLocation(g *game, playerLocation location) location {
	goalLocation := playerLocation

	x := rand.Intn(g.maxX + 1)
	if x > g.maxX {
		x = g.maxX
	}

	y := rand.Intn(g.maxY + 1)
	if y > g.maxY {
		y = g.maxY
	}

	randomLocation := location{x: x, y: y}

	if playerLocation.x < randomLocation.x {
		goalLocation.x++
	} else if playerLocation.x > randomLocation.x {
		goalLocation.x--
	}

	if playerLocation.y < randomLocation.y {
		goalLocation.y++
	} else if playerLocation.y > randomLocation.y {
		goalLocation.y--
	}

	return goalLocation
}
