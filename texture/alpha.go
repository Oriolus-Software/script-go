package texture

type AlphaMode interface {
	alphaValue() any
}

type opaque struct{}

func (o opaque) alphaValue() any {
	return "Opaque"
}

type blend struct{}

func (b blend) alphaValue() any {
	return "Blend"
}

type mask struct {
	Mask float32
}

func (m mask) alphaValue() any {
	return m
}

var (
	AlphaOpaque = opaque{}
	AlphaBlend  = blend{}
)

func AlphaMask(m float32) AlphaMode {
	return mask{Mask: m}
}
