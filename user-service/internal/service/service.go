package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"user-service/internal/kafka"
	"user-service/internal/methods"
	"user-service/internal/models"
	"user-service/pkg/proto"
	"user-service/utils"

	"github.com/go-redis/redis/v8"
)

type UserService struct {
	db          *sql.DB
	RedisClient *redis.Client
}

func NewUserService(db *sql.DB, redisClient *redis.Client) *UserService {
	return &UserService{db: db, RedisClient: redisClient}
}

func (s *UserService) RegisterUser(ctx context.Context, username, email, password, name, phone, code string) (models.User, error) {
	newUser := models.User{
		Username: username,
		Email:    email,
		Password: password,
		Name:     name,
		Phone:    phone,
	}

	var loginIdentifier string
	if email != "" {
		loginIdentifier = email
	} else if phone != "" {
		loginIdentifier = phone
	} else {
		return models.User{}, errors.New("login identifier is required")
	}

	queryCode := "INSERT INTO codes (login_identifier, code) VALUES ($1, $2) RETURNING *"
	var message models.Code

	err := s.db.QueryRowContext(ctx, queryCode, loginIdentifier, code).Scan(&message.ID, &message.LoginIdentifier, &message.Code, &message.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("error inserting code: %v", err)
	}

	queryCheck := `SELECT * FROM users WHERE username = $1`
	var existingUser models.User
	err = s.db.QueryRowContext(ctx, queryCheck, newUser.Username).Scan(&existingUser)
	if err != nil && err != sql.ErrNoRows {
		return models.User{}, fmt.Errorf("error checking user existence: %v", err)
	} else if err == nil {
		return models.User{}, errors.New("user already exists")
	}

	queryInsert := `INSERT INTO users (username, email, password, name, phone) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	err = s.db.QueryRowContext(ctx, queryInsert, newUser.Username, newUser.Email, newUser.Password, newUser.Name, newUser.Phone).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.Name,
		&newUser.Phone,
		&newUser.Password,
		&newUser.CreatedAt,
		&newUser.Bio,
		&newUser.IsPrivate,
		&newUser.Followers,
		&newUser.Following,
		&newUser.Tweets,
		&newUser.BlockedUsers,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("error inserting user: %v", err)
	}

	return newUser, nil
}

func (s *UserService) LoginUser(ctx context.Context, loginIdentifier string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT *
		FROM users 
		WHERE username = $1 OR email = $1 OR phone = $1`

	err := s.db.QueryRowContext(ctx, query, loginIdentifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("error querying user: %v", err)
	}
	user.FollowerCount = int32(len(user.Followers))
	user.FollowingCount = int32(len(user.Following))
	user.TweetCount = int32(len(user.Tweets))
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64, userIDToGet int64) (models.User, error) {
	user := models.User{}
	query := `
		SELECT * FROM users 
		WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Tweets,
		&user.Followers,
		&user.Following,
		&user.BlockedUsers,
	)
	if userIDToGet != 0 {
		if !utils.InSlice(user.BlockedUsers, int32(userIDToGet)) {
			return models.User{}, errors.New("user is not blocked")
		}
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("error querying user: %v", err)
	}
	user.FollowerCount = int32(len(user.Followers))
	user.FollowingCount = int32(len(user.Following))
	user.TweetCount = int32(len(user.Tweets))
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT * FROM users 
		WHERE email = $1`
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Tweets,
		&user.Followers,
		&user.Following,
		&user.BlockedUsers,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("error querying user: %v", err)
	}
	user.FollowerCount = int32(len(user.Followers))
	user.FollowingCount = int32(len(user.Following))
	user.TweetCount = int32(len(user.Tweets))
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT * FROM users 
		WHERE username = $1`
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Tweets,
		&user.Followers,
		&user.Following,
		&user.BlockedUsers,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("error querying user: %v", err)
	}
	user.FollowerCount = int32(len(user.Followers))
	user.FollowingCount = int32(len(user.Following))
	user.TweetCount = int32(len(user.Tweets))
	return user, nil
}

func (s *UserService) GetUserByPhone(ctx context.Context, phone string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT * FROM users 
		WHERE phone = $1`
	err := s.db.QueryRowContext(ctx, query, phone).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Tweets,
		&user.Followers,
		&user.Following,
		&user.BlockedUsers,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("error querying user: %v", err)
	}
	user.FollowerCount = int32(len(user.Followers))
	user.FollowingCount = int32(len(user.Following))
	user.TweetCount = int32(len(user.Tweets))
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	if user.Email == "" && user.Phone == "" {
		return models.User{}, errors.New("either email or phone must be provided")
	}

	var existingUsername string
	err := s.db.QueryRowContext(ctx, "SELECT username FROM users WHERE username = $1 AND id != $2", user.Username, user.ID).Scan(&existingUsername)
	if err == nil {
		return models.User{}, errors.New("username already exists")
	}

	query := "UPDATE users SET "
	params := []interface{}{}
	setClauses := []string{}

	if user.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", len(params)+1))
		params = append(params, user.Username)
	}
	if user.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", len(params)+1))
		params = append(params, user.Email)
	}
	if user.Phone != "" {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", len(params)+1))
		params = append(params, user.Phone)
	}
	if user.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(params)+1))
		params = append(params, user.Name)
	}
	setClauses = append(setClauses, fmt.Sprintf("bio = $%d", len(params)+1))
	params = append(params, user.Bio)
	setClauses = append(setClauses, fmt.Sprintf("is_private = $%d", len(params)+1))
	params = append(params, user.IsPrivate)

	setClauses = append(setClauses, fmt.Sprintf("WHERE id = $%d", len(params)+1))
	params = append(params, user.ID)

	query += strings.Join(setClauses, ", ") + " RETURNING *"

	err = s.db.QueryRowContext(ctx, query, params...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("error updating user: %v", err)
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	query := "UPDATE users SET password = $1 WHERE id = $2 AND password = $3"
	hashedPassword, err := utils.HashPassword(oldPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	_, err = s.db.ExecContext(ctx, query, newPassword, id, hashedPassword)
	if err != nil {
		return fmt.Errorf("error changing password: %v", err)
	}
	return nil
}

func (s *UserService) ForgotPassword(ctx context.Context, newPassword, loginIdentifier string) error {
	query := "UPDATE users SET password = $1 WHERE email = $2 OR phone = $2"
	_, err := s.db.ExecContext(ctx, query, newPassword, loginIdentifier)
	if err != nil {
		return fmt.Errorf("error changing password: %v", err)
	}
	return nil
}

func (s *UserService) AddAvatar(ctx context.Context, id int64, avatarURL string) error {
	query := "INSERT INTO avatars (user_id, url) VALUES ($1, $2)"
	_, err := s.db.ExecContext(ctx, query, id, avatarURL)
	if err != nil {
		return fmt.Errorf("error adding avatar: %v", err)
	}
	return nil
}

func (s *UserService) GetAvatar(ctx context.Context, id int64) (string, error) {
	query := "SELECT url FROM avatars WHERE user_id = $1"
	var avatarURL string
	err := s.db.QueryRowContext(ctx, query, id).Scan(&avatarURL)
	if err != nil {
		return "", fmt.Errorf("error getting avatar: %v", err)
	}
	return avatarURL, nil
}

func (s *UserService) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	query := "UPDATE avatars SET url = $1 WHERE user_id = $2"
	_, err := s.db.ExecContext(ctx, query, avatarURL, id)
	if err != nil {
		return fmt.Errorf("error updating avatar: %v", err)
	}
	return nil
}

