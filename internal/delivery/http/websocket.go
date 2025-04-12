package http

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		origin := r.Header.Get("Origin")
// 		return origin == config.Cors.AllowedOrigin
// 	},
// }

// type websocketController struct {
// 	websocketUsecase usecase.IWebsocketUsecase
// 	sessionUsecase   usecase.ISessionUsecase
// }

// func NewWebsocketController(r *mux.Router, websocketUsecase usecase.IWebsocketUsecase, sessionUsecase usecase.ISessionUsecase) {
// 	controller := &websocketController{
// 		websocketUsecase: websocketUsecase,
// 		sessionUsecase:   sessionUsecase,
// 	}
// 	r.Handle("/ws", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.WebsocketConnection))).Methods(http.MethodGet)
// }

// func (c *websocketController) WebsocketConnection(w http.ResponseWriter, r *http.Request) {

// }
