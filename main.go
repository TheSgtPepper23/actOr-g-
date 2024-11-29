package main

import (
	_ "embed"
	"fmt"

	"github.com/TheSgtPepper23/actOrg/models"
	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed Orbitron-Regular.ttf
var fontData []byte

func main() {
	fmt.Println(fontData)
	font := rl.LoadFontFromMemory(".ttf", fontData, 32, nil)
	defer rl.UnloadFont(font)
	config := models.NewConfig()
	rl.InitWindow(config.WIDTH, config.HEIGHT, "ActOr(g)")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		config.Update()
		config.Draw()
	}
}
