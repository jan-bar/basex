package basex

import (
	"bytes"
	"errors"
	"math/big"
)

// Encoding is a custom base encoding defined by an alphabet.
// It should bre created using NewEncoding function
type Encoding struct {
	base        *big.Int
	alphabet    []rune
	alphabetMap map[rune]int
}

// NewEncoding returns a custom base encoder defined by the alphabet string.
// The alphabet should contain non-repeating characters.
// Ordering is important.
// Example alphabets:
//   - base2: 01
//   - base16: 0123456789abcdef
//   - base32: 0123456789ABCDEFGHJKMNPQRSTVWXYZ
//   - base62: 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
func NewEncoding(alphabet string) (*Encoding, error) {
	runes := []rune(alphabet)
	runeMap := make(map[rune]int)

	for i := 0; i < len(runes); i++ {
		if _, ok := runeMap[runes[i]]; ok {
			return nil, errors.New("ambiguous alphabet")
		}
		runeMap[runes[i]] = i
	}

	return &Encoding{
		base:        big.NewInt(int64(len(runes))),
		alphabet:    runes,
		alphabetMap: runeMap,
	}, nil
}

// Encode function receives a byte slice and encodes it to a string using the alphabet provided
func (e *Encoding) Encode(source []byte) string {
	if len(source) == 0 {
		return ""
	}

	var (
		res bytes.Buffer
		k   = 0
	)
	for ; source[k] == 0 && k < len(source)-1; k++ {
		res.WriteRune(e.alphabet[0])
	}

	var (
		mod big.Int

		sourceInt = new(big.Int).SetBytes(source)
	)
	for sourceInt.Uint64() > 0 {
		sourceInt.DivMod(sourceInt, e.base, &mod)
		res.WriteRune(e.alphabet[mod.Uint64()])
	}

	var (
		buf = res.Bytes()
		j   = len(buf) - 1
	)
	for k < j {
		// Reverse bytes
		buf[k], buf[j] = buf[j], buf[k]
		k++
		j--
	}
	return string(buf)
}

// Decode function decodes a string previously obtained from Encode, using the same alphabet and returns a byte slice
// In case the input is not valid an arror will be returned
func (e *Encoding) Decode(source string) ([]byte, error) {
	if len(source) == 0 {
		return []byte{}, nil
	}

	var (
		data = []rune(source)
		dest = big.NewInt(0)
	)
	for i := 0; i < len(data); i++ {
		value, ok := e.alphabetMap[data[i]]
		if !ok {
			return nil, errors.New("non Base Character")
		}
		dest.Mul(dest, e.base)
		if value > 0 {
			dest.Add(dest, big.NewInt(int64(value)))
		}
	}

	k := 0 // leading zeros
	for ; data[k] == e.alphabet[0] && k < len(data)-1; k++ {
	}
	buf := dest.Bytes()
	res := make([]byte, k, k+len(buf))
	return append(res, buf...), nil
}
