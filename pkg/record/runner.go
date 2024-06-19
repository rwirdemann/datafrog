package record

import (
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mysql"
)

// Runner runs the recorder for the given channel.
type Runner struct {
	testname   string
	channel    df.Channel
	channelLog df.Log
	writer     df.TestWriter
	recorder   *Recorder
	done       chan struct{}
	stopped    chan struct{}
}

// NewRunner creates a new runner for recording interactions of the given
// channel.
func NewRunner(testname string, channel df.Channel, logFactory df.LogFactory) *Runner {
	return &Runner{testname: testname, channel: channel, channelLog: logFactory.Create(channel.Log)}
}

// Start starts a new recorder as go routine.
func (r *Runner) Start() error {
	var err error
	r.writer, err = df.NewFileTestWriter(r.testname)
	if err != nil {
		return err
	}

	r.recorder = NewRecorder(r.channel, mysql.Tokenizer{}, r.channelLog, r.writer, &df.UTCTimer{}, r.testname, df.GoogleUUIDProvider{})
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

	// wait till recorder has been finished gracefully
	<-r.stopped

	// close log file and writer
	r.channelLog.Close()
	_ = r.writer.Close()
}

// Testcase returns the testcase.
func (r *Runner) Testcase() df.Testcase {
	return r.recorder.testcase
}
