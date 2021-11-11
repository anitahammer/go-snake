package main

import (
	"testing"

	"github.com/faiface/pixel"
)

func TestInitSnake(t *testing.T) {
	specifiedGameSettings := GameSettings{
		CellSize:   20.0,
		GridWidth:  40.0,
		GridHeight: 40.0,
	}

	sut := initSnake(specifiedGameSettings)

	if len(sut.Segments) != 1 {
		t.Errorf("Expected a snake of length 1, actually got %d", len(sut.Segments))
	}

	if !sut.DirectionVector.Eq(pixel.V(0.0, specifiedGameSettings.CellSize)) {
		t.Errorf("Incorrect default movement vector on snake object %f - %f", sut.DirectionVector.Y, pixel.V(0.0, specifiedGameSettings.CellSize).Y)
	}

	if sut.Color != NOKIA_FOREGROUND_COLOR {
		t.Errorf("Unexpected foreground color on Snake object")
	}
}

func TestMoveSnake(t *testing.T) {
	specifiedGameSettings := GameSettings{
		CellSize:   20.0,
		GridWidth:  40.0,
		GridHeight: 40.0,
	}

	sut := initSnake(specifiedGameSettings)
	sut.Segments[0] = pixel.V(2, 2)
	sut.Segments = append(sut.Segments, pixel.V(4, 4))
	sut.Segments = append(sut.Segments, pixel.V(6, 6))

	sut.snakeMove()

	if !sut.Segments[0].Eq(pixel.V(2, 2).Add(sut.DirectionVector)) {
		t.Errorf("Incorrectly computed new snake (head) position after move")
	}

	if !sut.Segments[1].Eq(pixel.V(2, 2)) {
		t.Errorf("Incorrectly computed new snake (tail propagate) position after move")
	}

	if !sut.Segments[2].Eq(pixel.V(4, 4)) {
		t.Errorf("Incorrectly computed new snake (tail propagate) position after move")
	}
}

func TestGetSnakeDirection(t *testing.T) {
	specifiedGameSettings := GameSettings{
		CellSize:   20.0,
		GridWidth:  40.0,
		GridHeight: 40.0,
	}

	sut := initSnake(specifiedGameSettings)
	sut.Segments[0] = pixel.V(2, 2)
	sut.setSnakeDirection(Up, specifiedGameSettings)

	if sut.getSnakeDirection() != Up {
		t.Errorf("Incorrect snake direction")
	}

}
