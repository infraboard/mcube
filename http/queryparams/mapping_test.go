package queryparams_test

import (
	"testing"
	"time"

	"github.com/infraboard/mcube/http/queryparams"
	"github.com/stretchr/testify/assert"
)

func TestMappingDefault(t *testing.T) {
	var s struct {
		Int   int    `form:",default=9"`
		Slice []int  `form:",default=9"`
		Array [1]int `form:",default=9"`
	}
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.Int)
	assert.Equal(t, []int{9}, s.Slice)
	assert.Equal(t, [1]int{9}, s.Array)
}

func TestMappingSkipField(t *testing.T) {
	var s struct {
		A int
	}
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 0, s.A)
}

func TestMappingIgnoreField(t *testing.T) {
	var s struct {
		A int `form:"A"`
		B int `form:"-"`
	}
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"A": {"9"}, "B": {"9"}}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.B)
}

func TestMappingUnexportedField(t *testing.T) {
	var s struct {
		A int `form:"a"`
		b int `form:"b"`
	}
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"a": {"9"}, "b": {"9"}}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.b)
}

func TestMappingPrivateField(t *testing.T) {
	var s struct {
		f int `form:"field"`
	}
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"field": {"6"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, int(0), s.f)
}

func TestMappingUnknownFieldType(t *testing.T) {
	var s struct {
		U uintptr
	}

	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"U": {"unknown"}}, "form")
	assert.Error(t, err)
	assert.Equal(t, queryparams.ErrUnknownType, err)
}

func TestMappingURI(t *testing.T) {
	var s struct {
		F int `uri:"field"`
	}
	err := queryparams.MapURI(&s, map[string][]string{"field": {"6"}})
	assert.NoError(t, err)
	assert.Equal(t, int(6), s.F)
}

func TestMappingForm(t *testing.T) {
	var s struct {
		F int `json:"field"`
	}
	err := queryparams.MapFormJSON(&s, map[string][]string{"field": {"6"}})
	assert.NoError(t, err)
	assert.Equal(t, int(6), s.F)
}

func TestMapFormWithTag(t *testing.T) {
	var s struct {
		F int `externalTag:"field"`
	}
	err := queryparams.MapFormWithTag(&s, map[string][]string{"field": {"6"}}, "externalTag")
	assert.NoError(t, err)
	assert.Equal(t, int(6), s.F)
}

func TestMappingTime(t *testing.T) {
	var s struct {
		Time      time.Time
		LocalTime time.Time `time_format:"2006-01-02"`
		ZeroValue time.Time
		CSTTime   time.Time `time_format:"2006-01-02" time_location:"Asia/Shanghai"`
		UTCTime   time.Time `time_format:"2006-01-02" time_utc:"1"`
	}

	var err error
	time.Local, err = time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)

	err = queryparams.MapFormJSON(&s, map[string][]string{
		"Time":      {"2019-01-20T16:02:58Z"},
		"LocalTime": {"2019-01-20"},
		"ZeroValue": {},
		"CSTTime":   {"2019-01-20"},
		"UTCTime":   {"2019-01-20"},
	})
	assert.NoError(t, err)

	assert.Equal(t, "2019-01-20 16:02:58 +0000 UTC", s.Time.String())
	assert.Equal(t, "2019-01-20 00:00:00 +0100 CET", s.LocalTime.String())
	assert.Equal(t, "2019-01-19 23:00:00 +0000 UTC", s.LocalTime.UTC().String())
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", s.ZeroValue.String())
	assert.Equal(t, "2019-01-20 00:00:00 +0800 CST", s.CSTTime.String())
	assert.Equal(t, "2019-01-19 16:00:00 +0000 UTC", s.CSTTime.UTC().String())
	assert.Equal(t, "2019-01-20 00:00:00 +0000 UTC", s.UTCTime.String())

	// wrong location
	var wrongLoc struct {
		Time time.Time `time_location:"wrong"`
	}
	err = queryparams.MapFormJSON(&wrongLoc, map[string][]string{"Time": {"2019-01-20T16:02:58Z"}})
	assert.Error(t, err)

	// wrong time value
	var wrongTime struct {
		Time time.Time
	}
	err = queryparams.MapFormJSON(&wrongTime, map[string][]string{"Time": {"wrong"}})
	assert.Error(t, err)
}

func TestMappingTimeDuration(t *testing.T) {
	var s struct {
		D time.Duration
	}

	// ok
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"D": {"5s"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 5*time.Second, s.D)

	// error
	err = queryparams.MappingByPtr(&s, queryparams.FormSource{"D": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingSlice(t *testing.T) {
	var s struct {
		Slice []int `form:"slice,default=9"`
	}

	// default value
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{}, "form")
	assert.NoError(t, err)
	assert.Equal(t, []int{9}, s.Slice)

	// ok
	err = queryparams.MappingByPtr(&s, queryparams.FormSource{"slice": {"3", "4"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4}, s.Slice)

	// error
	err = queryparams.MappingByPtr(&s, queryparams.FormSource{"slice": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingArray(t *testing.T) {
	var s struct {
		Array [2]int `form:"array,default=9"`
	}

	// wrong default
	err := queryparams.MappingByPtr(&s, queryparams.FormSource{}, "form")
	assert.Error(t, err)

	// ok
	err = queryparams.MappingByPtr(&s, queryparams.FormSource{"array": {"3", "4"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, [2]int{3, 4}, s.Array)

	// error - not enough vals
	err = queryparams.MappingByPtr(&s, queryparams.FormSource{"array": {"3"}}, "form")
	assert.Error(t, err)

	// error - wrong value
	err = queryparams.MappingByPtr(&s, queryparams.FormSource{"array": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingStructField(t *testing.T) {
	var s struct {
		J struct {
			I int
		}
	}

	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"J": {`{"I": 9}`}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 9, s.J.I)
}

func TestMappingMapField(t *testing.T) {
	var s struct {
		M map[string]int
	}

	err := queryparams.MappingByPtr(&s, queryparams.FormSource{"M": {`{"one": 1}`}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"one": 1}, s.M)
}

func TestMappingIgnoredCircularRef(t *testing.T) {
	type S struct {
		S *S `form:"-"`
	}
	var s S

	err := queryparams.MappingByPtr(&s, queryparams.FormSource{}, "form")
	assert.NoError(t, err)
}