func (s *UserService) DeleteAvatar(ctx context.Context, id int64) error {
	query := "DELETE FROM avatars WHERE user_id = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting avatar: %v", err)
	}
	return nil
}

func (s *UserService) AddTweet(ctx context.Context, id int64, tweetID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error adding tweet: %v", err)
	}

	user.Tweets = append(user.Tweets, int32(tweetID))
	user.TweetCount += 1

	query = "UPDATE users SET tweets = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.Tweets, id)
	if err != nil {
		return fmt.Errorf("error updating user tweets: %v", err)
	}
	return nil
}

func (s *UserService) RemoveTweet(ctx context.Context, id int64, tweetID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error removing tweet: %v", err)
	}

	user.Tweets = utils.RemoveElement(user.Tweets, int32(tweetID))
	user.TweetCount -= 1

	query = "UPDATE users SET tweets = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.Tweets, id)
	if err != nil {
		return fmt.Errorf("error updating user tweets: %v", err)
	}
	return nil
}

func (s *UserService) AddFollower(ctx context.Context, id int64, followerID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error adding follower: %v", err)
	}

	user.Followers = append(user.Followers, int32(followerID))
	user.FollowerCount += 1

	query = "UPDATE users SET followers = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.Followers, id)
	if err != nil {
		return fmt.Errorf("error updating user followers: %v", err)
	}
	return nil
}

func (s *UserService) RemoveFollower(ctx context.Context, id int64, followerID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error removing follower: %v", err)
	}
	user.Followers = utils.RemoveElement(user.Followers, int32(followerID))
	user.FollowerCount -= 1

	query = "UPDATE users SET followers = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.Followers, id)
	if err != nil {
		return fmt.Errorf("error updating user followers: %v", err)
	}
	return nil
}

