package font

import (
	"github.com/oriolus-software/script-go/assets"
	"github.com/oriolus-software/script-go/internal/ffi"
)

type BitmapFont struct {
	BitmapFontProperties
	ContentId assets.ContentId
}

// / Properties of a bitmap font.
type BitmapFontProperties struct {
	/// The horizontal distance between letters.
	HorizontalDistance int32 `msgpack:"horizontal_distance"`
	/// The vertical size of the font.
	VerticalSize int32 `msgpack:"vertical_size"`
	/// The letters in the font.
	Letters map[string]FontLetter `msgpack:"letters"`
}

// / A letter in a bitmap font.
type FontLetter struct {
	/// The character represented by the letter.
	Character string `msgpack:"character"`
	/// The start of the letter in the texture.
	Start uint32 `msgpack:"start"`
	/// The width of the letter in the texture.
	Width uint32 `msgpack:"width"`
}

var bitmapFontCache = make(map[assets.ContentId]*BitmapFont)

func LoadBitmapFontProperties(contentId assets.ContentId) (*BitmapFont, bool) {
	if props, ok := bitmapFontCache[contentId]; ok {
		return props, true
	}

	ret := bitmapFontProperties(ffi.Serialize(contentId).ToPacked())

	if ret == 0 {
		return nil, false
	}

	props := ffi.Deserialize[BitmapFontProperties](ret)
	font := &BitmapFont{
		BitmapFontProperties: props,
		ContentId:            contentId,
	}
	bitmapFontCache[contentId] = font
	return font, true
}

func (font *BitmapFont) TextLen(text string, letterSpacing int) int {
	ret := textLen(ffi.Serialize(font.ContentId).ToPacked(), ffi.Serialize(text).ToPacked(), letterSpacing)

	if ret == -1 {
		panic("bitmap font properties should be loaded and text_len should return valid value")
	}

	return ret
}

//go:wasm-module font
//export bitmap_font_properties
func bitmapFontProperties(font uint64) uint64

//go:wasm-module font
//export text_len
func textLen(font, text uint64, letterSpacing int) int
