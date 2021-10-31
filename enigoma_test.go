package enigoma_test

import (
	"testing"

	"github.com/matsuyoshi30/enigoma"
)

func TestEnigoma(t *testing.T) {
	tests := []struct {
		pt  string
		m   []byte
		err bool
	}{
		{
			pt: "hello world",
		},
		{
			pt: "hello world with custom table",
			m:  stupidTable(),
		},
		{
			pt:  "Invalid Plain TEXT",
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.pt, func(t *testing.T) {
			e := enigoma.NewEnigoma(tt.m)

			c, err := e.Encrypt(tt.pt)
			if tt.err {
				if err == nil {
					t.Fatal("want error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("want no error but got %v", err)
			}
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
		ret[i-97] = byte(i) - 32
	}

	return ret
}
