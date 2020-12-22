package writers

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/base"
)

type tcpWriter struct {
	debug bool

	base.Writer

	listener net.Listener
	quitCh   chan bool
	//connsWG  sync.WaitGroup

	mu       sync.RWMutex
	conns    map[int]net.Conn
	connsIdx int

	worker *base.Worker
}

// NewCallbackAsyncWriter returns a callbackAsyncWriter.
func NewTCPWriter(address string, port int, debug bool) (sparalog.Writer, error) {
	if address == "0.0.0.0" {
		address = ""
	}

	address = fmt.Sprintf("%s:%d", address, port)

	netAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, fmt.Errorf("Invalid address: %s", err)
	}

	listener, err := net.Listen("tcp", netAddr.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to create listener: %s", err)
	}

	w := tcpWriter{
		debug: debug,

		listener: listener,
		quitCh:   make(chan bool),
	}

	w.worker = base.NewWorker(&w, 100)

	go w.serve()

	return &w, nil
}

func (w *tcpWriter) Close() {
	w.worker.Close(1)

	close(w.quitCh)
	//w.connsWG.Wait()
}

func (w *tcpWriter) Write(item sparalog.Item) {
	w.worker.Enqueue(item)
}

func (w *tcpWriter) ProcessQueueItem(item sparalog.Item) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, conn := range w.conns {
		_, err := conn.Write([]byte(item.ToString(true, true)))
		if err != nil {
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "error writing to client: ", err)
			}
		}
	}
}

func (srv *tcpWriter) serve() {
	for {
		select {

		case <-srv.quitCh:
			if srv.debug {
				srv.Feedback(sparalog.DebugLevel, "shutting down")
			}
			srv.listener.Close()
			return

		default:
			if srv.debug {
				srv.Feedback(sparalog.DebugLevel, "Listening for clients")
			}

			//srv.listener.SetDeadline(time.Now().Add(1e9))
			conn, err := srv.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				if srv.debug {
					srv.Feedback(sparalog.DebugLevel, "Failed to accept connection:", err)
				}
			}

			srv.mu.Lock()
			defer srv.mu.Unlock()

			srv.conns[srv.connsIdx] = conn
			srv.connsIdx++

			//srv.connsWG.Add(1)
			go func() {
				srv.handleConnection(conn, 0)

				srv.mu.Lock()
				defer srv.mu.Unlock()

				delete(srv.conns, srv.connsIdx)

				//srv.connsWG.Done()
			}()
		}
	}
}

func (srv *tcpWriter) handleConnection(conn net.Conn, id int) {
	if srv.debug {
		srv.Feedback(sparalog.DebugLevel, "Accepted connection from", conn.RemoteAddr())
	}

	defer func() {
		if srv.debug {
			srv.Feedback(sparalog.DebugLevel, "Closing connection from", conn.RemoteAddr())
		}
		conn.Close()
	}()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err != io.EOF {
				if srv.debug {
					srv.Feedback(sparalog.DebugLevel, "continue reading")
				}
				continue
			}

			if srv.debug {
				srv.Feedback(sparalog.DebugLevel, "Read error:", err.Error())
			}
			return
		}

		if n == 0 {
			if srv.debug {
				srv.Feedback(sparalog.DebugLevel, "Close connection from", conn.RemoteAddr())
			}
			return
		}
	}
}
