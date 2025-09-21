package texture

import (
	"github.com/oriolus-software/script-go/assets"
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/lmath"
)

type CreationOptions struct {
	Width   int  `msgpack:"width"`
	Height  int  `msgpack:"height"`
	MipMaps bool `msgpack:"mipmaps"`
}

type Texture uint32

type Pixel struct {
	R uint8 `msgpack:"r"`
	G uint8 `msgpack:"g"`
	B uint8 `msgpack:"b"`
}

type Color struct {
	R uint8 `msgpack:"r"`
	G uint8 `msgpack:"g"`
	B uint8 `msgpack:"b"`
	A uint8 `msgpack:"a"`
}

func Create(opts CreationOptions) Texture {
	return Texture(create(ffi.Serialize(opts).ToPacked()))
}

func (t Texture) Dispose() {
	dispose(uint32(t))
}

func (t Texture) GetPixel(x, y int) Pixel {
	return getPixel(uint32(t), x, y)
}

func (t Texture) ApplyTo(target string) {
	applyTo(uint32(t), ffi.Serialize(target).ToPacked())
}

func (t Texture) Clear(color Color) {
	t.addAction(struct {
		Clear Color
	}{
		Clear: color,
	})
}

func (t Texture) DrawPixels(pixels []DrawPixel) {
	t.addAction(struct {
		DrawPixels []DrawPixel
	}{
		DrawPixels: pixels,
	})
}

func (t Texture) DrawRect(start, end lmath.UVec2, color Color) {
	type drawRect struct {
		Start lmath.UVec2 `msgpack:"start"`
		End   lmath.UVec2 `msgpack:"end"`
		Color Color       `msgpack:"color"`
	}

	t.addAction(struct {
		DrawRect drawRect
	}{
		DrawRect: drawRect{
			Start: start,
			End:   end,
			Color: color,
		},
	})
}

type DrawTextOptions struct {
	Font          assets.ContentId `msgpack:"font"`
	Text          string           `msgpack:"text"`
	TopLeft       lmath.UVec2      `msgpack:"top_left"`
	LetterSpacing uint32           `msgpack:"letter_spacing"`
	FullColor     *Color           `msgpack:"full_color"`
	AlphaMode     AlphaMode        `msgpack:"alpha_mode"`
}

func (t Texture) DrawText(options *DrawTextOptions) {
	type drawText struct {
		Font          assets.ContentId `msgpack:"font"`
		Text          string           `msgpack:"text"`
		TopLeft       lmath.UVec2      `msgpack:"top_left"`
		LetterSpacing uint32           `msgpack:"letter_spacing"`
		FullColor     *Color           `msgpack:"full_color"`
		AlphaMode     any              `msgpack:"alpha_mode"`
	}

	var alphaMode any
	if options.AlphaMode != nil {
		alphaMode = options.AlphaMode.alphaValue()
	} else {
		alphaMode = "Opaque"
	}

	t.addAction(struct {
		DrawText drawText
	}{
		DrawText: drawText{
			Font:          options.Font,
			Text:          options.Text,
			TopLeft:       options.TopLeft,
			LetterSpacing: options.LetterSpacing,
			FullColor:     options.FullColor,
			AlphaMode:     alphaMode,
		},
	})
}

func (t Texture) DrawScriptTexture(src Texture, options DrawTextureOptions) {
	type drawScriptTexture struct {
		Handle  uint32             `msgpack:"handle"`
		Options DrawTextureOptions `msgpack:"options"`
	}

	t.addAction(struct {
		DrawScriptTexture drawScriptTexture
	}{
		DrawScriptTexture: drawScriptTexture{
			Handle:  uint32(src),
			Options: options,
		},
	})
}

func (t Texture) Flush() {
	flushActions(uint32(t))
}

func (t Texture) addAction(action any) {
	addAction(uint32(t), ffi.Serialize(action).ToPacked())
}

type DrawTextureOptions struct {
	SourceRect lmath.Rectangle
	TargetRect lmath.Rectangle
}

type DrawPixel struct {
	Pos   lmath.UVec2 `msgpack:"pos"`
	Color Color       `msgpack:"color"`
}

//go:wasm-module textures
//export create
func create(opts uint64) uint32

//go:wasm-module textures
//export dispose
func dispose(texture uint32)

//go:wasm-module textures
//export add_action
func addAction(texture uint32, action uint64)

//go:wasm-module textures
//export get_pixel
func getPixel(texture uint32, x, y int) Pixel

//go:wasm-module textures
//export flush_actions
func flushActions(texture uint32)

//go:wasm-module textures
//export apply_to
func applyTo(texture uint32, target uint64)
