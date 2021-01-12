package rejson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const (
	testUserJSON      = `{"first_name":"John","last_name":"Do","age":18,"married":false,"graduated":true}`
	testUserNestJSON  = `{"code":0,"msg": null,"data":{"name":"John"}}`
	testUserArrayJSON = `{"code":0,"msg":"ok","users":[{"name":"Han"},{"name":"Alex"}]}`
)

type user struct {
	FirstName string `jsonp:"first_name"`
	LastName  string `jsonp:"last_name"`
	FullName  string `jsonp:"func:ParseFullName"`
	Age       int    `jsonp:"age"`
	Married   bool   `jsonp:"married"`
	Graduated bool   `jsonp:"graduated"`
	Empty     string `jsonp:"-"`
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
		Name string `jsonp:"name"`
	}

	t.Run("Struct", func(t *testing.T) {
		resp := struct {
			Code int    `jsonp:"code"`
			Msg  string `jsonp:"msg"`
			Data user   `jsonp:"data"`
		}{}
		err := Unmarshal(testUserNestJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "John", resp.Data.Name)
	})

	t.Run("Pointer of struct", func(t *testing.T) {
		resp := struct {
			Code int    `jsonp:"code"`
			Msg  string `jsonp:"msg"`
			Data *user  `jsonp:"data"`
		}{}
		err := Unmarshal(testUserNestJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "John", resp.Data.Name)
	})
}

func TestUnmarshalJSONArray(t *testing.T) {
	type user struct {
		Name string `jsonp:"name"`
	}
	t.Run("Slice", func(t *testing.T) {
		resp := struct {
			Code  int    `jsonp:"code"`
			Msg   string `jsonp:"msg"`
			Users []user `jsonp:"users"`
		}{}
		err := Unmarshal(testUserArrayJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Users)
		assert.Len(t, resp.Users, 2)
	})

	t.Run("Pointer of slice", func(t *testing.T) {
		resp := struct {
			Users *[]user `jsonp:"users"`
		}{}
		err := Unmarshal(testUserArrayJSON, &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Users)
		assert.Len(t, *resp.Users, 2)
	})

	t.Run("Slice of number", func(t *testing.T) {
		resp := struct {
			Nums []int `jsonp:"nums"`
		}{}

		assert.NoError(t, Unmarshal(`{"nums":[1,2,3]}`, &resp))
		assert.Len(t, resp.Nums, 3)
		for i := 0; i < 3; i++ {
			assert.Equal(t, i+1, resp.Nums[i])
		}
	})

	t.Run("Slice of string", func(t *testing.T) {
		resp := struct {
			Names []string `jsonp:"names"`
		}{}

		assert.NoError(t, Unmarshal(`{"names":["a","b","c"]}`, &resp))
		assert.Len(t, resp.Names, 3)
	})
}

func TestUnmarshalNumber(t *testing.T) {
	t.Run("Float64", func(t *testing.T) {
		resp := struct {
			Money float64 `jsonp:"money"`
		}{}

		assert.NoError(t, Unmarshal(`{"money":3.2}`, &resp))
		assert.Equal(t, 3.2, resp.Money)
	})

	t.Run("Float32", func(t *testing.T) {
		resp := struct {
			Money float32 `jsonp:"money"`
		}{}

		assert.NoError(t, Unmarshal(`{"money":3.2}`, &resp))
		assert.Equal(t, float32(3.2), resp.Money)
	})

	t.Run("Int32", func(t *testing.T) {
		resp := struct {
			Money int32 `jsonp:"money"`
		}{}

		assert.NoError(t, Unmarshal(`{"money":3}`, &resp))
		assert.Equal(t, int32(3), resp.Money)
	})
}
