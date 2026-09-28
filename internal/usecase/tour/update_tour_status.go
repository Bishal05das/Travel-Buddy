package tourusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type updateTourStatusUseCase struct {
	txManager   port.TxManager
	repo        port.TourRepository
	bookingRepo port.BookingRepository
	paymentRepo port.PaymentRepository
}

func NewUpdateTourStatusUseCase(txManager port.TxManager, repo port.TourRepository, bookingRepo port.BookingRepository, paymentRepo port.PaymentRepository) port.UpdateTourStatus {
	return &updateTourStatusUseCase{
		txManager:   txManager,
		repo:        repo,
		bookingRepo: bookingRepo,
		paymentRepo: paymentRepo,
	}
}

// Execute changes the tour status. Cancelling a tour also cancels its
// pending and confirmed bookings (reason "tour_cancelled"), returns their
// seats and fails their unverified payments, all in one transaction.
// Verified payments stay "success" and are refunded outside the system.
func (uc *updateTourStatusUseCase) Execute(ctx context.Context, actor domain.Actor, tourID uuid.UUID, status string) (*domain.TourStatusChange, error) {
	result := &domain.TourStatusChange{TourID: tourID, Status: status}

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Locking the tour first serialises this with booking creation and
		// booking status changes, which lock the tour before any booking.
		tour, err := uc.repo.GetByIDForUpdate(txCtx, tourID)
		if err != nil {
			return err
		}
		scope := actor.AgencyScope()
		if scope != nil && tour.AgencyID != *scope {
			return errors.New("tour not found")
		}
		if err := domain.CheckTourStatus(tour.Status, status); err != nil {
			return err
		}

		if status == domain.TourCancelled {
			ids, seats, err := uc.bookingRepo.CancelActiveForTour(txCtx, tourID, domain.CancelledWithTour)
			if err != nil {
				return err
			}
			if err := uc.paymentRepo.SetStatusForBookings(txCtx, ids, "pending", "failed"); err != nil {
				return err
			}
			if seats > 0 {
				if err := uc.repo.UpdateAvailableSeats(txCtx, tourID, tour.AvailableSeat+seats); err != nil {
					return err
				}
			}
			result.CancelledBookings = len(ids)
		}

		return uc.repo.UpdateTourStatus(txCtx, tourID, status, scope)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
