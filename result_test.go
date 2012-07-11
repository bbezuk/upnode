package upnode

import (
	"testing"
)

func TestQueueString(t *testing.T) {
	task := DummyTask()

	samples := make(map[string]Result)

	samples["10123$$0.000000$$10$$1341575529$$10$$1$$OK$$65021$$1032131231$$options#kwd=123$$1"] = Result{true, "OK", 1341575529, 200, "test", *new(Timeline), 0, "", ""}

	for key, val := range samples {

		recievedString := val.QueueString(10, &task)

		if recievedString != key {
			t.Logf("Recieved string: %s => does not match expected: %s", recievedString, key)
			t.Fail()
		}

	}

}

func TestQueueStringOld(t *testing.T) {
	task := DummyTask()

	samples := make(map[string]Result)

	samples["10123$$0.000000$$-1$$10$$1341575529$$10$$1$$OK$$65021$$1032131231$$options#kwd=123$$0$$1"] = Result{true, "OK", 1341575529, 200, "test", *new(Timeline), 0, "", ""}

	for key, val := range samples {

		recievedString := val.QueueStringOld(10, &task)

		if recievedString != key {
			t.Logf("Recieved string: %s => does not match expected: %s", recievedString, key)
			t.Fail()
		}

	}

}
