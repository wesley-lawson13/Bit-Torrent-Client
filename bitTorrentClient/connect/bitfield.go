package connect

type Bitfield []byte

func (bf Bitfield) hasPiece(i int) bool {

    byteIndex, byteOffset := i / 8, i % 8
    if (bf[byteIndex] >> (7-byteOffset))&1 == 1 {
        return true
    }

    return false
}

func (bf Bitfield) setPiece(i int) {

    byteIndex, byteOffset := i / 8, i % 8
    bf[byteIndex] |= 1 << (7-byteOffset)
}
