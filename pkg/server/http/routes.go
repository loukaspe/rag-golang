package http

import (
	"github.com/loukaspe/rag-golang/internal/core/services"
	http2 "github.com/loukaspe/rag-golang/internal/handlers/http"
	chatSessions2 "github.com/loukaspe/rag-golang/internal/handlers/http/chatSessions"
	"github.com/loukaspe/rag-golang/internal/repositories"

	"github.com/loukaspe/rag-golang/pkg/auth"
	"net/http"
	"os"
)

//	@title			RAG in Golang
//	@version		1.0
//	@description	Experimentation with RAG in Golang using OpenAI, Pinecone, and MCP.

//	@host		localhost:8080
//	@BasePath	/

//	@contact.name	Loukas Peteinaris
//	@contact.url	loukas.peteinaris@gmail.com

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Header value should be in the form of `Bearer <JWT access token>`

// @accept		json
// @produce	json
func (s *Server) initializeRoutes() {
	// health check
	healthCheckHandler := http2.NewHealthCheckHandler(s.DB)
	s.router.HandleFunc("/health-check", healthCheckHandler.HealthCheckController).Methods("GET")

	mcpSSEServer := s.mcpServer.InitialiseSSEServer()

	s.router.HandleFunc("/mcp", mcpSSEServer.ServeHTTP)

	// auth
	jwtMechanism := auth.NewAuthMechanism(
		os.Getenv("JWT_SECRET_KEY"),
		os.Getenv("JWT_SIGNING_METHOD"),
	)
	jwtService := services.NewJwtService(jwtMechanism)
	jwtMiddleware := http2.NewAuthenticationMw(jwtMechanism)
	jwtHandler := http2.NewJwtClaimsHandler(jwtService, s.logger)

	s.router.HandleFunc("/token", jwtHandler.JwtTokenController).Methods(http.MethodPost)

	protected := s.router.PathPrefix("/").Subrouter()
	protected.Use(jwtMiddleware.AuthenticationMW)

	chatSessionRepository := repositories.NewChatSessionRepository(s.DB)
	chatSessionService := services.NewChatSessionService(s.logger, chatSessionRepository)
	messageRepository := repositories.NewMessageRepository(s.DB)
	messageService := services.NewMessageService(s.logger, messageRepository, chatSessionRepository, s.embedder, s.pineconeVectorDB, s.openAIClient)

	createChatSessionHandler := chatSessions2.NewCreateUserChatSessionHandler(chatSessionService, s.logger)
	getChatSessionHandler := chatSessions2.NewGetChatSessionHandler(chatSessionService, s.logger)
	sendMessageHandler := chatSessions2.NewSendMessageHandler(messageService, s.logger)
	submitFeedbackHandler := chatSessions2.NewSubmitFeedbackHandler(messageService, s.logger)

	protected.HandleFunc("/users/{user_id}/chat-sessions", createChatSessionHandler.CreateUserChatSessionController).Methods("POST")
	protected.HandleFunc("/users/{user_id}/chat-sessions", getChatSessionHandler.GetUserChatSessionsController).Methods("GET")
	protected.HandleFunc("/users/{user_id}/chat-sessions/{session_id}/messages", sendMessageHandler.SendMessageController).Methods("POST")
	protected.HandleFunc("/users/{user_id}/chat-sessions/{session_id}/messages/{message_id}/feedback", submitFeedbackHandler.SubmitFeedbackController).Methods("POST")

	protected.HandleFunc("/chat-sessions/{session_id}", getChatSessionHandler.GetChatSessionController).Methods("GET")

}
