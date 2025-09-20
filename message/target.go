package message

var (
	Myself Target = myself{}
	Parent Target = parent{}
)

type Target interface {
	toMessageTarget() any
}

type myself struct{}

func (m myself) toMessageTarget() any {
	return "Myself"
}

type parent struct{}

func (p parent) toMessageTarget() any {
	return "Parent"
}

type ChildByIndex uint32

func (c ChildByIndex) toMessageTarget() any {
	return struct {
		ChildByIndex uint32
	}{
		ChildByIndex: uint32(c),
	}
}

type Cockpit uint8

func (c Cockpit) toMessageTarget() any {
	return struct {
		Cockpit uint8
	}{
		Cockpit: uint8(c),
	}
}

type Broadcast struct {
	AcrossCouplings bool `msgpack:"across_couplings"`
	IncludeSelf     bool `msgpack:"include_self"`
}

func (b Broadcast) toMessageTarget() any {
	return struct {
		Broadcast Broadcast
	}{
		Broadcast: b,
	}
}

type AcrossCoupling struct {
	// Coupling is the name of the coupling to send the message across.
	// Either "front" or "rear".
	Coupling string `msgpack:"coupling"`
	Cascade  bool   `msgpack:"cascade"`
}

func (a AcrossCoupling) toMessageTarget() any {
	return struct {
		AcrossCoupling AcrossCoupling
	}{
		AcrossCoupling: a,
	}
}
