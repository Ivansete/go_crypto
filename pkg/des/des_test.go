package des

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Feistel(t *testing.T) {
	// https://page.math.tu-berlin.de/~kant/teaching/hess/krypto-ws2006/des.htm
	var input_key uint64 = 0b00010011_00110100_01010111_01111001_10011011_10111100_11011111_11110001
	var subkeys = GenerateSubkeys(input_key)
	var expected uint64 = 0b0010001101001010101010011011101100000000000000000000000000000000
	var r0 uint64 = 0b11110000101010101111000010101010_00000000_00000000_00000000_00000000
	assert.Equal(t, expected, Feistel(r0, subkeys[0]))
}

func Test_InitialPermutation(t *testing.T) {
	// Checking first bit of the output
	var input uint64 = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000000
	var res = InitialPermutation(input)
	var expected uint64 = 0x80_00_00_00_00_00_00_00
	res &= expected
	assert.Exactly(t, expected, res)

	// Checking one bit from the middle of the output
	input = 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	res = InitialPermutation(input)
	expected = 0b00000000_00000000_00000000_00000000_00000001_00000000_00000000_00000000
	res &= expected
	assert.Exactly(t, expected, res)

	// Checking the last bit of the output
	input = 0b00000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	res = InitialPermutation(input)
	expected = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001
	res &= expected
	assert.Exactly(t, expected, res)
}

func Test_FinalPermutation(t *testing.T) {
	// Checking first bit of the output
	var input uint64 = 0b00000000_00000000_00000000_00000000_00000001_00000000_00000000_00000000
	var res = FinalPermutation(input)
	var expected uint64 = 0x80_00_00_00_00_00_00_00
	res &= expected
	assert.Exactly(t, expected, res)

	// Checking one bit from the middle of the output
	input = 0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	res = FinalPermutation(input)
	expected = 0b00000000_00000000_00000000_01000000_00000000_00000000_00000000_00000000
	res &= expected
	assert.Exactly(t, expected, res)

	// Checking the last bit of the output
	input = 0b00000000_00000000_00000000_10000000_00000000_00000000_00000000_00000000
	res = FinalPermutation(input)
	expected = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001
	res &= expected
	assert.Exactly(t, expected, res)
}

func Test_Initial_Final_Permutation(t *testing.T) {
	var input uint64 = 0b00100100_01000100_00100100_11110000_00011101_01111110_10000000_00001011
	assert.Exactly(t, input, FinalPermutation(InitialPermutation(input)))
}

func Test_Expansion(t *testing.T) {
	// https://page.math.tu-berlin.de/~kant/teaching/hess/krypto-ws2006/des.htm
	var input uint64 = 0b1111000010101010111100001010101000000000000000000000000000000000
	var expected uint64 = 0b011110100001010101010101011110100001010101010101_00000000_00000000
	assert.Equal(t, expected, Expansion(input))
}

func Test_SBox(t *testing.T) {
	assert.Equal(t, uint8(14), S_Box(0, 0, 0))
	assert.Equal(t, uint8(13), S_Box(0, 15, 3))
	assert.Equal(t, uint8(11), S_Box(7, 15, 3))
}

func Test_PermutedChoice1(t *testing.T) {
	// https://page.math.tu-berlin.de/~kant/teaching/hess/krypto-ws2006/des.htm
	var input uint64 = 0b00010011_00110100_01010111_01111001_10011011_10111100_11011111_11110001
	var value = PermutedChoice1(input)
	var expected uint64 = 0b1111000011001100101010101111010101010110011001111000111100000000

	assert.Exactly(t, expected, value)
}

func Test_PermutedChoice2(t *testing.T) {
	var input uint64 = 0b00000000_00000100_00000000_00000000_00000000_00000000_00000000_00000000
	var res = PermutedChoice2(input)
	var expected uint64 = 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000

	assert.Exactly(t, expected, res)
}

func Test_Subkeys_generation(t *testing.T) {
	// https://page.math.tu-berlin.de/~kant/teaching/hess/krypto-ws2006/des.htm
	var input_key uint64 = 0b00010011_00110100_01010111_01111001_10011011_10111100_11011111_11110001
	var subkeys = GenerateSubkeys(input_key)
	var expected uint64 = 0b1100101100111101100010110000111000010111111101010000000000000000
	assert.Equal(t, expected, subkeys[15])
}

func Test_Mask(t *testing.T) {
	assert.Equal(t, uint64(0), GenerateMask(0))
	assert.Equal(t, uint64(0x80_00_00_00_00_00_00_00), GenerateMask(1))
	assert.Equal(t, uint64(0xFFFFFFF000000000), GenerateMask(28))
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), GenerateMask(64))
}

func Test_LeftShift(t *testing.T) {
	var expected uint64 = 0x2000001000000000
	var input uint64 = 0x90_00_00_00_00_00_00_00
	assert.Equal(t, expected, LeftShift(input, 1, 28))

	expected = 0x4000002000000000
	input = 0x90_00_00_00_00_00_00_00
	assert.Equal(t, expected, LeftShift(input, 2, 28))
}

func Test_Encrypt(t *testing.T) {
	// https://page.math.tu-berlin.de/~kant/teaching/hess/krypto-ws2006/des.htm
	var input uint64 = 0x0123456789ABCDEF
	var key uint64 = 0x133457799BBCDFF1
	assert.Equal(t, uint64(0x85E813540F0AB405), Encrypt(input, key))
}

func Test_Encrypt_Decrypt(t *testing.T) {
	var val uint64 = 0x1
	var key uint64 = 0x3b3898371520f75e
	assert.Equal(t, uint64(0x1), Decrypt(Encrypt(val, key), key))
}
