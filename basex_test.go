package basex

import (
	"crypto/rand"
	"strconv"
	"testing"

	basex "github.com/jan-bar/basex/eknkc_basex"
)

func test(t *testing.T, alphabet string) {
	eknkc_basex, err := basex.NewEncoding(alphabet)
	if err != nil {
		t.Fatal(err)
	}

	my_basex, err := NewEncoding(alphabet)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 16)
	for i := 0; i < 10240; i++ {
		_, _ = rand.Read(buf)
		s1 := eknkc_basex.Encode(buf)
		s2 := my_basex.Encode(buf)

		if s1 != s2 {
			t.Fatalf("Encode: buf:% x,s1:%s,s2:%s", buf, s1, s2)
		}
		for _, v1 := range s1 {
			if !strconv.IsPrint(v1) {
				t.Fatalf("s1:%s have not print", s1)
			}
		}
		for _, v2 := range s2 {
			if !strconv.IsPrint(v2) {
				t.Fatalf("s2:%s have not print", s2)
			}
		}

		b1, err := eknkc_basex.Decode(s1)
		if err != nil {
			t.Fatal(err)
		}

		b2, err := my_basex.Decode(s2)
		if err != nil {
			t.Fatal(err)
		}

		if string(b1) != string(b2) || string(buf) != string(b1) || string(buf) != string(b2) {
			t.Fatalf("Decode: buf:% x,b1:%s,b2:%s", buf, b1, b2)
		}
	}
}

func TestBasex(t *testing.T) {
	// test base2
	test(t, "01")
	// test base16
	test(t, "0123456789abcdef")
	// test base32
	test(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")
	// tst base64
	test(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
}
