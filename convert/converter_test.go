package convert

import (
	"testing"
	"time"

	"github.com/visheratin/url-short/config"
	"github.com/visheratin/url-short/storage"
)

var testlinks = [][2]string{
	[2]string{"https://golang.org/", "AAAAAA"},
	[2]string{"https://github.com/", "AAAAAa"},
	[2]string{"https://golang.org/", "AAAAAA"},
	[2]string{"https://www.reddit.com", "AAAAAB"},
}

func TestLoad(t *testing.T) {
	c := NewConverter(6, nil)
	for _, test := range testlinks {
		code, err := c.Load(test[0])
		if err != nil {
			panic(err)
		}
		if code != test[1] {
			t.Error(
				"For", test[0],
				"expected", test[1],
				"got", code,
			)
		}
	}
	for _, test := range testlinks {
		input := c.Extract(test[1])
		if input != test[0] {
			t.Error(
				"For", test[1],
				"expected", test[0],
				"got", input,
			)
		}
	}
}

func BenchmarkLoadWithStorage(b *testing.B) {
	config.Init("../config.json")
	storage, _ := storage.Instance()
	c := NewConverter(6, storage)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		str := time.Now().String()
		b.StartTimer()
		c.Load(str)
	}
}

func BenchmarkLoad(b *testing.B) {
	c := NewConverter(6, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		str := time.Now().String()
		b.StartTimer()
		c.Load(str)
	}
}
