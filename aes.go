package aes

import (
	"encoding/binary"
)

func ffAdd(a, b byte) byte {
	return a ^ b
}

func xtime(b byte) byte {
	// check the 7th bit, if it is set then we will have to reduce after bit shifting
	shouldReduce := (b & 0x80) != 0

	a := b << 1

	if shouldReduce {
		a = a ^ 0x1b
	}

	return a
}

func ffMultiply(a, b byte) byte {
	result := byte(0x00)
	mask := byte(0x01)
	intermediate := a

	for mask > 0 {
		// mask b to see if it has a bit set at each position
		if (mask & b) != 0 {
			// if the bit is set then we want to XOR that position'state intermediate value to the result
			result = result ^ intermediate
		}

		// calculate a new intermediate based on the old one
		intermediate = xtime(intermediate)

		// shift the mask to the next position
		mask = mask << 1
	}

	return result
}

func subWord(word uint32) uint32 {
	var output uint32
	mask := uint32(0x0000000f)

	for i := 0; i < 4; i++ {
		col := word & mask
		word = word >> 4

		row := word & mask
		word = word >> 4

		output += uint32(sbox[row][col]) << (8 * uint(i))
	}

	return output
}

func rotWord(word uint32) uint32 {
	temp := word >> 24
	word = word << 8
	return word + temp
}

func keyExpansion(key []byte) []uint32 {
	Nk := len(key) / 4
	Nr := Nk + 6

	w := make([]uint32, 4*(Nr+1))

	for i := 0; i < Nk; i++ {
		j := 4 * i
		w[i] = binary.BigEndian.Uint32(key[j : j+4])
	}

	for i := Nk; i < len(w); i++ {
		temp := w[i-1]

		if (i % Nk) == 0 {
			temp = subWord(rotWord(temp)) ^ rcon[i/Nk]
		} else if Nk > 6 && (i%Nk) == 4 {
			temp = subWord(temp)
		}

		w[i] = w[i-Nk] ^ temp
	}

	return w
}

func subBytes(state [][]byte) [][]byte {
	for i, substate := range state {
		for j, cell := range substate {
			col := cell & 0x0f
			cell = cell >> 4
			row := cell & 0x0f
			state[i][j] = sbox[row][col]
		}
	}
	return state
}

func shiftRows(state [][]byte) [][]byte {
	for i, row := range state {
		temp := make([]byte, len(row))
		copy(temp, row)
		for j := 0; j < 4; j++ {
			state[i][j] = temp[(j+i)%4]
		}
	}
	return state
}

func mixColumns(state [][]byte) [][]byte {
	sp := makeEmptyState()

	for col := 0; col < 4; col++ {
		sp[0][col] = ffMultiply(0x02, state[0][col]) ^ ffMultiply(0x03, state[1][col]) ^ state[2][col] ^ state[3][col]
		sp[1][col] = state[0][col] ^ ffMultiply(0x02, state[1][col]) ^ ffMultiply(0x03, state[2][col]) ^ state[3][col]
		sp[2][col] = state[0][col] ^ state[1][col] ^ ffMultiply(0x02, state[2][col]) ^ ffMultiply(0x03, state[3][col])
		sp[3][col] = ffMultiply(0x03, state[0][col]) ^ state[1][col] ^ state[2][col] ^ ffMultiply(0x02, state[3][col])
	}

	return sp
}

func addRoundKey(state [][]byte, w []uint32) [][]byte {

	for col := 0; col < 4; col++ {
		keyfragment := make([]byte, 4)
		binary.BigEndian.PutUint32(keyfragment, w[col])
		for row := 0; row < 4; row++ {
			state[row][col] = state[row][col] ^ keyfragment[row]
		}
	}
	return state
}

func invSubBytes(state [][]byte) [][]byte {
	for i, substate := range state {
		for j, cell := range substate {
			col := cell & 0x0f
			cell = cell >> 4
			row := cell & 0x0f
			state[i][j] = invsbox[row][col]
		}
	}
	return state
}

func invShiftRows(state [][]byte) [][]byte {
	for i, row := range state {
		temp := make([]byte, len(row))
		copy(temp, row)
		for j := 0; j < 4; j++ {
			state[i][(j+i)%4] = temp[j]
		}
	}
	return state
}

func invMixColumns(s [][]byte) [][]byte {
	sp := makeEmptyState()

	for col := 0; col < 4; col++ {
		sp[0][col] = ffMultiply(0x0e, s[0][col]) ^ ffMultiply(0x0b, s[1][col]) ^ ffMultiply(0x0d, s[2][col]) ^ ffMultiply(0x09, s[3][col])
		sp[1][col] = ffMultiply(0x09, s[0][col]) ^ ffMultiply(0x0e, s[1][col]) ^ ffMultiply(0x0b, s[2][col]) ^ ffMultiply(0x0d, s[3][col])
		sp[2][col] = ffMultiply(0x0d, s[0][col]) ^ ffMultiply(0x09, s[1][col]) ^ ffMultiply(0x0e, s[2][col]) ^ ffMultiply(0x0b, s[3][col])
		sp[3][col] = ffMultiply(0x0b, s[0][col]) ^ ffMultiply(0x0d, s[1][col]) ^ ffMultiply(0x09, s[2][col]) ^ ffMultiply(0x0e, s[3][col])
	}

	return sp
}

func cipher(in []byte, w []uint32) []byte {
	state := toState(in)

	Nr := (len(w) - 1) / 4

	state = addRoundKey(state, w[:4])

	for i := 1; i <= Nr; i++ {
		state = subBytes(state)
		state = shiftRows(state)

		if i != Nr {
			state = mixColumns(state)
		}

		state = addRoundKey(state, w[i*4:(i+1)*4])
	}

	return fromState(state)
}

func inverseCipher(in []byte, w []uint32) []byte {
	state := toState(in)

	Nr := (len(w) - 1) / 4

	state = addRoundKey(state, w[Nr*4:(Nr+1)*4])

	for round := Nr - 1; round >= 0; round-- {
		state = invShiftRows(state)
		state = invSubBytes(state)
		state = addRoundKey(state, w[round*4:(round+1)*4])
		if round != 0 {
			state = invMixColumns(state)
		}
	}

	return fromState(state)
}

func makeEmptyState() [][]byte {
	return [][]byte{{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00}}
}

func toState(in []byte) [][]byte {
	s := makeEmptyState()

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			s[row][col] = in[row+4*col]
		}
	}

	return s
}

func fromState(s [][]byte) []byte {
	out := make([]byte, 16)
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			out[row+4*col] = s[row][col]
		}
	}
	return out
}
