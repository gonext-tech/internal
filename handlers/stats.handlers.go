package handlers

import (
	"strconv"
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type StatsService interface {
	GetALL() ([]models.Stats, error)
	GetMonthly(month, year string) (models.Stats, error)
	Create(models.Stats) (models.Stats, error)
	Update(models.Stats) (models.Stats, error)
	Delete(models.Stats) error
}

type StatsHandler struct {
	StatsServices StatsService
}

func NewStatsHandler(ss StatsService) *StatsHandler {
	return &StatsHandler{
		StatsServices: ss,
	}
}

func (ss *StatsHandler) DashboardPage(c echo.Context) error {
	isError = false
	month := c.QueryParam("month")
	year := c.QueryParam("year")
	today := time.Now()
	if month == "" {
		month = strconv.Itoa(int(today.Month()))
	}
	if year == "" {
		year = strconv.Itoa(today.Year())
	}
	stats, err := ss.StatsServices.GetMonthly(month, year)
	if err != nil {
		setFlashmessages(c, "error", errorMsg)
	}
	var params models.ParamResponse
	if month != "" {
		params.Month = month
	}
	if year != "" {
		params.Year = year
	}

	titlePage := "Dashboard"
	return renderView(c, subscription_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		subscription_views.List(titlePage, stats, params),
	))
}

func (ss *StatsHandler) GetStatsApi(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	_, ok := sess.Values[user_id_key].(uint)
	if !ok {
		setFlashmessages(c, "error", "user is not authenticated")
		response := map[string]interface{}{"message": "something went wrong"}
		return c.JSON(400, response)
	}
	stats, err := ss.StatsServices.GetALL()
	if err != nil {
		setFlashmessages(c, "error", "can't fetch stats")
		response := map[string]interface{}{"message": "can't fetch stats"}
		return c.JSON(400, response)
	}
	return c.JSON(200, stats)
}
