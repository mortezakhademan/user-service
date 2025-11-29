package genopenapi_test

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFormatValidate(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Format Format
		Valid  bool
	}{{
		Format: FormatJSON,
		Valid:  true,
	}, {
		Format: FormatYAML,
		Valid:  true,
	}, {
		Format: Format("unknown"),
		Valid:  false,
	}, {
		Format: Format(""),
		Valid:  false,
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(string(tc.Format), func(t *testing.T) {
			t.Parallel()

			err := tc.Format.Validate()
			switch {
			case tc.Valid && err != nil:
				t.Fatalf("expect no validation error, got: %s", err)
			case !tc.Valid && err == nil:
				t.Fatal("expect validation error, got nil")
			}
		})
	}
}

func TestFormatEncode(t *testing.T) {
	t.Parallel()

	type contentDecoder interface {
		Decode(v interface{}) error
	}

	testCases := [...]struct {
		Format     Format
		NewDecoder func(r io.Reader) contentDecoder
	}{{
		Format: FormatJSON,
		NewDecoder: func(r io.Reader) contentDecoder {
			return json.NewDecoder(r)
		},
	}, {
		Format: FormatYAML,
		NewDecoder: func(r io.Reader) contentDecoder {
			return yaml.NewDecoder(r)
		},
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(string(tc.Format), func(t *testing.T) {
			t.Parallel()

			expParams := map[string]string{
				"hello": "world",
			}

			var buf bytes.Buffer
			enc, err := tc.Format.NewEncoder(&buf)
			if err != nil {
				t.Fatalf("expect no encoder creating error, got: %s", err)
			}

			err = enc.Encode(expParams)
			if err != nil {
				t.Fatalf("expect no encoding error, got: %s", err)
			}

			gotParams := make(map[string]string)
			err = tc.NewDecoder(&buf).Decode(&gotParams)
			if err != nil {
				t.Fatalf("expect no decoding error, got: %s", err)
			}

			if !reflect.DeepEqual(expParams, gotParams) {
				t.Fatalf("expected: %+v, actual: %+v", expParams, gotParams)
			}
		})
	}
}
