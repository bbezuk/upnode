package upnode

import (
	"database/sql"
	"fmt"
	"github.com/alphazero/Go-Redis"
	_ "github.com/bmizerany/pq"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const ListLen = 10

func RandomStatus(stamp int64) (string, bool) {
	seed := rand.Intn(10)

	if seed%2 == 0 {
		return fmt.Sprintf("%d", stamp), true
	}

	return "0", false
}

func TestIsOnline(t *testing.T) {
	//connecting to redis server

	spec := redis.DefaultSpec()
	rClient, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		t.Log("Could not establish redis connection")
		t.FailNow()
	}
	defer rClient.Quit()

	//getting current timestamp 

	now := time.Now().Unix()

	//filling dummy node statuses

	for i := uint16(1); i < 7; i++ {
		var val string

		if i%2 == 0 {
			val = fmt.Sprintf("%d", now)
		} else {
			val = "0"
		}

		rClient.Set(fmt.Sprintf("node-%d-status", i), []byte(val))
	}

	// evaluation of each result
	for i := uint16(1); i < 7; i++ {
		eval := IsOnline(rClient, i)
		if i%2 != 0 && eval {
			t.Logf("Iterator %d returned true but should be false", i)
			t.Fail()
		}
		if i%2 == 0 && !eval {
			t.Logf("Iterator %d returned false but should be true", i)
			t.Fail()
		}
	}

	//cleanup of dummy node statuses

	for i := uint16(1); i < 7; i++ {
		rClient.Del(fmt.Sprintf("node-%d-status", i))
	}

	t.Log("IsOnline test is done")
}

func TestOfflineList(t *testing.T) {
	var expectedList []uint16

	//connecting to redis server

	spec := redis.DefaultSpec()
	rClient, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		t.Log("Could not establish redis connection")
		t.FailNow()
	}
	defer rClient.Quit()

	// preparing seed for random results, and
	rand.Seed(time.Now().Unix())

	now := time.Now().Unix()

	// appeding own node id first
	expectedList = append(expectedList, uint16(NODE_ID))

	// generating values for redis and creating expected list
	for i := uint16(1); i < ListLen; i++ {
		val, eval := RandomStatus(now)

		if !eval {
			expectedList = append(expectedList, i)
		}

		rClient.Set(fmt.Sprintf("node-%d-status", i), []byte(val))
	}

	recievedList := OfflineList(rClient, []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9})

	expectedListString := fmt.Sprintf("%v", expectedList)
	recievedListString := fmt.Sprintf("%v", recievedList)

	if recievedListString != expectedListString {
		t.Logf("Expected list is not the same as recievedList\nExpected: %v\nRecieved: %v\n", expectedList, recievedList)
		t.Fail()
	}

	for i := uint16(1); i < ListLen; i++ {
		rClient.Del(fmt.Sprintf("node-%d-status", i))
	}

	t.Log("Offline list testing done")
}

type CoupleTest struct {
	stamp    int64
	interval int64
	node     int64
	mode     uint16
	expStr   string
	expSeed  int64
	recStr   string
	recSeed  int64
}

func TestGenerateCouple(t *testing.T) {
	test := CoupleTest{1337328172, 1, 10, 1, "2-10", 52, "", -1}
	CoupleEval(t, &test)

	t.Log("Generate couple testing done")
}

func CoupleEval(t *testing.T, test *CoupleTest) {
	test.recStr, test.recSeed = GenerateCouple(&test.stamp, &test.interval, &test.node, test.mode)

	if !strings.EqualFold(test.recStr, test.expStr) {
		t.Logf("String returned %s\texpected was -> %s\t Args=> %d, %d, %d, %d", test.expStr, test.recStr, test.stamp, test.interval, test.node, test.mode)
		t.FailNow()
	}
	if test.recSeed != test.expSeed {
		t.Logf("Seed returned %d\texpected was -> %d\t Args=> %d, %d, %d, %d", test.expSeed, test.recSeed, test.stamp, test.interval, test.node, test.mode)
		t.FailNow()
	}

}

func BenchmarkGenerateCouple(b *testing.B) {
	b.StopTimer()

	test := CoupleTest{1337328172, 1, 10, 1, "2-10", 52, "", -1}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		test.recStr, test.recSeed = GenerateCouple(&test.stamp, &test.interval, &test.node, test.mode)
	}
}

func SingleIFATest(t *testing.T, fbList []uint16, offList []uint16, prediction bool) {
	eval := IsFallbackAssigned(fbList, offList)

	if !eval && prediction {
		t.Logf("Fallback not recognized even it should be, fallbacklist is: %v, and offlineList is %v", fbList, offList)
		t.FailNow()
	}
	if eval && !prediction {
		t.Logf("Fallback recognized but it shouldn't, fallbacklist is: %v, and offlineList is %v", fbList, offList)
		t.FailNow()
	}

}

func TestIsFallbackAssigned(t *testing.T) {
	fbList := []uint16{12, NODE_ID, 11, 14}
	offList := []uint16{12, 11}

	SingleIFATest(t, fbList, offList, true)

	fbList = []uint16{12, 5, NODE_ID, 11}
	offList = []uint16{12, 11}

	SingleIFATest(t, fbList, offList, false)

	t.Log("IsFallbackAssigned testing done")

}

func BenchmarkIsFallbackAssigned(b *testing.B) {
	b.StopTimer()

	fbList := []uint16{12, 14, 10, 11}
	offList := []uint16{12, 11}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		IsFallbackAssigned(fbList, offList)
	}
}

