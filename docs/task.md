---
layout: default
title: Task Scheduling
nav_order: 10
---

# Task Scheduling

`github.com/mechanical-lich/mlge/task`

A priority-based task scheduler for AI and game entities with escalation, cooldowns, and proximity-based assignment.

## TaskAction

```go
type TaskAction string
```

String identifier for task types (e.g., `"gather"`, `"build"`, `"patrol"`).

## Task

```go
type Task struct {
    X, Y, Z         int
    Action           TaskAction
    Data             any
    Escalated        bool
    Completed        bool
    InProgress       bool
    ManuallyStopped  bool
    Created          time.Time
}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Start()` | Marks the task as in-progress |
| `Stop()` | Marks the task as manually stopped and not in-progress |
| `ReQueue()` | Resets the task to be available for assignment again |
| `Complete()` | Marks the task as completed |

## TaskScheduler

```go
type TaskScheduler struct{}
```

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddTask` | `(task *Task)` | Add a task to the scheduler |
| `GetTasks` | `() []*Task` | Get all tasks |
| `RemoveTask` | `(task *Task)` | Remove a specific task |
| `ClearTasks` | `()` | Remove all tasks |
| `SortTasks` | `()` | Sort by escalation then creation time; removes completed |
| `Count` | `() int` | Number of tasks in the queue |
| `GetNextTask` | `(allowed ...TaskAction) *Task` | Get and start the next available task |
| `GetClosestNextTask` | `(x, y, z int, allowed ...TaskAction) *Task` | Get the closest available task to a position |
| `PeekNextTask` | `(allowed ...TaskAction) *Task` | Preview next task without starting it |
| `PeekClosestNextTask` | `(x, y, z int, allowed ...TaskAction) *Task` | Preview closest task without starting it |

## Usage

```go
scheduler := &task.TaskScheduler{}

// Add tasks
scheduler.AddTask(&task.Task{
    X: 10, Y: 5, Z: 0,
    Action: "gather",
    Data:   &GatherData{ResourceType: "wood"},
})

scheduler.AddTask(&task.Task{
    X: 20, Y: 15, Z: 0,
    Action: "build",
    Data:   &BuildData{StructureType: "wall"},
})

// Sort to prioritize escalated and oldest tasks
scheduler.SortTasks()

// Get next task for a worker that can gather or build
t := scheduler.GetNextTask("gather", "build")
if t != nil {
    // Worker processes the task
    // ...
    t.Complete()
}

// Get closest task to a worker's position
t = scheduler.GetClosestNextTask(12, 8, 0, "gather")
```

## Task Lifecycle

1. **Created** — Task is added to the scheduler
2. **Started** — `GetNextTask`/`GetClosestNextTask` calls `Start()` automatically
3. **In Progress** — Worker is executing the task
4. **Completed** — Worker calls `Complete()`, task is removed on next `SortTasks()`
5. **Stopped** — Worker calls `Stop()` if interrupted; can be re-queued with `ReQueue()`

## Escalation

Tasks can be marked as `Escalated = true` to give them higher priority. `SortTasks()` places escalated tasks before non-escalated ones.
