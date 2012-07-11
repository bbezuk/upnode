package upnode

import (
	"github.com/nutrun/lentil"
	"testing"
)

func TestPutAndReserve(t *testing.T) {
	bClient, err := QueueConnect("tcp", "127.0.0.1:11300")
	if err != nil {
		t.Logf("Could not connect to beanstalk server with code: %s", err)
		t.FailNow()
	}

	defer QueueDisconnect(bClient)

	var testTube string = "test_tube"
	var testString string = "12345test"

	id, err := QueuePut(bClient, testTube, 10, 0, 120, []byte(testString))
	if err != nil {
		t.Logf("Could not put job in tube, error: %s\n", err)
		t.Fail()
	}

	id, body, err := QueueReserve(bClient, 120, testTube)
	if err != nil {
		t.Logf("Could not read job from tube: %s \n", err)
		t.Fail()
	}

	if testString != string(body) {
		t.Logf("Id: %d, Recieved string: %s is not equal to original string %s\n", id, body, testString)
	}

}

func BenchmarkPuttingJobs(b *testing.B) {
	b.StopTimer()

	bClient, err := QueueConnect("tcp", "127.0.0.1:11300")
	if err != nil {
		b.Logf("Could not connect to beanstalk server with code: %s", err)
		b.FailNow()
	}

	defer QueueDisconnect(bClient)

	var testString string = "{14113 http://www.musicload.de 10 200 OK ms:=http://www.musicload.de#options 10 3 98018 [14 12] 1341217746}"
	var testTube string = "test_tube"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = QueuePut(bClient, testTube, 10, 0, 120, []byte(testString))
	}
}

func BenchmarkReservingJobs(b *testing.B) {
	b.StopTimer()

	bClient, err := QueueConnect("tcp", "127.0.0.1:11300")
	if err != nil {
		b.Logf("Could not connect to beanstalk server with code: %s", err)
		b.FailNow()
	}

	defer QueueDisconnect(bClient)

	var testTube string = "test_tube"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = QueueReserve(bClient, 120, testTube)
	}

}

func BenchmarkPuttingJobsOld(b *testing.B) {
	b.StopTimer()
	client, err := lentil.Dial("127.0.0.1:11300")
	if err != nil {
		b.Logf("Beanstalkd connection error: ")
		b.FailNow()
	}

	defer client.Quit()

	var testString string = "{14113 http://www.musicload.de 10 200 OK ms:=http://www.musicload.de#options 10 3 98018 [14 12] 1341217746}"
	var testTube string = "test_tube"
	err = client.Use(testTube)
	if err != nil {
		b.Logf("Trouble with selecting correct tube in beanstalk %s\n", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.Put(10, 0, 120, []byte(testString))

	}
}

func BenchmarkReservingJobsOld(b *testing.B) {
	b.StopTimer()

	client, err := lentil.Dial("127.0.0.1:11300")
	if err != nil {
		b.Logf("Beanstalkd connection error: ")
		b.FailNow()
	}

	defer client.Quit()

	var testTube string = "test_tube"

	_, err = client.Watch(testTube)
	if err != nil {
		b.Logf("Trouble with selecting correct tube in beanstalk %s\n", err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.Reserve()
	}

}
