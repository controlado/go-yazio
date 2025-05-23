package assert

import (
	"encoding/json"
	"io"
	"reflect"
	"slices"
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("\ngot %v\nwant %v", got, want)
	}
}

func EqualSlices[S ~[]T, T comparable](t *testing.T, got, want S) {
	t.Helper()

	if !slices.Equal(got, want) {
		t.Fatalf("\ngot %v\nwant %v", got, want)
	}
}

func DeepEqual[T any](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %+v\nwant %+v", got, want)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("\nwant nil err\ngot: %v", err)
	}
}

func WantErr(t *testing.T, want bool, err error) {
	t.Helper()

	if want {
		if err == nil {
			t.Fatal("\nwant err\ngot nil")
		}
		t.SkipNow()
	} else {
		if err != nil {
			t.Fatalf("\nwant nil err\ngot: %v", err)
		}
	}
}

func NotNil(t *testing.T, got any) {
	t.Helper()

	if got == nil {
		t.Fatalf("\nwant non-nil\ngot nil")
	}
}

func WriteDTO(t *testing.T, w io.Writer, i any) {
	t.Helper()

	err := json.
		NewEncoder(w).
		Encode(i)

	NoError(t, err)
}

func DecodeDTO(t *testing.T, r io.Reader, o any) {
	t.Helper()

	err := json.
		NewDecoder(r).
		Decode(o)

	NoError(t, err)
}

func ToJSON(t *testing.T, r io.Reader) (out map[string]any) {
	t.Helper()
	DecodeDTO(t, r, &out)
	return out
}
