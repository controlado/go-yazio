package yazio

import (
	"fmt"
	"sync"
	"time"

	"github.com/controlado/go-yazio/internal/application"
)

type Token struct {
	sync.RWMutex
	expiresAt time.Time
	access    string
	refresh   string
}

func (t *Token) String() string {
	var tokenStatus = "Valid"

	if t.IsExpired() {
		tokenStatus = "Expired"
	}

	return fmt.Sprintf("Token(%v)", tokenStatus)
}

func (t *Token) Update(newToken application.Token) {
	t.Lock()
	defer t.Unlock()

	t.expiresAt = newToken.ExpiresAt()
	t.access = newToken.Access()
	t.refresh = newToken.Refresh()
}

func (t *Token) ExpiresAt() time.Time {
	t.RLock()
	defer t.RUnlock()
	return t.expiresAt
}

func (t *Token) Refresh() string {
	t.RLock()
	defer t.RUnlock()
	return t.refresh
}

func (t *Token) Access() string {
	t.RLock()
	defer t.RUnlock()
	return t.access
}

func (t *Token) Bearer() string {
	t.RLock()
	defer t.RUnlock()
	return fmt.Sprintf("Bearer %s", t.access)
}

// IsExpired reports whether the access token held
// by t has already expired relative to the current
// time.
func (t *Token) IsExpired() bool {
	timeNow := time.Now()

	t.RLock()
	defer t.RUnlock()

	return t.expiresAt.Before(timeNow)
}
