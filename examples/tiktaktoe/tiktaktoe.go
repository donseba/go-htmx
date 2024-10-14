package main

import (
	"github.com/donseba/go-htmx"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type (
	App struct {
		htmx *htmx.HTMX
	}

	Game struct {
		Board [3][3]string
		Turn  string
		Mu    sync.Mutex
	}
)

var game = Game{
	Board: [3][3]string{},
	Turn:  "X",
}

func main() {
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()

	htmx.UseTemplateCache = false

	mux.Handle("GET /", http.HandlerFunc(app.Home))
	mux.Handle("PUT /set/{row}/{column}", http.HandlerFunc(app.Set))
	mux.Handle("GET /reset", http.HandlerFunc(app.Reset))

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	game.Mu.Lock()
	defer game.Mu.Unlock()

	h := a.htmx.NewHandler(w, r)

	page := htmx.NewComponent("board.gohtml").SetData(map[string]any{
		"game":   &game,
		"winner": checkWinner(game.Board),
	}).Wrap(mainContent(), "board")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		log.Printf("error rendering page: %v", err.Error())
	}
}

func (a *App) Set(w http.ResponseWriter, r *http.Request) {
	game.Mu.Lock()
	defer game.Mu.Unlock()

	h := a.htmx.NewHandler(w, r)

	row, _ := strconv.Atoi(r.PathValue("row"))
	column, _ := strconv.Atoi(r.PathValue("column"))

	game.Board[row][column] = game.Turn
	if game.Turn == "X" {
		game.Turn = "O"
	} else {
		game.Turn = "X"
	}

	page := htmx.NewComponent("board.gohtml").SetData(map[string]any{
		"game":   &game,
		"winner": checkWinner(game.Board),
	}).Wrap(mainContent(), "board")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		log.Printf("error rendering page: %v", err.Error())
	}
}

func (a *App) Reset(w http.ResponseWriter, r *http.Request) {
	game.Mu.Lock()
	defer game.Mu.Unlock()

	h := a.htmx.NewHandler(w, r)

	game.Board = [3][3]string{}
	game.Turn = "X"

	page := htmx.NewComponent("board.gohtml").SetData(map[string]any{
		"game": &game,
	}).Wrap(mainContent(), "board")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		log.Printf("error rendering page: %v", err.Error())
	}
}

func checkWinner(board [3][3]string) string {
	// Check rows, columns, and diagonals for a winner
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][0] != "" {
			return board[i][0]
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[0][i] != "" {
			return board[0][i]
		}
	}

	if board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != "" {
		return board[0][0]
	}

	if board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != "" {
		return board[0][2]
	}

	return ""
}

func mainContent() htmx.RenderableComponent {
	return htmx.NewComponent("index.gohtml")
}
