package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
)

//PUBLIC orbs variable
var PUBLIC = sdk.Export(register, verify, search)

//SYSTEM orbs variable
var SYSTEM = sdk.Export(_init)

var dataKey = []byte("HASH_KEYS")

type searchHit struct {
	Hash   string                 `json:"hash"`
	Score  float64                `json:"score"`
	Source map[string]interface{} `json:"source"`
}

func _init() {
}

func register(phash string, meta string) string {
	//step 1: verify if meta is valid JSON
	var data map[string]interface{}
	err := json.Unmarshal([]byte(meta), &data)
	if err != nil {
		panic(fmt.Sprintf("Invalid JSON: %s", meta))
	}

	//step 2: convert meta JSON to a string
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("JSON stringify error: %s", meta))
	}

	//step 3: calculate real hash
	hash := fmt.Sprintf("%x", md5.Sum(bytes))

	//step 4: check for uniqueness
	key := []byte(phash)
	hs := state.ReadString(key)
	if hs != "" {
		keys := strings.Split(hs, ",")
		for _, k := range keys {
			if k == hash {
				panic(fmt.Sprintf("%s already exists in %s", k, phash))
			}
		}

		keys = append(keys, hash)
		hs = strings.Join(keys, ",")
	} else {
		hs = hash
	}

	//step 5: save phash
	state.WriteString(key, hs)

	//step 6: save hash
	state.WriteBytes([]byte(hash), bytes)

	//step 7: modify phash collection
	s := state.ReadString(dataKey)
	if s == "" {
		s = phash
	} else {
		s = fmt.Sprintf("%s,%s", s, phash)
	}
	state.WriteString(dataKey, s)

	return hs
}

func verify(phash string) string {
	key := []byte(phash)
	s := state.ReadString(key)
	if s == "" {
		panic(fmt.Sprintf("%s does not exists", phash))
	}
	return s
}

func search(phash string, minScore uint64) string {
	s := state.ReadString(dataKey)
	if s == "" {
		return "[]"
	}

	var (
		min  float64
		keys = make(map[string]float64)
		hits []searchHit
	)

	if minScore == 0 || minScore > 100 {
		min = 0.5
	} else {
		min = float64(minScore) / 100
	}

	key := []byte(phash)

	phs := strings.Split(s, ",")
	for _, ph := range phs {
		score := hamming(key, []byte(ph))
		if score >= min {
			keys[ph] = score
		}
	}

	for k, v := range keys {
		hs := state.ReadString([]byte(k))
		if hs != "" {
			hashes := strings.Split(hs, ",")
			for _, hash := range hashes {
				meta := state.ReadString([]byte(hash))
				var jo map[string]interface{}
				err := json.Unmarshal([]byte(meta), &jo)
				if err == nil {
					hit := searchHit{
						Hash:   hash,
						Score:  v,
						Source: jo,
					}
					hits = append(hits, hit)
				}
			}
		}
	}

	sort.SliceStable(hits, func(i, j int) bool {
		return hits[i].Score > hits[j].Score
	})

	bytes, err := json.Marshal(hits)
	if err != nil {
		return "[]"
	}

	return string(bytes)
}

func hamming(txt1, txt2 []byte) float64 {
	switch bytes.Compare(txt1, txt2) {
	case 0: // txt1 == txt2
	case 1: // txt1 > txt2
		temp := make([]byte, len(txt1))
		copy(temp, txt2)
		txt2 = temp
	case -1: // txt1 < txt2
		temp := make([]byte, len(txt2))
		copy(temp, txt1)
		txt1 = temp
	}
	if len(txt1) != len(txt2) {
		panic("Undefined for sequences of unequal length")
	}
	count := 0
	for idx, b1 := range txt1 {
		b2 := txt2[idx]
		xor := b1 ^ b2 // 1 if bits are different
		//
		// bit count (number of 1)
		// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetNaive
		//
		// repeat shifting from left to right (divide by 2)
		// until all bits are zero
		for x := xor; x > 0; x >>= 1 {
			// check if lowest bit is 1
			if int(x&1) == 1 {
				count++
			}
		}
	}
	if count == 0 {
		// similarity is 1 for equal texts.
		return 1
	}
	return float64(1) / float64(count)
}
