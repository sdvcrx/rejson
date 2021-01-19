package rejson

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

var (
	ErrFieldCannotSet    = errors.New("field cannot set")
	ErrUnexpectJSONValue = errors.New("Unexpect JSON value")
	ErrUnknownFieldType  = errors.New("Unknown field type")
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
		tag := parseTag(fieldType.Tag.Get(tagName))

		switch tag.Type {
		case tagTypePath:
			// json path value
			val := r.Get(tag.Value)

			// set value
			if err := setField(structValue.Field(i), val); err != nil {
				return err
			}
		case tagTypeIgnore:
			// do nothing
		case tagTypeEmpty:
			// do nothing
		case tagTypeFunc:
			callFunc(r, tag.Value, structValue)
		default:
			return fmt.Errorf("%w: %s", ErrUnknownTag, tag.Type)
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

func setFieldStringOrNumber(field reflect.Value, val gjson.Result) error {
	fieldType := field.Kind()
	switch fieldType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(val.Int())
	case reflect.Float32, reflect.Float64:
		field.SetFloat(val.Float())
	case reflect.String:
		field.SetString(val.String())
	default:
		return fmt.Errorf("setFieldStringOrNumber: %w %v", ErrUnknownFieldType, fieldType)
	}
	return nil
}

func setFieldBool(field reflect.Value, val gjson.Result) error {
	if field.Kind() == reflect.Bool {
		field.SetBool(val.Bool())
	}
	return nil
}

func setFieldObject(field reflect.Value, val gjson.Result) error {
	var v reflect.Value
	if field.Type().Kind() == reflect.Ptr {
		v = reflect.New(field.Type().Elem())
	} else {
		v = reflect.New(field.Type())
	}
	if err := unmarshalResult(val, v.Interface()); err != nil {
		return err
	}

	if field.Kind() == reflect.Ptr {
		// field type is *Entity
		field.Set(v)
	} else {
		// field type is Entity
		field.Set(v.Elem())
	}
	return nil
}

func setFieldArray(field reflect.Value, val gjson.Result) error {
	arr := val.Array()
	arrLength := len(arr)

	var arrVal reflect.Value
	if field.Type().Kind() == reflect.Ptr {
		arrVal = reflect.MakeSlice(field.Type().Elem(), arrLength, arrLength)
	} else {
		arrVal = reflect.MakeSlice(field.Type(), arrLength, arrLength)
	}

	for i := 0; i < arrLength; i++ {
		v := arr[i]

		setField(arrVal.Index(i), v)
	}
	if field.Kind() == reflect.Ptr {
		// Users *[]user `rejson:"users"`
		fieldVal := reflect.New(field.Type().Elem())
		fieldVal.Elem().Set(arrVal)
		field.Set(fieldVal)
	} else {
		field.Set(arrVal)
	}
	return nil
}

func setField(field reflect.Value, val gjson.Result) error {
	if !field.CanSet() {
		return fmt.Errorf("setField: %w %v", ErrFieldCannotSet, field.Type())
	}
	switch val.Type {
	case gjson.Number:
		return setFieldStringOrNumber(field, val)
	case gjson.String:
		return setFieldStringOrNumber(field, val)
	case gjson.Null:
		// not set by default
		return nil
	case gjson.False:
		return setFieldBool(field, val)
	case gjson.True:
		return setFieldBool(field, val)
	case gjson.JSON:
		if val.IsObject() {
			return setFieldObject(field, val)
		} else if val.IsArray() {
			return setFieldArray(field, val)
		}
	}
	// should not run to here
	return fmt.Errorf("setField: %w %v", ErrUnexpectJSONValue, val)
}

func (u *unmarshal) Unmarshal(v interface{}) error {
	return unmarshalResult(u.r, v)
}

func Unmarshal(jsonString string, v interface{}) error {
	u := newUnmarshal(jsonString)
	return u.Unmarshal(v)
}
