package bookingusecase

import (
	"context"
	"errors"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type createbookingusecase struct {
	txManager   port.TxManager
	bookingRepo port.BookingRepository
	tourRepo    port.TourRepository
	paymentRepo port.PaymentRepository
}

func NewCreateBookingUseCase(txManager port.TxManager, bookingRepo port.BookingRepository, tourRepo port.TourRepository, paymentRepo port.PaymentRepository) *createbookingusecase {
	return &createbookingusecase{
		txManager:   txManager,
		bookingRepo: bookingRepo,
		tourRepo:    tourRepo,
		paymentRepo: paymentRepo,
	}
}

func (uc *createbookingusecase) Execute(ctx context.Context, req *domain.BookingRequest, userID *uuid.UUID, memberID *uuid.UUID) (*domain.BookingResponse, error) {

	var response *domain.BookingResponse

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		//starts by row level locking
		tour, err := uc.tourRepo.GetByIDForUpdate(txCtx, req.TourID)
		if err != nil {
			return  err
		}
		if tour.Status != "open" {
			return  errors.New("tour is not open for booking")
		}
		if tour.AvailableSeat <= req.NumberOfPeople {
			return  errors.New("Not enough seats available")
		}
		if time.Now().After(tour.LastEnrollmentDate) {
			return  errors.New("enrollment deadline has passed")
		}
		//calculate price with discount
		totalPrice := tour.Price

		var customerID uuid.UUID

		if userID != nil {
			customerID, err = uc.bookingRepo.GetOrCreateCustomerByUser(txCtx, *userID)
			if err != nil {
				return  err
			}
		} else if memberID != nil {
			if req.CustomerName == "" || req.CustomerEmail == "" || req.CustomerPhone == 0 {
				return  errors.New("customer details required for guest booking")
			}
			customer := &domain.Customer{
				Name:  req.CustomerName,
				Email: req.CustomerEmail,
				Phone: req.CustomerPhone,
			}
			if err := uc.bookingRepo.CreateCustomer(txCtx, customer); err != nil {
				return  err
			}
			customerID = customer.CustomerID
		} else {
			return errors.New("either user or member must be sppecified")
		}

		booking := &domain.Booking{
			CustomerID:     customerID,
			TourID:         req.TourID,
			NumberOfPeople: req.NumberOfPeople,
			TotalPrice:     totalPrice,
			Status:         "pending",
		}
		if userID != nil {
			booking.UserID = userID
		}
		if memberID != nil {
			booking.MemberID = memberID
		}
	
		if err := uc.bookingRepo.Create(txCtx, booking); err != nil {
			return err
		}
		//update available seats
		newSeats := tour.AvailableSeat - req.NumberOfPeople
		if err := uc.tourRepo.UpdateAvailableSeats(txCtx, tour.TourID, newSeats); err != nil {
			return err
		}

		//create payment
		payment := &domain.Payment{
			BookingID:     booking.BookingID,
			Amount:        totalPrice,
			Method:        req.Method,
			TransactionID: req.TransactionId,
		}
		err = uc.paymentRepo.Create(txCtx, payment)
		if err != nil {
			return err
		}

		response, err = uc.bookingRepo.GetByID(txCtx, booking.BookingID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil

}
