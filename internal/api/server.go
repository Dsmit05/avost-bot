package api

import (
	"github.com/Dsmit05/avost-bot/docs"
	"github.com/Dsmit05/avost-bot/internal/api/controllers"
	"github.com/Dsmit05/avost-bot/internal/api/middlewares"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/pkg/stat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	r                    *gin.Engine
	addr                 string
	controllerBot        *controllers.Bot
	controllerRepository *controllers.Repository
	protected            *middlewares.ProtectedMiddleware
}

func NewServer(cache RepositoryI, bot BotI, crypto cryptographyI, stats *stat.Statistic, addr string) *Server {
	gin.SetMode(gin.ReleaseMode)

	return &Server{
		r:                    gin.New(),
		addr:                 addr,
		controllerBot:        controllers.NewBot(bot, stats),
		controllerRepository: controllers.NewRepository(cache),
		protected:            middlewares.NewProtectedMiddleware(crypto),
	}
}

func (s *Server) AddV1(basePath string) *Server {
	logger.Info("Server.AddV1()", "add api v1")

	v1 := s.r.Group(basePath)
	v1.Use(cors.Default())

	controlBot := v1.Group("/bot")
	controlBot.Use(s.protected.AuthMiddleware)

	{
		controlBot.POST("/send", s.controllerBot.SendMsgUser)
		controlBot.POST("/sends", s.controllerBot.SendAllMsgUser)
		controlBot.POST("/replay", s.controllerBot.ReplayMsgUser)
		controlBot.GET("/user/:id", s.controllerBot.GetUserInfo)
		controlBot.GET("/stat", s.controllerBot.GetStat)
	}

	controlRepository := v1.Group("/repository")
	controlRepository.Use(s.protected.AuthMiddleware)

	{
		controlRepository.GET("/users", s.controllerRepository.GetUsers)
		controlRepository.GET("/admins", s.controllerRepository.GetUsersAdmins)
		controlRepository.PATCH("/role", s.controllerRepository.UpdateUserRole)
		controlRepository.PATCH("/sub", s.controllerRepository.UpdateUserSub)
		controlRepository.GET("/users-file", s.controllerRepository.GetUsersFile)
	}

	return s
}

func (s *Server) AddSwagger(host string) *Server {
	docs.SwaggerInfo.Host = host

	s.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return s
}

func (s *Server) Start() {
	if err := s.r.Run(s.addr); err != nil {
		logger.Error("server start", err)
	}
}
