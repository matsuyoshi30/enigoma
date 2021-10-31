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
	s1 *Scrumble
	s2 *Scrumble
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

	s1 := &Scrumble{t: t}
	s2 := &Scrumble{t: t}

	return &Enigoma{s1: s1, s2: s2}
}

// Encrypt ...
func (e *Enigoma) Encrypt(pt string) string {
	ot1 := e.s1.t
	ot2 := e.s2.t

	var ct strings.Builder
	for _, t := range pt {
		so1 := e.s1.ptoc(byte(t))
		so2 := e.s2.ptoc(so1)

		fmt.Fprintf(&ct, "%s", string(so2))

		e.s1.rotate()
		if e.s1.fullRotated() {
			e.s2.rotate()
		}
	}
	e.s1.t, e.s1.ra = ot1, 0
	e.s2.t, e.s2.ra = ot2, 0

	return strings.ToUpper(ct.String())
}

// Decrypt ...
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range strings.ToLower(ct) {
		so2 := e.s2.ctop(byte(t))
		so1 := e.s1.ctop(so2)

		fmt.Fprintf(&pt, "%s", string(so1))
		e.s1.rotate()
		if e.s1.fullRotated() {
			e.s2.rotate()
		}
	}

	return pt.String()
}

type Scrumble struct {
	t  [26]byte
	ra int
}

// rotate
// A B C D ... Z
// to
// B C D E ... A
func (s *Scrumble) rotate() {
	top := s.t[0]
	copy(s.t[0:25], s.t[1:])
	s.t[25] = top

	s.ra = (s.ra + 1) % 26
}

func (s *Scrumble) fullRotated() bool {
	return s.ra == 0
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
