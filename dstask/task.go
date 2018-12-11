package dstask

// main task data structures

import (
	"sort"
	"time"
)

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"
	STATUS_RECURRING = "recurring" // tentative

	GIT_REPO   = "~/.dstask/"
	CACHE_FILE = "~/.cache/dstask/completion_cache.gob"

	// filter: P1 P2 etc
	PRIORITY_CRITICAL = "P1"
	PRIORITY_HIGH     = "P2"
	PRIORITY_NORMAL   = "P3"
	PRIORITY_LOW      = "P4"
)

// for import (etc) it's necessary to have full context
var ALL_STATUSES = []string{
	STATUS_PENDING,
	STATUS_ACTIVE,
	STATUS_RESOLVED,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_SOMEDAY,
	STATUS_RECURRING,
}

// for most operations, it's not necessary or desirable to load the expensive resolved tasks
var NORMAL_STATUSES = []string{
	STATUS_PENDING,
	STATUS_ACTIVE,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_SOMEDAY,
	STATUS_RECURRING,
}

type SubTask struct {
	Summary  string
	Resolved bool
}

type Task struct {
	// not stored in file -- rather filename and directory
	uuid   string
	status string
	// used to determine if an unlink should happen if status changes
	originalFilepath string

	// concise representation of task
	Summary string
	// task in more detail, only if necessary
	Description string
	Tags        []string
	Project     string
	// see const.go for PRIORITY_ strings
	Priority    string
	DelegatedTo string
	Subtasks    []SubTask
	Comments    []string
	// uuids of tasks that this task depends on
	// blocked status can be derived.
	// TODO possible filter: :blocked. Also, :overdue
	Dependencies []string

	Created  time.Time
	Modified time.Time
	Resolved time.Time
	Due      time.Time
}

type TaskSet struct {
	Tasks           []Task
	CurrentContext  string
	knownUuids      map[string]bool
	GitRepoLocation string
}

// Call before addressing and display. Sorts by status then UUID.
func (ts *TaskSet) SortTaskList() {
	sort.Slice(ts.Tasks, func(i, j int) bool {
		ti := ts.Tasks[i]
		tj := ts.Tasks[j]

		// TODO define precedent of statuses
		if ti.status == tj.status {
			return ti.uuid < tj.uuid
		} else {
			return ti.status < tj.status
		}
	})
}

// add a task, but only if it has a new uuid. Return true if task was added.
func (ts *TaskSet) MaybeAddTask(task Task) bool {
	if ts.knownUuids[task.uuid] {
		// load tasks, do not overwrite
		return false
	}

	ts.knownUuids[task.uuid] = true
	ts.Tasks = append(ts.Tasks, task)
	return true
}

// filter should be set before loading any data. The filter can be used to
// optimise a bit -- eg when listing, completed tasks should not be shown so we
// can avoid loading them. However when importing, it is important to load all
// tasks for full context.
type TaskFilter struct {
	Text     string
	Tags     []string
	Antitags []string
	Project  string
	Priority int
}

//func (ts *TaskSet) filter(filter *TaskFilter) TaskSet {
//
//}
//
//func (t *Task) Save() error {
//
//}
