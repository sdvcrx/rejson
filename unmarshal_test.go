package rejson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const (
	testUserJSON     = `{"first_name":"John","last_name":"Do","age":18}`
	testUserNestJSON = `{"code":0,"msg": null,"data":{"name":"John"}}`
)

type user struct {
	FirstName string `jsonp:"first_name"`
	LastName  string `jsonp:"last_name"`
	FullName  string `jsonp:"func:ParseFullName"`
	Age       int    `jsonp:"age"`
	Empty     string `jsonp:"-"`
}

func (usr *user) ParseFullName(jsonReader *gjson.Result) {
	usr.FullName = usr.FirstName + " " + usr.LastName
}

func TestUnmarshalJSONSimple(t *testing.T) {
	c := NewConverter(testUserJSON)

	u := &user{}
	err := c.Unmarshal(u)
	assert.NoError(t, err)
	t.Logf("%+v", u)
	assert.Equal(t, u.FirstName, "John")
	assert.Equal(t, u.Age, 18)
	assert.Equal(t, "John Do", u.FullName)
}

func TestUnmarshalJSONNest(t *testing.T) {
	cvt := NewConverter(testUserNestJSON)

	type user struct {
		Name string `jsonp:"name"`
	}
	type response struct {
		Code int    `jsonp:"code"`
		Msg  string `jsonp:"msg"`
		Data user   `jsonp:"data"`
	}
	resp := &response{}
	err := cvt.Unmarshal(resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "John", resp.Data.Name)

	type response2 struct {
		Code int    `jsonp:"code"`
		Msg  string `jsonp:"msg"`
		Data *user  `jsonp:"data"`
	}
	resp2 := &response2{}
	err = cvt.Unmarshal(resp2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2.Data)
	assert.Equal(t, "John", resp2.Data.Name)
}
