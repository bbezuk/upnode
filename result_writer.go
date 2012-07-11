package upnode

import (
	"fmt"
	"github.com/kr/beanstalk"
	"log"
)

type Writable struct {
	result *Result
	task   *Task
}

func ResultWriter(processId uint16, inbox chan *Writable, conn *beanstalk.Conn, tubePrefix string) {
	for w := range inbox {
		fmt.Printf("Got task from tube, writing\n")
		var tube string

		if w.task.oneTime == 1 || w.task.push > 0 || !w.result.success {
			tube = fmt.Sprintf("%sreport_tube", tubePrefix)
		} else {
			tube = fmt.Sprintf("%sresult_tube", tubePrefix)
		}

		_, err := QueuePut(conn, tube, uint32(w.task.priority), 0, 10, []byte(w.result.QueueStringOld(processId, w.task)))
		if err != nil {
			log.Printf("Could not put result in tube: %s\n result: %v\n error: %s\n", tube, w.result, err)
		}
	}

	fmt.Printf("Exiting writer\n")
}
