package verify

import (
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	log "github.com/sirupsen/logrus"
)

// Runner runs the verifier for the given channel.
type Runner struct {
	testname   string
	channel    df.Channel
	config     df.Config
	channelLog df.Log
	repository df.TestRepository
	verifier   *Verifier
	done       chan struct{}
	stopped    chan struct{}
}

// NewRunner creates a new runner for verifying interactions of the given
// channel.
func NewRunner(testname string, channel df.Channel, config df.Config, log df.Log, repository df.TestRepository) *Runner {
	return &Runner{testname: testname, channel: channel, config: config, channelLog: log, repository: repository}
}

// Start starts a new verifier as go routine.
func (r *Runner) Start() error {
	tc, err := r.repository.Get(r.testname)
	if err != nil {
		return nil
	}

	r.verifier = NewVerifier(r.config, r.channel, r.repository, mysql.Tokenizer{}, r.channelLog, tc, &df.UTCTimer{}, r.testname)
	r.done = make(chan struct{})
	r.stopped = make(chan struct{})
	go r.verifier.Start(r.done, r.stopped)
	return nil
}

// Stop stops the verification by closing the done channel, that is checked by the
// verifier for its termination. Closes also the channels log file and test
// writer.
func (r *Runner) Stop() error {
	// tell verifier that verification has been finished
	close(r.done)
	log.Printf("vrunner: waiting for stopped channel to be closed")

	// wait till verifier has been finished gracefully
	<-r.stopped
	log.Printf("vrunner: stopped channel closed")

	// close log file
	r.channelLog.Close()

	return nil
}

// Testcase returns the testcase.
func (r *Runner) Testcase() df.Testcase {
	return r.verifier.testcase
}
