package meater

import (
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

// isNil .
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// timeTrack .
func timeTrack(name string, log *logrus.Logger) {
	start := time.Now()
	log.Debugf("%s took %v\n", name, time.Since(start))
}
