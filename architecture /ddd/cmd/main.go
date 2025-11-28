package main

import (
	"fmt"
	"github.com/Nuyoahch/gopulse/architecture /ddd/internal/task/app"
)

func main() {
}

func mustTask(t *app.TaskService, task *app.TaskService) *app.TaskService {
	return t
}

func printTask(t interface {
	ID() string
	Name() string
	Status() interface{}
}) {
	// 为了简单，这里你可以直接在 main.go 里写：
	// func printTask(t *domain.Task) {...}
	fmt.Printf("Task{id=%s, name=%s, status=%v}\n", t.ID(), t.Name(), t.Status())
}
