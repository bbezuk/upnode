package upnode

import (
	"github.com/kr/beanstalk"
	"time"
)

func QueueConnect(protocol, address string) (*beanstalk.Conn, error) {
	conn, err := beanstalk.Dial(protocol, address)
	return conn, err
}

func QueueDisconnect(conn *beanstalk.Conn) {
	conn.Close()
}

func QueueUse(conn *beanstalk.Conn, name string) *beanstalk.Tube {
	tube := &beanstalk.Tube{conn, name}
	return tube
}

func QueuePut(conn *beanstalk.Conn, tube string, priority uint32, delay, ttr time.Duration, data []byte) (id uint64, err error) {
	t := QueueUse(conn, tube)
	id, err = t.Put(data, priority, delay, ttr)

	return
}

func QueueWatch(conn *beanstalk.Conn, name ...string) *beanstalk.TubeSet {
	return beanstalk.NewTubeSet(conn, name...)
}

func QueueReserve(conn *beanstalk.Conn, timeout time.Duration, tubes ...string) (id uint64, body []byte, err error) {
	tSet := QueueWatch(conn, tubes...)

	return tSet.Reserve(timeout)
}

func QueueDelete(conn *beanstalk.Conn,id uint64) (err error) {
	return conn.Delete(id)
}
