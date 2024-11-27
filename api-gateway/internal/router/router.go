package router

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/config"
	commentclients "api-gateway/internal/clients/comment_clients"
	directclients "api-gateway/internal/clients/direct_clients"
	likeclients "api-gateway/internal/clients/like_clients"
	tweetclients "api-gateway/internal/clients/tweet_clients"
	userclients "api-gateway/internal/clients/user_clients"
	commenthandler "api-gateway/internal/handlers/comment-handlers"
	directhandler "api-gateway/internal/handlers/direct-handlers"
	likehandler "api-gateway/internal/handlers/like-handlers"
	tweethandler "api-gateway/internal/handlers/tweet-handlers"
	userhandler "api-gateway/internal/handlers/user-handlers"

	"api-gateway/internal/jwt"

	middleware "api-gateway/internal/rate-limiting"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *http.Server {
	userHandlers := &userhandler.UserHandler{UserService: userclients.DialUserGrpc()}
	tweetClient := tweetclients.DialTweetGrpc()
	tweetHandlers := &tweethandler.TweetHandler{TweetService: tweetClient}

	likeClient := likeclients.DialLikeGrpc()
	likeHandler := &likehandler.LikeHandler{LikeService: likeClient}

	commentClient := commentclients.DialCommentGrpc()
	commentHandlers := &commenthandler.CommentHandler{CommentService: commentClient}

	directClient := directclients.DialDirectGrpc()
	directHandlers := &directhandler.DirectHandler{DirectService: directClient}

	// Rate limiter o‘rnatish
	userRateLimiter := middleware.NewRateLimiter(2, time.Minute)
	tweetRateLimiter := middleware.NewRateLimiter(2, time.Minute)
	likeRateLimiter := middleware.NewRateLimiter(2, time.Minute)
	commentRateLimiter := middleware.NewRateLimiter(2, time.Minute)
	directRateLimiter := middleware.NewRateLimiter(2, time.Minute)

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("Registered route: GET /swagger/*any")

	// userRoutes
	userRoutes := router.Group("/users")
	userRoutes.Use(userRateLimiter.Limit())
	{
		userRoutes.POST("/register", userHandlers.Register)
		log.Println("Registered route: POST /users/register")
		userRoutes.POST("/login", userHandlers.Login)
		log.Println("Registered route: POST /users/login")
		userRoutes.GET("/:id", jwt.Protected(), userHandlers.GetUser)
		log.Println("Registered route: GET /users/:id")
		userRoutes.PUT("/:id", jwt.Protected(), userHandlers.UpdateUser)
		log.Println("Registered route: PUT /users/:id")
		userRoutes.DELETE("/:id", jwt.Protected(), userHandlers.DeleteUser)
		log.Println("Registered route: DELETE /users/:id")
		userRoutes.PUT("/:id/password", jwt.Protected(), userHandlers.ChangePassword)
		log.Println("Registered route: PUT /users/:id/password")
		userRoutes.PUT("/:id/avatar", jwt.Protected(), userHandlers.UpdateAvatar)
		log.Println("Registered route: PUT /users/:id/avatar")
		userRoutes.POST("/:id/avatar", jwt.Protected(), userHandlers.AddAvatar)
		log.Println("Registered route: POST /users/:id/avatar")
		userRoutes.DELETE("/:id/avatar", jwt.Protected(), userHandlers.DeleteAvatar)
		log.Println("Registered route: DELETE /users/:id/avatar")
		userRoutes.PUT("/:id/follow", jwt.Protected(), userHandlers.AddFollowing)
		log.Println("Registered route: PUT /users/:id/follow")
		userRoutes.PUT("/:id/unfollow", jwt.Protected(), userHandlers.RemoveFollowing)
		log.Println("Registered route: PUT /users/:id/unfollow")
		userRoutes.PUT("/:id/follower", jwt.Protected(), userHandlers.AddFollower)
		log.Println("Registered route: PUT /users/:id/follower")
		userRoutes.PUT("/:id/unfollower", jwt.Protected(), userHandlers.RemoveFollower)
		log.Println("Registered route: PUT /users/:id/unfollower")
	}

	// tweetRoutes
	tweetRoutes := router.Group("/tweets")
	tweetRoutes.Use(tweetRateLimiter.Limit())
	{
		tweetRoutes.POST("", jwt.Protected(), tweetHandlers.CreateTweet)
		log.Println("Registered route: POST /tweets")
		tweetRoutes.GET("/user/:user_id", jwt.Protected(), tweetHandlers.GetTweetsByUser)
		log.Println("Registered route: GET /tweets/user/:user_id")
		tweetRoutes.GET("/:id", jwt.Protected(), tweetHandlers.GetTweetByID)
		log.Println("Registered route: GET /tweets/:id")
		tweetRoutes.PUT("/:id", jwt.Protected(), tweetHandlers.UpdateTweet)
		log.Println("Registered route: PUT /tweets/:id")
		tweetRoutes.DELETE("/:id", jwt.Protected(), tweetHandlers.DeleteTweet)
		log.Println("Registered route: DELETE /tweets/:id")
		tweetRoutes.PUT("/:id/like", jwt.Protected(), tweetHandlers.AddLike)
		log.Println("Registered route: PUT /tweets/:id/like")
		tweetRoutes.DELETE("/:id/like", jwt.Protected(), tweetHandlers.RemoveLike)
		log.Println("Registered route: DELETE /tweets/:id/like")
		tweetRoutes.POST("/:id/comment", jwt.Protected(), tweetHandlers.AddComment)
		log.Println("Registered route: POST /tweets/:id/comment")
		tweetRoutes.DELETE("/:id/comment/:comment_id", jwt.Protected(), tweetHandlers.RemoveComment)
		log.Println("Registered route: DELETE /tweets/:id/comment/:comment_id")
		tweetRoutes.PUT("/:id/share", jwt.Protected(), tweetHandlers.AddShare)
		log.Println("Registered route: PUT /tweets/:id/share")
		tweetRoutes.PUT("/:id/save", jwt.Protected(), tweetHandlers.SaveTweet)
		log.Println("Registered route: PUT /tweets/:id/save")
		tweetRoutes.DELETE("/:id/save", jwt.Protected(), tweetHandlers.RemoveSave)
		log.Println("Registered route: DELETE /tweets/:id/save")
	}

	// likeRoutes
	likeRoutes := router.Group("/likes")
	likeRoutes.Use(likeRateLimiter.Limit())
	{
		likeRoutes.POST("/tweet", jwt.Protected(), likeHandler.CreateLikeTweet)
		log.Println("Registered route: POST /likes/tweet")
		likeRoutes.DELETE("/tweet", jwt.Protected(), likeHandler.DeleteLikeTweet)
		log.Println("Registered route: DELETE /likes/tweet")
		likeRoutes.GET("/tweet/:tweet_id", jwt.Protected(), likeHandler.GetLikesTweet)
		log.Println("Registered route: GET /likes/tweet/:tweet_id")
		likeRoutes.GET("/tweet/user/:user_id", jwt.Protected(), likeHandler.GetLikeTweetByUser)
		log.Println("Registered route: GET /likes/tweet/user/:user_id")
		likeRoutes.POST("/comment", jwt.Protected(), likeHandler.CreateLikeComment)
		log.Println("Registered route: POST /likes/comment")
		likeRoutes.DELETE("/comment", jwt.Protected(), likeHandler.DeleteLikeComment)
		log.Println("Registered route: DELETE /likes/comment")
		likeRoutes.GET("/comment/:comment_id", jwt.Protected(), likeHandler.GetLikesComment)
		log.Println("Registered route: GET /likes/comment/:comment_id")
		likeRoutes.GET("/comment/user/:user_id", jwt.Protected(), likeHandler.GetLikeCommentByUser)
		log.Println("Registered route: GET /likes/comment/user/:user_id")
	}

	// commentRoutes
	commentRoutes := router.Group("/comments")
	commentRoutes.Use(commentRateLimiter.Limit())
	{
		commentRoutes.POST("", jwt.Protected(), commentHandlers.CreateComment)
		log.Println("Registered route: POST /comments")
		commentRoutes.GET("/tweet/:tweet_id", jwt.Protected(), commentHandlers.GetCommentsByTweetID)
		log.Println("Registered route: GET /comments/tweet/:tweet_id")
		commentRoutes.GET("/:id", jwt.Protected(), commentHandlers.GetComment)
		log.Println("Registered route: GET /comments/:id")
		commentRoutes.DELETE("/:id", jwt.Protected(), commentHandlers.DeleteComment)
		log.Println("Registered route: DELETE /comments/:id")
	}

	// directRoutes
	directRoutes := router.Group("/directs")
	directRoutes.Use(directRateLimiter.Limit())
	{
		directRoutes.POST("", jwt.Protected(), directHandlers.CreateDirectMessage)
		log.Println("Registered route: POST /directs")
		directRoutes.GET("", jwt.Protected(), directHandlers.GetDirectMessages)
		log.Println("Registered route: GET /directs")
		directRoutes.GET("/:id", jwt.Protected(), directHandlers.GetDirectMessageByID)
		log.Println("Registered route: GET /directs/:id")
		directRoutes.DELETE("/:id", jwt.Protected(), directHandlers.DeleteDirectMessage)
		log.Println("Registered route: DELETE /directs/:id")
	}

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router = gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ... (existing route setup code)

	server := &http.Server{
		Addr:    "localhost:" + conf.SERVER_PORT, // Remove "localhost" to bind to all interfaces
		Handler: router,
	}

	// Start the HTTPS server in a goroutine
	go func() {
		log.Printf("Starting HTTPS server on port %s", conf.SERVER_PORT)
		if err := server.ListenAndServeTLS("./tls/items.pem", "./tls/items-key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run HTTPS server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return server
}

func GracefulShutdown(srv *http.Server, logger *log.Logger) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	<-shutdownCh
	logger.Println("Shutdown signal received, initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Server shutdown encountered an error: %v", err)
	} else {
		logger.Println("Server gracefully stopped")
	}

	select {
	case <-shutdownCtx.Done():
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			logger.Println("Shutdown deadline exceeded, forcing server to stop")
		}
	default:
		logger.Println("Shutdown completed within the timeout period")
	}
}
