// Package convert performs shortening of an input string into a code of a fixed length.
package convert

import (
	"strings"
	"sync"

	"github.com/visheratin/url-short/log"
	"github.com/visheratin/url-short/storage"
)

// alphabet is a set of characters that can be used in the code.
const alphabet = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

/*
Converter contains maps of matching between input strings and their codes.
links: map of matching between codes and input strings;
revLinks: map of matching between codes and input strings, used for search
of repeating strings;
positions: slice containing current position of code literals in the alphabet;
m: mutex used for blocking position setting procedure;
storage: connector for reading and writing data into a permanent storage.
*/
type Converter struct {
	links     sync.Map
	revLinks  sync.Map
	positions []int
	m         sync.RWMutex
	storage   storage.Storage
}

// NewConverter creates a Converter instance using specified length
// and storage. If storage is nil, input-code pairs will not be saved
// after the program finishing.
func NewConverter(length int, storage storage.Storage) *Converter {
	pos := make([]int, length)
	converter := &Converter{
		positions: pos,
		storage:   storage,
	}
	converter.loadLinks()
	return converter
}

// loadLinks extracts existing input-code pairs from the permanent
// storage and sets a proper position for the converter.
func (conv *Converter) loadLinks() {
	if conv.storage == nil {
		return
	}
	changed := false
	existLinks, err := conv.storage.LoadAll()
	if err != nil {
		log.Log().Error.Println(err)
		return
	}
	for _, pair := range existLinks {
		conv.links.Store(pair[0], pair[1])
		conv.revLinks.Store(pair[1], pair[0])
		pos := make([]int, len(conv.positions))
		larger := true
		for i := 0; i < len(pair[0]); i++ {
			idx := strings.IndexByte(alphabet, pair[0][i])
			if idx == -1 {
				break
			}
			pos[i] = idx
			if idx < conv.positions[i] {
				larger = false
				break
			}
		}
		if larger {
			conv.positions = pos
			changed = true
		}
	}
	if changed {
		conv.getPositions()
	}
}

// Load generates a code for the input string, writes it to
// the matching maps and to the permanent storage. If a code
// for the input string exists, Load will return it instead
// of generating a new one.
func (conv *Converter) Load(input string) (string, error) {
	val, ok := conv.revLinks.Load(input)
	if ok {
		return val.(string), nil
	}
	conv.m.Lock()
	code, err := conv.getCode()
	conv.m.Unlock()
	if err != nil {
		log.Log().Error.Println(err)
		return "", err
	}
	if conv.storage != nil {
		err = conv.storage.Store(code, input)
		if err != nil {
			log.Log().Error.Println(err)
		}
	}
	conv.links.Store(code, input)
	conv.revLinks.Store(input, code)
	return code, err
}

// Extract gets an input string for the specified code.
// If there is no input string for the code, Extract will
// return empty string.
func (conv *Converter) Extract(code string) string {
	val, ok := conv.links.Load(code)
	if !ok {
		return ""
	}
	return val.(string)
}

// getCode takes alphabet characters for the current converter
// position and combines them into the code.
func (conv *Converter) getCode() (string, error) {
	positions := conv.getPositions()
	var sb strings.Builder
	for _, pos := range positions {
		err := sb.WriteByte(alphabet[pos])
		if err != nil {
			log.Log().Error.Println(err)
			return "", err
		}
	}
	return sb.String(), nil
}

// getPositions shifts extracts current positions of the
// converter and increases current positions by 1.
func (conv *Converter) getPositions() []int {
	res := make([]int, len(conv.positions))
	copy(res, conv.positions)
	for i := (len(conv.positions) - 1); i >= 0; i-- {
		if conv.positions[i] == len(alphabet)-1 {
			conv.positions[i] = 0
		} else {
			conv.positions[i]++
			break
		}
	}
	return res
}
