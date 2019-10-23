package starGo

import (
	"crypto/rand"
	"math"
)

func Int64() int64 {
	result, err := rand.Int(rand.Reader, maxBigInt64Edge)
	if err != nil {
		return math.MaxInt64
	}

	return result.Int64()
}

func Int64n(n int64) int64 {
	if n <= 0 {
		return 0
	}
	if n&(n-1) == 0 {
		return Int64() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := Int64()
	for v > max {
		v = Int64()
	}

	return v % n
}

func Uint64() uint64 {
	return uint64(Int64())>>31 | uint64(Int64())<<32
}

func Uint32() uint32 {
	return uint32(Int64() >> 31)
}

func Int32() int32 {
	return int32(Int64() >> 32)
}

func Int32n(n int32) int32 {
	if n <= 0 {
		return 0
	}
	if n&(n-1) == 0 {
		return Int32() & (n - 1)
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := Int32()
	for v > max {
		v = Int32()
	}
	return v % n
}

func Int() int {
	u := uint(Int64())
	return int(u << 1 >> 1)
}

func IntN(n int) int {
	if n <= 0 {
		return 0
	}
	if n <= 1<<31-1 {
		return int(Int32n(int32(n)))
	}
	return int(Int64n(int64(n)))
}

func Shuffle(n int, swap func(i, j int)) {
	if n <= 0 {
		return
	}
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(Int64n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(Int32n(int32(i + 1)))
		swap(i, j)
	}
}
