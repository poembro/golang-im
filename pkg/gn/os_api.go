package gn

type netpoll interface {
	accept() (nfd int, addr string, err error)
	closeFD(fd int) error
	getEvents(msec int) ([]event, error)
	closeFDRead(fd int) error
}
