package starGo

import (
	"crypto/rand"
	"math"
)

func RandInt64() int64 {
	result, err := rand.Int(rand.Reader, maxBigInt64Edge)
	if err != nil {
		return math.MaxInt64
	}

	return result.Int64()
}

func RandInt64n(n int64) int64 {
	if n <= 0 {
		return 0
	}
	if n&(n-1) == 0 {
		return RandInt64() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := RandInt64()
	for v > max {
		v = RandInt64()
	}

	return v % n
}

func RandUint64() uint64 {
	return uint64(RandInt64())>>31 | uint64(RandInt64())<<32
}

func RandUint32() uint32 {
	return uint32(RandInt64() >> 31)
}

func RandInt32() int32 {
	return int32(RandInt64() >> 32)
}

func RandInt32n(n int32) int32 {
	if n <= 0 {
		return 0
	}
	if n&(n-1) == 0 {
		return RandInt32() & (n - 1)
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := RandInt32()
	for v > max {
		v = RandInt32()
	}
	return v % n
}

func RandInt() int {
	u := uint(RandInt64())
	return int(u << 1 >> 1)
}

func RandIntN(n int) int {
	if n <= 0 {
		return 0
	}
	if n <= 1<<31-1 {
		return int(RandInt32n(int32(n)))
	}
	return int(RandInt64n(int64(n)))
}

func Shuffle(n int, swap func(i, j int)) {
	if n <= 0 {
		return
	}
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(RandInt64n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(RandInt32n(int32(i + 1)))
		swap(i, j)
	}
}
