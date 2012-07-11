package upnode

import (
	"log"
	"testing"
	"time"
)

func TestResultWriter(t *testing.T) {
	task := DummyTask()

	testChan := make(chan *Writable)

	wQueue, err := QueueConnect("tcp", "127.0.0.1:11300")
	if err != nil {
		log.Fatalf("Could not initialize wQueue: %s", err)
	}
	defer QueueDisconnect(wQueue)

	rQueue, err := QueueConnect("tcp", "127.0.0.1:11300")
	if err != nil {
		log.Fatalf("Could not initialize rQueue: %s", err)
	}
	defer QueueDisconnect(rQueue)

	tests := make(map[int]Writable)

	tests[0] = Writable{&Result{true, "OK", 1341920367, 200, "", Timeline{0.0, 0.0, 0.0, 0.0, 0.0, 0.760, 0.0}, 1, "", ""}, &task}
	tests[1] = Writable{&Result{true, "OK", 1341920368, 200, "", Timeline{0.0, 0.0, 0.0, 0.0, 0.0, 0.770, 0.0}, 1, "", ""}, &task}

	go ResultWriter(10, testChan, wQueue, "test_")

	for _, val := range tests {
		testChan <- &val
	}

	close(testChan)

	tSet := QueueWatch(rQueue, "test_report_tube")

	for i := 0; i < len(tests); i++ {
		id, body, err := tSet.Reserve(5 * time.Second)
		if err != nil {
			t.Logf("Could not get result from queue: %s", err)
			t.Fail()
		}

		t.Logf("Id: %d, Body: %s", id, body)

		err = QueueDelete(rQueue, id)
		if err != nil {
			t.Logf("Could not delete result in queue: %s", err)
			t.Fail()
		}
	}

}
