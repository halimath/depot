package depot

import (
	"testing"
	"time"
)

type testDef struct {
	input    interface{}
	ok       bool
	expected interface{}
}

type getter func(v *Values, key string) (interface{}, bool)

func TestValuesGetString(t *testing.T) {
	tests := []testDef{
		{"foo", true, "foo"},
		{[]byte("foo"), true, "foo"},
		{17, false, ""},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetString(key)
	})
}

func TestValuesGetBytes(t *testing.T) {
	tests := []testDef{
		{"foo", false, []byte{}},
		{[]byte("foo"), true, []byte("foo")},
		{17, false, []byte{}},
		{nil, true, []byte{}},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetBytes(key)
	})
}

func TestValuesGetTime(t *testing.T) {
	now := time.Now()
	tests := []testDef{
		{now, true, now},
		{17, false, time.Time{}},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetTime(key)
	})
}

func TestValuesGetBool(t *testing.T) {
	tests := []testDef{
		{"foo", false, false},
		{true, true, true},
		{false, true, false},
		{int64(1), true, true},
		{int64(0), true, false},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetBool(key)
	})
}

func TestValuesGetInt(t *testing.T) {
	tests := []testDef{
		{int64(17), true, 17},
		{17.2, false, 0},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetInt(key)
	})
}
func TestValuesGetInt8(t *testing.T) {
	tests := []testDef{
		{int64(17), true, int8(17)},
		{17.2, false, int8(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetInt8(key)
	})
}
func TestValuesGetInt16(t *testing.T) {
	tests := []testDef{
		{int64(17), true, int16(17)},
		{17.2, false, int16(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetInt16(key)
	})
}

func TestValuesGetInt32(t *testing.T) {
	tests := []testDef{
		{int64(17), true, int32(17)},
		{17.2, false, int32(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetInt32(key)
	})
}

func TestValuesGetInt64(t *testing.T) {
	tests := []testDef{
		{int64(17), true, int64(17)},
		{17.2, false, int64(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetInt64(key)
	})
}

func TestValuesGetFloat32(t *testing.T) {
	tests := []testDef{
		{float64(17.2), true, float32(17.2)},
		{17, false, float32(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetFloat32(key)
	})
}

func TestValuesGetFloat64(t *testing.T) {
	tests := []testDef{
		{float64(17.2), true, float64(17.2)},
		{17, false, float64(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetFloat64(key)
	})
}

func TestValuesGetUInt(t *testing.T) {
	tests := []testDef{
		{int64(17), true, uint(17)},
		{17.2, false, uint(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetUInt(key)
	})
}

func TestValuesGetUInt8(t *testing.T) {
	tests := []testDef{
		{int64(17), true, uint8(17)},
		{17.2, false, uint8(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetUInt8(key)
	})
}
func TestValuesGetUInt16(t *testing.T) {
	tests := []testDef{
		{int64(17), true, uint16(17)},
		{17.2, false, uint16(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetUInt16(key)
	})
}

func TestValuesGetUInt32(t *testing.T) {
	tests := []testDef{
		{int64(17), true, uint32(17)},
		{17.2, false, uint32(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetUInt32(key)
	})
}

func TestValuesGetUInt64(t *testing.T) {
	tests := []testDef{
		{int64(17), true, uint64(17)},
		{17.2, false, uint64(0)},
	}

	runTests(t, tests, func(v *Values, key string) (interface{}, bool) {
		return v.GetUInt64(key)
	})
}

func runTests(t *testing.T, tests []testDef, getter getter) {
	vals := Values{}
	_, ok := getter(&vals, "key")
	if ok {
		t.Errorf("<missing key>: expected to get not-ok but got ok")
	}

	for _, test := range tests {
		vals := Values{
			"key": test.input,
		}
		actual, ok := getter(&vals, "key")
		if ok != test.ok {
			t.Errorf("%#v: expected %v but got %v", test.input, test.ok, ok)
		}

		if _, ok := test.expected.([]byte); !ok {
			if actual != test.expected {
				t.Errorf("%#v: expected %v but got %v", test.input, test.expected, actual)
			}
		}
	}
}
