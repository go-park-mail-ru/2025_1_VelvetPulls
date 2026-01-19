// usecase/guest.go
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// IGuestModule - интерфейс usecase для работы с гостями (GuestModule)
type IGuestModule interface {
	CreateGuest(ctx context.Context, req *model.CreateGuestRequest) (*model.Guest, error)
	GetGuest(ctx context.Context, guestID uuid.UUID) (*model.Guest, error)
	GetGuestBySession(ctx context.Context, sessionID string) (*model.Guest, error)
	UpdateGuest(ctx context.Context, guestID uuid.UUID, req *model.UpdateGuestRequest) (*model.Guest, error)
	DeleteGuest(ctx context.Context, guestID uuid.UUID) error
	ConnectToChat(ctx context.Context, guestID uuid.UUID, chatID uuid.UUID) error
	SetGuestName(ctx context.Context, guestID uuid.UUID, name string) error
	UpdateGuestActivity(ctx context.Context, guestID uuid.UUID) error
	CleanupExpiredGuests(ctx context.Context, expirationHours int) error
}

type GuestModule struct {
	guestMapper repository.IGuestMapper
	nc        *nats.Conn
}

func NewGuestModule(guestMapper repository.IGuestMapper, nc *nats.Conn) IGuestModule {
	return &GuestModule{
		guestMapper: guestMapper,
		nc:        nc,
	}
}

func (uc *GuestModule) CreateGuest(ctx context.Context, req *model.CreateGuestRequest) (*model.Guest, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("CreateGuest start", zap.String("temporaryID", req.TemporaryID))

	// Валидация
	if err := req.Validate(); err != nil {
		logger.Error("CreateGuest: validation failed", zap.Error(err))
		return nil, err
	}

	existingGuest, err := uc.guestMapper.GetGuestByTemporaryID(ctx, req.TemporaryID)
	if err == nil && existingGuest != nil {
		logger.Info("CreateGuest: guest already exists, returning existing", 
			zap.String("guestID", existingGuest.ID.String()))
		return existingGuest, nil
	}

	guest := &model.Guest{
		ID:          uuid.New(),
		TemporaryID: req.TemporaryID,
		Status:      model.GuestStatusPending,
		CreatedAt:   time.Now(),
	}

	createdID, err := uc.guestMapper.CreateGuest(ctx, guest)
	if err != nil {
		logger.Error("CreateGuest: failed to create guest in repository", zap.Error(err))
		return nil, err
	}

	createdGuest, err := uc.guestMapper.GetGuestByID(ctx, createdID)
	if err != nil {
		logger.Error("CreateGuest: failed to get created guest", zap.Error(err))
		return nil, err
	}

	uc.publishGuestEvent(ctx, "guest.created", createdGuest)

	metrics.IncBusinessOp("create_guest")
	logger.Info("CreateGuest: success", zap.String("guestID", createdGuest.ID.String()))
	return createdGuest, nil
}

func (uc *GuestModule) GetGuest(ctx context.Context, guestID uuid.UUID) (*model.Guest, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetGuest", zap.String("guestID", guestID.String()))

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		logger.Error("GetGuest: failed to get guest", zap.Error(err))
		return nil, err
	}

	go uc.updateGuestActivity(context.Background(), guestID)

	metrics.IncBusinessOp("get_guest")
	return guest, nil
}

func (uc *GuestModule) GetGuestBySession(ctx context.Context, sessionID string) (*model.Guest, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetGuestBySession", zap.String("sessionID", sessionID))

	guest, err := uc.guestMapper.GetGuestBySession(ctx, sessionID)
	if err != nil {
		logger.Error("GetGuestBySession: failed to get guest", zap.Error(err))
		return nil, err
	}

	go uc.updateGuestActivity(context.Background(), guest.ID)

	metrics.IncBusinessOp("get_guest_by_session")
	return guest, nil
}

func (uc *GuestModule) UpdateGuest(ctx context.Context, guestID uuid.UUID, req *model.UpdateGuestRequest) (*model.Guest, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("UpdateGuest", zap.String("guestID", guestID.String()))

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		logger.Error("UpdateGuest: guest not found", zap.Error(err))
		return nil, err
	}

	if req.Status != "" {
		guest.Status = req.Status
	}
	if req.SessionID != "" {
		guest.SessionID = req.SessionID
	}

	if err := uc.guestMapper.UpdateGuest(ctx, guest); err != nil {
		logger.Error("UpdateGuest: failed to update guest", zap.Error(err))
		return nil, err
	}

	uc.publishGuestEvent(ctx, "guest.updated", guest)

	metrics.IncBusinessOp("update_guest")
	return guest, nil
}

