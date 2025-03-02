package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// Game constants
const (
	width        = 40
	height       = 15
	initialSize  = 3
	aspectRatio  = 1.8
	baseSpeed    = 100
	sidebarWidth = 20 // Width of the sidebar

	// Food timer constants
	minFoodTime     = 50  // Minimum ticks food stays on screen
	maxFoodTime     = 150 // Maximum ticks food stays on screen
	foodRespawnTime = 20  // Ticks to wait before spawning new food
)

// Food types and values
var (
	foodSymbols = []rune{'🍆', '🍗', '🧀', '🍬'}
	foodValues  = []int{1, 3, 5, 7}
)

// Cell symbols
const (
	symbolBorderHorizontal  = '━'
	symbolBorderVertical    = '┃'
	symbolBorderTopLeft     = '┏'
	symbolBorderTopRight    = '┓'
	symbolBorderBottomLeft  = '┗'
	symbolBorderBottomRight = '┛'
	symbolSnakeHead         = '▣'
	symbolSnakeBody         = '◼'
	symbolEmptyCell         = '⬚' // New symbol for empty cells in the game field
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
	snake              []Point
	food               Point
	foodType           int // Index of current food type in foodSymbols
	direction          Direction
	score              int
	highScore          int
	gameOver           bool
	foodTimer          int  // Countdown until food disappears
	foodVisible        bool // Is food currently visible?
	foodRespawnCounter int  // Countdown until next food appears
}

