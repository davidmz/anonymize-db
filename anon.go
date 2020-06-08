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

func anonUUIDs(str string) string {
	return uuid4Re.ReplaceAllStringFunc(str, func(uStr string) string {
		ub, _ := uuid.FromString(uStr)
		mac := getMacHash()
		mac.Write(ub.Bytes())
		u2, _ := uuid.FromBytes(mac.Sum(nil)[:16])
		u2.SetVersion(uuid.V4)
		u2.SetVariant(uuid.VariantRFC4122)
		return u2.String()
	})
}

var wordsMap = make(map[string]string)

func anonWord(str string) string {
	if w, ok := wordsMap[str]; ok {
		return w
	}
	for {
		w := lorem.Word(5, 12)
		if !createdStrings[w] {
			createdStrings[w] = true
			wordsMap[str] = w
			return w
		}
	}
}

func anonEmail(str string) string {
	if w, ok := wordsMap[str]; ok {
		return w
	}
	for {
		w := lorem.Email()
		if !createdStrings[w] {
			createdStrings[w] = true
			wordsMap[str] = w
			return w
		}
	}
}

var createdStrings = make(map[string]bool)

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
