package depot

import "time"

// Values contains the persistent column values for an entity either after reading
// the values from the database to re-create the entity value or to persist the
// entity's values in the database (either for insertion or update).
type Values map[string]interface{}

func (v Values) GetTime(key string) (time.Time, bool) {
	val, ok := v[key]
	if !ok {
		return time.Time{}, false
	}

	if val == nil {
		return time.Time{}, true
	}

	switch x := val.(type) {
	case time.Time:
		return x, ok
	default:
		return time.Time{}, false
	}
}

func (v Values) GetBytes(key string) ([]byte, bool) {
	val, ok := v[key]
	if !ok {
		return nil, false
	}

	if val == nil {
		return nil, true
	}

	switch x := val.(type) {
	case []byte:
		return x, ok
	default:
		return nil, false
	}
}

func (v Values) GetBool(key string) (bool, bool) {
	val, ok := v[key]
	if !ok {
		return false, false
	}

	switch x := val.(type) {
	case bool:
		return x, ok
	default:
		return false, false
	}
}

func (v Values) GetFloat32(key string) (float32, bool) {
	val, ok := v[key]
	if !ok {
		return 0.0, false
	}

	switch x := val.(type) {
	case float64:
		return float32(x), ok
	default:
		return 0.0, false
	}
}

func (v Values) GetFloat64(key string) (float64, bool) {
	val, ok := v[key]
	if !ok {
		return 0.0, false
	}

	switch x := val.(type) {
	case float64:
		return x, ok
	default:
		return 0.0, false
	}
}

func (v Values) GetInt(key string) (int, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return int(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetInt8(key string) (int8, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return int8(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetInt16(key string) (int16, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return int16(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetInt32(key string) (int32, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return int32(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetInt64(key string) (int64, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return x, ok
	default:
		return 0, false
	}
}

func (v Values) GetUInt(key string) (uint, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return uint(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetUInt8(key string) (uint8, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return uint8(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetUInt16(key string) (uint16, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return uint16(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetUInt32(key string) (uint32, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return uint32(x), ok
	default:
		return 0, false
	}
}

func (v Values) GetUInt64(key string) (uint64, bool) {
	val, ok := v[key]
	if !ok {
		return 0, false
	}

	switch x := val.(type) {
	case int64:
		return uint64(x), ok
	default:
		return 0, false
	}
}

// GetString returns the names value converted to a string.
func (v Values) GetString(key string) (string, bool) {
	val, ok := v[key]
	if !ok {
		return "", false
	}

	switch x := val.(type) {
	case string:
		return x, ok
	case []byte:
		return string(x), ok
	default:
		return "", false
	}
}
