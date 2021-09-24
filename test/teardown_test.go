package test

import (
	"testing"

	"niktj777/teardown/teardown"

	"github.com/stretchr/testify/require"
)

func TestTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddTeardown("", func() {
		tdcounter++
		require.Equal(t, 3, tdcounter)
	})

	teardown.AddTeardown("", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AddTeardown("", func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.Teardown("")

	require.Equal(t, 3, tdcounter)

	teardown.VerifyTeardown(t)
}

func TestNamedTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddTeardown("", func() {
		tdcounter++
		require.Equal(t, 3, tdcounter)
	})

	teardown.AddTeardown("", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AddTeardown("", func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.AddTeardown("other", func() {
		tdcounter++
		require.Equal(t, 5, tdcounter)
	})

	teardown.AddTeardown("other", func() {
		tdcounter++
		require.Equal(t, 4, tdcounter)
	})

	teardown.Teardown("")
	require.Equal(t, 3, tdcounter)

	teardown.Teardown("other")
	require.Equal(t, 5, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that an unconditional  DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestUnconditionalDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", true, func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 3, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.Teardown("name")

	require.Equal(t, 3, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that an unconditionally skipped DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestSkippedDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", false, func() {
		require.FailNow(t, "This diagnostic teardown should not have been run")
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.Teardown("name")

	require.Equal(t, 2, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that an unconditional  DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestUnconditionalFuncDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", func() bool { return true }, func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 3, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.Teardown("name")

	require.Equal(t, 3, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that an unconditional  DiagnosticTeardown is *always* executed *before* any other teardown for the same name */
func TestSkippedFuncDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", func() bool { return false }, func() {
		require.FailNow(t, "This diagnostic teardown should not have been run")
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.Teardown("name")

	require.Equal(t, 2, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *always* executed *before* any other teardown for the same name if the passed testing.T has Failed */
func TestTFailedDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	tt := new(testing.T)

	teardown.AddDiagnosticTeardown("name", tt, func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 3, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	tt.Fail()

	teardown.Teardown("name")

	require.Equal(t, 3, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *not* executed if the passed in T has *not* Failed */
func TestTNotFailedDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", t, func() {
		require.FailNow(t, "This diagnostic teardown should not have been called since T has not failed")
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.Teardown("name")

	require.Equal(t, 2, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *not* executed if the conditional is nil */
func TestNilDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", nil, func() {
		require.FailNow(t, "Diagnostic teardown should not be run if the conditional is nil")
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.Teardown("name")

	require.Equal(t, 2, tdcounter)

	teardown.VerifyTeardown(t)
}

/* verify that a DiagnosticTeardown is *always* executed if the ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS is true */
func TestUnconditionalEnvVarDiagnosticTeardown(t *testing.T) {
	tdcounter := 0

	teardown.AddDiagnosticTeardown("name", false, func() {
		tdcounter++
		require.Equal(t, 1, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 3, tdcounter)
	})

	teardown.AddTeardown("name", func() {
		tdcounter++
		require.Equal(t, 2, tdcounter)
	})

	teardown.AlwaysRunDiagnosticTeardowns = true
	defer func() { teardown.AlwaysRunDiagnosticTeardowns = false }()

	teardown.Teardown("name")

	require.Equal(t, 3, tdcounter)

	teardown.VerifyTeardown(t)
}
