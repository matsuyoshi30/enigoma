package enigoma_test

import (
	"testing"

	"github.com/matsuyoshi30/enigoma"
)

func TestEnigoma(t *testing.T) {
	tests := []struct {
		desc string
		pt   string
		m    []byte
	}{
		{
			desc: "simple",
			pt:   "hello world",
		},
		{
			desc: "custom table",
			pt:   "hello world with custom table",
			m:    stupidTable(),
		},
		{
			desc: "long text",
			pt: "a long time ago in a galaxy far far away " +
				"it is a period of civil war " +
				"rebel spaceships striking from a hidden base " +
				"have won their first victory against the evil galactic empire",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			e := enigoma.NewEnigoma(tt.m)

			c := e.Encrypt(tt.pt)
			t.Log(c)
			a := e.Decrypt(c)

			if tt.pt != a {
				t.Errorf("want '%s' but got '%s'", tt.pt, a)
			}
		})
	}
}

func stupidTable() []byte {
	ret := make([]byte, 26)

	for i := 'a'; i <= 'z'; i++ {
		ret[i-97] = byte(i)
	}

	return ret
}
