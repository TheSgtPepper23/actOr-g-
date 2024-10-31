package models

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type WindowConfig struct {
	TaskBuffer    *Task
	TextBuffer    string
	SelectedTasks []*Task
	NullArea      *rl.Rectangle
	Lines         []rl.Vector2
	Tasks         []*Task
	WIDTH         int32
	HEIGHT        int32
}

func NewConfig() WindowConfig {
	return WindowConfig{
		WIDTH:      800,
		HEIGHT:     600,
		TaskBuffer: nil,
	}
}

func (w *WindowConfig) Update() {
	freeMouse := w.NullArea == nil || !rl.CheckCollisionPointRec(rl.GetMousePosition(), *w.NullArea)

	for _, task := range w.Tasks {
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), task.GetRect()) {
			freeMouse = false
			if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
				w.Lines = append(w.Lines, task.GetCenter())
				w.SelectedTasks = append(w.SelectedTasks, task)
			}
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && freeMouse {
		task := &Task{
			Shape: rl.GetMousePosition(),
		}
		baseRect := task.GetRect()

		// The position of the rectangle is modified so its never rendered outside the window.
		// with the camera implementation this code should change to get the position relative to the
		// camera i guess
		if baseRect.X+baseRect.Width > float32(w.WIDTH) {
			baseRect.X = baseRect.X - ((baseRect.X + baseRect.Width) - float32(w.WIDTH))
		}

		if baseRect.Y+baseRect.Height > float32(w.HEIGHT) {
			baseRect.Y = baseRect.Y - ((baseRect.Y + baseRect.Height) - float32(w.HEIGHT))
		}

		w.NullArea = &baseRect
		w.TaskBuffer = task
	}
}

func (w *WindowConfig) clearBuffers() {
	w.SelectedTasks = make([]*Task, 0)
	w.Lines = make([]rl.Vector2, 0)
	w.TaskBuffer = nil
	w.TextBuffer = ""
	w.NullArea = nil
}

func (w *WindowConfig) CreateNewTask() {
	resp := gui.WindowBox(*w.NullArea, "New task")
	var padding float32 = 20
	gui.TextBox(
		rl.NewRectangle(w.NullArea.X+padding, w.NullArea.Y+padding*1.5, w.NullArea.Width-padding*2, 30),
		&w.TextBuffer,
		150,
		true,
	)
	okButton := gui.Button(
		rl.NewRectangle(w.NullArea.X+w.NullArea.Width-padding-50, w.NullArea.Y+padding*3.5, 50, 30),
		"Ok",
	)
	cancelButton := gui.Button(
		rl.NewRectangle(w.NullArea.X+w.NullArea.Width-padding*2-100, w.NullArea.Y+padding*3.5, 50, 30),
		"Cancel",
	)

	if resp || cancelButton {
		w.clearBuffers()
	}

	if okButton {
		if len(w.SelectedTasks) != 0 {
			for _, task := range w.SelectedTasks {
				w.TaskBuffer.PreviousTasks = append(w.TaskBuffer.PreviousTasks, task)
			}
		}
		w.TaskBuffer.Title = w.TextBuffer
		w.Tasks = append(w.Tasks, w.TaskBuffer)
		w.clearBuffers()
	}
}

func (w *WindowConfig) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	if w.TaskBuffer == nil {
		for _, line := range w.Lines {
			rl.DrawLineBezier(line, rl.GetMousePosition(), 4, rl.Blue)
		}
	} else {
		for _, line := range w.Lines {
			rl.DrawLineBezier(line, w.TaskBuffer.GetCenter(), 4, rl.Blue)
		}
	}

	// Draw lines
	for _, rect := range w.Tasks {
		for _, parent := range rect.PreviousTasks {
			rl.DrawLineBezier(rect.GetCenter(), parent.GetCenter(), 4, rl.Red)
		}
	}

	// Draw ready tasks
	for _, rect := range w.Tasks {
		rect.Draw()
	}

	// Draw task Buffer
	if w.TaskBuffer != nil {
		w.CreateNewTask()
	}

	rl.EndDrawing()
}
