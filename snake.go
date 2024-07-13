package main

import (
	"container/list"
	"image/color"

	"github.com/faiface/pixel"
)

type Snake struct {
	Segments   []pixel.Vec
	Directions *list.List
	Color      color.Color
}

type Direction int64

const (
	Up Direction = iota
	Down
	Left
	Right
	None
)

func initSnake(gameSettings GameSettings) Snake {
	snake := Snake{
		Segments:   make([]pixel.Vec, 1, 10),
		Directions: list.New(),
		Color:      NOKIA_FOREGROUND_COLOR,
	}
	snake.Directions.PushFront(Up)
	snake.growSnake(gameSettings)
	return snake
}

func (snake *Snake) setSnakeDirection(direction Direction) {
	snake.Directions.Front().Value = direction
}

func GetDirectionVector(direction Direction, gameSettings GameSettings) pixel.Vec {
	switch direction {
	case Up:
		return pixel.V(0, gameSettings.CellSize)
	case Down:
		return pixel.V(0, -gameSettings.CellSize)
	case Left:
		return pixel.V(-gameSettings.CellSize, 0)
	case Right:
		return pixel.V(gameSettings.CellSize, 0)
	default:
		return pixel.V(0, 0)
	}
}

func (snake Snake) getHeadDirection() Direction {
	return snake.Directions.Front().Value.(Direction)
}

func (snake Snake) getTailDirection() Direction {
	return snake.Directions.Back().Value.(Direction)
}

func (snake *Snake) growSnake(gameSettings GameSettings) {
	tail := snake.Segments[len(snake.Segments)-1]
	snake.Segments = append(snake.Segments, tail)
	snake.Directions.PushBack(snake.getTailDirection())
}

func (snake *Snake) snakeMove(gameSettings GameSettings) {
	var newSegments = make([]pixel.Vec, len(snake.Segments))
	newSegments[0] = snake.Segments[0].Add(GetDirectionVector(snake.getHeadDirection(), gameSettings))

	for i := 1; i < len(snake.Segments); i++ {
		newSegments[i] = snake.Segments[i-1]
	}

	snake.Segments = newSegments

	snake.Directions.PushFront(snake.getHeadDirection())
	snake.Directions.Remove(snake.Directions.Back())
}

func (snake Snake) snakeSelfIntersect() bool {
	for _, segment := range snake.Segments[1 : len(snake.Segments)-1] {
		if segment == snake.Segments[0] {
			return true
		}
	}
	return false
}
