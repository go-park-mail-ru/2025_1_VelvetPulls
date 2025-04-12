package usecase

type IWebsocketUsecase interface {
}

type WebsocketUsecase struct {
}

func NewWebsocketUsecase() IWebsocketUsecase {
	return &MessageUsecase{}
}
