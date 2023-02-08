package public

import (
	"github.com/dwnGnL/pg-contests/internal/repository"
	"net/http"
	"strconv"

	apiModels "github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func (ws publicHandler) wsContest(c *gin.Context) {
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
	tokenDetails, err := ws.jwtClient.ExtractTokenMetadata("Bearer" + req.Token)
	if err != nil {
		conn.WriteMessage(websocket.CloseInternalServerErr, []byte("token not valid"))
		return
	}
	contest, err := app.CheckAndReturnContestByUserID(contestID, tokenDetails.ID)
	if err != nil {
		goerrors.Log().Print("CheckAndReturnContestByUserID:", err)
		conn.WriteMessage(websocket.CloseInternalServerErr, []byte(err.Error()))
		conn.Close()
		return
	}
	if *contest.IsEnd {
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
				goerrors.Log().WithError(err).Error("ReadJSON error")
				//conn.WriteJSON()  // any model
				continue
			}
			//записать ответ на текущий вопрос в бд
			//посчитать время ответа на текущий вопрос оно должно быть от 0 до question.time
			curTime, err := app.CalculateTimeForQuestion(contestID, req.QuestionID)
			if err != nil {
				goerrors.Log().WithError(err).Error("CalculateTimeForQuestion error")
				continue
			}
			if curTime < 0 || curTime > contest.Questions[req.QuestionID].Time {
				//conn.WriteJSON()  // any model
				// не успел ответить
				return nil
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
				//conn.WriteJSON()  // any model
				goerrors.Log().WithError(err).Error("SubmitAnswer error")
				continue
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
