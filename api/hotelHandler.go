package api

import (
	"fmt"

	"github.com/GoDev/Hotel-reservatrion/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	roomStore  db.RoomStore
	hotelStore db.HotelStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}

	fmt.Println(qparams)

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
