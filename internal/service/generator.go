package service

import (
	"errors"
	"math/big"
	"strings"

	"crypto/rand"
	mrand "math/rand"

	"github.com/google/uuid"
)

const (
	numberCharset = "0123456789"
	stringCharset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
)

var (
	ErrInvalidLength  = errors.New("invalid length")
	ErrInvalidCharset = errors.New("invalid charset")
	ErrInvalidType    = errors.New("invalid type")
)

const (
	stringLike   = "string"
	num          = "numder"
	guid         = "guid"
	alphanumeric = "alphanumeric"
	specified    = "specified"
)

type GenerateOptions struct {
	generationType string

	charset          string
	charsetSpecified bool

	length          int
	lengthSpecified bool
}

func NewGenerateOptions(generationType string) GenerateOptions {
	return GenerateOptions{
		generationType:   generationType,
		charset:          "",
		charsetSpecified: false,
		length:           0,
		lengthSpecified:  false,
	}
}

func (g *GenerateOptions) SpecifyCharset(charset string) {
	g.charset = charset
	g.charsetSpecified = true
}

func (g *GenerateOptions) SpecifyLength(length int) {
	g.length = length
	g.lengthSpecified = true
}

// func NewGenerator() Generator {

// }
type Generator struct{}

func (g *Generator) Generate(opts GenerateOptions) (string, error) {
	length := 1
	if opts.lengthSpecified {
		length = opts.length
	} else {
		length = 1 + mrand.Intn(256)
	}

	if length < 1 || length > 256 {
		return "", ErrInvalidLength
	}

	switch opts.generationType {
	case stringLike:
		return g.generateWithCharset(stringCharset, length), nil

	case num:
		return g.generateNum(length)

	case guid:
		return strings.ToUpper(uuid.NewString()), nil

	case alphanumeric:
		return g.generateWithCharset(stringCharset+numberCharset, length), nil

	case specified:
		if !opts.charsetSpecified || opts.charset == "" {
			return "", ErrInvalidCharset
		}

		reducedCharset := g.reduceToUnique(opts.charset)

		return g.generateWithCharset(reducedCharset, length), nil
	}

	return "", ErrInvalidType
}

func (g *Generator) reduceToUnique(str string) string {
	m := make(map[rune]struct{}, len(str))
	for _, r := range str {
		m[r] = struct{}{}
	}
	outputRunes := make([]rune, 0, len(m))
	// reduced but shuffled
	for r := range m {
		outputRunes = append(outputRunes, r)
	}

	return string(outputRunes)
}

func (g *Generator) generateWithCharset(charset string, length int) string {
	runeset := []rune(charset)

	outputRunes := make([]rune, length)
	for i := range outputRunes {
		outputRunes[i] = runeset[mrand.Intn(len(runeset))]
	}

	return string(outputRunes)
}

func (g *Generator) generateNum(length int) (string, error) {
	base := big.NewInt(10)
	base = base.Exp(base, big.NewInt(int64(length)), nil)

	val, err := rand.Int(rand.Reader, base)
	if err != nil {
		return "", err
	}

	return val.String(), nil
}
