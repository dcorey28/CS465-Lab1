package aes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFFAdd(t *testing.T) {
	assert.Equal(t, byte(0xd4), ffAdd(0x57, 0x83))
}

func TestXtime(t *testing.T) {
	assert.Equal(t, byte(0xae), xtime(0x57))
	assert.Equal(t, byte(0x47), xtime(0xae))
	assert.Equal(t, byte(0x8e), xtime(0x47))
	assert.Equal(t, byte(0x07), xtime(0x8e))
}

func TestFFMultiply(t *testing.T) {
	assert.Equal(t, byte(0xfe), ffMultiply(0x57, 0x13))
}

func TestSubWord(t *testing.T) {
	assert.Equal(t, uint32(0x63cab704), subWord(0x00102030))
	assert.Equal(t, uint32(0x0953d051), subWord(0x40506070))
	assert.Equal(t, uint32(0xcd60e0e7), subWord(0x8090a0b0))
	assert.Equal(t, uint32(0xba70e18c), subWord(0xc0d0e0f0))
}

func TestRotWord(t *testing.T) {
	assert.Equal(t, uint32(0xcf4f3c09), rotWord(0x09cf4f3c))
	assert.Equal(t, uint32(0x6c76052a), rotWord(0x2a6c7605))
}

func TestKeyExpansion(t *testing.T) {
	key := []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}

	expandedKey := []uint32{0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x09cf4f3c,
		0xa0fafe17, 0x88542cb1, 0x23a33939, 0x2a6c7605,
		0xf2c295f2, 0x7a96b943, 0x5935807a, 0x7359f67f,
		0x3d80477d, 0x4716fe3e, 0x1e237e44, 0x6d7a883b,
		0xef44a541, 0xa8525b7f, 0xb671253b, 0xdb0bad00,
		0xd4d1c6f8, 0x7c839d87, 0xcaf2b8bc, 0x11f915bc,
		0x6d88a37a, 0x110b3efd, 0xdbf98641, 0xca0093fd,
		0x4e54f70e, 0x5f5fc9f3, 0x84a64fb2, 0x4ea6dc4f,
		0xead27321, 0xb58dbad2, 0x312bf560, 0x7f8d292f,
		0xac7766f3, 0x19fadc21, 0x28d12941, 0x575c006e,
		0xd014f9a8, 0xc9ee2589, 0xe13f0cc8, 0xb6630ca6}

	assert.Equal(t, expandedKey, keyExpansion(key))
}

func TestComponents(t *testing.T) {
	state := [][]byte{{0x19, 0xa0, 0x9a, 0xe9},
		{0x3d, 0xf4, 0xc6, 0xf8},
		{0xe3, 0xe2, 0x8d, 0x48},
		{0xbe, 0x2b, 0x2a, 0x08}}

	subExpected := [][]byte{{0xd4, 0xe0, 0xb8, 0x1e},
		{0x27, 0xbf, 0xb4, 0x41},
		{0x11, 0x98, 0x5d, 0x52},
		{0xae, 0xf1, 0xe5, 0x30}}

	shiftExpected := [][]byte{{0xd4, 0xe0, 0xb8, 0x1e},
		{0xbf, 0xb4, 0x41, 0x27},
		{0x5d, 0x52, 0x11, 0x98},
		{0x30, 0xae, 0xf1, 0xe5}}

	mixExpected := [][]byte{{0x04, 0xe0, 0x48, 0x28},
		{0x66, 0xcb, 0xf8, 0x06},
		{0x81, 0x19, 0xd3, 0x26},
		{0xe5, 0x9a, 0x7a, 0x4c}}

	roundExpected := [][]byte{{0xa4, 0x68, 0x6b, 0x02},
		{0x9c, 0x9f, 0x5b, 0x6a},
		{0x7f, 0x35, 0xea, 0x50},
		{0xf2, 0x2b, 0x43, 0x49}}

	w := []uint32{0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x09cf4f3c,
		0xa0fafe17, 0x88542cb1, 0x23a33939, 0x2a6c7605,
		0xf2c295f2, 0x7a96b943, 0x5935807a, 0x7359f67f,
		0x3d80477d, 0x4716fe3e, 0x1e237e44, 0x6d7a883b,
		0xef44a541, 0xa8525b7f, 0xb671253b, 0xdb0bad00,
		0xd4d1c6f8, 0x7c839d87, 0xcaf2b8bc, 0x11f915bc,
		0x6d88a37a, 0x110b3efd, 0xdbf98641, 0xca0093fd,
		0x4e54f70e, 0x5f5fc9f3, 0x84a64fb2, 0x4ea6dc4f,
		0xead27321, 0xb58dbad2, 0x312bf560, 0x7f8d292f,
		0xac7766f3, 0x19fadc21, 0x28d12941, 0x575c006e,
		0xd014f9a8, 0xc9ee2589, 0xe13f0cc8, 0xb6630ca6}

	assert.Equal(t, subExpected, subBytes(state), "subBytes failed")
	assert.Equal(t, shiftExpected, shiftRows(subExpected), "shiftRows failed")
	assert.Equal(t, mixExpected, mixColumns(shiftExpected), "mixColumns failed")
	assert.Equal(t, roundExpected, addRoundKey(mixExpected, w[4:8]), "addRoundKey failed")
}

