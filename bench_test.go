package rejson

import (
	"encoding/json"
	"testing"
)

func BenchmarkUnmarshalReJSON(b *testing.B) {
	type user struct {
		FirstName string `jsonp:"first_name"`
		LastName  string `jsonp:"last_name"`
		Age       int    `jsonp:"age"`
		Empty     string `jsonp:"-"`
	}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		u := &user{}
		Unmarshal(testUserJSON, u)
	}
}

func BenchmarkUnmarshalEncodingJSON(b *testing.B) {
	type user struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
		Empty     string `json:"-"`
	}
	data := []byte(testUserJSON)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		u := &user{}
		json.Unmarshal(data, u)
	}
}
