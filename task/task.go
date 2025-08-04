package task

import (
	"slices"
	"sort"
	"time"

	"github.com/mechanical-lich/mlge/utility"
)

type TaskAction string

// Default Task actions
const PickupAction TaskAction = "pickup"                  // Pickup item at tile
const ScoutAction TaskAction = "scout"                    // Non-aggressive move to location
const AggressiveMoveAction TaskAction = "aggressive_move" // Move to location but attack things in sight along the way
const AttackAction TaskAction = "attack"                  // Attack whatever is in square if possible (moves to it)
const BuildAction TaskAction = "build"                    // Build at tile
const DigAction TaskAction = "dig"                        // Digs up a tile if possible.
const HuntAction TaskAction = "hunt"                      // Hunt for food
const ButcherAction TaskAction = "butcher"                // Butcher an animal

type Task struct {
	X               int
	Y               int
	Z               int
	Action          TaskAction
	Data            any
	Escalated       bool
	Completed       bool
	InProgress      bool
	Created         time.Time
	ManuallyStopped bool
}

// Mark a task as in progress
func (t *Task) Start() {
	t.InProgress = true
	t.ManuallyStopped = false
}

// If a task can't be completed call Stop to put it back in the queue
func (t *Task) Stop() {
	t.InProgress = false
	t.Escalated = false
	t.ManuallyStopped = true
	t.Created = time.Now()
}

// ReQueue puts the task back in the queue without the restart delay
func (t *Task) ReQueue() {
	t.InProgress = false
	t.Escalated = false
	t.Created = time.Now()
}

// Mark a task as complete.  Completed tasks get automatically cleaned up.
func (t *Task) Complete() {
	//log.Printf("Task to %s completed at %s at [%d,%d,%d]", t.Action, t.Created.String(), t.X, t.Y, t.Z)
	t.Completed = true
}

/*
The task scheduler's job is to store all current tasks and give out tasks based on
priority and age.

Tasks that get started then stopped get put on the bottom of the queue.
Tasks that are completed get cleaned up automatically.
*/
type TaskScheduler struct {
	tasks []*Task
}

// Add a new task
func (ts *TaskScheduler) AddTask(task *Task) {
	ts.tasks = append(ts.tasks, task)
}

// Get all tasks
func (ts *TaskScheduler) GetTasks() []*Task {
	return ts.tasks
}

// Removes a task from the queue
func (ts *TaskScheduler) RemoveTask(task *Task) {
	for i, t := range ts.tasks {
		if t == task {
			ts.tasks = append(ts.tasks[:i], ts.tasks[i+1:]...)
			break
		}
	}
}

// Removes all currently scheduled tasks
func (ts *TaskScheduler) ClearTasks() {
	ts.tasks = []*Task{}
}

// Sort tasks by creation date by prioritize escalated tasks.
func (ts *TaskScheduler) SortTasks() {
	sort.Slice(ts.tasks, func(i, j int) bool {
		if ts.tasks[i].Escalated != ts.tasks[j].Escalated {
			return ts.tasks[i].Escalated
		}
		return ts.tasks[i].Created.Before(ts.tasks[j].Created)
	})

	// Clean up the completed tasks
	for i := len(ts.tasks) - 1; i >= 0; i-- {
		if ts.tasks[i].Completed {
			ts.tasks = append(ts.tasks[:i], ts.tasks[i+1:]...)
		}
	}
}

// Get the next task.   If allowed list is empty any task is picked.
// The task will be marked as "started" before it is returned.  It is up to the caller to mark it as completed
// or to stop it if it can't be completed.
func (ts *TaskScheduler) GetNextTask(allowed_tasks ...TaskAction) *Task {
	t := ts.PeekNextTask(allowed_tasks...)
	if t != nil {
		t.Start()
	}

	return t
}

func (ts *TaskScheduler) GetClosestNextTask(x, y, z int, allowed_tasks ...TaskAction) *Task {
	t := ts.PeekClosestNextTask(x, y, z, allowed_tasks...)
	if t != nil {
		t.Start()
	}

	return t
}

func (ts *TaskScheduler) PeekClosestNextTask(x, y, z int, allowed_tasks ...TaskAction) *Task {
	ts.SortTasks()
	if len(ts.tasks) == 0 {
		return nil
	}

	now := time.Now()
	var closestTask *Task
	minDist := int(^uint(0) >> 1) // Max int

	for _, task := range ts.tasks {
		if (len(allowed_tasks) == 0 || slices.Contains(allowed_tasks, task.Action)) &&
			!task.Completed && !task.InProgress {

			// Only allow manually stopped tasks if 5 seconds have passed since their creation
			if task.ManuallyStopped && now.Sub(task.Created) < 5*time.Second {
				continue
			}

			// Manhattan distance
			dist := utility.Abs(task.X-x) + utility.Abs(task.Y-y) + utility.Abs(task.Z-z)

			// If this task is closer, or if it's the same distance but higher priority (earlier in sorted list)
			if closestTask == nil || dist < minDist {
				closestTask = task
				minDist = dist
			}
		}
	}

	return closestTask
}

// PeekNextTask returns the next task without marking it as started.
func (ts *TaskScheduler) PeekNextTask(allowed_tasks ...TaskAction) *Task {
	ts.SortTasks()
	if len(ts.tasks) == 0 {
		return nil
	}

	now := time.Now()

	for i := 0; i < len(ts.tasks); i++ {
		task := ts.tasks[i]

		if (len(allowed_tasks) == 0 || slices.Contains(allowed_tasks, task.Action)) &&
			!task.Completed && !task.InProgress {

			// Only allow manually stopped tasks if 5 seconds have passed since their creation
			if task.ManuallyStopped && now.Sub(task.Created) < 5*time.Second {
				continue
			}

			return task
		}
	}

	return nil
}

func (ts *TaskScheduler) Count() int {
	return len(ts.tasks)
}
