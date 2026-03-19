package main

import (
	"net/http"

	"github.com/faissalmaulana/go-approve/cmd/handlers"
	"github.com/faissalmaulana/go-approve/cmd/middleware"
	"github.com/faissalmaulana/go-approve/internal/db"
	approvalroomrepository "github.com/faissalmaulana/go-approve/internal/repository/approvalRoom"
	approvalroomapprover "github.com/faissalmaulana/go-approve/internal/repository/approvalRoomApprover"
	"github.com/faissalmaulana/go-approve/internal/repository/blocklisttoken"
	filemetadata "github.com/faissalmaulana/go-approve/internal/repository/fileMetadata"
	requestreview "github.com/faissalmaulana/go-approve/internal/repository/requestReview"
	"github.com/faissalmaulana/go-approve/internal/repository/transactions"
	"github.com/faissalmaulana/go-approve/internal/repository/user"
	approvalroom "github.com/faissalmaulana/go-approve/internal/service/approvalRoom"
	"github.com/faissalmaulana/go-approve/internal/service/auth"
	"github.com/faissalmaulana/go-approve/internal/service/jwtfx"
	requestReviewService "github.com/faissalmaulana/go-approve/internal/service/requestReview"
	userService "github.com/faissalmaulana/go-approve/internal/service/user"
	"github.com/faissalmaulana/go-approve/internal/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		jwtfx.Module,
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			auth.New,
			db.New,
			fx.Annotate(transactions.New, fx.As(new(transactions.DatabaseTransaction))),
			fx.Annotate(user.New, fx.As(new(user.UserStorage))),
			fx.Annotate(approvalroomrepository.New, fx.As(new(approvalroomrepository.ApprovalRoomStorage))),
			fx.Annotate(requestreview.New, fx.As(new(requestreview.RequestReviewStorage))),
			fx.Annotate(filemetadata.New, fx.As(new(filemetadata.FileMetadataStorage))),
			fx.Annotate(approvalroomapprover.New, fx.As(new(approvalroomapprover.ApprovalRoomApproverStorage))),
			blocklisttoken.New,
			userService.New,
			requestReviewService.New,
			approvalroom.New,
			handlers.NewHealthHandler,
			handlers.NewRegisterHandler,
			handlers.NewLoginHandler,
			handlers.NewUserProfileHandler,
			handlers.NewLogoutHandler,
			handlers.NewCreateApprovalRoomHandler,
			handlers.NewGetUsersByUsernameHandler,
			handlers.NewGetApprovalRoomByIdHandler,
			handlers.NewUpdateApprovalStatusHandler,
			handlers.NewConfirmRequestReviewHandler,
			handlers.NewGetRequestReviewHandler,
			middleware.NewAuthMiddleware,
			zap.NewProduction,
			validator.New,
			utils.NewSugaredErrorMessageValidator,
		),
		fx.Invoke(func(*gorm.DB) {}),
		fx.Invoke(func(*http.Server) {})).Run()
}
