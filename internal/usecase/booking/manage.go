package bookingusecase

import (
	"context"
	"errors"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

const (
	defaultPageLimit = 10
	maxPageLimit     = 100
)

var (
	ErrForbiddenScope      = errors.New("not allowed to access these bookings")
	ErrInvalidStatusFilter = errors.New("status must be one of: pending, confirmed, cancelled, completed")
)

// checkScope makes sure the actor only reaches their own bookings: members
// their agency's, users their own. Super users may use any scope.
func checkScope(actor domain.Actor, scope domain.BookingScope) error {
	switch {
	case actor.IsSuper():
		return nil
	case actor.Role == domain.RoleMember:
		if scope.UserID == nil && scope.AgencyID != nil && actor.AgencyID != nil && *scope.AgencyID == *actor.AgencyID {
			return nil
		}
	case actor.Role == domain.RoleUser:
		if scope.AgencyID == nil && scope.UserID != nil && *scope.UserID == actor.ID {
			return nil
		}
	}
	return ErrForbiddenScope
}

type listBookingsUseCase struct {
	repo port.BookingRepository
}

func NewListBookingsUseCase(repo port.BookingRepository) port.ListBookings {
	return &listBookingsUseCase{repo: repo}
}

func (uc *listBookingsUseCase) Execute(ctx context.Context, actor domain.Actor, f domain.BookingFilter) (*util.PaginationData, error) {
	if err := checkScope(actor, f.BookingScope); err != nil {
		return nil, err
	}
	if f.Status != "" {
		switch f.Status {
		case domain.BookingPending, domain.BookingConfirmed, domain.BookingCancelled, domain.BookingCompleted:
		default:
			return nil, ErrInvalidStatusFilter
		}
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = defaultPageLimit
	}
	if f.Limit > maxPageLimit {
		f.Limit = maxPageLimit
	}

	bookings, total, err := uc.repo.List(ctx, f)
	if err != nil {
		return nil, err
	}
	return &util.PaginationData{
		Data: bookings,
		Meta: util.Meta{
			Page:       f.Page,
			Limit:      f.Limit,
			TotalCount: total,
			TotalPage:  (total + f.Limit - 1) / f.Limit,
		},
	}, nil
}

type getBookingUseCase struct {
	repo port.BookingRepository
}

func NewGetBookingUseCase(repo port.BookingRepository) port.GetBooking {
	return &getBookingUseCase{repo: repo}
}

func (uc *getBookingUseCase) Execute(ctx context.Context, actor domain.Actor, scope domain.BookingScope, bookingID uuid.UUID) (*domain.BookingResponse, error) {
	if err := checkScope(actor, scope); err != nil {
		return nil, err
	}
	return uc.repo.GetByID(ctx, bookingID, scope)
}

// statusChanger applies booking status changes together with their side
// effects on seats and payments, in one transaction.
type statusChanger struct {
	txManager   port.TxManager
	bookingRepo port.BookingRepository
	tourRepo    port.TourRepository
	paymentRepo port.PaymentRepository
	now         func() time.Time
}

// change locks the tour and then the booking, lets check veto the change,
// then updates the status. A cancellation returns the seats to the tour and
// records reason. Every path (booking creation, booking status changes and
// tour cancellation) locks the tour before any booking, so they cannot
// deadlock.
func (s *statusChanger) change(ctx context.Context, bookingID uuid.UUID, scope domain.BookingScope, to, reason string, check func(*domain.LockedBooking) error) (*domain.BookingResponse, error) {
	var response *domain.BookingResponse
	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		current, err := s.bookingRepo.GetByID(txCtx, bookingID, scope)
		if err != nil {
			return err
		}
		tour, err := s.tourRepo.GetByIDForUpdate(txCtx, current.TourID)
		if err != nil {
			return err
		}
		// Re-read under lock: the status may have changed while we waited.
		booking, err := s.bookingRepo.GetForUpdate(txCtx, bookingID, scope)
		if err != nil {
			return err
		}
		if err := domain.CheckBookingTransition(booking.Status, to); err != nil {
			return err
		}
		if check != nil {
			if err := check(booking); err != nil {
				return err
			}
		}

		if to == domain.BookingCancelled {
			if err := s.tourRepo.UpdateAvailableSeats(txCtx, tour.TourID, tour.AvailableSeat+booking.NumberOfPeople); err != nil {
				return err
			}
		} else {
			reason = ""
		}
		if err := s.bookingRepo.UpdateStatus(txCtx, bookingID, to, reason); err != nil {
			return err
		}

		// Confirming a booking means the agency verified the payment.
		// Cancelling voids a payment that was never verified; a verified
		// payment keeps "success" and is refunded outside the system.
		switch to {
		case domain.BookingConfirmed:
			err = s.paymentRepo.SetStatusForBookings(txCtx, []uuid.UUID{bookingID}, "pending", "success")
		case domain.BookingCancelled:
			err = s.paymentRepo.SetStatusForBookings(txCtx, []uuid.UUID{bookingID}, "pending", "failed")
		}
		if err != nil {
			return err
		}

		response, err = s.bookingRepo.GetByID(txCtx, bookingID, domain.BookingScope{})
		return err
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

type updateBookingStatusUseCase struct {
	statusChanger
}

func NewUpdateBookingStatusUseCase(txManager port.TxManager, bookingRepo port.BookingRepository, tourRepo port.TourRepository, paymentRepo port.PaymentRepository) port.UpdateBookingStatus {
	return &updateBookingStatusUseCase{statusChanger{txManager, bookingRepo, tourRepo, paymentRepo, time.Now}}
}

func (uc *updateBookingStatusUseCase) Execute(ctx context.Context, actor domain.Actor, agencyID, bookingID uuid.UUID, status string) (*domain.BookingResponse, error) {
	scope := domain.BookingScope{AgencyID: &agencyID}
	if err := checkScope(actor, scope); err != nil {
		return nil, err
	}
	return uc.change(ctx, bookingID, scope, status, domain.CancelledByAgency, nil)
}

var ErrTourAlreadyStarted = errors.New("bookings can only be cancelled before the tour starts")

type cancelMyBookingUseCase struct {
	statusChanger
}

func NewCancelMyBookingUseCase(txManager port.TxManager, bookingRepo port.BookingRepository, tourRepo port.TourRepository, paymentRepo port.PaymentRepository) port.CancelMyBooking {
	return &cancelMyBookingUseCase{statusChanger{txManager, bookingRepo, tourRepo, paymentRepo, time.Now}}
}

func (uc *cancelMyBookingUseCase) Execute(ctx context.Context, actor domain.Actor, bookingID uuid.UUID) (*domain.BookingResponse, error) {
	scope := domain.BookingScope{UserID: &actor.ID}
	if err := checkScope(actor, scope); err != nil {
		return nil, err
	}
	return uc.change(ctx, bookingID, scope, domain.BookingCancelled, domain.CancelledByCustomer, func(b *domain.LockedBooking) error {
		if !uc.now().Before(b.TourStartDate) {
			return ErrTourAlreadyStarted
		}
		return nil
	})
}
