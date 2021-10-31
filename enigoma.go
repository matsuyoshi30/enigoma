package enigoma

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

// Enigoma ...
type Enigoma struct {
	Scrumble
}

// NewEnigoma ...
func NewEnigoma(m []byte) *Enigoma {
	var t [26]byte

	if m == nil || !checkTable(m) {
		log.Printf("create table for substitution")
		t = createTable()
	} else {
		copy(t[:], m[:26])
	}

	return &Enigoma{
		Scrumble: Scrumble{
			t: t,
		},
	}
}

// Encrypt ...
func (e *Enigoma) Encrypt(pt string) (string, error) {
	ot := e.t

	var ct strings.Builder
	for _, t := range pt {
		if t == ' ' {
			fmt.Fprintf(&ct, "%s", " ")
		} else if t < 'a' || 'z' < t {
			return "", fmt.Errorf("only 'a' to 'z' in input text")
		} else {
			fmt.Fprintf(&ct, "%s", string(e.ptoc(byte(t))))
		}
		e.rotate()
	}
	e.t = ot

	return ct.String(), nil
}

// Decrypt ...
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range ct {
		if t == ' ' {
			fmt.Fprintf(&pt, "%s", " ")
		} else {
			fmt.Fprintf(&pt, "%s", string(e.ctop(byte(t))))
		}
		e.rotate()
	}

	return pt.String()
}

type Scrumble struct {
	t [26]byte
}

// rotate
// A B C D ... Z
// to
// B C D E ... A
func (s *Scrumble) rotate() {
	top := s.t[0]
	copy(s.t[0:25], s.t[1:])
	s.t[25] = top
}

func (s *Scrumble) ptoc(b byte) byte {
	if b < 'a' || 'z' < b {
		panic("invalid input byte")
	}

	return s.t[int(b-97)]
}

func (s *Scrumble) ctop(b byte) byte {
	if b < 'A' || 'Z' < b {
		panic("invalid input byte")
	}

	return byte('a' + s.indexAt(b))
}

func (s *Scrumble) indexAt(b byte) int {
	for i, elem := range s.t {
		if elem == b {
			return i
		}
	}

	return -1
}

func checkTable(m []byte) bool {
	if len(m) != 26 {
		return false
	}

	exists := make(map[byte]bool)
	for _, v := range m {
		if v < 'A' || 'Z' < v {
			return false
		}

		if exists[v] {
			return false
		}
		exists[v] = true
	}

	return true
}

func createTable() [26]byte {
	ret := [26]byte{}

	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range []byte(alpha) {
		rand.Seed(time.Now().UnixNano())

		a := rand.Intn(len(alpha))
		ret[i] = alpha[a]
		alpha = fmt.Sprintf("%s%s", alpha[0:a], alpha[a+1:])
	}

	return ret
}
