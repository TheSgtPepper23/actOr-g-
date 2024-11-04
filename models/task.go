package models

import rl "github.com/gen2brain/raylib-go/raylib"

type Task struct {
	Title         string
	PreviousTasks []*Task
	Shape         rl.Vector2
	Size          int32
	Completed     bool
	Dragging      bool
}

const padding int32 = 5

func (t *Task) ReadyToGo() bool {
	isReady := true
	for _, task := range t.PreviousTasks {
		if !task.Completed {
			isReady = false
			break
		}
	}
	return isReady
}

func (t Task) GetTextSize() int32 {
	return rl.MeasureText(t.Title, 20)
}

func (t Task) GetRect() rl.Rectangle {
	if t.Title == "" {
		return rl.NewRectangle(t.Shape.X, t.Shape.Y, 300, 110)
	}
	return rl.NewRectangle(t.Shape.X, t.Shape.Y, float32(t.GetTextSize()+padding*2), 50)
}

func (t Task) GetCenter() rl.Vector2 {
	rect := t.GetRect()
	return rl.NewVector2(t.Shape.X+rect.Width/2, t.Shape.Y+rect.Height/2)
}

func (t Task) Draw() {
	rl.DrawRectangleRec(t.GetRect(), rl.LightGray)
	rl.DrawText(t.Title, int32(t.Shape.X)+padding, int32(t.Shape.Y)+padding, 20, rl.Black)
}
