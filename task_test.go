// test case for task object

package upnode

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"strings"
	"testing"
)

func DummyTask() Task {
	var t Task
	t.requestId = 10123
	t.address = "www.example.com"
	t.checkType = 10
	t.expectedResult = "200 OK"
	t.options = "options#kwd=123"
	t.priority = 10
	t.push = 1
	t.oneTime = 0
	t.seq = 65021
	t.fallBack = make([]uint16, 2)
	t.fallBack[0] = 11
	t.fallBack[1] = 12
	t.taskStamp = 1032131231

	return t
}

func TestCreate(t *testing.T) {
	task := DummyTask()

	expected := "10123$$www.example.com$$10$$200 OK$$options#kwd=123$$10$$0$$1$$65021$$1032131231"
	created := task.Create()

	if strings.EqualFold(created, expected) != true {
		t.Logf("Create Output string was not good\nExpected\t:\t%s\nGot\t\t:\t%s", expected, created)
	}

	t.Log("TestCreate testing done")

}

func TestCreateOld(t *testing.T) {
	task := DummyTask()

	expected := "10123$$www.example.com$$10$$200 OK$$options#kwd=123$$10$$0$$1$$65021$$1032131231"
	created := task.CreateOld()

	if strings.EqualFold(created, expected) != true {
		t.Logf("CreateOld Output string was not good\nExpected\t:\t%s\nGot\t\t:\t%s", expected, created)
	}

	t.Log("TestCreateOld testing done")
}

func BenchmarkCreate(b *testing.B) {
	task := DummyTask()

	for i := 0; i < b.N; i++ {
		task.Create()
	}
}

func BenchmarkCreateOld(b *testing.B) {
	b.StopTimer()

	task := DummyTask()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		task.CreateOld()
	}
}

func TestGetTaskListFromDb(t *testing.T) {
	db, err := sql.Open("postgres", PostgresTest)
	if err != nil {
		t.Log("Cannot continue with test, database connection not established")
		t.FailNow()
	}
	defer db.Close()

	err = CreateTableRequest(db)
	if err != nil {
		t.Log("Cannot continue with test, table not created")
		t.FailNow()
	}

	for i := 0; i < 60; i++ {

		_, err = db.Exec(fmt.Sprintf("INSERT INTO \"Request\" VALUES(%d,%d,%d,'Address for check: %d','200 OK','options',%d,%d,%d,%d,%d,%d,%d,%s,%s,%d)", i, 1, 10, i, 10, 0, 0, 0, 0, 0, i, "ARRAY[10,11,12]", fmt.Sprintf("ARRAY['2-%d','1-11']", NODE_ID), 0))

		if err != nil {
			t.Log("Dummy data fill query failed: \n%s\n", err)
			t.Fail()
		}
	}

	st := new(Statements)
	st.PrepareAll(db)

	for i := 0; i < 60; i++ {
		rows, err := st.GetMainSql().Query(fmt.Sprintf("2-%d", NODE_ID), i, 1)
		if err != nil {
			t.Log("Could not execute main query\n%s\n", err)
			t.Fail()
		}

		var node uint16 = NODE_ID
		var stamp int64 = 1338975395

		testList := GetTasklistFromDb(rows, &node, []uint16{11, 12, 13}, &stamp)

		if testList[0].requestId != uint32(i) {
			t.Logf("Request with id: %d does not match expected task", i)
			t.Fail()
		}
	}

	err = DropTableRequest(db)
	if err != nil {
		t.Log("Cannot continue with test, table not dropped")
		t.FailNow()
	}

	t.Log("TestGetTaskListFromDb testing done")
}

func TestDecodePQIntArray(t *testing.T) {
	var rawBytes []byte = []byte("{12,13,14,15}")
	expectedList := []uint16{12, 13, 14, 15}

	recievedList := DecodePQIntArray(rawBytes)

	for i := 0; i < len(expectedList); i++ {
		if recievedList[i] != expectedList[i] {
			t.Logf("Recieved list is not the same on %d. member, Recieved was %d, and expected is %d", i+1, recievedList[i], expectedList[i])
			t.FailNow()
		}
	}

	t.Log("TestDecodePQIntArray testing done")

}
