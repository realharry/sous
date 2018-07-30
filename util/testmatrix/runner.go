package testmatrix

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

// Runner runs tests defined in a Matrix.
type Runner struct {
	t                  *testing.T
	matrix             Matrix
	testNames          map[string]struct{}
	testNamesMu        sync.RWMutex
	testNamesPassed    map[string]struct{}
	testNamesPassedMu  sync.Mutex
	testNamesSkipped   map[string]struct{}
	testNamesSkippedMu sync.Mutex
	testNamesFailed    map[string]struct{}
	testNamesFailedMu  sync.Mutex
	parent             *Supervisor
}

func (pf *Runner) recordTestStarted(t *testing.T) {
	t.Helper()
	name := t.Name()
	pf.testNamesMu.Lock()
	defer pf.testNamesMu.Unlock()
	if _, ok := pf.testNames[name]; ok {
		t.Fatalf("duplicate test name: %q", name)
	}
	pf.testNames[name] = struct{}{}
}

// Test is a test.
type Test func(*testing.T, Context)

// FixtureFactory generates Fixtures from test and combination.
type FixtureFactory func(*testing.T, Scenario) Fixture

// Fixture Teardown func is called after each test has finished.
type Fixture interface {
	Teardown(*testing.T)
}

// Context is passed to each test case.
type Context struct {
	Scenario Scenario
	// F is the fixture returned from FixtureFactory.
	F interface{}
}

// Run is analogous to *testing.T.Run, but takes a method that includes a
// Context as well as *testing.T. Run runs the defined test with all possible
// matrix combinations in parallel.
func (pf *Runner) Run(name string, test Test) {
	for _, c := range pf.matrix.scenarios() {
		c := c
		pf.t.Run(c.String()+"/"+name, func(t *testing.T) {
			t.Parallel()
			pf.parent.wg.Add(1)
			f := pf.parent.fixtureFactory(t, c)
			defer func() {
				defer pf.parent.wg.Done()
				pf.recordTestStatus(t)
				f.Teardown(t)
			}()
			pf.recordTestStarted(t)
			test(t, Context{Scenario: c, F: f})
		})
	}
}

func (pf *Runner) recordTestStatus(t *testing.T) {
	t.Helper()
	name := t.Name()
	pf.testNamesMu.RLock()
	_, started := pf.testNames[name]
	pf.testNamesMu.RUnlock()

	statusString := "UNKNOWN"
	status := &statusString
	defer func() { rtLog("Finished running %s: %s", name, *status) }()

	if !started {
		t.Fatalf("test %q reported as finished, but not started", name)
		*status = "ERROR: Not Started"
		return
	}
	switch {
	default:
		*status = "PASSED"
		pf.testNamesPassedMu.Lock()
		pf.testNamesPassed[name] = struct{}{}
		pf.testNamesPassedMu.Unlock()
		return
	case t.Skipped():
		*status = "SKIPPED"
		pf.testNamesSkippedMu.Lock()
		pf.testNamesSkipped[name] = struct{}{}
		pf.testNamesSkippedMu.Unlock()
		return
	case t.Failed():
		*status = "FAILED"
		pf.testNamesFailedMu.Lock()
		pf.testNamesFailed[name] = struct{}{}
		pf.testNamesFailedMu.Unlock()
		return
	}
}

func testNamesSlice(m map[string]struct{}) []string {
	var s, i = make([]string, len(m)), 0
	for n := range m {
		s[i] = n
		i++
	}
	return s
}

// Quiet causes less output to be produced if set to true.
// You must not change the value of Quiet after calling Init.
var Quiet bool

func rtLog(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func (pf *Runner) printSummary() (total, passed, skipped, failed, missing []string) {
	t := pf.t
	t.Helper()
	total = testNamesSlice(pf.testNames)
	passed = testNamesSlice(pf.testNamesPassed)
	skipped = testNamesSlice(pf.testNamesSkipped)
	failed = testNamesSlice(pf.testNamesFailed)

	if !Quiet {
		summary := fmt.Sprintf("%s summary: %d failed; %d skipped; %d passed (total %d)",
			t.Name(), len(failed), len(skipped), len(passed), len(total))
		t.Log(summary)
		fmt.Fprintln(os.Stdout, summary)
	}

	missingCount := len(total) - (len(passed) + len(failed) + len(skipped))
	if missingCount != 0 {
		for t := range pf.testNamesPassed {
			delete(pf.testNames, t)
		}
		for t := range pf.testNamesSkipped {
			delete(pf.testNames, t)
		}
		for t := range pf.testNamesFailed {
			delete(pf.testNames, t)
		}
		for t := range pf.testNames {
			missing = append(missing, t)
		}
	}
	return total, passed, skipped, failed, missing
}
