package enigoma

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Enigoma represents Enigma machine.
type Enigoma struct {
	p *PlugBoard
	s *Scramble
}

// NewEnigoma creates new enigoma instance.
//
// Enigoma has three scrambles (with conversion information for at least all
// alphabets), so it needs the key information that defines the scramble
// settings and initial state. Enigma also has a plugboard that can be
// configured by the operator, so we also need that information.
//
// If given information for the scramble is invalid, create it internally.
func NewEnigoma(m1, m2, m3 []byte, k1, k2, k3 byte, pb *PlugBoard) *Enigoma {
	var t1, t2, t3 [26]byte

	if !checkTable(m1) {
		t1 = createTable()
	} else {
		copy(t1[:], m1[:26])
	}
	t1 = shift(t1, k1)

	if !checkTable(m2) {
		t2 = createTable()
	} else {
		copy(t2[:], m2[:26])
	}
	t2 = shift(t2, k2)

	if !checkTable(m3) {
		t3 = createTable()
	} else {
		copy(t3[:], m3[:26])
	}
	t3 = shift(t3, k3)

	return &Enigoma{
		p: pb,
		s: NewScramble(t1, NewScramble(t2, NewScramble(t3, NewReflector(13)))),
	}
}

// Encrypt encrypts plain text to cipher.
func (e *Enigoma) Encrypt(pt string) string {
	_s := e.s.copyScramble()

	var ct strings.Builder
	for _, t := range strings.ToLower(pt) {
		ec := e.p.exchange(byte(t)) // exchange by plugboard
		o := e.convert(ec)          // convert via s1 -> s2 -> s3 -> reflector -> s3 -> s2 -> s1
		r := e.p.exchange(o)        // exchange by plugboard again
		fmt.Fprintf(&ct, "%s", string(r))

		e.s.Rotate()
	}
	e.s = _s

	return strings.ToUpper(ct.String())
}

// Decrypt decrypts cipher to plain text.
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range strings.ToLower(ct) {
		ec := e.p.exchange(byte(t)) // exchange by plugboard
		o := e.convert(ec)          // convert via s1 -> s2 -> s3 -> reflector -> s3 -> s2 -> s1
		r := e.p.exchange(o)        // exchange by plugboard again
		fmt.Fprintf(&pt, "%s", string(r))

		e.s.Rotate()
	}

	return pt.String()
}

func (e *Enigoma) convert(b byte) byte {
	return e.s.ctop(e.s.ptoc(b))
}

// Rotor is the interface that can rotate.
type Rotor interface {
	Rotate()
}

// Scramble implement a polyalphabetic substitution cipher that provides Enigma's security.
type Scramble struct {
	t  [26]byte
	ra int

	n Rotor
}

// NewScramble creates new scramble instance.
func NewScramble(t [26]byte, next Rotor) *Scramble {
	s := &Scramble{t: t}
	if next != nil {
		s.n = next
	}

	return s
}

func (s *Scramble) ptoc(b byte) byte {
	if refl, ok := s.n.(*Reflector); ok {
		return refl.Reflect(s.t, s._ptoc(b))
	}

	ns := s.n.(*Scramble)
	return ns.ptoc(s._ptoc(b))
}

func (s *Scramble) ctop(b byte) byte {
	if ns, ok := s.n.(*Scramble); ok {
		return s._ctop(ns.ctop(b))
	}

	return s._ctop(b)
}

// Rotate of Scramble rotates scrambles.
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

func (s *Scramble) _ptoc(b byte) byte {
	if b < 'a' || 'z' < b {
		return b
	}

	return s.t[int(b-97)]
}

func (s *Scramble) _ctop(b byte) byte {
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

// Reflector implements reflection.
type Reflector struct {
	gap int
}

// NewReflector creates new reflector instance.
func NewReflector(g int) *Reflector {
	if g < 1 || 25 < g {
		rand.Seed(time.Now().UnixNano())
		g = rand.Intn(24) + 1
	}

	return &Reflector{gap: g}
}

// Reflect reflects one byte to another.
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

// Rotate of Reflector does nothing.
func (r *Reflector) Rotate() {}

// PlugBoard permitted variable wiring that could be reconfigured by the operator.
type PlugBoard struct {
	m map[byte]byte
}

// NewPlugBoard creates new plugboard instance.
func NewPlugBoard() *PlugBoard {
	return &PlugBoard{m: make(map[byte]byte)}
}

// AddExchange add a wire.
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
	if len(m) < 26 {
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
