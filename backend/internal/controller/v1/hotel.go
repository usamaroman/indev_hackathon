package v1

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/middleware"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/request"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/response"
	_ "github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/service"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type hotelRoutes struct {
	log   *slog.Logger
	valid *validator.Validate

	hotelService service.Hotel
}

func newHotelRoutes(log *slog.Logger, g *gin.RouterGroup, hotelService service.Hotel, authMiddleware *middleware.AuthMiddleware) {
	log = log.With(slog.String("component", "hotel routes"))

	v := validator.New()

	r := &hotelRoutes{
		log:          log,
		valid:        v,
		hotelService: hotelService,
	}

	g.POST("/rooms", authMiddleware.PublicMiddleware(), r.getAvailableRooms)
	g.GET("/rooms/:id", authMiddleware.HotelsOnly(), r.getRoomByID)
	g.POST("/rooms/reserve", authMiddleware.CustomersOnly(), r.reserveRoom)
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

	hasRsv, err := r.hotelService.RoomHasReservations(c, roomID)
	if err != nil {
		r.log.Error("failed to get if room has reservations", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.RoomInfo{
		Room:     room,
		IsActive: hasRsv,
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
