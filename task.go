// Task definition and functions for manipulation

package upnode

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Task struct {
	requestId      uint32   // database id of check
	address        string   // url or ip if check
	checkType      uint8    // type of check
	expectedResult string   // expected result of check
	options        string   //all non standard options for check
	priority       uint8    // priority compared to other checks
	push           uint8    // push flag for check, important to pass for result processing
	oneTime        uint8    // flag to signal if this request is one time or not, important for tube select in writing result
	seq            uint64   // cardinal number of task in relation to its check
	fallBack       []uint16 //list of nodes ordered by fallback priority 
	taskStamp      int64    // time stamp for task
}

//main function to create string that is sent to beanstalkd queue
func (t *Task) Create() string {
	return fmt.Sprintf("%d$$%s$$%d$$%s$$%s$$%d$$%d$$%d$$%d$$%d",
		t.requestId,
		t.address,
		t.checkType,
		t.expectedResult,
		t.options,
		t.priority,
		t.oneTime,
		t.push,
		t.seq,
		t.taskStamp)
}

//old protocol function to create string that is sent to beanstalkd queue
//it has extra field that use to be one_time , it is deprecated from protocol, but it has to remain until worker processes are rewritte as it is expected
func (t *Task) CreateOld() string {
	return fmt.Sprintf("%d$$%s$$%d$$%s$$%s$$%d$$%d$$%d$$%d$$%d",
		t.requestId,
		t.address,
		t.checkType,
		t.expectedResult,
		t.options,
		t.priority,
		t.oneTime,
		t.push,
		t.seq,
		t.taskStamp)
}

func GetTasklistFromDb(rows *sql.Rows, node *uint16, nodeList []uint16, stamp *int64) (tasklist []Task) {
	for rows.Next() {
		var t Task

		var temp []byte

		rows.Scan(&t.requestId, &t.address, &t.checkType, &t.expectedResult, &t.options, &t.priority, &t.push, &t.seq, &temp)

		t.oneTime = 0;

		//fmt.Printf("Task recieved: %v\n", t)

		if *node == NODE_ID || IsFallbackAssigned(t.fallBack, nodeList) {
			t.fallBack = DecodePQIntArray(temp)
			t.taskStamp = *stamp
			tasklist = append(tasklist, t)
		}
	}

	return
}

func DecodePQIntArray(raw []byte) (intList []uint16) {
	intList = nil
	nodes := strings.Split(strings.Trim(string(raw), "{}"), ",")

	for i := 0; i < len(nodes); i++ {
		temp, _ := strconv.ParseUint(nodes[i], 10, 16)

		intList = append(intList, uint16(temp))
	}

	return
}
