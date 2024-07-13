package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
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
var MUSTARD_COLOR = color.RGBA{0xff, 0xdb, 0x58, 0xff}
var DARK_MUSTARD_COLOR = color.RGBA{0xa8, 0x89, 0x05, 0xff}

type GameSettings struct {
	CellSize       float64
	GridWidth      float64
	GridHeight     float64
	FoodSpawnCount int
}

type GameState struct {
	Score              int
	Speed              float64
	NextMoveOperations *list.List
	IsGameOver         bool
	Debug              bool
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

func drawSnake(imd *imdraw.IMDraw, gameState GameState, gameSettings GameSettings, snake Snake, progress float64) {
	if gameState.Debug {
		imd.Color = colornames.Salmon
		for _, segment := range snake.Segments {
			imd.Push(segment.Add(pixel.V(+1, +1)))
			imd.Push(segment.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
			imd.Rectangle(0)
		}
	}

	progress = math.Min(progress, 1)
	headProgress := progress - math.Sin(progress*2*math.Pi)/10
	tailProgress := progress + math.Sin(progress*2*math.Pi)/5
	imd.Color = snake.Color

	headDirection := snake.getHeadDirection()
	head := snake.Segments[0]
	switch headDirection {
	case Up:
		imd.Push(head.Add(pixel.V(1, 1)))
		imd.Push(head.Add(pixel.V(gameSettings.CellSize-1, 1+headProgress*(gameSettings.CellSize-2))))
	case Down:
		imd.Push(head.Add(pixel.V(1, 1+(1-headProgress)*(gameSettings.CellSize-2))))
		imd.Push(head.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
	case Left:
		imd.Push(head.Add(pixel.V(1+(1-headProgress)*(gameSettings.CellSize-2), 1)))
		imd.Push(head.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
	case Right:
		imd.Push(head.Add(pixel.V(1, 1)))
		imd.Push(head.Add(pixel.V(1+headProgress*(gameSettings.CellSize-2), gameSettings.CellSize-1)))
	}
	imd.Rectangle(0)

	for _, segment := range snake.Segments[1 : len(snake.Segments)-1] {
		imd.Push(segment.Add(pixel.V(1, 1)))
		imd.Push(segment.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
		imd.Rectangle(0)
	}

	tailDirection := snake.getTailDirection()
	tail := snake.Segments[len(snake.Segments)-1]
	switch tailDirection {
	case Up:
		imd.Push(tail.Add(pixel.V(1, 1+tailProgress*(gameSettings.CellSize-2))))
		imd.Push(tail.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
	case Down:
		imd.Push(tail.Add(pixel.V(1, 1)))
		imd.Push(tail.Add(pixel.V(gameSettings.CellSize-1, 1+(1-tailProgress)*(gameSettings.CellSize-2))))
	case Left:
		imd.Push(tail.Add(pixel.V(1, 1)))
		imd.Push(tail.Add(pixel.V(1+(1-tailProgress)*(gameSettings.CellSize-2), gameSettings.CellSize-1)))
	case Right:
		imd.Push(tail.Add(pixel.V(1+tailProgress*(gameSettings.CellSize-2), 1)))
		imd.Push(tail.Add(pixel.V(gameSettings.CellSize-1, gameSettings.CellSize-1)))
	}
	imd.Rectangle(0)
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
			for _, segment := range snake.Segments[1:] {
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
	lastDirection := snake.getHeadDirection()
	if gameState.NextMoveOperations.Len() > 0 {
		lastDirection = gameState.NextMoveOperations.Back().Value.(Direction)
	}

	setDirection := None
	if (win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyW)) && lastDirection != Down {
		setDirection = Up
	} else if (win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyS)) && lastDirection != Up {
		setDirection = Down
	} else if (win.JustPressed(pixelgl.KeyLeft) || win.JustPressed(pixelgl.KeyA)) && lastDirection != Right {
		setDirection = Left
	} else if (win.JustPressed(pixelgl.KeyRight) || win.JustPressed(pixelgl.KeyD)) && lastDirection != Left {
		setDirection = Right
	}

	if setDirection != None && setDirection != lastDirection {
		gameState.NextMoveOperations.PushBack(setDirection)
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
		Score:              0,
		Speed:              200.0,
		NextMoveOperations: list.New(),
		IsGameOver:         false,
		Debug:              false,
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	basicTxt := text.New(pixel.V(16, 500), atlas)
	basicTxt.Color = colornames.Red
	fmt.Fprintln(basicTxt, "GAME OVER.")

	scoreText := text.New(pixel.V(16, 8), atlas)
	scoreText.Color = DARK_MUSTARD_COLOR

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

	for !win.Closed() {
		handleKeyboardEvents(win, &gameState, snake)

		millisFromUpdate := time.Now().UnixMilli() - lastUpdate

		imd := imdraw.New(nil)
		imd.Color = colornames.Blueviolet
		drawGrid(imd, gameSettings)
		drawFood(imd, gameSettings, snakeFood)
		drawSnake(imd, gameState, gameSettings, snake, float64(millisFromUpdate)/gameState.Speed)

		win.Clear(NOKIA_BACKGROUND_COLOR)
		imd.Draw(win)

		if gameState.IsGameOver {
			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 5))
		}

		scoreText.Clear()
		fmt.Fprintln(scoreText, "SCORE", gameState.Score)
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 3))

		if millisFromUpdate > int64(gameState.Speed) && !gameState.IsGameOver {

			if gameState.NextMoveOperations.Len() > 0 {
				snake.setSnakeDirection(gameState.NextMoveOperations.Front().Value.(Direction))
				gameState.NextMoveOperations.Remove(gameState.NextMoveOperations.Front())
			}

			snake.snakeMove(gameSettings)

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

			updateConsumeFoodOnIntersect(&gameSettings, &snake, &gameState, &snakeFood)

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
