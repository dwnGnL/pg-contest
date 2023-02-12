package public

import (
	"sync"
	"time"

	"github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gorilla/websocket"
)

type subscribeSwitcher struct {
	event       <-chan models.WsResponse
	subscribers *subscribers
}

const writeWait = 10 * time.Second

func (s subscribeSwitcher) ReceiveEvent() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case resp, ok := <-s.event:
			if !ok {
				break
			}
			s.subscribers.Each(func(_ int, conn *websocket.Conn) {
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := conn.WriteJSON(resp)
				if err != nil {
					goerrors.Log().Warnf("writeJson err:%s", err.Error())
				}
				if resp.ContestStatus == models.End {
					conn.Close()
				}
			})
			if resp.ContestStatus == models.End {
				return
			}
		case <-ticker.C:
			s.subscribers.Each(func(_ int, conn *websocket.Conn) {
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			})
		}

	}
}

func (s *subscribers) Add(conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()
	s.wsConnections = append(s.wsConnections, conn)
}

func (s *subscribers) Each(fn func(i int, conn *websocket.Conn)) {
	s.Lock()
	defer s.Unlock()
	for i, v := range s.wsConnections {
		fn(i, v)
	}
}

type subscribers struct {
	sync.RWMutex
	wsConnections []*websocket.Conn
}
