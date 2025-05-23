package assert

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"
)

func Equal[T comparable](t testing.TB, got, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("\ngot %v\nwant %v", got, want)
	}
}

func EqualSlicesItems[S ~[]T, T comparable](t testing.TB, got, want S) {
	t.Helper()

	var (
		gotLength  = len(got)
		wantLength = len(want)
	)

	if gotLength != wantLength {
		t.Fatalf("\nlen(got)=%d\nlen(want)=%d", gotLength, wantLength)
	}

	gotItems := make(map[T]int, gotLength)
	for _, v := range got {
		gotItems[v]++
	}

	for _, v := range want {
		if gotItems[v] == 0 {
			t.Fatalf("\nitem %v not received\ngot %v\nwant %v", v, got, want)
		}
		gotItems[v]--
	}

	for k, n := range gotItems {
		if n != 0 {
			t.Fatalf("\nreceived %d of %v which is not wanted\ngot %v\nwant %v", n, k, got, want)
		}
	}
}

func DeepEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %+v\nwant %+v", got, want)
	}
}

func NoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("\nwant nil err\ngot: %v", err)
	}
}

func WantErr(t testing.TB, want bool, err error) {
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

func NotNil(t testing.TB, got any) {
	t.Helper()

	if got == nil {
		t.Fatalf("\nwant non-nil\ngot nil")
	}
}

func WriteDTO(t testing.TB, w io.Writer, i any) {
	t.Helper()

	err := json.
		NewEncoder(w).
		Encode(i)

	NoError(t, err)
}

func DecodeDTO(t testing.TB, r io.Reader, o any) {
	t.Helper()

	err := json.
		NewDecoder(r).
		Decode(o)

	NoError(t, err)
}

func ToJSON(t testing.TB, r io.Reader) (out map[string]any) {
	t.Helper()
	DecodeDTO(t, r, &out)
	return out
}
