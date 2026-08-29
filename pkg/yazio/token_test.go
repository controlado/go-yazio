package yazio

import (
	"fmt"
	"testing"

	"github.com/controlado/go-yazio/v3/internal/testutil/assert"
	"github.com/controlado/go-yazio/v3/internal/testutil/times"
	"github.com/google/uuid"
)

func TestToken_IsExpired(t *testing.T) {
	t.Parallel()

	var (
		randomAccessToken  = uuid.NewString()
		randomRefreshToken = uuid.NewString()

		tk = &Token{
			expiresAt: times.Past(), // invalid
			access:    randomAccessToken,
			refresh:   randomRefreshToken,
		}

		want = true
		got  = tk.IsExpired()
	)

	assert.Equal(t, got, want)
}

func TestToken_Bearer(t *testing.T) {
	t.Parallel()

	var (
		randomAccessToken = uuid.NewString()

		tk = &Token{
			expiresAt: times.Future(),
			access:    randomAccessToken,
		}

		want = fmt.Sprintf("Bearer %s", randomAccessToken)
		got  = tk.Bearer()
	)

	assert.Equal(t, got, want)
}

func TestToken_Refresh(t *testing.T) {
	t.Parallel()

	var (
		randomRefreshToken = uuid.NewString()

		tk = &Token{
			expiresAt: times.Future(),
			refresh:   randomRefreshToken,
		}
	)

	var (
		want = randomRefreshToken
		got  = tk.Refresh()
	)

	assert.Equal(t, got, want)
}

func TestToken_Access(t *testing.T) {
	t.Parallel()

	var (
		randomAccessToken = uuid.NewString()

		tk = &Token{
			expiresAt: times.Future(),
			access:    randomAccessToken,
		}
	)

	var (
		want = randomAccessToken
		got  = tk.Access()
	)

	assert.Equal(t, got, want)
}

func TestToken_Update(t *testing.T) {
	t.Parallel()

	var (
		newValidToken = func() *Token {
			return &Token{
				expiresAt: times.Future(),
				access:    uuid.NewString(),
				refresh:   uuid.NewString(),
			}
		}

		testBlocks = []struct {
			name     string
			wantErr  error
			newToken *Token
		}{
			{name: "using a valid new token", wantErr: nil, newToken: newValidToken()},
			{name: "using nil new token", wantErr: ErrTokenCannotBeNil, newToken: nil},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			tk := newValidToken()
			tkBearerSnapshot := tk.Bearer()

			err := tk.Update(tb.newToken)
			if tb.wantErr != nil {
				assert.Equal(t, err, tb.wantErr)
				assert.DeepEqual(t, tk.Bearer(), tkBearerSnapshot)
				return
			}
			assert.DeepEqual(t, tk, tb.newToken)
		})
	}
}

func TestToken_String(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name string
		tr   *Token
		want string
	}{
		{
			name: "expired token",
			tr:   &Token{expiresAt: times.Past()},
			want: "Token(Expired)",
		},
		{
			name: "expired token",
			tr:   &Token{expiresAt: times.Future()},
			want: "Token(Valid)",
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.tr.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
