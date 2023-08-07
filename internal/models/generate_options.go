package models

const (
	StringLike   = "string"
	Num          = "number"
	Guid         = "guid"
	Alphanumeric = "alphanumeric"
	Specified    = "specified"
)

type GenerateOptions struct {
	GenerationType string

	Charset          string
	CharsetSpecified bool

	Length          int
	LengthSpecified bool
}

// func NewGenerateOptions(generationType string) GenerateOptions {
// 	return GenerateOptions{
// 		GenerationType:   generationType,
// 		Charset:          "",
// 		CharsetSpecified: false,
// 		Length:           0,
// 		LengthSpecified:  false,
// 	}
// }

func (g *GenerateOptions) SpecifyCharset(charset string) {
	g.Charset = charset
	g.CharsetSpecified = true
}

func (g *GenerateOptions) SpecifyLength(length int) {
	g.Length = length
	g.LengthSpecified = true
}
