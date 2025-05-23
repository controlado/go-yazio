package assert

import "testing"

type mockedTest struct {
	testing.TB
	helper bool
	failed bool
}

func (mt *mockedTest) Helper()                   { mt.helper = true }
func (mt *mockedTest) Fatal(_ ...any)            { mt.failed = true }
func (mt *mockedTest) Fatalf(_ string, _ ...any) { mt.failed = true }
