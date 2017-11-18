package serial

import "testing"
import "time"

type TestStruct struct {
	StringField  string
	IntField     int
	BoolField    bool
	Float64Field float64
	TimeField    time.Time
}

func TestUnmarsal(t *testing.T) {
	s := []byte("string#12#true#13.45#1510986597")
	v := &TestStruct{}
	if err := Unmarshal(s, v, '#'); err != nil {
		t.Error(err)
	}
	if v.StringField != "string" {
		t.Errorf("string field error: %s!=%s", v.StringField, "string")
	}

	if v.BoolField != true {
		t.Errorf("bool field error: %t", v.BoolField)
	}

	if v.IntField != 12 {
		t.Errorf("int field error: %d!=%d", v.IntField, 12)
	}

	if v.Float64Field != 13.45 {
		t.Errorf("float64 field error: %f!=%f", v.Float64Field, 13.45)
	}

	if v.TimeField != time.Unix(1510986597, 0) {
		t.Errorf("time field error: %v!=%v", v.TimeField, time.Unix(1510986597, 0))
	}
}
func TestMarsal(t *testing.T) {
	t.Skipf("not implemented yet")
}
