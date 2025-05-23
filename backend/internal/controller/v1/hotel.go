package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/middleware"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/request"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/response"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/service"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/box"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type hotelRoutes struct {
	log   *slog.Logger
	valid *validator.Validate

	hotelService service.Hotel
	box          *box.Client
	cache        *redis.Client
}

func newHotelRoutes(log *slog.Logger, g *gin.RouterGroup, hotelService service.Hotel, authMiddleware *middleware.AuthMiddleware, b *box.Client, redisClient *redis.Client) {
	log = log.With(slog.String("component", "hotel routes"))

	v := validator.New()

	r := &hotelRoutes{
		log:          log,
		valid:        v,
		hotelService: hotelService,
		box:          b,
		cache:        redisClient,
	}

	g.POST("/rooms", authMiddleware.PublicMiddleware(), r.getAvailableRooms)
	g.GET("/rooms/:id", authMiddleware.HotelsOnly(), r.getRoomByID)
	g.POST("/rooms/reserve", authMiddleware.CustomersOnly(), r.reserveRoom)
	g.PATCH("/rooms/reservations/:id", authMiddleware.HotelsOnly(), r.updateReservationStatus)
	g.POST("/rooms/light", authMiddleware.CustomersOnly(), r.roomLights)
	g.GET("/rooms/token", authMiddleware.CustomersOnly(), r.getRoomToken)
	g.GET("/rooms/dump_token", r.getToken)
	g.GET("/rooms/reservations/confirmed", authMiddleware.HotelsOnly(), r.getConfirmedReservations)
	g.GET("/rooms/reservations/checkedin", authMiddleware.HotelsOnly(), r.getCheckedInReservations)
}

