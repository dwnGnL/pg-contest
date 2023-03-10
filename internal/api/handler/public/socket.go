package public

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/internal/service"

	"github.com/dwnGnL/pg-contests/internal/api/models"
	apiModels "github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/lib/goerrors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
const (
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func (ws publicHandler) wsContest(c *gin.Context) {
	goerrors.Log().Println("start socket")

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
	goerrors.Log().Println("start Upgrade")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		goerrors.Log().Print("upgrade:", err)
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}
	goerrors.Log().Println("read token")

	req := new(apiModels.WsRequest)
	err = conn.ReadJSON(req)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}
	goerrors.Log().Println("check token ")

	tokenDetails, err := ws.jwtClient.ExtractTokenMetadata("Bearer " + req.Token)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("token not valid"))
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}
	goerrors.Log().Println("CheckAndReturnContestByUserID")

	contest, err := app.CheckAndReturnContestByUserID(contestID, tokenDetails.ID)
	if err != nil {
		if errors.Is(err, service.SubscribeErr) {
			conn.WriteJSON(models.WsResponse{ErrorCode: 2, ErrorMess: err.Error()})
			conn.Close()
			return
		}
		goerrors.Log().Print("CheckAndReturnContestByUserID:", err)
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return
	}
	if *contest.IsEnd {
		conn.WriteJSON(app.Generate(contestID))
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	// ????????????
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	go func() {
		for {
			err := conn.ReadJSON(req)
			// if errors.Is(err,websocket.ErrBadHandshake)

			if err != nil {
				conn.WriteJSON(models.WsResponse{ErrorCode: 1, ErrorMess: "???????????? ???????????? ?????????????? " + err.Error()}) // any model
				goerrors.Log().WithError(err).Error("ReadJSON error")
				conn.Close()
				break
			}
			if req.AnswerID == 0 || req.QuestionID == 0 {
				continue
			}
			//???????????????? ?????????? ???? ?????????????? ???????????? ?? ????
			//?????????????????? ?????????? ???????????? ???? ?????????????? ???????????? ?????? ???????????? ???????? ???? 0 ???? question.time
			curTime, err := app.CalculateTimeForQuestion(contestID, req.QuestionID)
			if err != nil {
				conn.WriteJSON(models.WsResponse{ErrorCode: 1, ErrorMess: "?????????????????? ?????????????? ???????????????? " + err.Error()}) // any model
				goerrors.Log().WithError(err).Error("CalculateTimeForQuestion error")
				continue
			}
			var question *repository.Question
			for _, v := range contest.Questions {
				if v.ID == req.QuestionID {
					question = &v
				}
			}
			if question == nil {
				conn.WriteJSON(models.WsResponse{ErrorCode: 1, ErrorMess: "?????? ???????????? ?????????????? ?? ???????? ???????????????? "}) // any model
				goerrors.Log().Error("CalculateTimeForQuestion error")
				continue
			}
			if curTime+1 < 0 || curTime > question.Time {
				conn.WriteJSON(models.WsResponse{ErrorCode: 1, ErrorMess: "?????????? ?????????? ?????? ???? ?????????????? ?????? "}) // any model
				goerrors.Log().Error("CalculateTimeForQuestion error")
				continue
			}
			userAnswer := repository.UserAnswers{
				UserID:     tokenDetails.ID,
				ContestID:  contestID,
				QuestionID: req.QuestionID,
				AnswerID:   req.AnswerID,
				Time:       curTime,
			}
			err = app.SubmitAnswer(&userAnswer)
			if err != nil {
				conn.WriteJSON(models.WsResponse{ErrorCode: 1, ErrorMess: "SubmitAnswer error " + err.Error()}) // any model
				goerrors.Log().WithError(err).Error("SubmitAnswer error")
				continue
			}
		}
	}()

	// ????????????
	go func() {
		switcher, ok := ws.contestMap.Load(contestID)
		if ok && !switcher.End {
			conn.WriteJSON(app.Generate(contestID))
			switcher.subscribers.Add(conn)
			return
		}
		ch := app.GenerateAndProcessChan(contestID)
		subscriber := new(subscribers)
		subscriber.Add(conn)
		switcher = &subscribeSwitcher{
			event:       ch,
			subscribers: subscriber,
		}
		ws.contestMap.Store(contestID, switcher)
		go switcher.ReceiveEvent()
		return
	}()

}
