package driver

import (
	"fmt"
	"log"
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
// becomes full.spec.ts) and passing the playwright version to
// scripts/playwright.sh.
func (r PlaywrightRunner) Run(testname string) {
	if !r.Exists(testname) {
		log.Printf("PlaywrightRunner: test file '%s' not found", testname)
		return
	}

	fn := r.ToPlaywright(testname)
	log.Printf("driver: running test '%s'", fn)

	// change into playwright project directory
	_, err := exec.Command("cd", r.config.Playwright.BaseDir).Output()
	if err != nil {
		fmt.Printf("could not run command: %v", err)
		return
	}

	// run playwright test
	out, err := exec.Command("scripts/playwright.sh", fn).Output()
	if err != nil {
		fmt.Println("output: ", string(out))
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println("Output: ", string(out))
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
