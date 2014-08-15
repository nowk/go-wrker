package wrker

type Job interface {
	Do() error
}
