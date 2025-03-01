package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// Game constants
const (
	width       = 40
	height      = 20
	initialSize = 3
)

// Cell symbols
const (
	symbolBorderHorizontal  = '━'
	symbolBorderVertical    = '┃'
	symbolBorderTopLeft     = '┏'
	symbolBorderTopRight    = '┓'
	symbolBorderBottomLeft  = '┗'
	symbolBorderBottomRight = '┛'
	symbolSnakeHead         = '◉'
	symbolSnakeBody         = '■'
	symbolFood              = '★'
)

// Direction represents the snake's movement direction
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

// Point represents a position on the grid
type Point struct {
	X, Y int
}

// Game represents the state of the game
type Game struct {
	snake     []Point
	food      Point
	direction Direction
	score     int
	gameOver  bool
}

// Initialize a new game
func NewGame() *Game {
	g := &Game{
		snake:     make([]Point, initialSize),
		direction: Right,
	}

	// Initialize snake in the middle of the board
	for i := 0; i < initialSize; i++ {
		g.snake[i] = Point{
			X: width/2 - i,
			Y: height / 2,
		}
	}

	// Place initial food
	g.PlaceFood()

	return g
}

// Place food at a random location not occupied by the snake
func (g *Game) PlaceFood() {
	for {
		g.food = Point{
			X: rand.Intn(width),
			Y: rand.Intn(height),
		}

		// Check if food is on snake
		collision := false
		for _, p := range g.snake {
			if p.X == g.food.X && p.Y == g.food.Y {
				collision = true
				break
			}
		}

		if !collision {
			break
		}
	}
}

// Update game state
func (g *Game) Update() {
	if g.gameOver {
		return
	}

	// Calculate new head position
	head := g.snake[0]
	var newHead Point

	switch g.direction {
	case Up:
		newHead = Point{X: head.X, Y: head.Y - 1}
	case Right:
		newHead = Point{X: head.X + 1, Y: head.Y}
	case Down:
		newHead = Point{X: head.X, Y: head.Y + 1}
	case Left:
		newHead = Point{X: head.X - 1, Y: head.Y}
	}

	// Implement wraparound for walls
	if newHead.X < 0 {
		newHead.X = width - 1
	} else if newHead.X >= width {
		newHead.X = 0
	}

	if newHead.Y < 0 {
		newHead.Y = height - 1
	} else if newHead.Y >= height {
		newHead.Y = 0
	}

	// Check self collision
	for _, p := range g.snake {
		if p.X == newHead.X && p.Y == newHead.Y {
			g.gameOver = true
			return
		}
	}

	// Add new head to snake
	g.snake = append([]Point{newHead}, g.snake...)

	// Check food collision
	if newHead.X == g.food.X && newHead.Y == g.food.Y {
		g.score++
		g.PlaceFood()
	} else {
		// Remove tail if no food was eaten
		g.snake = g.snake[:len(g.snake)-1]
	}
}

// Draw the game
func (g *Game) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw border
	for i := 0; i < width+2; i++ {
		termbox.SetCell(i, 0, symbolBorderHorizontal, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(i, height+1, symbolBorderHorizontal, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i := 0; i < height+2; i++ {
		termbox.SetCell(0, i, symbolBorderVertical, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width+1, i, symbolBorderVertical, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.SetCell(0, 0, symbolBorderTopLeft, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(width+1, 0, symbolBorderTopRight, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, height+1, symbolBorderBottomLeft, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(width+1, height+1, symbolBorderBottomRight, termbox.ColorWhite, termbox.ColorDefault)

	// Draw score
	scoreStr := []rune(fmt.Sprintf("Score: %d", g.score))
	for i, ch := range scoreStr {
		termbox.SetCell(2+i, height+3, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw snake
	for i, p := range g.snake {
		symbol := symbolSnakeBody
		if i == 0 {
			// First segment is the head
			symbol = symbolSnakeHead
		}
		termbox.SetCell(p.X+1, p.Y+1, symbol, termbox.ColorGreen, termbox.ColorDefault)
	}

	// Draw food
	termbox.SetCell(g.food.X+1, g.food.Y+1, symbolFood, termbox.ColorRed, termbox.ColorDefault)

	// Game over message
	if g.gameOver {
		gameOverMsg := "Game Over! Press 'q' to quit or 'r' to restart."
		for i, ch := range []rune(gameOverMsg) {
			termbox.SetCell(width/2-len(gameOverMsg)/2+i, height/2, ch, termbox.ColorRed, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := NewGame()
	eventQueue := make(chan termbox.Event)

	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowUp:
					if game.direction != Down {
						game.direction = Up
					}
				case termbox.KeyArrowRight:
					if game.direction != Left {
						game.direction = Right
					}
				case termbox.KeyArrowDown:
					if game.direction != Up {
						game.direction = Down
					}
				case termbox.KeyArrowLeft:
					if game.direction != Right {
						game.direction = Left
					}
				case termbox.KeyEsc:
					return
				}

				if ev.Ch == 'q' {
					return
				} else if ev.Ch == 'r' && game.gameOver {
					game = NewGame()
				}
			}
		case <-ticker.C:
			game.Update()
			game.Draw()
		}
	}
}
