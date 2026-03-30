package dtos

type ContainerActionRequest struct {
	ContainerID string
	Action      string // restart, stop, start
}

type LogStreamRequest struct {
	ContainerID string
	Tail        int
	Follow      bool
}
