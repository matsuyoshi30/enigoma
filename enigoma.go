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
	s *Scrumble
}

// NewEnigoma ...
func NewEnigoma(m1, m2, m3 []byte) *Enigoma {
	var t1, t2, t3 [26]byte

	if !checkTable(m1) {
		log.Printf("create table for substitution")
		t1 = createTable()
	} else {
		copy(t1[:], m1[:26])
	}

	if !checkTable(m2) {
		log.Printf("create table for substitution")
		t2 = createTable()
	} else {
		copy(t2[:], m2[:26])
	}

	if !checkTable(m3) {
		log.Printf("create table for substitution")
		t3 = createTable()
	} else {
		copy(t3[:], m3[:26])
	}

	return &Enigoma{
		s: NewScrumble(t1, NewScrumble(t2, NewScrumble(t3, nil))),
	}
}

// Encrypt ...
func (e *Enigoma) Encrypt(pt string) string {
	_s := e.s.copyScrumble()

	var ct strings.Builder
	for _, t := range pt {
		fmt.Fprintf(&ct, "%s", string(e.s.PtoC(byte(t))))

		e.s.Rotate()
	}
	e.s = _s

	return strings.ToUpper(ct.String())
}

// Decrypt ...
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range strings.ToLower(ct) {
		fmt.Fprintf(&pt, "%s", string(e.s.CtoP(byte(t))))

		e.s.Rotate()
	}

	return pt.String()
}

type Scrumble struct {
	t  [26]byte
	ra int

	n *Scrumble
}

func NewScrumble(t [26]byte, next *Scrumble) *Scrumble {
	s := &Scrumble{t: t}
	if next != nil {
		s.n = next
	}

	return s
}

func (s *Scrumble) PtoC(b byte) byte {
	if s.n == nil {
		return s.ptoc(b)
	}

	return s.n.PtoC(s.ptoc(b))
}

func (s *Scrumble) CtoP(b byte) byte {
	if s.n == nil {
		return s.ctop(b)
	}

	return s.ctop(s.n.CtoP(b))
}

func (s *Scrumble) Rotate() {
	s.rotate()
	if s.fullRotated() && s.n != nil {
		s.n.Rotate()
	}
}

func (s *Scrumble) copyScrumble() *Scrumble {
	return NewScrumble(s.t, NewScrumble(s.n.t, NewScrumble(s.n.n.t, nil)))
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
