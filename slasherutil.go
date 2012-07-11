package upnode

import (
	"database/sql"
	"fmt"
	"github.com/alphazero/Go-Redis"
	_ "github.com/bmizerany/pq"
	"github.com/kr/beanstalk"
	"log"
	"time"
)

type Statements struct {
	mainSql     *sql.Stmt //statement to select all checks in mode 0
	mainSqlMode *sql.Stmt //statement to select all checks in mode 1
	updateSql   *sql.Stmt //statement to call when we want to increase seq counter
}

func (s *Statements) PrepareAll(db *sql.DB) {
	var err error

	s.mainSql, err = db.Prepare(`SELECT "requestId","address","type","expectedResult","options","priority","push","seq","fallbackList" FROM "Request" WHERE paused = 0 AND $1 = ANY("slotList") AND "plantedSeed" = $2 AND interval = $3 AND push IN (0,1) `)
	if err != nil {
		log.Fatal(fmt.Sprintf("\n1. prepare statement failed: %s\n", err))
	}

	s.mainSqlMode, err = db.Prepare(`SELECT "requestId","address","type","expectedResult","options","priority","push","seq","fallbackList" FROM "Request" WHERE paused = 0 AND $1 = ANY("slotList") AND "plantedSeed" % ($2) = ($3) AND interval = ($4) AND push IN (2,3) `)
	if err != nil {
		log.Fatal(fmt.Sprintf("\n2. prepare statement failed: %s\n", err))
	}
	s.updateSql, err = db.Prepare(`UPDATE "Request" SET "seq" = "seq" + 1 WHERE "requestId" = $1 `)
	if err != nil {
		log.Fatal(fmt.Sprintf("\n3. prepare statement failed failed: %s\n", err))
	}

}

func (s Statements) GetMainSql() *sql.Stmt {
	return s.mainSql
}

func (s Statements) GetMainSqlMode() *sql.Stmt {
	return s.mainSqlMode
}

func (s Statements) GetUpdateSql() *sql.Stmt {
	return s.updateSql
}

func OfflineList(redisClient redis.Client, nodes []uint16) (list []uint16) {
	list = append(list, uint16(NODE_ID))

	for _, val := range nodes {
		if !IsOnline(redisClient, val) {
			list = append(list, val)
		}
	}

	return
}

func IsOnline(redisClient redis.Client, node uint16) (eval bool) {
	eval = false
	value, err := redisClient.Get(fmt.Sprintf("node-%d-status", node))
	if err != nil {
		log.Println("Error on get node status\n")
		return
	}

	if value != nil {
		var val int64

		num, err := fmt.Sscanf(string(value), "%d", &val)
		if err != nil {
			log.Println("Error converting value as int")
		}

		if num == 0 {
			log.Println("Value read as status is not int\n")
			return
		}

		if time.Now().Unix() <= (val + FailedOffset) {
			eval = true
		}
	}

	return
}

func GenerateCouple(stamp, interval, node *int64, mode uint16) (desc string, seed int64) {

	var slot int64 = 0
	seed = 0

	effectiveInterval := *interval * 60

	if mode > 1 {
		effectiveInterval /= CycleSize
		slot = int64((*stamp % (*interval * 60)) / effectiveInterval)
		seed = *stamp % ((*interval * 60) / CycleSize)
	} else {
		slot = int64((*stamp % (effectiveInterval * CycleSize)) / effectiveInterval)
		seed = *stamp % effectiveInterval
	}

	return fmt.Sprintf("%d-%d", slot, *node), seed
}

func IsFallbackAssigned(fallbackList []uint16, offlineList []uint16) (isIt bool) {
	isIt = false

	var diff []uint16

	for i := 0; i < len(fallbackList); i++ {
		var j int
		counter := 0
		for j = 0; j < len(offlineList); j++ {
			if fallbackList[i] != offlineList[j] {
				counter++
			}
		}
		if counter == j {
			diff = append(diff, fallbackList[i])
		}
	}

	if len(diff) > 0 && diff[0] == NODE_ID {
		isIt = true
	}

	return
}

func GetNodes(db *sql.DB) (list []uint16) {
	rows, err := db.Query(fmt.Sprintf("SELECT id FROM  \"Cluster\" WHERE active = 1 AND \"group\" = %d ", NODE_GROUP))
	if err != nil {
		log.Fatal(fmt.Sprintf("Node select query failed: \n%s\n", err))
	}

	for rows.Next() {
		var id uint16
		rows.Scan(&id)
		list = append(list, id)
	}

	return
}

func GetIntervals(db *sql.DB) (list map[uint16]uint16) {
	rows, err := db.Query(fmt.Sprintf("SELECT id,value FROM  \"Interval\" ORDER BY \"value\" "))
	if err != nil {
		log.Fatal(fmt.Sprintf("Interval select query failed: \n%s\n", err))
	}

	list = make(map[uint16]uint16)

	for rows.Next() {
		var id, value uint16
		rows.Scan(&id, &value)
		list[id] = value
	}

	return
}

func GetTasks(st *Statements, stamp *int64, nodeList []uint16, intervals map[uint16]uint16) (tasklist []Task) {

	for i := 0; i < len(nodeList); i++ {

		for intervalId, interval := range intervals {
			partialList := GetPartialTask(st, stamp, &nodeList[i], nodeList, &intervalId, &interval)
			tasklist = append(tasklist, partialList...)
		}
	}

	return

}

func GetPartialTask(st *Statements, stamp *int64, node *uint16, nodeList []uint16, intervalId *uint16, interval *uint16) (tasklist []Task) {

	var tInterval, tNode int64

	tInterval = int64(*interval)
	tNode = int64(*node)

	slot, seed := GenerateCouple(stamp, &tInterval, &tNode, 0)

	//fmt.Printf("Mode 0, Getting slot: %s, and seed: %d, and interval: %d\n", slot, seed, tInterval)

	rows, err := st.GetMainSql().Query(slot, seed, *intervalId)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not execute main query\n%s\n", err))
	}

	tasklist = append(tasklist, GetTasklistFromDb(rows, node, nodeList, stamp)...)

	slot, seed = GenerateCouple(stamp, &tInterval, &tNode, 2)

	//fmt.Printf("Mode 1, Getting slot: %s, and seed: %d, and interval: %d\n", slot, seed, tInterval)

	rows, err = st.GetMainSqlMode().Query(slot, (60 * (*interval) / CycleSize), seed, *intervalId)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not execute main mode query\n%s\n", err))
	}

	tasklist = append(tasklist, GetTasklistFromDb(rows, node, nodeList, stamp)...)

	return
}

func AssignTasks(st *Statements, beanstalkClient *beanstalk.Conn, tasklist []Task, tube string, dispOnly *bool) {
	for _, task := range tasklist {
		if *dispOnly {
			fmt.Printf("Check: %d, Stamp: %d, seq: %d\n", task.requestId, task.taskStamp, task.seq)
			continue
		}

		_, err := QueuePut(beanstalkClient, fmt.Sprintf("%s%d", tube, task.checkType), uint32(task.priority), SLASHER_DELAY, 120, []byte(task.CreateOld()))
		if err != nil {
			log.Printf("Could not put job in tube, error: %s\n", err)
			continue
		}

		_, err = st.GetUpdateSql().Exec(task.requestId)
		if err != nil {
			log.Fatalf("Trouble updating sequence counter\n")
		}
	}

}
