package api

import (
	"github.com/GoDev/Hotel-reservatrion/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrNotfound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnAthorized()
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{Type: "msg", Msg: "updated"})
}

// when admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.Getbookings(c.Context(), bson.M{})
	if err != nil {
		return ErrNotfound("bookings")
	}

	return c.JSON(bookings)
}

// when user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return ErrNotfound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnAthorized()
	}
	return c.JSON(booking)
}
