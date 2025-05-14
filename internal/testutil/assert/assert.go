package assert

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("\ngot %v\nwant %v", got, want)
	}
}

func DeepEqual(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %+v\nwant %+v", got, want)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got an error: %v", err)
	}
}

func NotNil(t *testing.T, v any) {
	t.Helper()

	if v == nil {
		t.Fatalf("\nwant no-pointer\ngot nil")
	}
}

func DecodeDTO(t *testing.T, r io.Reader, o any) {
	t.Helper()

	err := json.NewDecoder(r).Decode(o)
	NoError(t, err)
}

func JSON(t *testing.T, r io.Reader) map[string]any {
	t.Helper()

	var out map[string]any
	err := json.NewDecoder(r).Decode(&out)
	NoError(t, err)

	return out
}

func Write(t *testing.T, w io.Writer, i any) {
	t.Helper()

	err := json.NewEncoder(w).Encode(i)
	NoError(t, err)
}
