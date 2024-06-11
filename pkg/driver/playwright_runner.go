package driver

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"os/exec"
	"strings"

	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/file"
)

// PlaywrightRunner runs TypeScript based playwright tests via npx (wrapped in
// scripts/playwright.sh).
type PlaywrightRunner struct {
	config df.Config
}

// NewPlaywrightRunner creates a new PlaywrightRunner using the given config.
func NewPlaywrightRunner(c df.Config) PlaywrightRunner {
	return PlaywrightRunner{config: c}
}

// Run runs testname by converting the name to its playwright format (full.json
// becomes full.spec.ts)
func (r PlaywrightRunner) Run(testname string) {
	if !r.Exists(testname) {
		log.Errorf("PlaywrightRunner: test file '%s' not found", testname)
		return
	}

	fn := r.ToPlaywright(testname)
	log.Printf("PlaywrightRunner: running test '%s' in '%s'", fn, r.config.Playwright.BaseDir)

	// run playwright test
	cmd := exec.Command("npx", "playwright", "test", fn, "--project=chromium")
	cmd.Dir = r.config.Playwright.BaseDir
	if err := cmd.Run(); err != nil {
		log.Errorf("error running command: %v", err)
	}
}

// Exists converts testname to its corresponding playwright format and returns
// true if the file exists in the playwright project test directory specified in
// config.
func (r PlaywrightRunner) Exists(testname string) bool {
	fn := r.ToPlaywright(testname)
	path := fmt.Sprintf("%s/%s/%s", r.config.Playwright.BaseDir, r.config.Playwright.TestDir, fn)
	if !file.Exists(path) {
		return false
	}
	return true
}

// ToPlaywright converts testname from datafrog to playwright format.
func (r PlaywrightRunner) ToPlaywright(testname string) string {
	return strings.Split(testname, ".")[0] + ".spec.ts"
}
