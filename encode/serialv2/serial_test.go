package serialv2

import "testing"
import "time"
import "log"

type TestStruct struct {
	StringField  string
	IntField     int
	BoolField    bool
	Float64Field float64
	TimeField    time.Time
}
type SubStruct struct {
	Name string
}
type TestStruct2 struct {
	StringField string
	Sub         []SubStruct
}

func TestUnmarsal(t *testing.T) {
	s := []byte("string;str1#str2#str3")
	v := &TestStruct2{}
	if err := Unmarshal(s, v, ';'); err != nil {
		t.Error(err)
	}
	if v.StringField != "string" {
		t.Errorf("string field error: %s!=%s", v.StringField, "string")
	}
	log.Printf("%v", v.Sub)

}
