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
func (e *Enigoma) Encrypt(pt string) string {
	ot := e.t

	var ct strings.Builder
	for _, t := range pt {
		fmt.Fprintf(&ct, "%s", string(e.ptoc(byte(t))))
		e.rotate()
	}
	e.t = ot

	return strings.ToUpper(ct.String())
}

// Decrypt ...
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range strings.ToLower(ct) {
		fmt.Fprintf(&pt, "%s", string(e.ctop(byte(t))))
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
		return b
	}

	return s.t[int(b-97)]
}

func (s *Scrumble) ctop(b byte) byte {
	if b < 'a' || 'z' < b {
		return b
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
		if v < 'a' || 'z' < v {
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

	// alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alpha := "abcdefghijklmnopqrstuvwxyz"
	for i := range []byte(alpha) {
		rand.Seed(time.Now().UnixNano())

		a := rand.Intn(len(alpha))
		ret[i] = alpha[a]
		alpha = fmt.Sprintf("%s%s", alpha[0:a], alpha[a+1:])
	}

	return ret
}
