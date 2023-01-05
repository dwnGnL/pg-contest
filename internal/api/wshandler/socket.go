package wshandler

import (
	"net/http"
	"strconv"

	apiModels "github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/lib/cachemap"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func (ws wsHandler) wsContest(c *gin.Context) {
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	contestID, err := strconv.ParseInt(c.Param("contestID"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		goerrors.Log().Print("upgrade:", err)
		conn.WriteMessage(websocket.CloseInternalServerErr, []byte(err.Error()))
		return
	}
	req := new(apiModels.WsRequest)
	err = conn.ReadJSON(req)
	if err != nil {
		conn.WriteMessage(websocket.CloseInternalServerErr, []byte(err.Error()))
		return
	}
	if req.Token != "token" {
		conn.WriteMessage(websocket.CloseInternalServerErr, []byte("token not valid"))
		return
	}
	contest, err := app.CheckAndReturnContestByUserID(contestID, 1)
	if err != nil {
		goerrors.Log().Print("CheckAndReturnContestByUserID:", err)
		conn.WriteMessage(websocket.CloseInternalServerErr, []byte(err.Error()))
		conn.Close()
		return
	}
	if contest.IsEnd {
		conn.WriteJSON(app.Generate(contestID))
		conn.Close()
		return
	}

	var group errgroup.Group
	// чтение
	group.Go(func() error {
		for {
			err := conn.ReadJSON(req)
			if err != nil {
				return err
			}

		}
	})

	// запись
	group.Go(func() error {
		switcher, ok := ws.contestMap.Load(contestID)
		if ok {
			conn.WriteJSON(app.Generate(contestID))
			switcher.subscribers.Add(conn)
			return nil
		}
		ch := app.GenerateAndProcessChan(contestID)
		subscriber := new(subscribers)
		subscriber.Add(conn)
		switcher = subscribeSwitcher{
			event:       ch,
			subscribers: subscriber,
		}
		ws.contestMap.Store(contestID, switcher)
		go switcher.ReceiveEvent()
		return nil
	})

	if err := group.Wait(); err != nil {
		goerrors.Log().WithError(err).Error("Stopping ws with error")
	}
}

func newWsHandler() *wsHandler {
	return &wsHandler{
		contestMap: cachemap.NewCacheMap[int64, subscribeSwitcher](),
	}
}

func GenRouting(r *gin.RouterGroup) {
	ws := newWsHandler()
	r.Any("/connect/:contestID", ws.wsContest)
}