// @Summary Получение доступных типов комнат
// @Description Получение доступных типов комнат
// @Tags отель
// @Accept json
// @Param input body request.GetRooms true "Тело запроса"
// @Produce json
// @Success 200 {object} response.GetAllRooms
// @Router /v1/hotel/rooms [post]
func (r *hotelRoutes) getAvailableRooms(c *gin.Context) {
	var req request.GetRooms

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	startDate, err := time.Parse("02.01.2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	endDate, err := time.Parse("02.01.2006", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rooms, err := r.hotelService.GetAvailableRooms(c, startDate, endDate)
	if err != nil {
		r.log.Error("failed to get available rooms", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.GetAllRooms{
		Rooms: rooms,
	})
}

// @Summary Получение информации по комнате для админки
// @Description Получение информации по комнате для админки
// @Security BearerAuth
// @Tags отель
// @Produce json
// @Param id path string true "Идентификатор комнаты"
// @Success 200 {object} response.RoomInfo
// @Router /v1/hotel/rooms/{id} [get]
func (r *hotelRoutes) getRoomByID(c *gin.Context) {
	roomID := c.Param("id")

	room, err := r.hotelService.GetRoomByID(c, roomID)
	if err != nil {
		r.log.Error("failed to get room", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	resStatus, err := r.hotelService.RoomReservationStatus(c, roomID)
	if err != nil {
		r.log.Error("failed to get room reservation status", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.RoomInfo{
		Room:   room,
		Status: resStatus,
	})
}

// @Summary Бронь
// @Description Бронь комнаты
// @Security BearerAuth
// @Tags отель
// @Accept json
// @Produce json
// @Param input body request.ReserveRoom true "Тело запроса"
// @Success 204
// @Router /v1/hotel/rooms/reserve [post]
func (r *hotelRoutes) reserveRoom(c *gin.Context) {
	var reservationDto request.ReserveRoom
	if err := c.ShouldBindJSON(&reservationDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	startDate, err := time.Parse("02.01.2006", reservationDto.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	endDate, err := time.Parse("02.01.2006", reservationDto.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDCtx)

	if err = r.hotelService.CreateReservation(c, &service.CreateReservationInput{
		UserID:    userID,
		RoomType:  reservationDto.RoomType,
		StartDate: startDate,
		EndDate:   endDate,
	}); err != nil {
		r.log.Error("failed to create reservation", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Свет
// @Description Тогл света в комнате
// @Security BearerAuth
// @Tags отель
// @Accept json
// @Produce json
// @Param input body request.Light true "Тело запроса"
// @Success 204
// @Router /v1/hotel/rooms/light [post]
func (r *hotelRoutes) roomLights(c *gin.Context) {
	var req request.Light

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := validator.New().Struct(req); err != nil {
		r.log.Error("failed to validate request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	state, err := strconv.ParseBool(req.State)
	if err != nil {
		r.log.Error("failed to parse state", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	lightName := r.box.GetBleName()

	rsv, err := r.hotelService.GetUserCurrentReservation(c, c.GetInt64(middleware.UserIDCtx))
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			r.log.Error("user does not have a reservation", logger.Error(err))
			c.JSON(http.StatusForbidden, gin.H{
				"error": "user does not have a reservation",
			})
			return
		}

		r.log.Error("failed to get user current reservation", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if rsv.RoomID != lightName {
		r.log.Error("user does not have a light", logger.Error(err))
		c.JSON(http.StatusForbidden, gin.H{
			"error": "user does not have a light",
		})
		return
	}

	switch state {
	case true:
		if err := r.box.LightOn(); err != nil {
			r.log.Error("failed to turn on the light", logger.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	case false:
		if err := r.box.LightOff(); err != nil {
			r.log.Error("failed to turn off the light", logger.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.Status(http.StatusNoContent)
}

// @Summary Обновление статуса брони
// @Description Обновление статуса брони
// @Security BearerAuth
// @Tags отель
// @Accept json
// @Produce json
// @Param id path string true "ID брони"
// @Param input body request.UpdateReservationStatus true "Тело запроса"
// @Success 204
// @Router /v1/hotel/rooms/reservations/{id} [patch]
func (r *hotelRoutes) updateReservationStatus(c *gin.Context) {
	reservationID := c.Param("id")

	var req request.UpdateReservationStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := validator.New().Struct(req); err != nil {
		r.log.Error("failed to validate request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := r.hotelService.UpdateReservationStatus(c, reservationID, types.StringToReservationType[req.Status]); err != nil {
		r.log.Error("failed to update reservation status", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Получение токена комнаты [ХАРДКОД]
// @Description Получение токена комнаты [ХАРДКОД]
// @Tags отель
// @Produce json
// @Success 200 {object} response.Token
// @Router /v1/hotel/rooms/dump_token [get]
func (r *hotelRoutes) getToken(c *gin.Context) {
	c.JSON(http.StatusOK, response.Token{Token: "tZAZwckrp2ZH4fCM"})
}

// @Summary Получение токена комнаты
// @Description Получение токена комнаты
// @Security BearerAuth
// @Tags отель
// @Accept json
// @Produce json
// @Success 200 {object} response.Token
// @Router /v1/hotel/rooms/token [get]
func (r *hotelRoutes) getRoomToken(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDCtx)

	res, err := r.hotelService.GetUserCurrentReservation(c, userID)
	if err != nil {
		r.log.Error("failed to get user current reservation", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	token, err := r.cache.Get(c, fmt.Sprint(res.RoomID)).Result()
	if err != nil {
		r.log.Error("failed to get room token", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.Token{Token: token})
}

// @Summary Получение бронирований, в которые еще не заехали
// @Description Получение бронирований, в которые еще не заехали
// @Security BearerAuth
// @Tags отель
// @Accept json
// @Produce json
// @Success 200 {object} []entity.ReservationInfo
// @Router /v1/hotel/rooms/reservations/confirmed [get]
func (r *hotelRoutes) getConfirmedReservations(c *gin.Context) {
	res, err := r.hotelService.GetReservationsByStatus(c, types.Confirmed.String())
	if err != nil {
		r.log.Error("failed to get confirmed reservations", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Получение бронирований, в которые заехали
// @Description Получение бронирований, в которые заехали
// @Security BearerAuth
// @Tags отель
// @Accept json
// @Produce json
// @Success 200 {object} []entity.ReservationInfo
// @Router /v1/hotel/rooms/reservations/confirmed [get]
func (r *hotelRoutes) getCheckedInReservations(c *gin.Context) {
	res, err := r.hotelService.GetReservationsByStatus(c, types.CheckedIn.String())
	if err != nil {
		r.log.Error("failed to get confirmed reservations", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}