func (s *UserService) AddFollowing(ctx context.Context, id int64, followingID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error adding following: %v", err)
	}
	user.Following = append(user.Following, int32(followingID))
	user.FollowingCount += 1
	err = kafka.PublishNotification(ctx, "follow_request_accepted", kafka.NotificationMessage{
		UserID:     int32(id),
		FollowerID: int32(followingID),
		Action:     "add",
		Message:    "Following added",
	})
	if err != nil {
		return fmt.Errorf("error accepting follow request: %v", err)
	}
	query = "UPDATE users SET following = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.Following, id)
	if err != nil {
		return fmt.Errorf("error updating user following: %v", err)
	}
	return nil
}

func (s *UserService) RemoveFollowing(ctx context.Context, id int64, followingID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error removing following: %v", err)
	}
	user.Following = utils.RemoveElement(user.Following, int32(followingID))
	user.FollowingCount -= 1

	query = "UPDATE users SET following = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.Following, id)
	if err != nil {
		return fmt.Errorf("error updating user following: %v", err)
	}
	return nil
}

func (s *UserService) AcceptFollowRequest(ctx context.Context, id int64, followerID int64) error {
	err := kafka.PublishNotification(ctx, "follow_request_accepted", kafka.NotificationMessage{
		UserID:     int32(id),
		FollowerID: int32(followerID),
		Action:     "accept",
		Message:    "Follow request accepted",
	})
	if err != nil {
		return fmt.Errorf("error accepting follow request: %v", err)
	}
	return nil
}

func (s *UserService) RejectFollowRequest(ctx context.Context, id int64, followerID int64) error {
	err := kafka.PublishNotification(ctx, "follow_request_rejected", kafka.NotificationMessage{
		UserID:     int32(id),
		FollowerID: int32(followerID),
		Action:     "reject",
		Message:    "Follow request rejected",
	})
	if err != nil {
		return fmt.Errorf("error rejecting follow request: %v", err)
	}
	return nil
}

func (s *UserService) SendFollowRequest(ctx context.Context, id int64, followerID int64) error {
	err := kafka.PublishNotification(ctx, "follow_request_sent", kafka.NotificationMessage{
		UserID:     int32(id),
		FollowerID: int32(followerID),
		Action:     "request",
		Message:    "Follow request sent",
	})
	if err != nil {
		return fmt.Errorf("error sending follow request: %v", err)
	}
	return nil
}

func (s *UserService) VerifyCode(ctx context.Context, loginIdentifier, message string) (models.Code, error) {
	query := "SELECT * FROM codes WHERE login_identifier = $1 AND code = $2 ORDER BY created_at DESC LIMIT 1 RETURNING *"
	var code models.Code
	err := s.db.QueryRowContext(ctx, query, loginIdentifier, message).Scan(
		&code.ID,
		&code.LoginIdentifier,
		&code.Code,
		&code.CreatedAt,
	)
	if err != nil {
		return models.Code{}, fmt.Errorf("error verifying code: %v", err)
	}
	return code, nil
}

