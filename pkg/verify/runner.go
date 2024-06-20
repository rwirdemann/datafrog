package verify

import (
	"fmt"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

// Runner runs the verifier for the given channel.
type Runner struct {
	testname   string
	channel    df.Channel
	config     df.Config
	channelLog df.Log
	repositoy  df.TestRepository
	writer     df.TestWriter
	verifier   *Verifier
	done       chan struct{}
	stopped    chan struct{}
}

// NewRunner creates a new runner for verifying interactions of the given
// channel.
func NewRunner(testname string, channel df.Channel, config df.Config, logFactory df.LogFactory, repository df.TestRepository) *Runner {
	return &Runner{testname: testname, channel: channel, config: config, channelLog: logFactory.Create(channel.Log), repositoy: repository}
}

// Start starts a new verifier as go routine.
func (r *Runner) Start() error {
	tc, err := r.repositoy.Get(r.testname)
	if err != nil {
		return nil
	}

	// create the test writer
	r.writer, err = df.NewFileTestWriter(fmt.Sprintf("%s.json.running", r.testname))
	if err != nil {
		return err
	}

	r.verifier = NewVerifier(r.config, r.channel, mysql.Tokenizer{}, r.channelLog, tc, r.writer, &df.UTCTimer{}, r.testname)
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

	// close log file and writer
	r.channelLog.Close()
	_ = r.writer.Close()

	// copy .running testfile to original file
	if err := copyFile(fmt.Sprintf("%s.json.running", r.testname), fmt.Sprintf("%s.json", r.testname)); err != nil {
		return err
	}

	// delete .running file
	if err := deleteFile(fmt.Sprintf("%s.json.running", r.testname)); err != nil {
		return err
	}

	return nil
}

// Testcase returns the testcase.
func (r *Runner) Testcase() df.Testcase {
	return r.verifier.testcase
}

var mutex = &sync.Mutex{}

func deleteFile(testname string) error {
	mutex.Lock()
	defer mutex.Unlock()
	log.Printf("vrunner: deleting test file %s", testname)
	return os.Remove(testname)
}

func copyFile(src string, dst string) error {
	mutex.Lock()
	defer mutex.Unlock()
	log.Printf("vrunner: copying %s to %s", src, dst)
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
	_, err = io.Copy(destination, source)
	return err
}
