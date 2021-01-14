package rejson

import (
	"encoding/json"
	"testing"

	"github.com/tidwall/gjson"
)

const (
	benchmarkJSON = `{"widget":{"debug":"on","window":{"title":"Sample Konfabulator Widget","name":"main_window","width":500,"height":500},"image":{"src":"Images/Sun.png","hOffset":250,"vOffset":250,"alignment":"center"},"text":{"data":"Click Here","size":36,"style":"bold","vOffset":100,"alignment":"center","onMouseUp":"sun1.opacity = (sun1.opacity / 100) * 90;"}}}`
)

type benchWidget struct {
	Name      string `rejson:"widget.window.name"`
	HOffset   int    `rejson:"widget.image.hOffset"`
	OnMouseUp string `rejson:"widget.text.onMouseUp"`
}

type benchStruct struct {
	Widget struct {
		Window struct {
			Name string `json:"name"`
		} `json:"window"`
		Image struct {
			HOffset int `json:"hOffset"`
		} `json:"image"`
		Text struct {
			OnMouseUp string `json:"onMouseUp"`
		} `json:"text"`
	} `json:"widget"`
}

func BenchmarkUnmarshalReJSON(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := benchWidget{}
		Unmarshal(benchmarkJSON, &w)
	}
}

func BenchmarkUnmarshalGJSONGet(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := benchWidget{}
		r := gjson.Parse(benchmarkJSON)
		w.Name = r.Get("widget.window.name").String()
		w.HOffset = int(r.Get("widget.image.hOffset").Int())
		w.OnMouseUp = r.Get("widget.text.onMouseUp").String()
	}
}

func BenchmarkUnmarshalEncodingJSON(b *testing.B) {
	data := []byte(benchmarkJSON)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bs := benchStruct{}
		json.Unmarshal(data, &bs)

		w := benchWidget{}
		w.Name = bs.Widget.Window.Name
		w.HOffset = bs.Widget.Image.HOffset
		w.OnMouseUp = bs.Widget.Text.OnMouseUp
	}
}
