package enigoma

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

// SimpleSubstitutionCipherer ...
type SimpleSubstitutionCipherer interface {
	Encrypt(string) (string, error)
	Decrypt(string) string
}

// Enigoma ...
type Enigoma struct {
	t [26]byte
}

var _ SimpleSubstitutionCipherer = (*Enigoma)(nil)

// NewEnigoma ...
func NewEnigoma(m []byte) *Enigoma {
	var t [26]byte

	if m == nil || !checkTable(m) {
		log.Printf("create table for substitution")
		t = createTable()
	} else {
		copy(t[:], m[:26])
	}

	return &Enigoma{t: t}
}

// Encrypt ...
func (e *Enigoma) Encrypt(pt string) (string, error) {
	var ct strings.Builder
	for _, t := range pt {
		if t == ' ' {
			fmt.Fprintf(&ct, "%s", " ")
		} else if t < 'a' || 'z' < t {
			return "", fmt.Errorf("only 'a' to 'z' in input text")
		} else {
			fmt.Fprintf(&ct, "%s", string(e.t[atoi(byte(t))]))
		}
	}

	return ct.String(), nil
}

// Decrypt ...
func (e *Enigoma) Decrypt(ct string) string {
	var pt strings.Builder
	for _, t := range ct {
		if t == ' ' {
			fmt.Fprintf(&pt, "%s", " ")
		} else {
			fmt.Fprintf(&pt, "%s", string(byte('a'+e.indexAt(byte(t)))))
		}
	}

	return pt.String()
}

func (e *Enigoma) indexAt(b byte) int {
	for i, elem := range e.t {
		if elem == b {
			return i
		}
	}

	return -1
}

func atoi(b byte) int {
	if 'a' <= b && b <= 'z' {
		return int(b - 97)
	}
	if 'A' <= b && b <= 'Z' {
		return int(b - 65)
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
