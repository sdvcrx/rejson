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
	type response struct {
		Code int    `jsonp:"code"`
		Msg  string `jsonp:"msg"`
		Data user   `jsonp:"data"`
	}
	resp := &response{}
	err := Unmarshal(testUserNestJSON, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "John", resp.Data.Name)

	type response2 struct {
		Code int    `jsonp:"code"`
		Msg  string `jsonp:"msg"`
		Data *user  `jsonp:"data"`
	}
	resp2 := &response2{}
	err = Unmarshal(testUserNestJSON, resp2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2.Data)
	assert.Equal(t, "John", resp2.Data.Name)
}

func TestUnmarshalJSONArray(t *testing.T) {
	type user struct {
		Name string `jsonp:"name"`
	}
	type response struct {
		Code  int    `jsonp:"code"`
		Msg   string `jsonp:"msg"`
		Users []user `jsonp:"users"`
	}

	resp := &response{}
	err := Unmarshal(testUserArrayJSON, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Users)
	assert.Len(t, resp.Users, 2)
}