func (s *UserService) BlockUser(ctx context.Context, id int64, blockedID int64) error {
	query := "SELECT * FROM users WHERE id = $1"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error blocking user: %v", err)
	}
	user.BlockedUsers = append(user.BlockedUsers, int32(blockedID))

	blockedUser, err := s.GetUser(ctx, blockedID, id)
	if err != nil {
		return fmt.Errorf("error blocking user: %v", err)
	}

	var wg sync.WaitGroup

	// Remove following and followers
	wg.Add(4) // 4 ta operatsiya uchun

	go func() {
		defer wg.Done()
		if utils.InSlice(user.Following, int32(blockedID)) {
			err = s.RemoveFollowing(ctx, id, blockedID)
			if err != nil {
				fmt.Printf("error removing following: %v\n", err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if utils.InSlice(user.Followers, int32(blockedID)) {
			err = s.RemoveFollower(ctx, id, blockedID)
			if err != nil {
				fmt.Printf("error removing follower: %v\n", err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if utils.InSlice(blockedUser.Following, int32(id)) {
			err = s.RemoveFollowing(ctx, blockedID, id)
			if err != nil {
				fmt.Printf("error removing following: %v\n", err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if utils.InSlice(blockedUser.Followers, int32(id)) {
			err = s.RemoveFollower(ctx, blockedID, id)
			if err != nil {
				fmt.Printf("error removing follower: %v\n", err)
			}
		}
	}()

	wg.Wait() // Barcha goroutines tugagach kuting

	// Tweets va commentlarni o'chirish
	tweetsConn := methods.ConnectTweet()
	tweets, err := tweetsConn.GetTweetsByUser(ctx, &proto.GetTweetsByUserRequest{
		UserId: int32(id),
	})
	if err != nil {
		return fmt.Errorf("error getting tweets: %v", err)
	}

	commentsConn := methods.ConnectComment()
	likesConn := methods.ConnectLike()

	var wgTweets sync.WaitGroup

	for _, tweet := range tweets.Tweets {
		wgTweets.Add(2) // Like va commentlarni o'chirish uchun 2 ta operatsiya

		go func(tweet *proto.Tweet) {
			defer wgTweets.Done()
			if utils.InSlice(tweet.Likes, int32(blockedID)) {
				_, err := likesConn.DeleteLikeTweet(ctx, &proto.DeleteLikeTweetRequest{
					TweetId: tweet.Id,
					UserId:  int32(blockedID),
				})
				if err != nil {
					fmt.Printf("error removing like: %v\n", err)
				}
			}
		}(tweet)

		go func(tweet *proto.Tweet) {
			defer wgTweets.Done()
			if utils.InSlice(tweet.Comments, int32(blockedID)) {
				_, err := commentsConn.DeleteComment(ctx, &proto.DeleteCommentRequest{
					TweetId: int64(tweet.Id),
					UserId:  int64(blockedID),
				})
				if err != nil {
					fmt.Printf("error removing comment: %v\n", err)
				}
			}
		}(tweet)
	}

	wgTweets.Wait() 

	for _, tweet := range tweets.Tweets {
		var wgComments sync.WaitGroup

		for _, com := range tweet.Comments {
			wgComments.Add(1)
			go func(com int64) {
				defer wgComments.Done()
				comment, err := commentsConn.GetComment(ctx, &proto.GetCommentRequest{
					Id:     com,
					UserId: int64(id),
				})
				if err != nil {
					fmt.Printf("error getting comment: %v\n", err)
					return
				}
				if comment.Comment.UserId == id {
					if utils.InSliceInt64(comment.Comment.Likes, int64(blockedID)) {
						_, err := likesConn.DeleteLikeComment(ctx, &proto.DeleteLikeCommentRequest{
							CommentId: int32(comment.Comment.Id),
							UserId:    int32(blockedID),
						})
						if err != nil {
							fmt.Printf("error removing comment like: %v\n", err)
						}
					}
				}
			}(int64(com))
		}
		wgComments.Wait() 
	}

	tweetsBlocked, err := tweetsConn.GetTweetsByUser(ctx, &proto.GetTweetsByUserRequest{
		UserId: int32(blockedID),
	})
	if err != nil {
		return fmt.Errorf("error getting tweets: %v", err)
	}

	for _, tweet := range tweetsBlocked.Tweets {
		wgTweets.Add(2) 

		go func(tweet *proto.Tweet) {
			defer wgTweets.Done()
			if utils.InSlice(tweet.Likes, int32(id)) {
				_, err := likesConn.DeleteLikeTweet(ctx, &proto.DeleteLikeTweetRequest{
					TweetId: tweet.Id,
					UserId:  int32(id),
				})
				if err != nil {
					fmt.Printf("error removing like: %v\n", err)
				}
			}
		}(tweet)

		go func(tweet *proto.Tweet) {
			defer wgTweets.Done()
			if utils.InSlice(tweet.Comments, int32(id)) {
				_, err := commentsConn.DeleteComment(ctx, &proto.DeleteCommentRequest{
					TweetId: int64(tweet.Id),
					UserId:  int64(id),
				})
				if err != nil {
					fmt.Printf("error removing comment: %v\n", err)
				}
			}
		}(tweet)
	}

	wgTweets.Wait() 

	query = "UPDATE users SET blocked_users = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.BlockedUsers, id)
	if err != nil {
		return fmt.Errorf("error updating user blocked users: %v", err)
	}
	return nil
}

func (s *UserService) UnblockUser(ctx context.Context, id int64, blockedID int64) error {
	query := "SELECT * FROM users WHERE id = $1 RETURNING *"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.CreatedAt,
		&user.Bio,
		&user.IsPrivate,
		&user.Followers,
		&user.Following,
		&user.Tweets,
		&user.BlockedUsers,
	)
	if err != nil {
		return fmt.Errorf("error unblocking user: %v", err)
	}

	user.BlockedUsers = utils.RemoveElement(user.BlockedUsers, int32(blockedID))

	query = "UPDATE users SET blocked_users = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, user.BlockedUsers, id)
	if err != nil {
		return fmt.Errorf("error updating user blocked users: %v", err)
	}
	return nil
}