func TestCipher(t *testing.T) {
	w := []uint32{0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x09cf4f3c,
		0xa0fafe17, 0x88542cb1, 0x23a33939, 0x2a6c7605,
		0xf2c295f2, 0x7a96b943, 0x5935807a, 0x7359f67f,
		0x3d80477d, 0x4716fe3e, 0x1e237e44, 0x6d7a883b,
		0xef44a541, 0xa8525b7f, 0xb671253b, 0xdb0bad00,
		0xd4d1c6f8, 0x7c839d87, 0xcaf2b8bc, 0x11f915bc,
		0x6d88a37a, 0x110b3efd, 0xdbf98641, 0xca0093fd,
		0x4e54f70e, 0x5f5fc9f3, 0x84a64fb2, 0x4ea6dc4f,
		0xead27321, 0xb58dbad2, 0x312bf560, 0x7f8d292f,
		0xac7766f3, 0x19fadc21, 0x28d12941, 0x575c006e,
		0xd014f9a8, 0xc9ee2589, 0xe13f0cc8, 0xb6630ca6}

	in := []byte{0x32, 0x43, 0xf6, 0xa8, 0x88, 0x5a, 0x30, 0x8d,
		0x31, 0x31, 0x98, 0xa2, 0xe0, 0x37, 0x07, 0x34}

	result := []byte{0x39, 0x25, 0x84, 0x1d, 0x02, 0xdc, 0x09, 0xfb,
		0xdc, 0x11, 0x85, 0x97, 0x19, 0x6a, 0x0b, 0x32}

	out := cipher(in, w)

	assert.Equal(t, result, out)
}

func TestInvSubBytes(t *testing.T) {
	input := [][]byte{{0xd4, 0xe0, 0xb8, 0x1e},
		{0x27, 0xbf, 0xb4, 0x41},
		{0x11, 0x98, 0x5d, 0x52},
		{0xae, 0xf1, 0xe5, 0x30}}

	output := [][]byte{{0x19, 0xa0, 0x9a, 0xe9},
		{0x3d, 0xf4, 0xc6, 0xf8},
		{0xe3, 0xe2, 0x8d, 0x48},
		{0xbe, 0x2b, 0x2a, 0x08}}

	assert.Equal(t, output, invSubBytes(input))
}

func TestInvShiftRows(t *testing.T) {
	input := [][]byte{{0xd4, 0xe0, 0xb8, 0x1e},
		{0xbf, 0xb4, 0x41, 0x27},
		{0x5d, 0x52, 0x11, 0x98},
		{0x30, 0xae, 0xf1, 0xe5}}

	output := [][]byte{{0xd4, 0xe0, 0xb8, 0x1e},
		{0x27, 0xbf, 0xb4, 0x41},
		{0x11, 0x98, 0x5d, 0x52},
		{0xae, 0xf1, 0xe5, 0x30}}

	assert.Equal(t, output, invShiftRows(input))
}

func TestInvMixColumns(t *testing.T) {
	input := [][]byte{{0x04, 0xe0, 0x48, 0x28},
		{0x66, 0xcb, 0xf8, 0x06},
		{0x81, 0x19, 0xd3, 0x26},
		{0xe5, 0x9a, 0x7a, 0x4c}}

	output := [][]byte{{0xd4, 0xe0, 0xb8, 0x1e},
		{0xbf, 0xb4, 0x41, 0x27},
		{0x5d, 0x52, 0x11, 0x98},
		{0x30, 0xae, 0xf1, 0xe5}}

	assert.Equal(t, output, invMixColumns(input))
}

