package service

import (
	"testing"

	"github.com/hablof/generate-random-value/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate_ErrorsOnly(t *testing.T) {
	tests := []struct {
		name string
		opts models.GenerateOptions
		// want    string
		wantErr error
	}{
		{
			name: "StringLike, no err",
			opts: models.GenerateOptions{
				GenerationType: models.StringLike,
			},
			wantErr: nil,
		},
		{
			name: "Number, no err",
			opts: models.GenerateOptions{
				GenerationType: models.Num,
			},
			wantErr: nil,
		},
		{
			name: "StringLike, with charset, no err",
			opts: models.GenerateOptions{
				GenerationType:   models.StringLike,
				Charset:          "test",
				CharsetSpecified: true,
			},
			wantErr: nil,
		},
		{
			name: "Number, with charset, no err",
			opts: models.GenerateOptions{
				GenerationType:   models.Num,
				Charset:          "test",
				CharsetSpecified: true,
			},
			wantErr: nil,
		},
		{
			name: "guid",
			opts: models.GenerateOptions{
				GenerationType: models.Guid,
			},
			wantErr: nil,
		},
		{
			name: "guid with ignored length",
			opts: models.GenerateOptions{
				GenerationType:  models.Guid,
				Length:          45,
				LengthSpecified: true,
			},
			wantErr: nil,
		},
		{
			name: "invalid charset",
			opts: models.GenerateOptions{
				GenerationType:   models.Specified,
				Charset:          "",
				CharsetSpecified: true,
			},
			wantErr: ErrInvalidCharset,
		},
		{
			name: "invalid charset (unspecified)",
			opts: models.GenerateOptions{
				GenerationType: models.Specified,
			},
			wantErr: ErrInvalidCharset,
		},
		{
			name: "invalid type",
			opts: models.GenerateOptions{
				GenerationType: "0",
			},
			wantErr: ErrInvalidType,
		},
		{
			name: "too small length",
			opts: models.GenerateOptions{
				GenerationType:  models.Alphanumeric,
				Length:          0,
				LengthSpecified: true,
			},
			wantErr: ErrInvalidLength,
		},
		{
			name: "too huge length",
			opts: models.GenerateOptions{
				GenerationType:  models.Alphanumeric,
				Length:          257,
				LengthSpecified: true,
			},
			wantErr: ErrInvalidLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Generator{}
			_, err := g.Generate(tt.opts)

			// assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGenerator_reduceToUnique(t *testing.T) {
	checkUnique := func(input string, output string) {
		outputRunes := make(map[rune]int)
		for _, r := range output {
			outputRunes[r]++
		}

		for _, r := range input {
			assert.Equal(t, 1, outputRunes[r])
		}
	}

	tests := []string{
		"test",
		"pohgfd",
		"vgf",
		"aaaaaaaaaaaaaaa",
		"pyfxtybt",
		"банальность",
		"поле )))))",
		"ХЪ{}1029403o5i",
	}
	g := Generator{}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			got := g.reduceToUnique(tt)
			checkUnique(tt, got)
		})
	}
}

func TestGenerator_generateNum(t *testing.T) {
	tests := []int{1, 2, 3, 4, 5, 100, 256, 256, 256, 150, 15, 45, 50}
	g := Generator{}
	for _, tt := range tests {
		got, err := g.generateNum(tt)
		assert.Equal(t, nil, err)
		assert.Equal(t, len(got), tt)
	}
}

func TestGenerator_generateWithCharset(t *testing.T) {
	checkIsInCharset := func(got, charset string) {
		charsetMap := make(map[rune]struct{})
		for _, r := range charset {
			charsetMap[r] = struct{}{}
		}
		for _, r := range got {
			_, ok := charsetMap[r]
			assert.Equal(t, true, ok)
		}
	}

	tests := []struct {
		name    string
		charset string
		length  int
	}{
		{
			name:    "1",
			charset: numberCharset,
			length:  15,
		},
		{
			name:    "2",
			charset: numberCharset,
			length:  50,
		},
		{
			name:    "3",
			charset: stringCharset,
			length:  50,
		},
		{
			name:    "4",
			charset: stringCharset + numberCharset,
			length:  50,
		},
		{
			name:    "5",
			charset: "kavabanga",
			length:  50,
		},
		{
			name:    "6",
			charset: "stringCharset+numberCharset",
			length:  50,
		},
	}
	g := Generator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.generateWithCharset(tt.charset, tt.length)
			checkIsInCharset(got, tt.charset)
			assert.Equal(t, tt.length, len(got))
		})
	}
}
