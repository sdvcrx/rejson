package rejson

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/tidwall/gjson"
)

type unmarshal struct {
	r gjson.Result
}

func newUnmarshal(jsonString string) *unmarshal {
	return &unmarshal{
		r: gjson.Parse(jsonString),
	}
}

func unmarshalResult(r gjson.Result, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("the value must be a non-nil pointer")
	}

	structValue := rv.Elem()
	numField := structValue.NumField()
	valueType := structValue.Type()

	for i := 0; i < numField; i++ {
		fieldType := valueType.Field(i)
		tag, err := parseTag(fieldType.Tag.Get(tagName))
		if err != nil {
			return fmt.Errorf("failed parseTag: %w", err)
		}

		switch tag.Type {
		case tagTypePath:
			// json path value
			val := r.Get(tag.Value)

			// set value
			setField(structValue.Field(i), val)
		case tagTypeIgnore:
			// do nothing
		case tagTypeEmpty:
			// do nothing
		case tagTypeFunc:
			callFunc(r, tag.Value, structValue)
		default:
			log.Println("unknown type: ", tag.Type)
		}
	}

	return nil
}

func callFunc(r gjson.Result, funcName string, structValue reflect.Value) {
	method := structValue.Addr().MethodByName(funcName)
	jsonResultType := reflect.TypeOf((*gjson.Result)(nil))
	if method.IsValid() && method.Type().NumIn() == 1 && method.Type().In(0) == jsonResultType {
		// TODO handle error
		method.Call([]reflect.Value{
			reflect.ValueOf(&r),
		})
	}
	// TODO print error log and return error
}

func setFieldStringOrNumber(field reflect.Value, val gjson.Result) {
	fieldType := field.Kind()
	switch fieldType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(val.Int())
	case reflect.String:
		field.SetString(val.String())
	}
}

func setFieldBool(field reflect.Value, val gjson.Result) {
	if field.Kind() == reflect.Bool {
		field.SetBool(val.Bool())
	}
}

func setFieldJSON(field reflect.Value, val gjson.Result) {
	var v reflect.Value
	if field.Type().Kind() == reflect.Ptr {
		v = reflect.New(field.Type().Elem())
	} else {
		v = reflect.New(field.Type())
	}
	if err := unmarshalResult(val, v.Interface()); err != nil {
		log.Fatalln(err)
	}

	if field.Kind() == reflect.Ptr {
		// field type is *Entity
		field.Set(v)
	} else {
		// field type is Entity
		field.Set(v.Elem())
	}
}

func setField(field reflect.Value, val gjson.Result) {
	if field.CanSet() {
		switch val.Type {
		case gjson.Number:
			setFieldStringOrNumber(field, val)
		case gjson.String:
			setFieldStringOrNumber(field, val)
		case gjson.Null:
			// not set by default
		case gjson.False:
			setFieldBool(field, val)
		case gjson.True:
			setFieldBool(field, val)
		case gjson.JSON:
			setFieldJSON(field, val)
		}
	} else {
		log.Printf("%+v cannot be setted", field)
	}
}

func (u *unmarshal) Unmarshal(v interface{}) error {
	return unmarshalResult(u.r, v)
}

func Unmarshal(jsonString string, v interface{}) error {
	u := newUnmarshal(jsonString)
	return u.Unmarshal(v)
}
