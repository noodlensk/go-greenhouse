// Package serial helps marshal/unmarshal data transfered via serial port
// Currently supports field types:
// 	- bool
//	- int
//	- string
//	- float64
//  - time.Time
package serial

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Unmarshal parses the Serial-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}, delimiter rune) error {
	parsed := strings.Split(string(data), string(delimiter))

	rv := reflect.ValueOf(v).Elem()
	if rv.NumField() != len(parsed) {
		return errors.Errorf("failed to parse: fields num mismatch: %d != %d", rv.NumField(), len(parsed))
	}
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Field(i)
		if rf.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if rf.CanSet() {
				switch rf.Kind() {
				case reflect.Int:
					x, err := strconv.ParseInt(parsed[i], 0, 64)
					if err != nil {
						return errors.Wrap(err, "failed to parse field")
					}
					if !rf.OverflowInt(x) {
						rf.SetInt(x)
					}
				case reflect.String:
					rf.SetString(parsed[i])
				case reflect.Float64:
					x, err := strconv.ParseFloat(parsed[i], 64)
					if err != nil {
						return errors.Wrap(err, "failed to parse field")
					}
					if !rf.OverflowFloat(x) {
						rf.SetFloat(x)
					}
				case reflect.Bool:
					if parsed[i] == "true" {
						rf.SetBool(true)
					} else if parsed[i] == "false" {
						rf.SetBool(false)
					} else {
						return errors.Errorf("failed to convert bool, got:%s", parsed[i])
					}
				case reflect.Struct:
					switch rf.Type() {
					case reflect.TypeOf(time.Time{}):
						x, err := strconv.ParseInt(parsed[i], 0, 64)
						if err != nil {
							return errors.Wrap(err, "failed to parse field")
						}
						rf.Set(reflect.ValueOf(time.Unix(x, 0)))
					}
				}
			}
		}
	}
	return nil
}
