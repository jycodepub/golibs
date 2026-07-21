package sqldb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// MapRows maps all remaining rows in rows to a slice of T.
// T must be a struct type or a pointer to a struct type.
// It does NOT close the rows; the caller is responsible for closing them.
func MapRows[T any](rows *sql.Rows) ([]T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	isPtr := false
	if t.Kind() == reflect.Pointer {
		isPtr = true
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("target type must be a struct or pointer to struct, got %s", t.Kind())
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	fieldMap := make(map[string][]int)
	buildFieldMap(t, nil, fieldMap)

	var result []T
	for rows.Next() {
		structVal := reflect.New(t).Elem()
		scanArgs := make([]any, len(columns))
		for i, col := range columns {
			if indexPath, ok := fieldMap[strings.ToLower(col)]; ok {
				scanArgs[i] = getFieldRef(structVal, indexPath).Addr().Interface()
			} else {
				scanArgs[i] = new(any)
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		var item T
		if isPtr {
			item = structVal.Addr().Interface().(T)
		} else {
			item = structVal.Interface().(T)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// MapRow maps the current row in rows to T.
// T must be a struct type or a pointer to a struct type.
// The caller must call rows.Next() before calling MapRow.
func MapRow[T any](rows *sql.Rows) (T, error) {
	var zero T
	t := reflect.TypeOf((*T)(nil)).Elem()
	isPtr := false
	if t.Kind() == reflect.Pointer {
		isPtr = true
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return zero, fmt.Errorf("target type must be a struct or pointer to struct, got %s", t.Kind())
	}

	structVal := reflect.New(t).Elem()
	columns, err := rows.Columns()
	if err != nil {
		return zero, err
	}

	fieldMap := make(map[string][]int)
	buildFieldMap(t, nil, fieldMap)

	scanArgs := make([]any, len(columns))
	for i, col := range columns {
		if indexPath, ok := fieldMap[strings.ToLower(col)]; ok {
			scanArgs[i] = getFieldRef(structVal, indexPath).Addr().Interface()
		} else {
			scanArgs[i] = new(any)
		}
	}

	if err := rows.Scan(scanArgs...); err != nil {
		return zero, err
	}

	if isPtr {
		return structVal.Addr().Interface().(T), nil
	}
	return structVal.Interface().(T), nil
}

// QueryRows executes a query and maps the resulting rows to a slice of T.
func QueryRows[T any](c *Client, query string, args ...any) ([]T, error) {
	rows, err := c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return MapRows[T](rows)
}

// QueryRow executes a query and maps the first resulting row to T.
// It returns sql.ErrNoRows if no rows were returned.
func QueryRow[T any](c *Client, query string, args ...any) (T, error) {
	var zero T
	rows, err := c.Query(query, args...)
	if err != nil {
		return zero, err
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return zero, err
		}
		return zero, sql.ErrNoRows
	}
	return MapRow[T](rows)
}

// buildFieldMap recursively maps struct field names/tags to their reflect index paths.
func buildFieldMap(t reflect.Type, indexPrefix []int, m map[string][]int) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// Skip unexported fields unless they are embedded/anonymous (which we traverse)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}

		indexPath := append(indexPrefix, i)

		// If it's an anonymous (embedded) field, traverse it if it's a struct (or pointer to struct)
		if f.Anonymous {
			ft := f.Type
			if ft.Kind() == reflect.Pointer {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				buildFieldMap(ft, indexPath, m)
				continue
			}
		}

		// Check db tag
		tag := f.Tag.Get("db")
		if tag == "-" {
			continue
		}
		if tag != "" {
			parts := strings.Split(tag, ",")
			name := strings.TrimSpace(parts[0])
			if name != "" {
				m[strings.ToLower(name)] = indexPath
			}
		} else {
			m[strings.ToLower(f.Name)] = indexPath
		}
	}
}

// getFieldRef retrieves a field's reflect.Value by its index path,
// initializing any intermediate nil pointers if they are embedded.
func getFieldRef(v reflect.Value, indexPath []int) reflect.Value {
	curr := v
	for _, idx := range indexPath {
		if curr.Kind() == reflect.Pointer {
			if curr.IsNil() {
				curr.Set(reflect.New(curr.Type().Elem()))
			}
			curr = curr.Elem()
		}
		curr = curr.Field(idx)
	}
	return curr
}
