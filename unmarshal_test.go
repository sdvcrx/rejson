package rejson

import (
	"reflect"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const (
	testUserJSON      = `{"first_name":"John","last_name":"Do","age":18,"married":false,"graduated":true}`
	testUserNestJSON  = `{"code":0,"msg": null,"data":{"name":"John"}}`
	testUserArrayJSON = `{"code":0,"msg":"ok","users":[{"name":"Han"},{"name":"Alex"}]}`
)

type user struct {
	FirstName string `rejson:"first_name"`
	LastName  string `rejson:"last_name"`
	FullName  string `rejson:"func:ParseFullName"`
	Age       int    `rejson:"age"`
	Married   bool   `rejson:"married"`
	Graduated bool   `rejson:"graduated"`
	Empty     string `rejson:"-"`
}

func (usr *user) ParseFullName(jsonReader *gjson.Result) {
	usr.FullName = usr.FirstName + " " + usr.LastName
}

func TestUnmarshalJSONSimple(t *testing.T) {
	u := &user{}
	err := Unmarshal(testUserJSON, u)
	assert.NoError(t, err)
	assert.Equal(t, u.FirstName, "John")
	assert.Equal(t, u.Age, 18)
	assert.Equal(t, "John Do", u.FullName)
	assert.Equal(t, false, u.Married)
	assert.Equal(t, true, u.Graduated)
}

func TestUnmarshalJSONNest(t *testing.T) {
	type user struct {
		Name string `rejson:"name"`
	}

	t.Run("Struct", func(t *testing.T) {
		resp := struct {
			Code int    `rejson:"code"`
			Msg  string `rejson:"msg"`
			Data user   `rejson:"data"`
		}{}
		err := Unmarshal(testUserNestJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "John", resp.Data.Name)
	})

	t.Run("Pointer of struct", func(t *testing.T) {
		resp := struct {
			Code int    `rejson:"code"`
			Msg  string `rejson:"msg"`
			Data *user  `rejson:"data"`
		}{}
		err := Unmarshal(testUserNestJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "John", resp.Data.Name)
	})
}

func TestUnmarshalJSONArray(t *testing.T) {
	type user struct {
		Name string `rejson:"name"`
	}
	t.Run("Slice", func(t *testing.T) {
		resp := struct {
			Code  int    `rejson:"code"`
			Msg   string `rejson:"msg"`
			Users []user `rejson:"users"`
		}{}
		err := Unmarshal(testUserArrayJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Users)
		assert.Len(t, resp.Users, 2)
	})

	t.Run("Pointer of slice", func(t *testing.T) {
		resp := struct {
			Users *[]user `rejson:"users"`
		}{}
		err := Unmarshal(testUserArrayJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Users)
		assert.Len(t, *resp.Users, 2)
	})

	t.Run("Slice of number", func(t *testing.T) {
		resp := struct {
			Nums []int `rejson:"nums"`
		}{}

		assert.NoError(t, Unmarshal(`{"nums":[1,2,3]}`, &resp))
		assert.Len(t, resp.Nums, 3)
		for i := 0; i < 3; i++ {
			assert.Equal(t, i+1, resp.Nums[i])
		}
	})

	t.Run("Slice of string", func(t *testing.T) {
		resp := struct {
			Names []string `rejson:"names"`
		}{}

		assert.NoError(t, Unmarshal(`{"names":["a","b","c"]}`, &resp))
		assert.Len(t, resp.Names, 3)
	})
}

func TestUnmarshalUnknownTag(t *testing.T) {
	d := struct {
		Nums []int `rejson:"test:test"`
	}{}
	err := Unmarshal(`{}`, &d)
	assert.True(t, errors.Is(err, ErrUnknownTag))
}

func TestUnmarshalNumber(t *testing.T) {
	t.Run("Float64", func(t *testing.T) {
		resp := struct {
			Money float64 `rejson:"money"`
		}{}

		assert.NoError(t, Unmarshal(`{"money":3.2}`, &resp))
		assert.Equal(t, 3.2, resp.Money)
	})

	t.Run("Float32", func(t *testing.T) {
		resp := struct {
			Money float32 `rejson:"money"`
		}{}

		assert.NoError(t, Unmarshal(`{"money":3.2}`, &resp))
		assert.Equal(t, float32(3.2), resp.Money)
	})

	t.Run("Int32", func(t *testing.T) {
		resp := struct {
			Money int32 `rejson:"money"`
		}{}

		assert.NoError(t, Unmarshal(`{"money":3}`, &resp))
		assert.Equal(t, int32(3), resp.Money)
	})
}

func TestUnmarshalResultError(t *testing.T) {
	jsonStr := gjson.Parse("\"123\"")
	var u *user
	err := unmarshalResult(jsonStr, u)
	assert.Error(t, err)
}

func TestSetField(t *testing.T) {
	jsonStr := gjson.Parse("\"123\"")

	t.Run("Expect return ErrFieldCannotSet", func(t *testing.T) {
		var s string
		err := setField(reflect.ValueOf(s), jsonStr)
		assert.True(t, errors.Is(err, ErrFieldCannotSet))
	})

	t.Run("Expect return ErrUnknownFieldType", func(t *testing.T) {
		var v bool
		err := setFieldStringOrNumber(reflect.ValueOf(v), jsonStr)
		assert.True(t, errors.Is(err, ErrUnknownFieldType))
	})
}
