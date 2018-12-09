package dstask

// see https://taskwarrior.org/docs/design/task.html

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type TwTime struct {
	time.Time
}

// TW time is stored in ISO 8601 basic format, which the default parser does
// not understand.
func (tt *TwTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == "null" {
		return nil
	}

	// drop quotes
	if len(s) > 2 {
		s = s[1:len(s)-1]
	}

	if len(s) == 16 {
		// convert from basic format to normal format which is RFC3339 compatible
		s = s[0:4]+"-"+s[4:6]+"-"+s[6:11]+":"+s[11:13]+":"+s[13:len(s)]
	}

	t, err := time.Parse(time.RFC3339, s)
	tt.Time = t

	if (err != nil) {
		return err
	}

	return nil
}

type TwAnnotation struct {
	Description string `json:"description"`
	Entry       string `json: entry`
}

type TwTask struct {
	Description string         `json:"description"`
	End         TwTime         `json:"end"`
	Entry       TwTime         `json: entry`
	Start       TwTime         `json: start`
	Modified    TwTime         `json: modified`
	Due         TwTime         `json: due`
	Status      string         `json: status`
	Project     string         `json: project`
	Priority    string         `json: priority`
	Depends     string         `json: depends`
	Tags        []string       `json: tags`
	Uuid        string         `json: uuid`
	Annotations []TwAnnotation `json:annotations`
}

var priorityMap = map[string]string{
	"H": PRIORITY_HIGH,
	"M": PRIORITY_NORMAL,
	"L": PRIORITY_LOW,
}

// convert a tw status into a dstask status
func (t *TwTask) ConvertStatus() string {
	if t.Start.Time.IsZero() {
		return STATUS_ACTIVE
	}

	switch t.Status {
	case "completed":
		return STATUS_RESOLVED
	case "deleted":
		return STATUS_RESOLVED
	case "waiting":
		return STATUS_PENDING
	case "recurring":
		// TODO -- implement reccurence
		//return STATUS_RECURRING
		return STATUS_RESOLVED
	default:
		return t.Status
	}
}

func (t *TwTask) ConvertAnnotations() []string {
	var comments []string

	for _, ann := range t.Annotations {
		comments = append(comments, ann.Description)
	}

	return comments
}

func (ts *TaskSet) ImportFromTaskwarrior() error {
	var twTasks []TwTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&twTasks)

	if err != nil {
		return err
	}

	for _, twTask := range twTasks {
		fmt.Println(twTask.Status, twTask.Description)
		ts.tasks = append(ts.tasks, Task{
			uuid:         twTask.Uuid,
			status:       twTask.ConvertStatus(),
			Summary:      twTask.Description,
			Tags:         twTask.Tags,
			Project:      twTask.Project,
			Priority:     priorityMap[twTask.Priority],
			Comments:     twTask.ConvertAnnotations(),
			Dependencies: strings.Split(twTask.Depends, ","),
			Created:      twTask.Entry.Time,
			Modified:     twTask.Modified.Time,
			//Resolved:     twTask.
			Due:          twTask.Due.Time,
		})
	}

	return nil
}