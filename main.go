package main

import (
	"github.com/TheSgtPepper23/actOrg/models"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	config := models.NewConfig()
	rl.InitWindow(config.WIDTH, config.HEIGHT, "ActOr(g)")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		config.Update()
		config.Draw()
	}
}
