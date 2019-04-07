package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

func serialPacket(elink Elink) []byte {
	header := ElinkHead{
		Flag: GlobalFlag,
		Len:  uint32(len(elink.Data)),
	}

	elink.ElinkHead = header

	data, _ := encodeHeader(header)
	res := make([]byte, 8+len(elink.Data))

	copy(res[:8], data)
	copy(res[8:], elink.Data)
	return res
}

func parseHeader(b []byte) (ElinkHead, error) {
	header := ElinkHead{}
	fmt.Println(hex.Dump(b))
	header.Flag = binary.BigEndian.Uint32(b[:4])
	header.Len = binary.BigEndian.Uint32(b[4:8])
	if header.Flag != GlobalFlag {
		errMsg := fmt.Sprintf("Parse Header Flag error: Recv Flag is %x, need Flag is %x\n", header.Flag, GlobalFlag)
		return ElinkHead{}, errors.New(errMsg)
	}
	return header, nil
}

func encodeHeader(header ElinkHead) ([]byte, error) {
	if header.Flag != GlobalFlag {
		errMsg := fmt.Sprintf("Encode Header Flag error: Current Flag is %x, need Flag is %x\n", header.Flag, GlobalFlag)
		return nil, errors.New(errMsg)
	}
	res := make([]byte, 8)
	binary.BigEndian.PutUint32(res[:4], header.Flag)
	binary.BigEndian.PutUint32(res[4:8], header.Len)
	return res, nil
}
