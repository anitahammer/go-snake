package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var NOKIA_FOREGROUND_COLOR = color.RGBA{0x43, 0x52, 0x3d, 0xFF}
var NOKIA_BACKGROUND_COLOR = color.RGBA{0xc7, 0xf0, 0xd8, 0xff}

type GameSettings struct {
	CellSize       float64
	GridWidth      float64
	GridHeight     float64
	FoodSpawnCount int
}

type GameState struct {
	Score             int
	Speed             float64
	NextMoveOperation Direction
	IsGameOver        bool
}

func drawGrid(imd *imdraw.IMDraw, gameSettings GameSettings) {
	for i := 0; i < int(gameSettings.GridWidth); i++ {
		for j := 0; j < int(gameSettings.GridHeight); j++ {
			imd.Push(pixel.V(float64(i)*gameSettings.CellSize, float64(j)*gameSettings.CellSize))
			imd.Color = pixel.RGB(255, 255, 255)
			imd.Push(pixel.V(float64(i)*gameSettings.CellSize+gameSettings.CellSize, float64(j)*gameSettings.CellSize+gameSettings.CellSize))
			imd.Color = pixel.RGB(255, 0, 255)
			imd.Rectangle(1)
		}
	}
}

func drawFood(imd *imdraw.IMDraw, gameSettings GameSettings, snakeFood []SnakeFood) {
	for _, food := range snakeFood {
		if food.Alive {
			imd.Color = colornames.Hotpink
			imd.Push(food.Position.Add(pixel.V(+1, +1)))
			imd.Color = colornames.Hotpink
			imd.Push(food.Position.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
			imd.Color = colornames.Hotpink
			imd.Rectangle(0)
		}
	}
}

func drawSnake(imd *imdraw.IMDraw, gameSettings GameSettings, snake Snake) {
	for _, segment := range snake.Segments {
		imd.Color = snake.Color
		imd.Push(segment.Add(pixel.V(+1, +1)))
		imd.Color = snake.Color
		imd.Push(segment.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
		imd.Color = snake.Color
		imd.Rectangle(0)
	}
}

type SnakeFood struct {
	Position pixel.Vec
	Alive    bool
	Color    color.Color
	Value    int
}

func updateConsumeFoodOnIntersect(gameSettings *GameSettings, snake *Snake, gameState *GameState, snakeFood *[]SnakeFood) {
	for i, food := range *snakeFood {
		if food.Alive {
			for _, segment := range snake.Segments {
				if segment.Eq(food.Position) {
					(*snake).growSnake(*gameSettings)
					(*gameState).Score += food.Value
					(*snakeFood)[i].Position = generateRandomPosition(*gameSettings)
					gameState.Speed -= 1
				}
			}
		}
	}
}

func generateRandomPosition(gameSettings GameSettings) pixel.Vec {
	return pixel.V(float64(rand.Intn(int(gameSettings.GridWidth)))*gameSettings.CellSize, float64(rand.Intn(int(gameSettings.GridHeight)))*gameSettings.CellSize)
}

func handleKeyboardEvents(win *pixelgl.Window, gameState *GameState, snake Snake) {
	if (win.Pressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW)) && snake.getSnakeDirection() != Down {
		gameState.NextMoveOperation = Up
	} else if (win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS)) && snake.getSnakeDirection() != Up {
		gameState.NextMoveOperation = Down
	} else if (win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA)) && snake.getSnakeDirection() != Right {
		gameState.NextMoveOperation = Left
	} else if (win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD)) && snake.getSnakeDirection() != Left {
		gameState.NextMoveOperation = Right
	}
}

func run() {
	gameSettings := GameSettings{
		CellSize:       40.0,
		GridWidth:      20.0,
		GridHeight:     20.0,
		FoodSpawnCount: 2,
	}
	cfg := pixelgl.WindowConfig{
		Title:  "go-snake",
		Bounds: pixel.R(0, 0, gameSettings.GridWidth*gameSettings.CellSize, gameSettings.GridHeight*gameSettings.CellSize),
		VSync:  true,
	}

	gameState := GameState{
		Score:             0,
		Speed:             200.0,
		NextMoveOperation: None,
		IsGameOver:        false,
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(16, 500), atlas)
	scoreText := text.New(pixel.V(16, 8), atlas)
	scoreText.Color = colornames.Yellow

	var snake = initSnake(gameSettings)
	var snakeFood = make([]SnakeFood, gameSettings.FoodSpawnCount)

	for i := 0; i < len(snakeFood); i++ {
		snakeFoodItem := SnakeFood{
			Position: generateRandomPosition(gameSettings),
			Alive:    true,
			Value:    1,
		}
		snakeFood[i] = snakeFoodItem
	}
	lastUpdate := time.Now().UnixMilli()
	basicTxt.Color = colornames.Red

	for !win.Closed() {
		handleKeyboardEvents(win, &gameState, snake)

		imd := imdraw.New(nil)
		imd.Color = colornames.Blueviolet
		drawGrid(imd, gameSettings)
		drawFood(imd, gameSettings, snakeFood)
		drawSnake(imd, gameSettings, snake)

		if gameState.IsGameOver {
			fmt.Fprintln(basicTxt, "GAME OVER.")
			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
		}

		scoreText.Clear()
		fmt.Fprintln(scoreText, "SCORE", gameState.Score)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 3))

		win.Clear(NOKIA_BACKGROUND_COLOR)
		imd.Draw(win)

		if time.Now().UnixMilli()-lastUpdate > int64(gameState.Speed) && !gameState.IsGameOver {

			snake.setSnakeDirection(gameState.NextMoveOperation, gameSettings)
			gameState.NextMoveOperation = None

			updateConsumeFoodOnIntersect(&gameSettings, &snake, &gameState, &snakeFood)
			snake.snakeMove()

			if snake.Segments[0].X >= cfg.Bounds.W() {
				snake.Segments[0].X = 0
			}

			if snake.Segments[0].X < 0 {
				snake.Segments[0].X = cfg.Bounds.W() - gameSettings.CellSize
			}

			if snake.Segments[0].Y >= cfg.Bounds.H() {
				snake.Segments[0].Y = 0
			}

			if snake.Segments[0].Y < 0 {
				snake.Segments[0].Y = cfg.Bounds.H() - gameSettings.CellSize
			}

			if snake.snakeSelfIntersect() {
				gameState.IsGameOver = true
				snake.Color = colornames.Red
			} else {
				snake.Color = NOKIA_FOREGROUND_COLOR
			}

			lastUpdate = time.Now().UnixMilli()
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
