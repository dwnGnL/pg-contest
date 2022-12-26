package wshandler

import (
	"sync"

	"github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/lib/cachemap"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gorilla/websocket"
)

type wsHandler struct {
	contestMap *cachemap.CacheMaper[int64, subscribeSwitcher]
}

type subscribeSwitcher struct {
	event       <-chan models.WsResponse
	subscribers *subscribers
}

func (s subscribeSwitcher) ReceiveEvent() {
	for {
		resp, ok := <-s.event
		if !ok {
			break
		}
		s.subscribers.Each(func(_ int, conn *websocket.Conn) {

			err := conn.WriteJSON(resp)
			if err != nil {
				goerrors.Log().Warnf("writeJson err:%w", err)
			}
			if resp.ContestStatus == models.End {
				conn.Close()
			}
		})
		if resp.ContestStatus == models.End {
			break
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
