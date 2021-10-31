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
	p *PlugBoard
	s *Scramble
}

// NewEnigoma ...
func NewEnigoma(m1, m2, m3 []byte, k1, k2, k3 byte, g int, pb *PlugBoard) *Enigoma {
	var t1, t2, t3 [26]byte

	if !checkTable(m1) {
		log.Printf("create table for substitution")
		t1 = createTable()
	} else {
		copy(t1[:], m1[:26])
	}
	t1 = shift(t1, k1)

	if !checkTable(m2) {
		log.Printf("create table for substitution")
		t2 = createTable()
	} else {
		copy(t2[:], m2[:26])
	}
	t2 = shift(t2, k2)

	if !checkTable(m3) {
		log.Printf("create table for substitution")
		t3 = createTable()
	} else {
		copy(t3[:], m3[:26])
	}
	t3 = shift(t3, k3)

	return &Enigoma{
		p: pb,
		s: NewScramble(t1, NewScramble(t2, NewScramble(t3, NewReflector(g)))),
	}
}

// Encrypt ...
func (e *Enigoma) Encrypt(pt string) string {
	_s := e.s.copyScramble()

	var ct strings.Builder
	for _, t := range pt {
		ec := e.p.exchange(byte(t))
		o := e.s.CtoP(e.s.PtoC(ec))
		result := e.p.exchange(o)
		fmt.Fprintf(&ct, "%s", string(result))

		e.s.Rotate()
	}
	e.s = _s

	return strings.ToUpper(ct.String())
}

// Decrypt ...
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range strings.ToLower(ct) {
		ec := e.p.exchange(byte(t))
		o := e.s.CtoP(e.s.PtoC(ec))
		result := e.p.exchange(o)
		fmt.Fprintf(&pt, "%s", string(result))

		e.s.Rotate()
	}

	return pt.String()
}

type Rotor interface {
	Rotate()
}

type Scramble struct {
	t  [26]byte
	ra int

	n Rotor
}

func NewScramble(t [26]byte, next Rotor) *Scramble {
	s := &Scramble{t: t}
	if next != nil {
		s.n = next
	}

	return s
}

func (s *Scramble) PtoC(b byte) byte {
	if refl, ok := s.n.(*Reflector); ok {
		return refl.Reflect(s.t, s.ptoc(b))
	}

	ns := s.n.(*Scramble)
	return ns.PtoC(s.ptoc(b))
}

func (s *Scramble) CtoP(b byte) byte {
	if ns, ok := s.n.(*Scramble); ok {
		return s.ctop(ns.CtoP(b))
	}

	return s.ctop(b)
}

func (s *Scramble) Rotate() {
	s.rotate()
	if s.fullRotated() && s.n != nil {
		s.n.Rotate()
	}
}

func (s *Scramble) copyScramble() *Scramble {
	ns := s.n.(*Scramble)
	nns := ns.n.(*Scramble)
	refl := nns.n.(*Reflector)

	return NewScramble(s.t, NewScramble(ns.t, NewScramble(nns.t, refl)))
}

// rotate
// A B C D ... Z
// to
// B C D E ... A
func (s *Scramble) rotate() {
	top := s.t[0]
	copy(s.t[0:25], s.t[1:])
	s.t[25] = top

	s.ra = (s.ra + 1) % 26
}

func (s *Scramble) fullRotated() bool {
	return s.ra == 0
}

func (s *Scramble) ptoc(b byte) byte {
	if b < 'a' || 'z' < b {
		return b
	}

	return s.t[int(b-97)]
}

func (s *Scramble) ctop(b byte) byte {
	if b < 'a' || 'z' < b {
		return b
	}

	return byte('a' + s.indexAt(b))
}

func (s *Scramble) indexAt(b byte) int {
	for i, elem := range s.t {
		if elem == b {
			return i
		}
	}

	return -1
}

type Reflector struct {
	gap int
}

func NewReflector(g int) *Reflector {
	if g < 1 || 25 < g {
		rand.Seed(time.Now().UnixNano())
		g = rand.Intn(24) + 1
	}

	return &Reflector{gap: g}
}

func (r *Reflector) Reflect(t [26]byte, b byte) byte {
	if b < 'a' || 'z' < b {
		return b
	}

	base := -1
	for i, elem := range t {
		if elem == b {
			base = i
			break
		}
	}
	if base == -1 {
		panic("cannot reflect")
	}

	return t[(base+r.gap)%26]
}

func (r *Reflector) Rotate() {}

type PlugBoard struct {
	m map[byte]byte
}

func NewPlugBoard() *PlugBoard {
	return &PlugBoard{m: make(map[byte]byte)}
}

func (p *PlugBoard) AddExchange(b1, b2 byte) {
	if _, exists := p.m[b1]; exists {
		delete(p.m, b1)
	}
	if _, exists := p.m[b2]; exists {
		delete(p.m, b2)
	}

	p.m[b1] = b2
	p.m[b2] = b1
}

func (p *PlugBoard) exchange(b byte) byte {
	e, ok := p.m[b]
	if ok {
		return e
	}

	return b
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

func shift(t [26]byte, k byte) [26]byte {
	top := -1
	for i, elem := range t {
		if elem == k {
			top = i
			break
		}
	}
	if top == -1 {
		panic("invalid key")
	}

	var ret [26]byte
	copy(ret[:], t[top:])
	copy(ret[len(t[top:]):], t[:top])
	return ret
}
