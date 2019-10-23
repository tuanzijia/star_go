package starGo

import (
	"bytes"
	"encoding/binary"
)

func Int32ToBytes(n int32, bigEndian bool) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bigEndian {
		binary.Write(bytesBuffer, binary.BigEndian, n)
	} else {
		binary.Write(bytesBuffer, binary.LittleEndian, n)
	}

	return bytesBuffer.Bytes()
}

func BytesToInt32(data []byte, bigEndian bool) int32 {
	bytesBuffer := bytes.NewBuffer(data)

	var result int32
	if bigEndian {
		binary.Read(bytesBuffer, binary.BigEndian, &result)
	} else {
		binary.Read(bytesBuffer, binary.LittleEndian, &result)
	}

	return result
}
