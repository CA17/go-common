package sqltype

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func NewJsonNullString(s string) JsonNullString {
	return JsonNullString{
		NullString: sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}

func NewJsonNullInt64(s int64, err error) JsonNullInt64 {
	return JsonNullInt64{
		NullInt64: sql.NullInt64{
			Int64: s,
			Valid: s > 0 && err == nil,
		},
	}
}

func NewJsonNullInt32(s int32, v bool) JsonNullInt32 {
	return JsonNullInt32{
		NullInt32: sql.NullInt32{
			Int32: s,
			Valid: v,
		},
	}
}
func NewJsonNullTime(s time.Time) JsonNullTime {
	return JsonNullTime{
		NullTime: sql.NullTime{
			Time:  s,
			Valid: true,
		},
	}
}

type JsonNullString struct {
	sql.NullString
}

//
// func(v *JsonNullString) String() string {
//     return v.NullString.String
// }

func (v *JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullTime struct {
	sql.NullTime
}

func (v *JsonNullTime) String() string {
	if v.Valid {
		return v.Time.Format("2006-01-02 15:04:05")
	} else {
		return ""
	}
}

func (v *JsonNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time.String())
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v *JsonNullInt64) String() string {
	return fmt.Sprintf("%d", v.Int64)
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullInt32 struct {
	sql.NullInt32
}

func (v *JsonNullInt32) String() string {
	return fmt.Sprintf("%d", v.Int32)
}

func (v JsonNullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt32) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int32
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int32 = *x
	} else {
		v.Valid = false
	}
	return nil
}
