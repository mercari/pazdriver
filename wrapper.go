package zapdriver

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// A Field can be anything that can be handled by zap.Any() https://pkg.go.dev/go.uber.org/zap#Any
type Field interface{}

type Logger struct {
	logger    *zap.Logger
	labelMap  map[string]string
	labelList []zap.Field
	fieldMap  map[string]Field
	fieldList []zap.Field
}

var defaultFields = []string{"caller", "context", "error", "message", "serviceContext", "stacktrace", "timestamp"}

func isDefaultKey(key string) bool {
	for _, k := range defaultFields {
		if k == key {
			return true
		}
	}
	return false
}

func NewLogger(serviceName string) (*Logger, error) {
	noName := false
	if len(serviceName) == 0 {
		serviceName = "unknown-service"
		noName = true
	}
	zapLogger, err := NewProductionWithCore(WrapCore(
		ReportAllErrors(true),
		ServiceName(serviceName),
	))
	if err != nil {
		return nil, err
	}
	if noName {
		err = fmt.Errorf("zapdriver.NewLogger, servicename not set")
	}
	return &Logger{logger: zapLogger, labelMap: make(map[string]string), fieldMap: make(map[string]Field)}, err
}

// NewLoggerWithKServiceName() uses the K_SERVICE environment variable, which is available by default in Cloud Function and Cloud Run
func NewLoggerWithKServiceName() (*Logger, error) {
	return NewLogger(os.Getenv("K_SERVICE"))
}

func copyLabelMap(l *Logger) map[string]string {
	newMap := make(map[string]string, len(l.labelMap))
	for key, value := range l.labelMap {
		newMap[key] = value
	}
	return newMap
}

func copyFieldMap(l *Logger) map[string]Field {
	newMap := make(map[string]Field, len(l.fieldMap))
	for key, value := range l.fieldMap {
		newMap[key] = value
	}
	return newMap
}

func copyList(origList []zap.Field) []zap.Field {
	newList := make([]zap.Field, 0, len(origList))
	for _, label := range origList {
		newList = append(newList, label)
	}
	return newList
}

func generateLabelList(labelMap map[string]string) []zap.Field {
	zapfields := make([]zap.Field, 0, len(labelMap))
	for key, value := range labelMap {
		zapfields = append(zapfields, Label(key, value))
	}
	return zapfields
}

func generateFieldList(fieldMap map[string]Field) []zap.Field {
	zapfields := make([]zap.Field, 0, len(fieldMap))
	for key, value := range fieldMap {
		zapfields = append(zapfields, zap.Any(key, value))
	}
	return zapfields
}

func (l *Logger) WithLabel(key string, value string) *Logger {
	labelMap := copyLabelMap(l)
	labelMap[key] = value
	labelList := generateLabelList(labelMap)
	return &Logger{logger: l.logger,
		labelMap:  labelMap,
		labelList: labelList,
		fieldMap:  copyFieldMap(l),
		fieldList: copyList(l.fieldList)}
}

func (l *Logger) WithLabels(labels map[string]string) *Logger {
	labelMap := copyLabelMap(l)
	for key, value := range labels {
		labelMap[key] = value
	}
	labelList := generateLabelList(labelMap)
	return &Logger{logger: l.logger,
		labelMap:  labelMap,
		labelList: labelList,
		fieldMap:  copyFieldMap(l),
		fieldList: copyList(l.fieldList)}
}

func (l *Logger) WithField(key string, value Field) *Logger {
	fieldMap := copyFieldMap(l)
	if !isDefaultKey(key) { //if it's a default key, we just don't overwrite it
		fieldMap[key] = value
	}
	fieldList := generateFieldList(fieldMap)
	return &Logger{logger: l.logger,
		labelMap:  copyLabelMap(l),
		labelList: copyList(l.labelList),
		fieldMap:  fieldMap,
		fieldList: fieldList}
}

func (l *Logger) WithFields(fields map[string]Field) *Logger {
	fieldMap := copyFieldMap(l)
	for key, value := range fields {
		if !isDefaultKey(key) { //if it's a default key, we just don't overwrite it
			fieldMap[key] = value
		}
	}
	fieldList := generateFieldList(fieldMap)
	return &Logger{logger: l.logger,
		labelMap:  copyLabelMap(l),
		labelList: copyList(l.labelList),
		fieldMap:  fieldMap,
		fieldList: fieldList}
}

func (l *Logger) Infof(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.Info(msg)
}

func (l *Logger) Errorf(err error, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.Error(msg, err)
}

func (l *Logger) Fatalf(err error, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.Fatal(msg, err)
}

func (l *Logger) Info(msg string) {
	fields := append(l.fieldList, Labels(l.labelList...))
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, err error) {
	fields := append(l.fieldList, zap.Error(err), Labels(l.labelList...))
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, err error) {
	fields := append(l.fieldList, zap.Error(err), Labels(l.labelList...))
	l.logger.Fatal(msg, fields...)
}