func TestInvCipher(t *testing.T) {
	w := []uint32{0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x09cf4f3c,
		0xa0fafe17, 0x88542cb1, 0x23a33939, 0x2a6c7605,
		0xf2c295f2, 0x7a96b943, 0x5935807a, 0x7359f67f,
		0x3d80477d, 0x4716fe3e, 0x1e237e44, 0x6d7a883b,
		0xef44a541, 0xa8525b7f, 0xb671253b, 0xdb0bad00,
		0xd4d1c6f8, 0x7c839d87, 0xcaf2b8bc, 0x11f915bc,
		0x6d88a37a, 0x110b3efd, 0xdbf98641, 0xca0093fd,
		0x4e54f70e, 0x5f5fc9f3, 0x84a64fb2, 0x4ea6dc4f,
		0xead27321, 0xb58dbad2, 0x312bf560, 0x7f8d292f,
		0xac7766f3, 0x19fadc21, 0x28d12941, 0x575c006e,
		0xd014f9a8, 0xc9ee2589, 0xe13f0cc8, 0xb6630ca6}

	in := []byte{0x39, 0x25, 0x84, 0x1d, 0x02, 0xdc, 0x09, 0xfb,
		0xdc, 0x11, 0x85, 0x97, 0x19, 0x6a, 0x0b, 0x32}

	expected := []byte{0x32, 0x43, 0xf6, 0xa8, 0x88, 0x5a, 0x30, 0x8d,
		0x31, 0x31, 0x98, 0xa2, 0xe0, 0x37, 0x07, 0x34}

	out := inverseCipher(in, w)

	assert.Equal(t, expected, out)
}

func TestCipher128(t *testing.T) {
	in := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	expected := []byte{0x69, 0xc4, 0xe0, 0xd8, 0x6a, 0x7b, 0x04, 0x30, 0xd8, 0xcd, 0xb7, 0x80, 0x70, 0xb4, 0xc5, 0x5a}

	w := keyExpansion(key)

	assert.Equal(t, expected, cipher(in, w))
}

func TestInvCipher128(t *testing.T) {
	expected := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	in := []byte{0x69, 0xc4, 0xe0, 0xd8, 0x6a, 0x7b, 0x04, 0x30, 0xd8, 0xcd, 0xb7, 0x80, 0x70, 0xb4, 0xc5, 0x5a}

	w := keyExpansion(key)

	assert.Equal(t, expected, inverseCipher(in, w))
}

func TestCipher192(t *testing.T) {
	in := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}
	expected := []byte{0xdd, 0xa9, 0x7c, 0xa4, 0x86, 0x4c, 0xdf, 0xe0, 0x6e, 0xaf, 0x70, 0xa0, 0xec, 0x0d, 0x71, 0x91}

	w := keyExpansion(key)

	assert.Equal(t, expected, cipher(in, w))
}

func TestInvCipher192(t *testing.T) {
	expected := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}
	in := []byte{0xdd, 0xa9, 0x7c, 0xa4, 0x86, 0x4c, 0xdf, 0xe0, 0x6e, 0xaf, 0x70, 0xa0, 0xec, 0x0d, 0x71, 0x91}

	w := keyExpansion(key)

	assert.Equal(t, expected, inverseCipher(in, w))
}

func TestCipher256(t *testing.T) {
	in := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	expected := []byte{0x8e, 0xa2, 0xb7, 0xca, 0x51, 0x67, 0x45, 0xbf, 0xea, 0xfc, 0x49, 0x90, 0x4b, 0x49, 0x60, 0x89}

	w := keyExpansion(key)

	assert.Equal(t, expected, cipher(in, w))
}

func TestInvCipher256(t *testing.T) {
	expected := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	in := []byte{0x8e, 0xa2, 0xb7, 0xca, 0x51, 0x67, 0x45, 0xbf, 0xea, 0xfc, 0x49, 0x90, 0x4b, 0x49, 0x60, 0x89}

	w := keyExpansion(key)

	assert.Equal(t, expected, inverseCipher(in, w))
}
