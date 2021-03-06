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
	base.Writer

	debug bool

	stateCb StateChangeCallback

	listener net.Listener
	quitCh   chan bool
	//connsWG  sync.WaitGroup

	mu        sync.RWMutex
	conns     map[int]net.Conn
	connsIdx  int
	connCount int

	worker *base.Worker
}

// StateChangeCallback is called (true) when first client connecting,
// and (false) when there are no more clients.
type StateChangeCallback func(bool)

// NewTCPWriter returns a tcpWriter.
func NewTCPWriter(address string, port int, debug bool, cb StateChangeCallback) (sparalog.Writer, error) {

	w := tcpWriter{
		debug:   debug,
		stateCb: cb,
		quitCh:  make(chan bool),
		conns:   make(map[int]net.Conn, 4),
	}

	w.worker = base.NewWorker(&w, 100)

	err := w.startServer(address, port)
	if err != nil {
		return nil, err
	}

	go w.serve()

	return &w, nil
}

func (w *tcpWriter) startServer(address string, port int) error {
	if address == "0.0.0.0" {
		address = ""
	}

	address = fmt.Sprintf("%s:%d", address, port)

	netAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return fmt.Errorf("Invalid address: %s", err)
	}

	listener, err := net.Listen("tcp", netAddr.String())
	if err != nil {
		return fmt.Errorf("Failed to create listener: %s", err)
	}
	w.listener = listener

	return nil
}

func (w *tcpWriter) Close() {
	w.worker.Close(1)

	close(w.quitCh)
	//w.connsWG.Wait()
}

func (w *tcpWriter) Write(item sparalog.Item) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.conns) == 0 {
		return
	}

	w.worker.Enqueue(item)
}

func (w *tcpWriter) ProcessQueueItem(item sparalog.Item) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	bb := []byte(item.ToString(true, true) + "\n")

	for _, conn := range w.conns {
		_, err := conn.Write(bb)
		if err != nil {
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "error writing to client: ", err)
			}
		}
	}
}

func (w *tcpWriter) serve() {
	for {
		select {

		case <-w.quitCh:
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "shutting down")
			}
			w.listener.Close()

			/*for _, conn := range w.conns {
				conn.Close()
			}*/

			return

		default:
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "Listening for clients")
			}

			//w.listener.SetDeadline(time.Now().Add(1e9))
			conn, err := w.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				if w.debug {
					w.Feedback(sparalog.DebugLevel, "Failed to accept connection:", err)
				}
			}

			// Add connection to the pool.
			w.mu.Lock()

			w.connsIdx++
			w.conns[w.connsIdx] = conn

			if w.debug {
				w.Feedback(sparalog.DebugLevel, "Accepting TCP client ", w.connsIdx)
			}

			w.connCount++

			w.mu.Unlock()

			if w.connCount == 1 {
				if w.stateCb != nil {
					w.stateCb(true)
				}
			}

			//w.connsWG.Add(1)
			go func(idx int) {
				w.handleConnection(conn, 0)

				// Delete connection from the pool.
				w.mu.Lock()

				delete(w.conns, idx)

				if w.debug {
					w.Feedback(sparalog.DebugLevel, "Dispose TCP client ", idx)
				}

				w.connCount--

				w.mu.Unlock()

				if w.connCount == 0 {
					if w.stateCb != nil {
						w.stateCb(false)
					}
				}

				//w.connsWG.Done()
			}(w.connsIdx)
		}
	}
}

func (w *tcpWriter) handleConnection(conn net.Conn, id int) {
	if w.debug {
		w.Feedback(sparalog.DebugLevel, "Accepted connection from ", conn.RemoteAddr())
	}

	defer func() {
		if w.debug {
			w.Feedback(sparalog.DebugLevel, "Closed connection from ", conn.RemoteAddr())
		}
		conn.Close()
	}()

	conn.Write([]byte("Logger streaming activated for TCP channel.\n"))

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		switch {
		case err != nil && (err != io.EOF):
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "Read error: ", err.Error())
			}

		case err == io.EOF:
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "continue reading")
			}
			fallthrough

		case n == 0:
			if w.debug {
				w.Feedback(sparalog.DebugLevel, "Closing connection from ", conn.RemoteAddr())
			}
			return
		}
	}
}
