package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/sse"
)

type (
	// ColorGame is a game where the player can toggle the color of an element in a large grid of elements.
	ColorGame struct {
		// Colors is the number of colors that can be toggled.
		Colors map[int]map[int]string
	}

	App struct {
		htmx *htmx.HTMX
	}
)

var (
	sseManager sse.Manager
	game       *ColorGame
)

func main() {
	// Create a new game with 10 rows, 10 columns, and 3 colors.
	game = &ColorGame{
		Colors: generateColoredGrid(100, 100),
	}

	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	sseManager = sse.NewManager(5)

	mux := http.NewServeMux()

	htmx.UseTemplateCache = false

	mux.Handle("GET /", http.HandlerFunc(app.Index))
	mux.Handle("POST /", http.HandlerFunc(app.ToggleColor))
	mux.Handle("GET /sse", http.HandlerFunc(app.SSE))

	err := http.ListenAndServe(":3456", mux)
	log.Fatal(err)
}

func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the color game",
		"Game": game,
	}

	page := htmx.NewComponent("index.gohtml").SetData(data).AddTemplateFunction("safe", func(s string) template.HTML {
		return template.HTML(s)
	})

	_, err := h.Render(r.Context(), page)
	if err != nil {
		fmt.Printf("error rendering page: %v", err.Error())
	}
}

type ToggleColorInput struct {
	X             int
	Y             int
	SelectedColor string
}

func (a *App) ToggleColor(w http.ResponseWriter, r *http.Request) {
	var input ToggleColorInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the color at the x and y coordinates.
	game.Colors[input.X][input.Y] = input.SelectedColor

	div := `<div hx-post='/' hx-ext='json-enc' id='color%d%d' sse-swap='color%d%d' hx-vals='{ "X":%d ,"Y":%d,"SelectedColor": %s }' style='background-color: {{ safe $color }}; width: 10px; height: 10px; float:left; cursor: pointer;' onclick='updateButtonColor(this)'></div>`

	// Send a message to all connected clients with the updated color.
	sseManager.Send(
		sse.NewMessage(fmt.Sprintf(div, input.X, input.Y, input.X, input.Y, input.X, input.Y, game.Colors[input.X][input.Y])).
			WithEvent(fmt.Sprintf("color%d%d", input.X, input.Y)),
	)

	_, _ = w.Write([]byte(div))
}

func (a *App) SSE(w http.ResponseWriter, r *http.Request) {
	cl := sse.NewClient(randStringRunes(10))

	fmt.Println("new client connected")
	fmt.Println(cl.ID)

	sseManager.Handle(w, r, cl)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateColoredGrid(rows, cols int) map[int]map[int]string {
	colors := make(map[int]map[int]string)

	for x := 0; x < rows; x++ {
		colors[x] = make(map[int]string)
		for y := 0; y < cols; y++ {
			colors[x][y] = randomColor()
		}
	}

	return colors
}

func randomColor() string {
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
