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
	epl map[byte]byte
	dcl map[byte]byte
}

var _ SimpleSubstitutionCipherer = (*Enigoma)(nil)

// NewEnigoma ...
func NewEnigoma(m map[byte]byte) *Enigoma {
	if m == nil || !checkTable(m) {
		log.Printf("create table for substitution")
		m = createTable()
	}

	return &Enigoma{epl: m, dcl: reverseTable(m)}
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
			fmt.Fprintf(&ct, "%s", string(e.epl[byte(t)]))
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
			fmt.Fprintf(&pt, "%s", string(e.dcl[byte(t)]))
		}
	}

	return pt.String()
}

func checkTable(m map[byte]byte) bool {
	if len(m) != 26 {
		return false
	}

	exists := make(map[byte]bool)
	for k, v := range m {
		if k < 'a' || 'z' < k {
			return false
		}

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

func reverseTable(m map[byte]byte) map[byte]byte {
	ret := make(map[byte]byte)
	for k, v := range m {
		ret[v] = k
	}

	return ret
}

func createTable() map[byte]byte {
	ret := make(map[byte]byte)

	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, c := range []byte(strings.ToLower(alpha)) {
		rand.Seed(time.Now().UnixNano())

		i := rand.Intn(len(alpha))
		ret[c] = alpha[i]
		alpha = fmt.Sprintf("%s%s", alpha[0:i], alpha[i+1:])
	}

	return ret
}
