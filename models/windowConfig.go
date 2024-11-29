package models

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type WindowConfig struct {
	TaskBuffer    *Task
	LocalFont     rl.Font
	NullArea      *rl.Rectangle
	TextBuffer    string
	SelectedTasks []*Task
	Lines         []rl.Vector2
	Tasks         []*Task
	Buttons       []string
	Camera        rl.Camera2D
	MenuRectangle rl.Rectangle
	WIDTH         int32
	HEIGHT        int32
}

func NewConfig() WindowConfig {
	tempButtons := []string{"Nuevo", "Abrir", "Cerrar", "Guardar", "Guardar como..."}
	return WindowConfig{
		WIDTH:         1280,
		HEIGHT:        720,
		TaskBuffer:    nil,
		Buttons:       tempButtons,
		MenuRectangle: rl.NewRectangle(0, 0, 1280, 60),
		Camera: rl.Camera2D{
			Offset:   rl.NewVector2(400, 300),
			Rotation: 0,
			Target:   rl.Vector2Zero(),
			Zoom:     1,
		},
	}
}

func (w *WindowConfig) Update() {
	// nullarea is a gui element so its rendered outside the 2d mode so the real mosuse position should be used
	freeMouse := (w.NullArea == nil || !rl.CheckCollisionPointRec(rl.GetMousePosition(), *w.NullArea)) && !rl.CheckCollisionPointRec(rl.GetMousePosition(), w.MenuRectangle)

	if freeMouse && rl.IsKeyDown(rl.KeySpace) && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		rl.SetMouseCursor(rl.MouseCursorResizeAll)
		mouseDelta := rl.GetMouseDelta()
		mouseDelta = rl.Vector2Scale(mouseDelta, -1.0*w.Camera.Zoom)
		w.Camera.Target = rl.Vector2Add(w.Camera.Target, mouseDelta)
	}

	if rl.IsKeyReleased(rl.KeySpace) || rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		mouseWorld := rl.GetScreenToWorld2D(rl.GetMousePosition(), w.Camera)
		w.Camera.Offset = rl.GetMousePosition()
		w.Camera.Target = mouseWorld
		w.Camera.Zoom += wheel * 0.325
		if w.Camera.Zoom < 0.100 {
			w.Camera.Zoom = 0.100
		}
	}

	for _, task := range w.Tasks {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) && task.Dragging {
			task.Shape = rl.GetScreenToWorld2D(rl.GetMousePosition(), w.Camera)
		}
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) && task.Dragging {
			task.Dragging = false
		}

		if rl.CheckCollisionPointRec(rl.GetScreenToWorld2D(rl.GetMousePosition(), w.Camera), task.GetRect()) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				task.Dragging = true
			}
			if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
				w.Lines = append(w.Lines, task.GetCenter())
				w.SelectedTasks = append(w.SelectedTasks, task)
			}
			freeMouse = false
		}
	}

	// Create a new TaskBuffer
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && freeMouse && !rl.IsKeyDown(rl.KeySpace) {
		task := &Task{
			Shape: rl.GetScreenToWorld2D(rl.GetMousePosition(), w.Camera),
		}

		baseRect := Task{
			Shape: rl.GetMousePosition(),
		}.GetRect()

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

	rl.BeginMode2D(w.Camera)
	if w.TaskBuffer == nil {
		for _, line := range w.Lines {
			rl.DrawLineBezier(line, rl.GetScreenToWorld2D(rl.GetMousePosition(), w.Camera), 4, rl.Blue)
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
	rl.EndMode2D()

	// Draw task Buffer
	if w.TaskBuffer != nil {
		w.CreateNewTask()
	}

	// Menu area

	gui.Panel(w.MenuRectangle, "")
	var buttonWide float32 = 100
	var buttonPadding float32 = 15
	for i, button := range w.Buttons {
		i := float32(i)
		gui.Button(rl.NewRectangle(buttonPadding+buttonWide*i+buttonPadding*i, buttonPadding, buttonWide, 30), button)
	}
	rl.EndDrawing()
}
