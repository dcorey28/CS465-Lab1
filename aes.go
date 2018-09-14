package aes

func ffAdd(a, b byte) byte {
	return 0x00
}

func xime(b byte) byte {
	return 0x00
}

func ffMultiply(a, b byte) byte {
	return 0x00
}

func subWord(word uint32) uint32 {
	return 0x000000
}

func rotWord(word uint32) uint32 {
	return 0x000000
}

func keyExpansion(key []byte) []uint32 {
	return []uint32{}
}

func subBytes(state [][]byte) [][]byte {
	return [][]byte{}
}

func shiftRows(state [][]byte) [][]byte {
	return [][]byte{}
}

func mixColumns(state [][]byte) [][]byte {
	return [][]byte{}
}

func addRoundKey(state [][]byte, w []uint32) [][]byte {
	return [][]byte{}
}
