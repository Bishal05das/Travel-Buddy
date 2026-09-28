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

func (uc *createbookingusecase) Execute(ctx context.Context, req *domain.BookingCommand) (*domain.BookingResponse, error) {

	var response *domain.BookingResponse

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		//starts by row level locking
		tour, err := uc.tourRepo.GetByIDForUpdate(txCtx, req.TourID)
		if err != nil {
			return err
		}
		if tour.Status != "open" {
			return errors.New("tour is not open for booking")
		}
		if req.NumberOfPeople < 1 {
			return errors.New("number of people must be at least 1")
		}
		if tour.AvailableSeat < req.NumberOfPeople {
			return errors.New("not enough seats available")
		}
		if time.Now().After(tour.LastEnrollmentDate) {
			return errors.New("enrollment deadline has passed")
		}
		//calculate price with discount for the whole group
		receivedPrice := req.TotalPrice
		calculatedPrice := tour.UnitPrice() * req.NumberOfPeople
		if receivedPrice != calculatedPrice {
			return errors.New("price mismatch, please check the price and try again")
		}

		var customerID uuid.UUID

		if req.UserID != nil {
			customerID, err = uc.bookingRepo.GetOrCreateCustomerByUser(txCtx, *req.UserID)
			if err != nil {
				return err
			}
		} else if req.MemberID != nil {
			if req.GuestInfo == nil || req.GuestInfo.Name == "" || req.GuestInfo.Email == "" || req.GuestInfo.Phone == "" {
				return errors.New("customer details required for guest booking")
			}
			customer := &domain.Customer{
				Name:  req.GuestInfo.Name,
				Email: req.GuestInfo.Email,
				Phone: req.GuestInfo.Phone,
			}
			if err := uc.bookingRepo.CreateCustomer(txCtx, customer); err != nil {
				return err
			}
			customerID = customer.CustomerID
		} else {
			return errors.New("either user or member must be specified")
		}

		booking := &domain.Booking{
			CustomerID:     customerID,
			TourID:         req.TourID,
			NumberOfPeople: req.NumberOfPeople,
			TotalPrice:     receivedPrice,
			Status:         "pending",
		}
		if req.UserID != nil {
			booking.UserID = req.UserID
		}
		if req.MemberID != nil {
			booking.MemberID = req.MemberID
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
			Amount:        receivedPrice,
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
