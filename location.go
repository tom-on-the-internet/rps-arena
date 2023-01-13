package main

import "strconv"

type location struct {
	x int
	y int
}

func (l location) String() string {
	return strconv.Itoa(l.x) + "," + strconv.Itoa(l.y)
}

func (l location) relativeDistance(otherLocation location) int {
	xDist := l.x - otherLocation.x
	if xDist < 0 {
		xDist *= -1
	}

	yDist := l.y - otherLocation.y
	if yDist < 0 {
		yDist *= -1
	}

	if xDist == 0 {
		xDist = 1
	}

	if yDist == 0 {
		yDist = 1
	}

	return xDist * yDist
}