func TestGetNodes(t *testing.T) {
	db, err := sql.Open("postgres", PostgresTest)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = CreateTableCluster(db)
	if err != nil {
		t.Log("Cannot continue with test, table 'Cluster' not created")
		t.FailNow()
	}

	_, err = db.Exec(`INSERT INTO "Cluster" VALUES 
		(10,'127.0.0.1',13339,'Bill','Clinton',4021637572967437,1,3),
		(11,'127.0.0.2',13339,'Bob','Marley',4021637572967437,0,2),
		(12,'127.0.0.3',13339,'Sterling','Archer',4021637572967437,1,2),
		(13,'127.0.0.1',13339,'Bill','Clinton',4021637572967437,0,3),
		(14,'127.0.0.2',13339,'Bob','Marley',4021637572967437,1,4),
		(15,'127.0.0.3',13339,'Sterling','Archer',4021637572967437,1,2)
	`)

	if err != nil {
		t.Logf("Dummy data fill query failed: \n%s\n", err)
		t.Fail()
	}

	nodesExpected := []uint16{12, 15}
	nodesRecieved := GetNodes(db)

	nodesExpectedString := fmt.Sprintf("%v", nodesExpected)
	nodesRecievedString := fmt.Sprintf("%v", nodesRecieved)

	if nodesExpectedString != nodesRecievedString {
		t.Logf("Nodes expected were:\t %v\nNodes recieved are:\t %v", nodesExpected, nodesRecieved)
		t.Fail()
	}

	err = DropTableCluster(db)
	if err != nil {
		t.Log("Cannot continue with test, table 'Cluster' not dropped")
		t.FailNow()
	}

	t.Log("TestGetNodes testing done")
}

func TestGetIntervals(t *testing.T) {
	db, err := sql.Open("postgres", PostgresTest)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = CreateTableIncident(db)
	if err != nil {
		t.Log("Cannot continue with test, table 'Incident' not created")
		t.FailNow()
	}

	err = InsertDefaultIntervals(db)
	if err != nil {
		t.Log("Cannot continue with test, default intervals not filled")
		t.Fail()
	}

	intervalsExpected := map[uint16]uint16{1: 1, 2: 2, 3: 5, 4: 10, 5: 20, 6: 30, 7: 60}
	intervalsRecieved := GetIntervals(db)

	var diff bool = false

	for key, _ := range intervalsExpected {
		if _, ok := intervalsRecieved[key]; !ok {
			diff = true
			break
		}
	}

	if diff {
		t.Logf("Intervals expected were: ", intervalsExpected, "\nIntervals recieved are: ", intervalsRecieved, "\n")
		t.Fail()
	}

	err = DropTableIncident(db)
	if err != nil {
		t.Log("Cannot continue with test, table 'Incident' not dropped")
		t.FailNow()
	}

	t.Log("GetIntervals() testing done")

}

func TestGetPartialTask(t *testing.T) {
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
		if i == 49 {
			continue
		}

		_, err = db.Exec(fmt.Sprintf("INSERT INTO \"Request\" VALUES(%d,%d,%d,'Address for check: %d','200 OK','options',%d,%d,%d,%d,%d,%d,%d,%s,%s,%d)", i, 1, 10, i, 10, 0, 0, 0, 0, 0, i, "ARRAY[10,11,12,14]", fmt.Sprintf("ARRAY['2-%d','1-11']", NODE_ID), 0))

		if err != nil {
			t.Log("Dummy data fill query failed: \n%s\n", err)
			t.FailNow()
		}
	}

	statements := new(Statements)

	statements.PrepareAll(db)

	var node, interval, intervalId uint16 = NODE_ID, 1, 1

	var i int64

	for i = 120; i < 180; i++ {
		testList := GetPartialTask(statements, &i, &node, []uint16{12, 13, 14}, &intervalId, &interval)

		if len(testList) > 0 {
			if testList[0].requestId != uint32(i-120) {
				t.Logf("Request Id is not same as expected, expected is %d, and recieved is %d", i-120, testList[0].requestId)
				t.Fail()
			}
		} else {
			if i != 169 {
				t.Logf("List with no members recieved, but exepcted members at index: %d", i)
				t.Fail()
			}
		}
	}

	err = DropTableRequest(db)
	if err != nil {
		t.Log("Cannot continue with test, table not dropped")
		t.FailNow()
	}

	t.Log("TestGetTaskListFromDb testing done")
}

func TestAssignTasks(t *testing.T) {
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

	statements := new(Statements)

	statements.PrepareAll(db)

	bClient, err := QueueConnect()
	if err != nil {
		t.Logf("Could not connect to beanstalk server with code: %s", err)
		t.FailNow()
	}

	defer QueueDisconnect(bClient)

	var testTube string = "test_tube_"

	task := dummyTask()
	taskList := make([]Task, 1)

	expectedString := task.CreateOld()

	taskList = append(taskList, task)

	dispOnly := false

	AssignTasks(statements, bClient, taskList, testTube, &dispOnly)

	_, body, err := QueueReserve(bClient, 120, fmt.Sprintf("%s10", testTube))
	if err != nil {
		t.Logf("Could not read job from tube: %s \n", err)
		t.Fail()
	}

	recievedString := string(body)

	if expectedString != recievedString {
		t.Logf("Recieved string from beanstalkd does not match expected string. Expected: %s, and recieved: %s", expectedString, recievedString)
		t.Fail()
	}

	err = DropTableRequest(db)
	if err != nil {
		t.Log("Cannot continue with test, table not dropped")
		t.FailNow()
	}

	t.Log("TestAssignTasks testing done")
}
