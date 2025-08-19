package msgpack

// MessagePack format constants
const (
	// fixint
	PositiveFixintMax = 0x7f
	NegativeFixintMin = 0xe0

	// nil
	Nil = 0xc0

	// bool
	False = 0xc2
	True  = 0xc3

	// bin
	Bin8  = 0xc4
	Bin16 = 0xc5
	Bin32 = 0xc6

	// float
	Float32 = 0xca
	Float64 = 0xcb

	// uint
	Uint8  = 0xcc
	Uint16 = 0xcd
	Uint32 = 0xce
	Uint64 = 0xcf

	// int
	Int8  = 0xd0
	Int16 = 0xd1
	Int32 = 0xd2
	Int64 = 0xd3

	// fixstr
	FixstrMask = 0xa0
	FixstrMax  = 0x1f

	// str
	Str8  = 0xd9
	Str16 = 0xda
	Str32 = 0xdb

	// array
	FixarrayMask = 0x90
	FixarrayMax  = 0x0f
	Array16      = 0xdc
	Array32      = 0xdd

	// map
	FixmapMask = 0x80
	FixmapMax  = 0x0f
	Map16      = 0xde
	Map32      = 0xdf
)
