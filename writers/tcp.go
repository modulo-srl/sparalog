package writers

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/modulo-srl/sparalog/logs"
)

type TcpWriter struct {
	Writer

	debug bool

	stateCb StateChangeCallback

	address string
	port    int

	listener net.Listener
	quitCh   chan bool
	//connsWG  sync.WaitGroup

	mu          sync.RWMutex
	conns       map[int]net.Conn
	connsIdx    int
	connCount   int
	pause       bool
	filter      string
	filterInput bool
}

// StateChangeCallback is called (true) when first client connecting,
// and (false) when there are no more clients.
type StateChangeCallback func(bool)

// Ritorna un nuovo TCP Writer.
// - debug: se attivo logga in feedback eventuali errori e cambi di stato.
func NewTCPWriter(address string, port int, debug bool, cb StateChangeCallback) (*TcpWriter, error) {

	w := TcpWriter{
		debug:   debug,
		address: address,
		port:    port,
		stateCb: cb,
		quitCh:  make(chan bool),
		conns:   make(map[int]net.Conn, 4),
	}

	return &w, nil
}

func (w *TcpWriter) startServer(address string, port int) error {
	if address == "0.0.0.0" {
		address = ""
	}

	address = fmt.Sprintf("%s:%d", address, port)

	netAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return fmt.Errorf("invalid address: %s", err)
	}

	listener, err := net.Listen("tcp", netAddr.String())
	if err != nil {
		return fmt.Errorf("failed to create listener: %s", err)
	}
	w.listener = listener

	return nil
}

func (w *TcpWriter) Start() error {
	err := w.startServer(w.address, w.port)
	if err != nil {
		return err
	}

	go w.serve()

	w.StartQueue(100, w.onQueueItem)

	return nil
}

func (w *TcpWriter) Stop() {
	w.StopQueue(1)

	close(w.quitCh)
	//w.connsWG.Wait()
}

func (w *TcpWriter) Write(item *logs.Item) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.conns) == 0 {
		return
	}

	w.Enqueue(item)
}

func (w *TcpWriter) onQueueItem(item *logs.Item) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.pause {
		return nil
	}

	s := item.ToString(true, true) + "\n"

	if w.filter != "" {
		i := strings.Index(s, w.filter)
		if i < 0 {
			return nil
		}

		s = s[0:i] + "\033[91m" + s[i:i+len(w.filter)] + "\033[39m" + s[i+len(w.filter):]
	}

	bb := []byte(s)
	for _, conn := range w.conns {
		_, err := conn.Write(bb)
		if err != nil {
			if w.debug {
				w.Feedback(logs.DebugLevel, "error writing to client: ", err)
			}
		}
	}

	return nil
}

func (w *TcpWriter) serve() {
	for {
		select {

		case <-w.quitCh:
			if w.debug {
				w.Feedback(logs.DebugLevel, "shutting down")
			}
			w.listener.Close()

			/*for _, conn := range w.conns {
				conn.Close()
			}*/

			return

		default:
			if w.debug {
				w.Feedback(logs.DebugLevel, "Listening for clients")
			}

			//w.listener.SetDeadline(time.Now().Add(1e9))
			conn, err := w.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				if w.debug {
					w.Feedback(logs.DebugLevel, "Failed to accept connection:", err)
				}
			}

			// Add connection to the pool.
			w.mu.Lock()

			w.connsIdx++
			w.conns[w.connsIdx] = conn

			if w.debug {
				w.Feedback(logs.DebugLevel, "Accepting TCP client ", w.connsIdx)
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
				w.handleConnection(conn)

				// Delete connection from the pool.
				w.mu.Lock()

				delete(w.conns, idx)

				if w.debug {
					w.Feedback(logs.DebugLevel, "Dispose TCP client ", idx)
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

func (w *TcpWriter) handleConnection(conn net.Conn) {
	if w.debug {
		w.Feedback(logs.DebugLevel, "Accepted connection from ", conn.RemoteAddr())
	}

	defer func() {
		if w.debug {
			w.Feedback(logs.DebugLevel, "Closed connection from ", conn.RemoteAddr())
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
				w.Feedback(logs.DebugLevel, "Read error: ", err.Error())
			}

		case err == io.EOF:
			if w.debug {
				w.Feedback(logs.DebugLevel, "Continue reading")
			}
			fallthrough

		case n == 0:
			if w.debug {
				w.Feedback(logs.DebugLevel, "Closing connection from ", conn.RemoteAddr())
			}
			return
		}

		s := string(buf[:n-1])
		switch {
		case w.filterInput:
			f := strings.TrimSpace(s)

			w.mu.Lock()
			w.filter = f
			w.pause = false
			w.mu.Unlock()

			if f == "" {
				conn.Write([]byte("Show all logs\n"))
			} else {
				conn.Write([]byte("Show logs for \"" + f + "\"\n"))
			}

			w.filterInput = false

		case s == "filter":
			w.mu.Lock()
			w.pause = true
			w.mu.Unlock()

			conn.Write([]byte("\nFilter: "))
			w.filterInput = true
		}
	}
}
