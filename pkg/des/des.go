package des

func Permutation(input uint64, permutationIndexes []uint64) uint64 {
	var maskInput uint64 = 0x80_00_00_00_00_00_00_00

	var ret uint64 = 0

	for i := 0; i < len(permutationIndexes); i++ {
		if input&(maskInput>>(permutationIndexes[i]-1)) > 0 {
			ret |= maskInput >> i
		}
	}

	return ret
}

func InitialPermutation(input uint64) uint64 {
	var permutationIndexes = []uint64{58, 50, 42, 34, 26, 18, 10, 2, 60, 52, 44, 36, 28, 20, 12, 4, 62, 54, 46, 38, 30, 22, 14, 6, 64, 56, 48, 40, 32, 24, 16, 8, 57, 49, 41, 33, 25, 17, 9, 1, 59, 51, 43, 35, 27, 19, 11, 3, 61, 53, 45, 37, 29, 21, 13, 5, 63, 55, 47, 39, 31, 23, 15, 7}
	return Permutation(input, permutationIndexes)
}

func FinalPermutation(input uint64) uint64 {
	var permutationIndexes = []uint64{40, 8, 48, 16, 56, 24, 64, 32, 39, 7, 47, 15, 55, 23, 63, 31, 38, 6, 46, 14, 54, 22, 62, 30, 37, 5, 45, 13, 53, 21, 61, 29, 36, 4, 44, 12, 52, 20, 60, 28, 35, 3, 43, 11, 51, 19, 59, 27, 34, 2, 42, 10, 50, 18, 58, 26, 33, 1, 41, 9, 49, 17, 57, 25}
	return Permutation(input, permutationIndexes)
}

// The expansion function is interpreted as for the initial and
// final permutations. Note that some bits from the input are duplicated at
// the output; e.g. the fifth bit of the input is duplicated in both the
// sixth and eighth bit of the output. Thus, the 32-bit half-block is expanded to 48 bits.
func Expansion(input uint64) uint64 {
	var expansionIndexes = []uint64{32, 1, 2, 3, 4, 5, 4, 5, 6, 7, 8, 9, 8, 9, 10, 11, 12, 13, 12, 13, 14, 15, 16, 17, 16, 17, 18, 19, 20, 21, 20, 21, 22, 23, 24, 25, 24, 25, 26, 27, 28, 29, 28, 29, 30, 31, 32, 1}
	return Permutation(input, expansionIndexes)
}

// Shuffles the bits of a 32-bit half-block
func Shuffle(input uint64) uint64 {
	var index = []uint64{16, 7, 20, 21, 29, 12, 28, 17, 1, 15, 23, 26, 5, 18, 31, 10, 2, 8, 24, 14, 32, 27, 3, 9, 19, 13, 30, 6, 22, 11, 4, 25}
	return Permutation(input, index)
}

func PermutedChoice1(input uint64) uint64 {
	var index_perm = []uint64{57, 49, 41, 33, 25, 17, 9, 1, 58, 50, 42, 34, 26, 18, 10, 2, 59, 51, 43, 35, 27, 19, 11, 3, 60, 52, 44, 36, 63, 55, 47, 39, 31, 23, 15, 7, 62, 54, 46, 38, 30, 22, 14, 6, 61, 53, 45, 37, 29, 21, 13, 5, 28, 20, 12, 4}
	return Permutation(input, index_perm)
}

// This permutation selects the 48-bit subkey for each
// round from the 56-bit key-schedule state.
// This permutation will ignore 8 bits below:
// Permuted Choice 2 "PC-2" Ignored bits 9, 18, 22, 25, 35, 38, 43, 54.
func PermutedChoice2(input uint64) uint64 {
	var permutIndexes = []uint64{14, 17, 11, 24, 1, 5, 3, 28, 15, 6, 21, 10, 23, 19, 12, 4, 26, 8, 16, 7, 27, 20, 13, 2, 41, 52, 31, 37, 47, 55, 30, 40, 51, 45, 33, 48, 44, 49, 39, 56, 34, 53, 46, 42, 50, 36, 29, 32}
	return Permutation(input, permutIndexes)
}

