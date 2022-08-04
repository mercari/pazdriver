package pazdriver

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorReport(t *testing.T) {
	t.Parallel()

	got := ErrorReport(runtime.Caller(0)).Interface.(*reportContext)

	assert.Contains(t, got.ReportLocation.File, "pazdriver/report_test.go")
	assert.Equal(t, "13", got.ReportLocation.Line)
	assert.Contains(t, got.ReportLocation.Function, "pazdriver.TestErrorReport")
}

func TestNewReportContext(t *testing.T) {
	t.Parallel()

	got := newReportContext(runtime.Caller(0))

	assert.Contains(t, got.ReportLocation.File, "pazdriver/report_test.go")
	assert.Equal(t, "23", got.ReportLocation.Line)
	assert.Contains(t, got.ReportLocation.Function, "pazdriver.TestNewReportContext")
}
