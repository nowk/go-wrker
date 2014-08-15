# go-wrker

[![Build Status](https://travis-ci.org/nowk/go-wrker.svg?branch=master)](https://travis-ci.org/nowk/go-wrker)

Go worker

## Example

    type PassValue struct {
      Value  string
      Return chan string
    }

    func (p PassValue) Do() error {
      p.Return <- p.Value
      return nil
    }

    func main() {
      jobs := make(chan wrker.Job, 10)
      res := make(chan string)

      pool := wrker.NewPool(3)
      errs := wrker.Dispatch(jobs, pool)

      go func() {
        for {
          select {
          case str := <-res:
            log.Printf("received %s", str)
          case err := <-errs:
            log.Printf("error: %v", err)
          }
        }
      }()

      for i := 0; i < 10; i++ {
        jobs <- PassValue{fmt.Sprintf("-- %d", i), res}
      }

      sig := make(chan os.Signal, 2)
      signal.Notify(sig, os.Interrupt, os.Kill)
      <-sig
      <-sig
    }

