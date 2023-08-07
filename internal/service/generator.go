package service

import (
	"errors"
	"math/big"
	"strings"

	"crypto/rand"
	mrand "math/rand"

	"github.com/google/uuid"
	"github.com/hablof/generate-random-value/internal/models"
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

// func NewGenerator() Generator {

// }
type Generator struct{}

func (g *Generator) Generate(opts models.GenerateOptions) (string, error) {
	length := 1
	if opts.LengthSpecified {
		length = opts.Length
	} else {
		length = 1 + mrand.Intn(256)
	}

	if length < 1 || length > 256 {
		return "", ErrInvalidLength
	}

	switch opts.GenerationType {
	case models.StringLike:
		return g.generateWithCharset(stringCharset, length), nil

	case models.Num:
		return g.generateNum(length)

	case models.Guid:
		return strings.ToUpper(uuid.NewString()), nil

	case models.Alphanumeric:
		return g.generateWithCharset(stringCharset+numberCharset, length), nil

	case models.Specified:
		if !opts.CharsetSpecified || opts.Charset == "" {
			return "", ErrInvalidCharset
		}

		reducedCharset := g.reduceToUnique(opts.Charset)

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
	min := big.NewInt(10)
	max := big.NewInt(10)
	max = max.Exp(max, big.NewInt(int64(length)), nil)
	min = min.Exp(min, big.NewInt(int64(length-1)), nil)
	intervalRange := max.Sub(max, min)

	delta, err := rand.Int(rand.Reader, intervalRange)
	if err != nil {
		return "", err
	}

	return delta.Add(min, delta).String(), nil
}
