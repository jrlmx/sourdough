package main

import "fmt"

type cleanUpTask struct {
	task      func() error
	onSuccess bool
}

type cleanUpManager struct {
	tasks []cleanUpTask
}

func (cm *cleanUpManager) addTask(task func() error, onSuccess bool) {
	cm.tasks = append(cm.tasks, cleanUpTask{task, onSuccess})
}

func (cm *cleanUpManager) cleanUp(success bool) {
	for _, task := range cm.tasks {
		if success && !task.onSuccess {
			continue
		}
		if err := task.task(); err != nil {
			fmt.Println(err)
		}
	}
}

func NewCleanUpManager() *cleanUpManager {
	return &cleanUpManager{
		tasks: []cleanUpTask{},
	}
}
