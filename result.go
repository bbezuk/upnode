package upnode

import (
	"fmt"
)

type Result struct {
	success      bool
	msg          string
	timestamp    int64
	httpMsg      int32
	effUrl       string
	timeline     Timeline
	numRedirects uint16
	head         string
	body         string
}

func NewResult() *Result {
	return &Result{success: false, timeline: *new(Timeline), numRedirects: 0, httpMsg: -1}
}

func (r *Result) QueueString(processId uint16, task *Task) string {
	successString := "1"
	if !r.success {
		successString = "0"
	}

	return fmt.Sprintf("%d$$%f$$%d$$%d$$%d$$%s$$%s$$%d$$%d$$%s$$%d",
		task.requestId,
		(r.timeline.totalTime - r.timeline.nameLookupTime),
		processId,
		r.timestamp,
		task.priority,
		successString,
		r.msg,
		task.seq,
		task.taskStamp,
		task.options,
		task.push)
}

func (r *Result) QueueStringOld(processId uint16, task *Task) string {
	successString := "1"
	if !r.success {
		successString = "0"
	}

	return fmt.Sprintf("%d$$%f$$-1$$%d$$%d$$%d$$%s$$%s$$%d$$%d$$%s$$0$$%d",
		task.requestId,
		(r.timeline.totalTime - r.timeline.nameLookupTime),
		processId,
		r.timestamp,
		task.priority,
		successString,
		r.msg,
		task.seq,
		task.taskStamp,
		task.options,
		task.push)
}
