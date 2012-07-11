package upnode

import (
	"testing"
)

func TestDoHttpCheck(t *testing.T) {
	task := DummyTask()

	opts := Options{80, "default", []uint64{}, "", 3, false, false, true}

	result, err := DoHttpCheck(&task, &opts)

	if err != nil {
		t.Logf("Error doing http check: %s\n", err)
		t.Fail()
	}

	var eval string = "Fail"
	if result.success {
		eval = "Success"
	}

	t.Logf("Address: %s, %s, time: %f sec, http code : %d, eval: %s, num of redirects: %d\n", task.address, eval, (result.timeline.totalTime - result.timeline.nameLookupTime), result.httpMsg,result.msg, result.numRedirects)
}
