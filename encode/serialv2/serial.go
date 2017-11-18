// Package serialv2 helps marshal/unmarshal data transfered via serial port
// Currently supports field types:
// 	- bool
//	- int
//	- string
//	- float64
//  - time.Time
package serialv2

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// used for unmarshal data
var delemiters = []rune{
	';',
	'#',
}

// Unmarshal parses the Serial-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}, delimiter rune) error {
	rv := reflect.ValueOf(v)
	mmarshal(rv, string(data), delimiter)

	return nil
}
func mmarshal(rv reflect.Value, data string, delimiter rune) error {
	parsed := strings.Split(string(data), string(delimiter))
	log.Println(rv.Elem().Kind())
	switch rv.Kind() {
	case reflect.Struct:
		if rv.Elem().NumField() != len(parsed) {
			return errors.Errorf("failed to parse: fields num mismatch: %d != %d", rv.Elem().NumField(), len(parsed))
		}
		for i := 0; i < rv.Elem().NumField(); i++ {
			rf := rv.Elem().Field(i)
			if rf.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if rf.CanSet() {
					if rf.Kind() == reflect.Slice {
						mmarshal(rf, parsed[i], '#')
					} else {
						setVal(&rf, parsed[i], rf.Kind())
					}
				}
			}
		}
	case reflect.Slice:
		for i := range parsed {
			item := reflect.New(rv.Elem().Type().Elem())
			mmarshal(item, parsed[i], '#')
			log.Printf("f %+v\n", rv)
			log.Printf("s %+v\n", rv)
			log.Printf("d %+v\n", item)
			reflect.Append(rv, item)
		}
	}
	return nil
}
func setVal(rf *reflect.Value, value string, rType reflect.Kind) error {
	switch rType {
	case reflect.Int:
		x, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse field")
		}
		if !rf.OverflowInt(x) {
			rf.SetInt(x)
		}
	case reflect.String:
		rf.SetString(value)
	case reflect.Float64:
		x, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse field")
		}
		if !rf.OverflowFloat(x) {
			rf.SetFloat(x)
		}
	case reflect.Bool:
		if value == "true" {
			rf.SetBool(true)
		} else if value == "false" {
			rf.SetBool(false)
		} else {
			return errors.Errorf("failed to convert bool, got:%s", value)
		}
	case reflect.Struct:
		switch rf.Type() {
		case reflect.TypeOf(time.Time{}):
			x, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				return errors.Wrap(err, "failed to parse field")
			}
			rf.Set(reflect.ValueOf(time.Unix(x, 0)))
		}
	}
	return nil
}

// Marshal returns the Serial encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return []byte{}, errors.Errorf("not implemented yet")
}
