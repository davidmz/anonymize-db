package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"
	"regexp"

	lorem "github.com/drhodes/golorem"

	"github.com/gofrs/uuid"
)

// Matches only lowercase V4 uuids
var uuid4Re = regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}")

func anonUUID(str string) string {
	ub, _ := uuid.FromString(str)
	mac := getMacHash()
	mac.Write(ub.Bytes())
	u2, _ := uuid.FromBytes(mac.Sum(nil)[:16])
	u2.SetVersion(uuid.V4)
	u2.SetVariant(uuid.VariantRFC4122)
	return u2.String()
}

func anonAllUUIDs(str string) string {
	return uuid4Re.ReplaceAllStringFunc(str, anonUUID)
}

var anonWord = anonUniqString(func() string { return lorem.Word(5, 12) })
var anonEmail = anonUniqString(func() string { return lorem.Email() })

var createdStrings = make(map[string]bool)
var wordsMap = make(map[string]string)

func anonUniqString(gen func() string) func(string) string {
	return func(str string) string {
		if w, ok := wordsMap[str]; ok {
			return w
		}
		for {
			w := gen()
			if !createdStrings[w] {
				createdStrings[w] = true
				wordsMap[str] = w
				return w
			}
		}
	}
}

var _macHash hash.Hash

func getMacHash() hash.Hash {
	if _macHash == nil {
		key := make([]byte, 32)
		_, _ = io.ReadFull(rand.Reader, key)
		_macHash = hmac.New(sha256.New, key)
	} else {
		_macHash.Reset()
	}
	return _macHash
}
