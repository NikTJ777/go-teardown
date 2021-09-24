package teardown

import (
	"os"
	"strings"
	"testing"
)

/** Lists of the teardown and diagnostic teardown funcs */
var teardownLists = make(map[string][]func())
var diagnosticTeardownLists = make(map[string][]func())

const ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS = "ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS"

/** Exported var - initialised from the EnvVar, but can be reset in code if desired */
var AlwaysRunDiagnosticTeardowns = strings.EqualFold(os.Getenv(ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS), "true")

/**
 * add a teardown function to the named list - for deferred execution.
 *
 * The teardown functions are called in reverse order of insertion, by a call to Teardown(name).
 *
 * The typical idiom is:
 * <pre>
 *   teardown.AddTeardown("DATABASE", func() { ...})
 *   // possibly more teardown.AddTeardown("DATABASE", func() { ... })
 *   defer teardown.Teardown("DATABASE")
 * <pre>
 */
func AddTeardown(name string, teardownFunc func()) {
	teardownLists[name] = append(teardownLists[name], teardownFunc)
}

/**
 * Adds a teardown function to all named teardown lists - for deferred execution.
 *
 * The teardown functions are called in reverse order of insertion, by a call to Teardown(name).
 *
 */
func AddGlobalTeardown(teardownFunc func()) {
	for name := range teardownLists {
		AddTeardown(name, teardownFunc)
	}
}

/**
 * add a diagnostic teardown func to be called before any other teardowns in the named list - to aid diagnostics/debugging.
 * This allows a diagnostic teardown to do such things as:
 * <ul>
 *   <li>Generate logging and debug information immediately prior to resurce teardown
 *   <li>call time.Sleep() to allow inspection and/or debugging of the exit state before teardown.
 * <ul>
 *
 * NOTE: it is generally undesirable to add multiple diagnostic teardowns that sleep - so it would usually be best to
 * add any Sleep() debug teardown to the innermost teardown list.
 * Nonetheless, there are use-cases where multiple Sleep() teardowns are useful - to allow inspecting different
 * intermediate states.
 */
func AddDiagnosticTeardown(name string, condition interface{}, teardownFunc func()) {

	// the test for whether to run the diagnostic teardown must be executed at TEARDOWN time, not at DEFER time.
	// So, create a wrapper func that has the logic to determine whether to run the teardown func, and calls it conditionally.
	tdfunc := func() {
		shouldIdoIt := AlwaysRunDiagnosticTeardowns

		if !shouldIdoIt {
			switch c := condition.(type) {
			case *testing.T:
				shouldIdoIt = c.Failed()

			case func() bool:
				shouldIdoIt = c()

			case bool:
				shouldIdoIt = c

			default:
				shouldIdoIt = c != nil
			}
		}

		if shouldIdoIt {
			teardownFunc()
		}
	}

	// add the wrapper func to the diagnosticTeardown map
	diagnosticTeardownLists[name] = append(diagnosticTeardownLists[name], tdfunc)
}

/**
 * Adds a diagnostic teardown function to all named diagnostic teardown lists - for deferred execution.
 *
 * The teardown functions are called in reverse order of insertion, by a call to Teardown(name).
 *
 */
func AddGlobalDiagnosticTeardown(condition interface{}, teardownFunc func()) {
	for name := range diagnosticTeardownLists {
		AddDiagnosticTeardown(name, condition, teardownFunc)
	}
}

/**
 * Call the stored teardown functions in the named list, in the correct order (last-in-first-out)
 *
 * NOTE: Any DIAGNOSTIC teardowns - those added with AddDiagnosticTeardown() for this name - are called BEFORE any other teardowns for this name.
 *
 * The typical use of Teardown is with a deferred call:
 * defer teardown.Teardown("SOME NAME")
 * See: teardown.AddTeardown(); teardown.AddDiagnosticTeardown()
 */
func Teardown(name string) {
	// ensure both list and diagnostic list are removed.
	defer func() { delete(diagnosticTeardownLists, name) }()
	defer func() { delete(teardownLists, name) }()

	list := teardownLists[name]
	list = append(list, diagnosticTeardownLists[name]...) // append any diagnostic funcs - so they are called FIRST

	for x := len(list) - 1; x >= 0; x-- {
		list[x]()
	}
}

/**
* Verify all teardownLists have been executed already; and throw an error if not.
* Can be used to verify correct coding of a test that uses teardown - and to ensure eventual release of resources.
*
* NOTE: while the funcs are called in the correct order for each list,
* there can be NO guarantee that the lists are iterated in the correct order.
*
* This function MUST NOT be used as a replacement for calling teardown() at the correct point in the code.
 */
func VerifyTeardown(t *testing.T) {

	// ensure all funcs in all lists are released
	defer func() { teardownLists = make(map[string][]func()) }()
	defer func() { diagnosticTeardownLists = make(map[string][]func()) }()

	// append each diagnostic list to the corresponding (possibly empty) teardown list
	for name, list := range diagnosticTeardownLists {
		teardownLists[name] = append(teardownLists[name], list...)
	}

	// release all remaining resources - this is a "best effort" as the order of iterating the map is arbitrary
	uncleared := make([]string, 0)

	// make a "best-effort" at releasing all remaining resources
	for name, list := range teardownLists {
		uncleared = append(uncleared, name)

		for x := len(list) - 1; x >= 0; x-- {
			list[x]()
		}
	}

	if len(uncleared) > 0 && t != nil {
		t.Fatalf("Error - %d teardownLists were left uncleared: %s", len(uncleared), uncleared)
	}
}