func (uc *GuestModule) DeleteGuest(ctx context.Context, guestID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("DeleteGuest", zap.String("guestID", guestID.String()))

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		logger.Error("DeleteGuest: guest not found", zap.Error(err))
		return err
	}

	if err := uc.guestMapper.DeleteGuest(ctx, guestID); err != nil {
		logger.Error("DeleteGuest: failed to delete guest", zap.Error(err))
		return err
	}

	uc.publishGuestEvent(ctx, "guest.deleted", guest)

	metrics.IncBusinessOp("delete_guest")
	return nil
}

func (uc *GuestModule) ConnectToChat(ctx context.Context, guestID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("ConnectToChat", 
		zap.String("guestID", guestID.String()),
		zap.String("chatID", chatID.String()))

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		logger.Error("ConnectToChat: guest not found", zap.Error(err))
		return err
	}

	guest.Status = model.GuestStatusConnected
	if err := uc.guestMapper.UpdateGuest(ctx, guest); err != nil {
		logger.Error("ConnectToChat: failed to update guest status", zap.Error(err))
		return err
	}

	event := model.GuestEvent{
		Action: "guest.connected_to_chat",
		Guest:  *guest,
	}
	data, _ := json.Marshal(event)
	
	subject := fmt.Sprintf("guest.%s.events", guestID.String())
	if err := uc.nc.Publish(subject, data); err != nil {
		logger.Error("ConnectToChat: failed to publish event", zap.Error(err))
	}

	logger.Info("ConnectToChat: guest connected to chat successfully")
	metrics.IncBusinessOp("connect_guest_to_chat")
	return nil
}

func (uc *GuestModule) SetGuestName(ctx context.Context, guestID uuid.UUID, name string) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("SetGuestName", 
		zap.String("guestID", guestID.String()),
		zap.String("name", name))

	if len(name) < 2 || len(name) > 50 {
		return ErrInvalidGuestName
	}

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		logger.Error("SetGuestName: guest not found", zap.Error(err))
		return err
	}

	guest.TemporaryID = name
	if err := uc.guestMapper.UpdateGuest(ctx, guest); err != nil {
		logger.Error("SetGuestName: failed to update guest name", zap.Error(err))
		return err
	}

	metrics.IncBusinessOp("set_guest_name")
	return nil
}

func (uc *GuestModule) UpdateGuestActivity(ctx context.Context, guestID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Debug("UpdateGuestActivity", zap.String("guestID", guestID.String()))

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		return err
	}

	if err := uc.guestMapper.UpdateGuest(ctx, guest); err != nil {
		logger.Error("UpdateGuestActivity: failed to update guest activity", zap.Error(err))
		return err
	}

	return nil
}

func (uc *GuestModule) CleanupExpiredGuests(ctx context.Context, expirationHours int) error {
	logger := zap.L()
	logger.Info("CleanupExpiredGuests started")

	// В реальном проекте здесь был бы запрос на удаление старых гостей
	// query := `DELETE FROM guests WHERE last_activity_at < NOW() - INTERVAL '$1 hours'`
	// Но в нашем интерфейсе репозитория нет такого метода, поэтому оставляем заглушку

	logger.Info("CleanupExpiredGuests completed")
	return nil
}

func (uc *GuestModule) publishGuestEvent(ctx context.Context, action string, guest *model.Guest) {
	logger := utils.GetLoggerFromCtx(ctx)
	
	event := model.GuestEvent{
		Action: action,
		Guest:  *guest,
	}
	
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error("publishGuestEvent: failed to marshal event", zap.Error(err))
		return
	}
	
	subject := fmt.Sprintf("guest.%s.events", guest.ID.String())
	if err := uc.nc.Publish(subject, data); err != nil {
		logger.Error("publishGuestEvent: failed to publish event", 
			zap.String("subject", subject), 
			zap.Error(err))
	}
}

func (uc *GuestModule) updateGuestActivity(ctx context.Context, guestID uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	guest, err := uc.guestMapper.GetGuestByID(ctx, guestID)
	if err != nil {
		zap.L().Warn("updateGuestActivity: failed to get guest", zap.Error(err))
		return
	}

	if err := uc.guestMapper.UpdateGuest(ctx, guest); err != nil {
		zap.L().Warn("updateGuestActivity: failed to update guest", zap.Error(err))
	}
}

var (
	ErrGuestNotFound        = fmt.Errorf("guest not found")
	ErrInvalidTemporaryID   = fmt.Errorf("invalid temporary ID")
	ErrInvalidGuestName     = fmt.Errorf("invalid guest name (must be 2-50 characters)")
	ErrGuestAlreadyExists   = fmt.Errorf("guest already exists")
	ErrGuestExpired         = fmt.Errorf("guest session expired")
)