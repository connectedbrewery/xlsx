package xlsx

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gobuffalo/nulls"
)

// WriteSlice writes an array to row r. Accepts a pointer to array type 'e',
// and writes the number of columns to write, 'cols'. If 'cols' is < 0,
// the entire array will be written if possible. Returns -1 if the 'e'
// doesn't point to an array, otherwise the number of columns written.
func (r *Row) WriteSlice(e interface{}, cols int) int {
	if cols == 0 {
		return cols
	}

	// make sure 'e' is a Ptr to Slice
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Ptr {
		return -1
	}

	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return -1
	}

	// it's a slice, so open up its values
	n := v.Len()
	if cols < n && cols > 0 {
		n = cols
	}

	var setCell func(reflect.Value)
	setCell = func(val reflect.Value) {
		switch t := val.Interface().(type) {
		case time.Time:
			cell := r.AddCell()
			cell.SetValue(t)
		case fmt.Stringer: // check Stringer first
			cell := r.AddCell()
			cell.SetString(t.String())
		case sql.NullString: // check null sql types nulls = ''
			cell := r.AddCell()
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.String)
			}
		case sql.NullBool:
			cell := r.AddCell()
			if cell.SetString(``); t.Valid {
				cell.SetBool(t.Bool)
			}
		case sql.NullInt64:
			cell := r.AddCell()
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Int64)
			}
		case sql.NullFloat64:
			cell := r.AddCell()
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Float64)
			}
		default:
			switch val.Kind() { // underlying type of slice
			case reflect.String, reflect.Int, reflect.Int8,
				reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float64, reflect.Float32:
				cell := r.AddCell()
				cell.SetValue(val.Interface())
			case reflect.Bool:
				cell := r.AddCell()
				cell.SetBool(t.(bool))
			case reflect.Interface:
				setCell(reflect.ValueOf(t))
			}
		}
	}

	var i int
	for i = 0; i < n; i++ {
		setCell(v.Index(i))
	}
	return i
}

// WriteStruct writes a struct to row r. Accepts a pointer to struct type
// 'e', and the number of columns to write, `cols`. If 'cols' is < 0,
// the entire struct will be written if possible. Returns -1 if the 'e'
// doesn't point to a struct, otherwise the number of columns written
func (r *Row) WriteStruct(e interface{}, cols int) (int, error) {
	if cols == 0 {
		return cols, nil
	}

	v := reflect.ValueOf(e).Elem()
	if v.Kind() != reflect.Struct {
		return 0, errNotStructPointer
	}

	n := v.NumField() // number of fields in struct
	if cols < n && cols > 0 {
		n = cols
	}

	var k int
	for i := 0; i < n; i, k = i+1, k+1 {
		field := v.Type().Field(i)
		idx := field.Tag.Get("xlsx")

		if idx == "-" {
			k-- // nothing set to reset to previous
			continue
		}

		pos, err := strconv.Atoi(idx)
		if err != nil {
			return 0, errInvalidTag
		}

		f := v.Field(i)
		switch t := f.Interface().(type) {
		case time.Time:
			cell := r.GetCell(pos)
			cell.SetValue(t)
		case fmt.Stringer: // check Stringer first
			cell := r.GetCell(pos)
			cell.SetString(t.String())
		case sql.NullString: // check null sql types nulls = ''
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.String)
			}
		case sql.NullBool:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetBool(t.Bool)
			}
		case sql.NullInt64:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Int64)
			}
		case sql.NullFloat64:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Float64)
			}
		case nulls.String: // check null sql types nulls = ''
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.String)
			}
		case nulls.Bool:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetBool(t.Bool)
			}
		case nulls.Int:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Int)
			}
		case nulls.Int64:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Int64)
			}
		case nulls.Float64:
			cell := r.GetCell(pos)
			if cell.SetString(``); t.Valid {
				cell.SetValue(t.Float64)
			}
		default:
			switch f.Kind() {
			case reflect.String, reflect.Int, reflect.Int8,
				reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float64, reflect.Float32:
				cell := r.GetCell(pos)
				cell.SetValue(f.Interface())
			case reflect.Bool:
				cell := r.GetCell(pos)
				cell.SetBool(t.(bool))
			default:
				k-- // nothing set so reset to previous
			}
		}
	}

	return k, nil
}
