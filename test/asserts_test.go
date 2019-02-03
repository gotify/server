package test_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
)

type obj struct {
	Test string
	ID   int
}

type fakeTesting struct {
	hasErrors bool
}

func (t *fakeTesting) Errorf(format string, args ...interface{}) {
	t.hasErrors = true
}

func Test_BodyEquals(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteString(`{"ID": 2, "Test": "asd"}`)

	fakeTesting := &fakeTesting{}

	test.BodyEquals(fakeTesting, &obj{ID: 2, Test: "asd"}, recorder)
	assert.False(t, fakeTesting.hasErrors)
}

func Test_BodyEquals_failing(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteString(`{"ID": 3, "Test": "asd"}`)

	fakeTesting := &fakeTesting{}

	test.BodyEquals(fakeTesting, &obj{ID: 2, Test: "asd"}, recorder)
	assert.True(t, fakeTesting.hasErrors)
}

func Test_UnreaableReader(t *testing.T) {
	_, err := ioutil.ReadAll(test.UnreadableReader())
	assert.Error(t, err)
}
