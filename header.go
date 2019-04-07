package main

type Elink struct {
	ElinkHead
	Data      []byte
	DeCryData []byte
}

type ElinkHead struct {
	Flag uint32
	Len  uint32
}
