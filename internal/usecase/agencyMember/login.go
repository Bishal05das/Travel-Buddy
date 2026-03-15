package memberusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type memberLoginUseCase struct {
	memberRepo port.AgencyMemberRepository
	cfg        *config.Config
}

func NewMemberLoginUseCase(memberRepo port.AgencyMemberRepository, cfg *config.Config) port.LoginUser {
	return &memberLoginUseCase{
		memberRepo: memberRepo,
		cfg:        cfg,
	}
}

func (uc *memberLoginUseCase) Execute(ctx context.Context, member *domain.ReqLogin) (*string, error) {
	mem, err := uc.memberRepo.FindMember(ctx, member.Email, member.Password)
	if mem == nil {
		// http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return nil, errors.New("Invalid Credentials")
	}
	accessToken, err := util.CreateJWT(uc.cfg.JWTSecretkey, util.Payload{
		UserID: mem.MemberID,
		Role:   "member",
		RoleID: &mem.RoleID,
	})
	if err != nil {
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	return &accessToken, nil
}
