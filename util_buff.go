package starGo

import (
	"bytes"
	"encoding/binary"
)

func Int32ToBytes(n int32, bigEndian bool) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bigEndian {
		_ = binary.Write(bytesBuffer, binary.BigEndian, n)
	} else {
		_ = binary.Write(bytesBuffer, binary.LittleEndian, n)
	}

	return bytesBuffer.Bytes()
}

func BytesToInt32(data []byte, bigEndian bool) int32 {
	bytesBuffer := bytes.NewBuffer(data)

	var result int32
	if bigEndian {
		_ = binary.Read(bytesBuffer, binary.BigEndian, &result)
	} else {
		_ = binary.Read(bytesBuffer, binary.LittleEndian, &result)
	}

	return result
}
