package assert

import (
	"testing"
)

func TestEqualSlicesItems(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name       string
		shouldFail bool
		got, want  []int
	}{
		{
			name: "wrong order with same items",
			got:  []int{1, 2, 3},
			want: []int{3, 2, 1},
		},
		{
			name:       "with different items",
			shouldFail: true,
			got:        []int{4, 5, 6},
			want:       []int{3, 2, 1},
		},
		{
			name:       "with different sizes",
			shouldFail: true,
			got:        []int{3, 2},
			want:       []int{3, 2, 1},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			mt := new(mockedTest)
			EqualSlicesItems(mt, tb.got, tb.want)

			if mt.failed != tb.shouldFail {
				t.Errorf("\nwant fail=%v\ngot fail=%v", tb.shouldFail, mt.failed)
			}
		})
	}
}
