package record

import (
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	log "github.com/sirupsen/logrus"
)

// Runner runs the recorder for the given channel.
type Runner struct {
	testname   string
	channel    df.Channel
	repository df.TestRepository
	channelLog df.Log
	recorder   *Recorder
	done       chan struct{}
	stopped    chan struct{}
}

// NewRunner creates a new runner for recording interactions of the given
// channel.
func NewRunner(testname string, channel df.Channel, repository df.TestRepository, log df.Log) *Runner {
	return &Runner{testname: testname, channel: channel, repository: repository, channelLog: log}
}

// Start starts a new recorder as go routine.
func (r *Runner) Start() error {
	r.recorder = NewRecorder(r.channel, mysql.Tokenizer{}, r.channelLog, &df.UTCTimer{}, r.testname, df.GoogleUUIDProvider{}, r.repository)
	r.done = make(chan struct{})
	r.stopped = make(chan struct{})
	go r.recorder.Start(r.done, r.stopped)
	return nil
}

// Stop stops the recording by closing the done channel, that is checked by the
// recorder for its termination. Closes also the channels log file and test
// writer.
func (r *Runner) Stop() {
	// tell recorder that recording has been finished
	close(r.done)
	log.Printf("rrunner: waiting for stopped channel to be closed")

	// wait till recorder has been finished gracefully
	<-r.stopped
	log.Printf("rrunner: stopped channel closed")

	// close log file
	r.channelLog.Close()
}

// Testcase returns the testcase.
func (r *Runner) Testcase() df.Testcase {
	return r.recorder.testcase
}
