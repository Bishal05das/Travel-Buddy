package booking

import (
	"context"
	"errors"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type createbookingusecase struct {
	bookingRepo port.BookingRepository
	tourRepo    port.TourRepository
}

func NewCreateBookingUseCase(bookingRepo port.BookingRepository, tourRepo port.TourRepository) *createbookingusecase {
	return &createbookingusecase{
		bookingRepo: bookingRepo,
		tourRepo:    tourRepo,
	}
}

func (uc *createbookingusecase) Execute(ctx context.Context, req *domain.BookingRequest, userID *uuid.UUID, memberID *uuid.UUID) (*domain.BookingResponse, error) {
	tour, err := uc.tourRepo.GetByID(ctx, req.TourID)
	if err != nil {
		return nil, err
	}
	if tour.Status != "open" {
		return nil, errors.New("tour is not open for booking")
	}
	if tour.AvailableSeat <= 0 {
		return nil, errors.New("Not enough seats available")
	}
	if time.Now().After(tour.LastEnrollmentDate) {
		return nil, errors.New("enrollment deadline has passed")
	}
	//calculate price with discount
	totalPrice := tour.Price

	var customerID uuid.UUID

	if userID != nil {
		customerID, err = uc.bookingRepo.GetOrCreateCustomerByUser(ctx, *userID)
		if err != nil {
			return nil, err
		}
	} else if memberID != nil {
		if req.CustomerName == "" || req.CustomerEmail == "" || req.CustomerPhone == 0 {
			return nil, errors.New("customer details required for guest booking")
		}
		customer := &domain.Customer{
			Name:  req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		}
		if err := uc.bookingRepo.CreateCustomer(ctx, customer); err != nil {
			return nil, err
		}
		customerID = customer.CustomerID
	} else {
		return nil, errors.New("either user or member must be sppecified")
	}

	booking := &domain.Booking{
		CustomerID:     customerID,
		UserID:         *userID,
		MemberID:       *memberID,
		TourID:         req.TourID,
		NumberOfPeople: req.NumberOfPeople,
		TotalPrice:     totalPrice,
		Status:         "pending",
	}
	if err := uc.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}
	//update available seats
	newSeats := tour.AvailableSeat - req.NumberOfPeople
	if err := uc.tourRepo.UpdateAvailableSeats(ctx, tour.TourID, newSeats); err != nil {
		return nil, err
	}

	return uc.bookingRepo.GetByID(ctx, booking.BookingID)

}
