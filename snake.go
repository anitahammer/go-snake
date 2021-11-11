package main

import (
	"image/color"

	"github.com/faiface/pixel"
)

type Snake struct {
	Segments        []pixel.Vec
	DirectionVector pixel.Vec
	Color           color.Color
}

type Direction int64

const (
	Up    Direction = 0
	Down  Direction = 1
	Left  Direction = 2
	Right Direction = 3
	None  Direction = 4
)

func initSnake(gameSettings GameSettings) Snake {
	return Snake{Segments: make([]pixel.Vec, 1, 10), DirectionVector: pixel.V(0.0, gameSettings.CellSize), Color: NOKIA_FOREGROUND_COLOR}
}

func (snake *Snake) setSnakeDirection(direction Direction, gameSettings GameSettings) {
	switch direction {
	case Up:
		snake.DirectionVector = pixel.V(0, gameSettings.CellSize)
	case Down:
		snake.DirectionVector = pixel.V(0, -gameSettings.CellSize)
	case Left:
		snake.DirectionVector = pixel.V(-gameSettings.CellSize, 0)
	case Right:
		snake.DirectionVector = pixel.V(gameSettings.CellSize, 0)
	}
}

func (snake Snake) getSnakeDirection() Direction {
	if snake.DirectionVector.X < 0 {
		return Left
	} else if snake.DirectionVector.X > 0 {
		return Right
	} else if snake.DirectionVector.Y < 0 {
		return Down
	} else {
		return Up
	}
}

func (snake *Snake) growSnake(gameSettings GameSettings) {
	tail := snake.Segments[len(snake.Segments)-1]
	(*snake).Segments = append(snake.Segments, tail.Add(snake.DirectionVector.Scaled(-1)))
}

func (snake *Snake) snakeMove() {
	var newSegments = make([]pixel.Vec, len(snake.Segments))
	newSegments[0] = snake.Segments[0].Add(snake.DirectionVector)

	for i := 1; i < len(snake.Segments); i++ {
		newSegments[i] = snake.Segments[i-1]
	}

	snake.Segments = newSegments
}

func (snake Snake) snakeSelfIntersect() bool {
	for i := 1; i < len(snake.Segments); i++ {
		if snake.Segments[i] == snake.Segments[0] {
			return true
		}
	}
	return false
}
