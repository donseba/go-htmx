package main

import "math/rand"

func moveSnake(game *SnakeGame) {
	head := game.Snake[0]
	newHead := Position{X: head.X + game.Dir.X, Y: head.Y + game.Dir.Y}

	// Wrap around the edges of the board
	if newHead.X < 0 {
		newHead.X = 19 // Move to the rightmost edge
	} else if newHead.X >= 20 {
		newHead.X = 0 // Move to the leftmost edge
	}

	if newHead.Y < 0 {
		newHead.Y = 19 // Move to the bottom edge
	} else if newHead.Y >= 20 {
		newHead.Y = 0 // Move to the top edge
	}

	// Check if the snake eats the food
	if newHead == game.Food {
		// Grow the snake by not removing the last part
		placeFood(game) // Place new food
	} else {
		// Move the snake by removing the tail
		game.Snake = game.Snake[:len(game.Snake)-1]
	}

	// Add the new head to the front of the snake
	game.Snake = append([]Position{newHead}, game.Snake...)

	// Update the board
	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = ""
		}
	}
	for _, pos := range game.Snake {
		game.Board[pos.X][pos.Y] = "S"
	}
	// Place food on the board
	game.Board[game.Food.X][game.Food.Y] = "F"
}

func placeFood(game *SnakeGame) {
	for {
		x := rand.Intn(20)
		y := rand.Intn(20)
		foodPos := Position{X: x, Y: y}

		// Ensure food is not placed on the snake
		occupied := false
		for _, pos := range game.Snake {
			if pos == foodPos {
				occupied = true
				break
			}
		}
		if !occupied {
			game.Food = foodPos
			game.Board[game.Food.X][game.Food.Y] = "F"
			break
		}
	}
}