// Create 16 subkeys, each of which is 48-bits long.
func GenerateSubkeys(key uint64) []uint64 {
	var ret = make([]uint64, 16)

	// In k we will have 56-bit permutation. Only the msb 56 bit are considered.
	var k = PermutedChoice1(key)

	var c0 = k & GenerateMask(28)
	var d0 = k << 28 & GenerateMask(28)

	var shift_left_per_round = []uint8{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}

	var ci = c0
	var di = d0
	for num_round := 0; num_round < 16; num_round++ {
		ci = LeftShift(ci,
			shift_left_per_round[num_round],
			28)

		di = LeftShift(di,
			shift_left_per_round[num_round],
			28)

		var ci_di = ci | (di >> 28)     // ci_di contains a 56-bit partial key
		var ki = PermutedChoice2(ci_di) // ki contains the 48-bit subkey

		ret[num_round] = ki
	}

	return ret
}

// / Generates a mask with the specified number of 1's on the msb position
func GenerateMask(num_bits uint8) uint64 {
	var ret uint64 = 0

	var partial_mask uint64 = 0x80_00_00_00_00_00_00_00
	for i := 0; i < int(num_bits); i++ {
		ret |= partial_mask >> i
	}

	return ret
}

func LeftShift(input uint64, shift_amount uint8, num_bits uint8) uint64 {
	var left = input << shift_amount
	var right = input >> (num_bits - shift_amount)

	return (left | right) & GenerateMask(num_bits)
}

func S_Box(sBoxNum int, x int, y int) uint8 {

	var S1 = [][]uint8{{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
		{0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
		{4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
		{15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13}}

	var S2 = [][]uint8{{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
		{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
		{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
		{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9}}

	var S3 = [][]uint8{{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
		{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
		{13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
		{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12}}

	var S4 = [][]uint8{{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
		{13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
		{10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
		{3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14}}

	var S5 = [][]uint8{{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
		{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
		{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
		{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3}}

	var S6 = [][]uint8{{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
		{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
		{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
		{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13}}

	var S7 = [][]uint8{{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
		{13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
		{1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
		{6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12}}

	var S8 = [][]uint8{{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
		{1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
		{7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
		{2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11}}

	var sBoxes = [][][]uint8{S1, S2, S3, S4, S5, S6, S7, S8}

	return sBoxes[sBoxNum][y][x]
}

// Operates on two blocks, one of 32 bits and one of
// 48 bits, and produces a block of 32 bits.
func Feistel(input uint64, subkey uint64) uint64 {
	// Expansion
	var expanded = Expansion(input)

	// Key mixing
	var mixed = expanded ^ subkey

	// Substitution
	var ret uint64 = 0
	for numBox := 0; numBox < 8; numBox++ {
		// Only consider the 6 bits for each box input
		var boxInput = mixed >> (58 - 6*numBox)
		boxInput &= 0x3F

		var x = (boxInput & 0x1E) >> 1
		var y = (boxInput&0x20)>>4 | (boxInput & 0x01)

		var boxOutput = S_Box(numBox, int(x), int(y))

		ret |= uint64(boxOutput) << (60 - (4 * numBox))
	}

	// Permutation
	var index_perm = []uint64{16, 7, 20, 21, 29, 12, 28, 17, 1, 15, 23, 26, 5, 18, 31, 10, 2, 8, 24, 14, 32, 27, 3, 9, 19, 13, 30, 6, 22, 11, 4, 25}
	return Permutation(ret, index_perm)
}

func Encrypt(input uint64, key uint64) uint64 {
	return EncryptOrDecrypt(input, key, true)
}

func Decrypt(input uint64, key uint64) uint64 {
	return EncryptOrDecrypt(input, key, false)
}

// Both encryption and decryption procedures share the
// same algorithm. It only varies in the order the partial
// keys are applied.
func EncryptOrDecrypt(input uint64, key uint64, encrypt bool) uint64 {
	input = InitialPermutation(input)

	var l0 uint64 = input & GenerateMask(32)
	var r0 uint64 = input << 32 & GenerateMask(32)

	var subkeys = GenerateSubkeys(key)

	var li = l0
	var ri = r0

	for numRound := range subkeys {
		var index = numRound
		if !encrypt {
			index = 15 - numRound
		}
		var kn = subkeys[index]

		var prevLi = li
		li = ri
		ri = prevLi ^ Feistel(ri, kn)
	}

	return FinalPermutation(ri | (li >> 32))
}
