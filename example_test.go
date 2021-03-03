package rejson_test

import (
	"fmt"

	"github.com/sdvcrx/rejson"
)

type User struct {
	Name string `rejson:"data.name"`
}

func Example() {
	jsonString := `{
    "code":0,
    "msg": null,
    "data":{"name":"John"}
  }`

	u := User{}
	err := rejson.Unmarshal(jsonString, &u)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User%+v", u)
	// Output: User{Name:John}
}
