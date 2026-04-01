package domain

import "errors"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

var ErrInvalidTransition = errors.New("invalid status transition")

var transitions = map[TaskStatus][]TaskStatus{
	StatusTodo:       {StatusInProgress},
	StatusInProgress: {StatusDone},
}

func (s TaskStatus) CanTransitionTo(target TaskStatus) bool {
	for _, allowed := range transitions[s] {
		if allowed == target {
			return true
		}
	}
	return false
}
