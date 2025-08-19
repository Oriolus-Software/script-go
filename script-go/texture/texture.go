package texture

import (
	"github.com/oriolus-software/script-go/assets"
	"github.com/oriolus-software/script-go/internal/ffi"
	"github.com/oriolus-software/script-go/mathf"
)

type CreationOptions struct {
	Width   int  `msgpack:"width"`
	Height  int  `msgpack:"height"`
	MipMaps bool `msgpack:"mipmaps"`
}

type Texture uint32
type Pixel struct {
	R uint8
	G uint8
	B uint8
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
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

func (t Texture) DrawRect(start, end mathf.UVec2, color Color) {
	type drawRect struct {
		Start mathf.UVec2
		End   mathf.UVec2
		Color Color
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

func (t Texture) DrawText(font assets.ContentId, text string, topLeft mathf.UVec2, letterSpacing uint32, fullColor *Color, alphaMode any) {
	var alpha any
	switch alphaMode := alphaMode.(type) {
	case string:
		switch alphaMode {
		case "opaque":
			alpha = "Opaque"
		case "blend":
			alpha = "Blend"
		default:
			panic("invalid alpha mode")
		}
	case float32:
		alpha = struct {
			Mask float32 `msgpack:"mask"`
		}{
			Mask: alphaMode,
		}
	default:
		panic("invalid alpha mode")
	}

	t.addAction(struct {
		DrawText drawText
	}{
		DrawText: drawText{
			Font:          font,
			Text:          text,
			TopLeft:       topLeft,
			LetterSpacing: letterSpacing,
			FullColor:     fullColor,
			Alpha:         alpha,
		},
	})
}

func (t Texture) DrawScriptTexture(src Texture, options drawTextureOptions) {
	t.addAction(struct {
		DrawScriptTexture drawScriptTexture
	}{
		DrawScriptTexture: drawScriptTexture{
			Handle:  uint32(src),
			Options: options,
		},
	})
}

func (t Texture) addAction(action any) {
	addAction(uint32(t), ffi.Serialize(action).ToPacked())
}

func (t Texture) Flush() {
	flushActions(uint32(t))
}

type drawText struct {
	Font          assets.ContentId
	Text          string
	TopLeft       mathf.UVec2
	LetterSpacing uint32
	FullColor     *Color
	Alpha         any
}

type drawTextureOptions struct {
	SourceRect *mathf.Rectangle
	TargetRect *mathf.Rectangle
}

type drawScriptTexture struct {
	Handle  uint32
	Options drawTextureOptions
}

type DrawPixel struct {
	Pos   mathf.UVec2
	Color Color
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