// Initialize a new game
func NewGame() *Game {
	g := &Game{
		snake:              make([]Point, initialSize),
		direction:          Right,
		score:              0,     // Explicitly initialize score to 0
		foodVisible:        false, // Start with no food
		foodRespawnCounter: 0,     // Spawn food immediately
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
	// Select random food type
	g.foodType = rand.Intn(len(foodSymbols))

	// Set a random timer for this food
	g.foodTimer = rand.Intn(maxFoodTime-minFoodTime) + minFoodTime

	// Make food visible
	g.foodVisible = true

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

	// Food timer management
	if g.foodVisible {
		// Countdown food timer
		g.foodTimer--
		if g.foodTimer <= 0 {
			// Food has disappeared
			g.foodVisible = false
			g.foodRespawnCounter = foodRespawnTime
		}
	} else {
		// Food is not visible, count down to respawn
		g.foodRespawnCounter--
		if g.foodRespawnCounter <= 0 {
			// Time to respawn food
			g.PlaceFood()
		}
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

	// Check food collision only if food is visible
	if g.foodVisible && newHead.X == g.food.X && newHead.Y == g.food.Y {
		// Award points based on food type
		pointsEarned := foodValues[g.foodType]
		g.score += pointsEarned

		// Flash score notification
		// (Could extend this in the future to show +N points briefly)

		// Update high score if current score is higher
		if g.score > g.highScore {
			g.highScore = g.score
		}

		// Place new food
		g.PlaceFood()
	} else {
		// Remove tail if no food was eaten
		g.snake = g.snake[:len(g.snake)-1]
	}
}

// Draw the game
func (g *Game) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Clear sidebar area explicitly to prevent artifacts
	clearSidebarArea()

	// Draw sidebar with minimal info
	drawSidebar(g)

	// Draw border with offset for sidebar
	for i := 0; i < width+2; i++ {
		termbox.SetCell(i+sidebarWidth, 0, symbolBorderHorizontal, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(i+sidebarWidth, height+1, symbolBorderHorizontal, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i := 0; i < height+2; i++ {
		termbox.SetCell(sidebarWidth, i, symbolBorderVertical, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width+sidebarWidth+1, i, symbolBorderVertical, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.SetCell(sidebarWidth, 0, symbolBorderTopLeft, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(width+sidebarWidth+1, 0, symbolBorderTopRight, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(sidebarWidth, height+1, symbolBorderBottomLeft, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(width+sidebarWidth+1, height+1, symbolBorderBottomRight, termbox.ColorWhite, termbox.ColorDefault)

	// Fill game field with empty cell symbols
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			termbox.SetCell(x+sidebarWidth+1, y+1, symbolEmptyCell, termbox.ColorDarkGray, termbox.ColorDefault)
		}
	}

	// Draw snake with offset for sidebar
	for i, p := range g.snake {
		symbol := symbolSnakeBody
		if i == 0 {
			// First segment is the head
			symbol = symbolSnakeHead
		}
		termbox.SetCell(p.X+sidebarWidth+1, p.Y+1, symbol, termbox.ColorGreen, termbox.ColorDefault)
	}

	// Draw food if visible, with color indicating timer
	if g.foodVisible {
		// Calculate color based on food timer
		var fg termbox.Attribute = termbox.ColorRed

		// Change color as timer runs down
		if g.foodTimer < minFoodTime/3 {
			fg = termbox.ColorRed | termbox.AttrBlink // Blinking when about to disappear
		} else if g.foodTimer < minFoodTime/2 {
			fg = termbox.ColorRed | termbox.AttrBold // Bold red when getting low
		}

		termbox.SetCell(g.food.X+sidebarWidth+1, g.food.Y+1, foodSymbols[g.foodType], fg, termbox.ColorDefault)
	}

	// Game over message (centered in game area)
	if g.gameOver {
		gameOverX := sidebarWidth + width/2
		gameOverMsg := "Game Over!\n\rPress 'q' to quit or 'r' to restart."
		scoreMsg := fmt.Sprintf("Final Score: %d", g.score)

		for i, ch := range []rune(gameOverMsg) {
			termbox.SetCell(gameOverX-len(gameOverMsg)/2+i, height/2, ch, termbox.ColorRed, termbox.ColorDefault)
		}

		for i, ch := range []rune(scoreMsg) {
			termbox.SetCell(gameOverX-len(scoreMsg)/2+i, height/2+1, ch, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

// Clear the entire sidebar area to prevent artifacts
func clearSidebarArea() {
	for y := 0; y < height+4; y++ { // +4 to include score area below game
		for x := 0; x < sidebarWidth; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

// Draw the sidebar with scores and food information
func drawSidebar(g *Game) {
	// Draw vertical separator line
	for i := 0; i < height+2; i++ {
		termbox.SetCell(sidebarWidth-1, i, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw minimal score display
	scoreStr := []rune(fmt.Sprintf("SCORE: %d", g.score))
	for i, ch := range scoreStr {
		termbox.SetCell(2+i, 2, ch, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
	}

	// Draw food value table header with minimal styling
	tableHeader := " "
	for i, ch := range []rune(tableHeader) {
		termbox.SetCell(sidebarWidth/2-len(tableHeader)/2+i, 5, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw food symbols and their values in a compact format
	for i := 0; i < len(foodSymbols); i++ {
		// Draw food symbol
		termbox.SetCell(4, 7+i, foodSymbols[i], termbox.ColorRed, termbox.ColorDefault)

		// Draw equals sign
		termbox.SetCell(6, 7+i, '=', termbox.ColorWhite, termbox.ColorDefault)

		// Draw points value
		valueStr := []rune(fmt.Sprintf("%d", foodValues[i]))
		for j := 0; j < len(valueStr); j++ {
			termbox.SetCell(8+j, 7+i, valueStr[j], termbox.ColorYellow, termbox.ColorDefault)
		}
	}
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

	// Initialize with horizontal speed (will be adjusted based on direction)
	updateInterval := time.Duration(baseSpeed) * time.Millisecond
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	highScore := 0 // Track high score across games

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				oldDirection := game.direction
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

				// If direction changed between horizontal/vertical, adjust the ticker
				if directionChanged(oldDirection, game.direction) {
					ticker.Stop()
					updateInterval = getUpdateInterval(game.direction)
					ticker = time.NewTicker(updateInterval)
				}

				if ev.Ch == 'q' {
					return
				} else if ev.Ch == 'r' && game.gameOver {
					// Preserve high score when starting a new game
					highScore = max(highScore, game.highScore)
					game = NewGame()
					game.highScore = highScore
				}
			}
		case <-ticker.C:
			game.Update()
			game.Draw()
		}
	}
}

// Helper function to check if direction changed between horizontal and vertical
func directionChanged(old, new Direction) bool {
	return (old == Up || old == Down) != (new == Up || new == Down)
}

// Get the appropriate update interval based on direction
func getUpdateInterval(dir Direction) time.Duration {
	if dir == Left || dir == Right {
		// Horizontal movement
		return time.Duration(baseSpeed) * time.Millisecond
	} else {
		// Vertical movement - adjust for aspect ratio
		return time.Duration(float64(baseSpeed)*aspectRatio) * time.Millisecond
	}
}

// Helper function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
