package task

type ticket struct {
	id   string
	exec func() error
}
