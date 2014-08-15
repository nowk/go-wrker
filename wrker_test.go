package wrker_test

import "errors"
import "testing"
import "github.com/bmizerany/assert"
import . "github.com/nowk/go-wrker"

type PassValue struct {
	Value  string
	Return chan string
}

func (p PassValue) Do() error {
	p.Return <- p.Value
	return nil
}

type PassErr struct {
	//
}

func (p PassErr) Do() error {
	return errors.New("Uh-oh!")
}

func TestWrkpool(t *testing.T) {
	jobs := make(chan Job, 10)

	pool := NewPool(3)
	errs := pool.Dispatch(jobs)

	res := make(chan string)

	jobs <- PassValue{"foo", res}
	jobs <- PassValue{"bar", res}
	jobs <- PassValue{"baz", res}
	jobs <- PassValue{"qux", res}
	jobs <- PassErr{}

	assert.Equal(t, <-res, "foo")
	assert.Equal(t, <-res, "bar")
	assert.Equal(t, <-res, "baz")
	assert.Equal(t, <-res, "qux")
	assert.Equal(t, <-errs, errors.New("Uh-oh!"))
}
