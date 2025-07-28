package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskScheduler_AddTask(t *testing.T) {
	ts := &TaskScheduler{}
	task := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now()}
	ts.AddTask(task)
	assert.Equal(t, 1, len(ts.tasks))
	assert.Equal(t, task, ts.tasks[0])
}

func TestTaskScheduler_GetTasks(t *testing.T) {
	ts := &TaskScheduler{}
	task := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now()}
	ts.AddTask(task)
	tasks := ts.GetTasks()
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, task, tasks[0])
}

func TestTaskScheduler_RemoveTask(t *testing.T) {
	ts := &TaskScheduler{}
	task := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now()}
	ts.AddTask(task)
	ts.RemoveTask(task)
	assert.Equal(t, 0, len(ts.tasks))
}

func TestTaskScheduler_ClearTasks(t *testing.T) {
	ts := &TaskScheduler{}
	task := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now()}
	ts.AddTask(task)
	ts.ClearTasks()
	assert.Equal(t, 0, len(ts.tasks))
}

func TestTaskScheduler_SortTasks(t *testing.T) {
	ts := &TaskScheduler{}
	task1 := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now()}
	task2 := &Task{X: 3, Y: 4, Action: ScoutAction, Created: time.Now().Add(-time.Second), Escalated: true}
	ts.AddTask(task1)
	ts.AddTask(task2)

	ts.SortTasks()
	assert.Equal(t, task2, ts.tasks[0])
	assert.Equal(t, task1, ts.tasks[1])
}

func TestTaskScheduler_SortTasks_Escalation(t *testing.T) {
	ts := &TaskScheduler{}
	task1 := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now(), Escalated: true}                   // New escalated task
	task2 := &Task{X: 3, Y: 4, Action: ScoutAction, Created: time.Now()}                                     // New unescalated task
	task3 := &Task{X: 3, Y: 6, Action: AttackAction, Created: time.Now().Add(-time.Second)}                  // Oldest non-escalated task
	task4 := &Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now().Add(-time.Second), Escalated: true} // Oldest escalated task

	ts.AddTask(task1)
	ts.AddTask(task2)
	ts.AddTask(task3)
	ts.AddTask(task4)

	ts.SortTasks()
	assert.Equal(t, task4, ts.tasks[0]) // Oldest escalated task should be first
	assert.Equal(t, task1, ts.tasks[1]) // Escalated tasks should come before unescalated tasks
	assert.Equal(t, task3, ts.tasks[2]) // Oldest non-escalated task should be third
	assert.Equal(t, task2, ts.tasks[3]) // Unescalated tasks should come last
}

func TestTaskScheduler_GetNextTask(t *testing.T) {
	ts := &TaskScheduler{}

	// Test case 1: No tasks available
	task := ts.GetNextTask()
	assert.Nil(t, task)

	// Add some tasks for further testing
	task1 := Task{X: 1, Y: 2, Action: PickupAction, Created: time.Now()}
	task2 := Task{X: 3, Y: 4, Action: ScoutAction, Created: time.Now().Add(-time.Second)}
	task3 := Task{X: 5, Y: 6, Action: AttackAction, Created: time.Now().Add(-time.Second * 2)}

	ts.AddTask(&task1)
	ts.AddTask(&task2)
	ts.AddTask(&task3)

	// Test case 2: Tasks with allowed actions
	task = ts.GetNextTask(PickupAction)
	assert.NotNil(t, task)
	assert.Equal(t, PickupAction, task.Action)

	// Test case 3: No allowed tasks
	task = ts.GetNextTask(DigAction)
	assert.Nil(t, task)

	// Test case 4: First task in list
	task = ts.GetNextTask()
	assert.NotNil(t, task)
	assert.Equal(t, AttackAction, task.Action)

	task.Complete()

	// Completed tasks don't get picked
	task = ts.GetNextTask()
	assert.False(t, task.Completed)
	assert.NotNil(t, task)
	assert.Equal(t, ScoutAction, task.Action)

	// New escalated task comes next
	task4 := &Task{X: 5, Y: 6, Action: AggressiveMoveAction, Created: time.Now(), Escalated: true}
	ts.AddTask(task4)
	task = ts.GetNextTask()
	assert.True(t, task.Escalated)
	assert.NotNil(t, task)
	assert.Equal(t, AggressiveMoveAction, task.Action)

}

func TestTaskScheduler_PeekNextTask_ManuallyStopped(t *testing.T) {
	ts := &TaskScheduler{}

	task1 := &Task{
		Action:          DigAction,
		Created:         time.Now().Add(-3 * time.Second),
		ManuallyStopped: true,
	}
	task2 := &Task{
		Action:          PickupAction,
		Created:         time.Now().Add(-6 * time.Second),
		ManuallyStopped: true,
	}
	task3 := &Task{
		Action:  ScoutAction,
		Created: time.Now().Add(-7 * time.Second),
	}

	ts.AddTask(task1)
	ts.AddTask(task2)
	ts.AddTask(task3)

	peeked := ts.PeekNextTask()
	assert.Equal(t, task3, peeked, "Expected normal task")

	task3.InProgress = true // Simulate task3 being in progress
	peeked = ts.PeekNextTask()
	assert.Equal(t, task2, peeked, "Expected manually stopped task2")
}

func TestTaskScheduler_PeekClosestNextTask_ManuallyStopped(t *testing.T) {
	ts := &TaskScheduler{}

	task1 := &Task{
		X:               1,
		Y:               1,
		Z:               1,
		Action:          PickupAction,
		Created:         time.Now().Add(-3 * time.Second),
		ManuallyStopped: true,
	}
	task2 := &Task{
		X:               2,
		Y:               2,
		Z:               2,
		Action:          PickupAction,
		Created:         time.Now().Add(-6 * time.Second),
		ManuallyStopped: true,
	}
	task3 := &Task{
		X:       0,
		Y:       0,
		Z:       0,
		Action:  PickupAction,
		Created: time.Now(),
	}

	ts.AddTask(task1)
	ts.AddTask(task2)
	ts.AddTask(task3)

	peeked := ts.PeekClosestNextTask(0, 0, 0)
	assert.Equal(t, task3, peeked, "Expected closest normal task")

	task3.InProgress = true
	peeked = ts.PeekClosestNextTask(0, 0, 0)
	assert.Equal(t, task2, peeked, "Expected manually stopped task2")
}
