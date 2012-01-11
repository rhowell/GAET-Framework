package gaet

import "appengine"

var TestState = []string{"pass", "fail"}

const (
    PASS = iota
    FAIL
)

type Test struct {
    Context appengine.Context
    Output string
    Status string
}

func (t *Test) Fail(output string) {
    t.Status = TestState[FAIL]
    t.Output = output
}

func (t *Test) Pass(output string) {
    t.Status = TestState[PASS]
    t.Output = output
}

func (t *Test) IsStatusSet() bool {
    return 0 != len(t.Status)
}

type TestListEntry struct {
    Name string
    Test func(*Test)
}
