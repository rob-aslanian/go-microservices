package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	companyadmin "gitlab.lan/Rightnao-site/microservices/news_feed/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
	"google.golang.org/grpc/metadata"
)

// AddPost ...
func (s Service) AddPost(ctx context.Context, p *post.Post) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddPost")
	defer span.Finish()

	err := p.Validate()
	if err != nil {
		return "", err
	}

	id := p.GenerateID()
	p.CreatedAt = time.Now()

	if p.GetCompanyID() != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			p.GetCompanyID(),
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return "", errors.New("not_allowed")
		}
	}

	p.Trim()
	p.FindHashtags()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = p.SetUserID(userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.repository.SavePost(ctx, p)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.mq.SendNewPostEvent(p)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, nil
}

// ChangePost ...
func (s Service) ChangePost(ctx context.Context, p *post.Post) error {
	span := s.tracer.MakeSpan(ctx, "ChangePost")
	defer span.Finish()

	err := p.ValidateText()
	if err != nil {
		return err
	}

	oldPost, err := s.repository.GetPostByID(ctx, p.GetID())
	if err != nil {
		return err
	}

	if p.GetCompanyID() != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			p.GetCompanyID(),
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
		if oldPost.GetCompanyID() != p.GetCompanyID() {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
		if oldPost.GetUserID() != userID {
			return errors.New("not_allowed")
		}
	}

	now := time.Now()
	p.ChangedAt = &now
	p.Trim()
	p.FindHashtags()

	err = s.repository.ChangePost(ctx, p)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemovePost ...
func (s Service) RemovePost(ctx context.Context, postID string, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemovePost")
	defer span.Finish()

	oldPost, err := s.repository.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	if companyID != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
		if oldPost.GetCompanyID() != companyID {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
		if oldPost.GetUserID() != userID {
			return errors.New("not_allowed")
		}
	}

	err = s.repository.RemovePost(ctx, postID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetNewsfeed ...
func (s Service) GetNewsfeed(ctx context.Context, id, companyID string, pinned bool, first, after uint32) (*post.Newsfeed, error) {
	span := s.tracer.MakeSpan(ctx, "GetNewsfeed")
	defer span.Finish()

	var feed *post.Newsfeed
	var err error

	myID := ""

	if companyID == "" {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
		myID = userID
	} else {
		// TODO: check if admin??
		myID = companyID
	}

	if id != "" {
		// get someone's newsfeed
		feed, err = s.repository.GetNewsfeed(ctx, myID, id, pinned, first, after)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
	} else {
		var isCompany bool
		if companyID != "" {
			isCompany = true
		}

		// get followings ids
		ids, err := s.networkRPC.GetFollowersIDs(ctx, myID, isCompany)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		ids = append(ids, myID) // append my id

		feed, err = s.repository.GetNewsfeedOfFollowings(ctx, myID, ids, first, after)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
	}

	// get amount of shares
	if feed != nil {
		postIDs := make([]string, 0, len(feed.Posts))
		for _, post := range feed.Posts {
			postIDs = append(postIDs, post.GetID())
		}

		shares, err := s.repository.GetAmountOfSharedPosts(ctx, postIDs)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		for _, post := range feed.Posts {
			post.SharesAmount = shares[post.GetID()]
		}
	}

	// // ---
	// for _, f := range feed.Posts {
	// 	if f.Liked != nil {
	// 		log.Println(*f.Liked)
	// 	}
	// }
	// // ---

	return feed, nil
}

// AddComment ...
func (s Service) AddComment(ctx context.Context, com *post.Comment) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddComment")
	defer span.Finish()

	err := com.Validate()
	if err != nil {
		return "", err
	}

	com.CreatedAt = time.Now()

	if com.GetCompanyID() != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			com.GetCompanyID(),
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return "", errors.New("not_allowed")
		}
	}

	com.Trim()
	id := com.GenerateID()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = com.SetUserID(userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.repository.AddComment(ctx, com)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.mq.SendNewCommentEvent(com)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, nil
}

// GetComments ...
func (s Service) GetComments(ctx context.Context, id string, companyID string, sort post.CommentSort, first, after uint32) (*post.Comments, error) {
	span := s.tracer.MakeSpan(ctx, "GetComments")
	defer span.Finish()

	var myID string
	if companyID == "" {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
		myID = userID
	} else {
		// TODO: check if admin??
		myID = companyID
	}

	com, err := s.repository.GetComments(ctx, myID, id, sort, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// get amount of replies
	if com != nil {
		commentIDs := make([]string, 0, len(com.Comments))
		for _, post := range com.Comments {
			commentIDs = append(commentIDs, post.GetID())
		}

		amountReplies, err := s.repository.GetAmountOfReplies(ctx, id, commentIDs)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		for i := range com.Comments {
			com.Comments[i].AmountOfReplies = amountReplies[com.Comments[i].GetID()]
		}
	}

	return com, nil
}

// RemoveComment ...
func (s Service) RemoveComment(ctx context.Context, postID, commentID, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveComment")
	defer span.Finish()

	var userID string

	if companyID != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	oldPost, err := s.repository.GetPostByID(ctx, postID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if oldPost.GetNewsFeedUserID() != userID ||
		oldPost.GetNewsFeedCompanyID() != companyID {
		// check owner
		if companyID != "" {
			if oldPost.GetCompanyID() != companyID {
				return errors.New("not_allowed")
			}
		} else {
			if oldPost.GetUserID() != userID {
				return errors.New("not_allowed")
			}
		}
	}

	err = s.repository.RemoveComment(ctx, postID, commentID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeComment ...
func (s Service) ChangeComment(ctx context.Context, com *post.Comment) error {
	span := s.tracer.MakeSpan(ctx, "ChangeComment")
	defer span.Finish()

	err := com.Validate()
	if err != nil {
		return err
	}

	now := time.Now()
	com.ChangedAt = &now
	com.Trim()

	var userID string

	if com.GetCompanyID() != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			com.GetCompanyID(),
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	oldPost, err := s.repository.GetPostByID(ctx, com.GetPostID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check owner
	if com.GetCompanyID() != "" {
		if oldPost.GetCompanyID() != com.GetCompanyID() {
			return errors.New("not_allowed")
		}
	} else {
		if oldPost.GetUserID() != userID {
			return errors.New("not_allowed")
		}
	}

	err = s.repository.ChangeComment(ctx, com)
	if err != nil {
		return err
	}

	return nil
}

// GetCommentReplies ...
func (s Service) GetCommentReplies(ctx context.Context, companyID, postID string, commentID string, first, after uint32) (*post.Comments, error) {
	span := s.tracer.MakeSpan(ctx, "GetCommentReplies")
	defer span.Finish()

	var myID string
	if companyID == "" {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
		myID = userID
	} else {
		// TODO: check if admin??
		myID = companyID
	}

	com, err := s.repository.GetCommentReplies(ctx, myID, postID, commentID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return com, nil
}

// GetSharedPosts ...
func (s Service) GetSharedPosts(ctx context.Context, companyID, id string, first, after uint32) (*post.Newsfeed, error) {
	span := s.tracer.MakeSpan(ctx, "GetSharedPosts")
	defer span.Finish()

	myID := ""

	if companyID == "" {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
		myID = userID
	} else {
		// TODO: check if admin??
		myID = companyID
	}

	feed, err := s.repository.GetSharedPosts(ctx, myID, id, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// get amount of shares
	if feed != nil {
		postIDs := make([]string, 0, len(feed.Posts))
		for _, post := range feed.Posts {
			postIDs = append(postIDs, post.GetID())
		}

		shares, err := s.repository.GetAmountOfSharedPosts(ctx, postIDs)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		for _, post := range feed.Posts {
			post.SharesAmount = shares[post.GetID()]
		}
	}

	return feed, nil
}

// GetPostByID ...
func (s Service) GetPostByID(ctx context.Context, postID string) (*post.Post, error) {
	span := s.tracer.MakeSpan(ctx, "GetPostByID")
	defer span.Finish()

	p, err := s.repository.GetPostByID(ctx, postID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// get amount of shares
	if p != nil {
		shares, err := s.repository.GetAmountOfSharedPosts(ctx, []string{p.GetID()})
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		p.SharesAmount = shares[p.GetID()]
	}

	return p, nil
}

// AddFile ...
func (s Service) AddFile(ctx context.Context, userID string, postID string, commentID, companyID string, f *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFile")
	defer span.Finish()

	id := f.GenerateID()

	if companyID != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return "", errors.New("not_allowed")
		}
	} else {
		// token := s.retriveToken(ctx)
		// var err error
		// userID, err = s.authRPC.GetUser(ctx, token)
		// if err != nil {
		// 	s.tracer.LogError(span, err)
		// 	return "", err
		// }
	}

	if commentID == "" {
		// get post by ID
		oldPost, err := s.repository.GetPostByID(ctx, postID)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}

		if len(oldPost.Files) > 10 {
			return "", errors.New("too_much_files")
		}

		// check owner
		if companyID != "" {
			if oldPost.GetCompanyID() != companyID {
				return "", errors.New("not_allowed")
			}
		} else {
			if oldPost.GetUserID() != userID {
				return "", errors.New("not_allowed")
			}
		}

		err = s.repository.AddFile(ctx, postID, f)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}
	} else {
		// get comment by ID
		com, err := s.repository.GetCommentByID(ctx, postID, commentID)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}

		if len(com.Files) > 10 {
			return "", errors.New("too_much_files")
		}

		// check owner
		if companyID != "" {
			if com.GetCompanyID() != companyID {
				return "", errors.New("not_allowed")
			}
		} else {
			if com.GetUserID() != userID {
				return "", errors.New("not_allowed")
			}
		}

		err = s.repository.AddFileInComment(ctx, postID, commentID, f)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}
	}

	return id, nil
}

// RemoveFile ...
func (s Service) RemoveFile(ctx context.Context, postID string, commentID, companyID string, fileID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFile")
	defer span.Finish()

	var userID string

	if companyID != "" {
		isAdmin := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	if commentID == "" {
		// get post by ID
		oldPost, err := s.repository.GetPostByID(ctx, postID)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}

		// check owner
		if companyID != "" {
			if oldPost.GetCompanyID() != companyID {
				return errors.New("not_allowed")
			}
		} else {
			if oldPost.GetUserID() != userID {
				return errors.New("not_allowed")
			}
		}

		err = s.repository.RemoveFile(ctx, postID, fileID)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	} else {
		// get comment by ID
		com, err := s.repository.GetCommentByID(ctx, postID, commentID)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}

		// check owner
		if companyID != "" {
			if com.GetCompanyID() != companyID {
				return errors.New("not_allowed")
			}
		} else {
			if com.GetUserID() != userID {
				return errors.New("not_allowed")
			}
		}

		err = s.repository.RemoveFileInComment(ctx, postID, commentID, fileID)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	return nil
}

// Like ...
func (s Service) Like(ctx context.Context, postID string, commentID string, like *post.Like) error {
	span := s.tracer.MakeSpan(ctx, "Like")
	defer span.Finish()

	switch like.Type {
	case "user":
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
		_ = like.SetID(userID)
	case "company":
		fallthrough
	case "organization":
		isAdmin := s.checkAdminLevel(
			ctx,
			like.GetID(),
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
	default:
		return errors.New("wrong_type")
	}

	like.CreatedAt = time.Now()

	if commentID != "" {
		err := s.repository.LikeComment(ctx, postID, commentID, like)
		if err != nil {
			return err
		}
	} else {
		err := s.repository.LikePost(ctx, postID, like)
		if err != nil {
			return err
		}
	}

	like.PostID = postID
	if commentID != "" {
		like.CommentID = &commentID
	}

	err := s.mq.SendNewLikeEvent(like)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// Unlike ...
func (s Service) Unlike(ctx context.Context, postID string, commentID string, id string) error {
	span := s.tracer.MakeSpan(ctx, "Unlike")
	defer span.Finish()

	var like *post.Like
	var err error

	if commentID == "" {
		like, err = s.repository.GetLikeInPostByID(ctx, postID, id)
		if err != nil {
			return err
		}
	} else {
		like, err = s.repository.GetLikeInCommentByID(ctx, postID, commentID, id)
		if err != nil {
			return err
		}
	}

	if like == nil {
		return nil
	}

	switch like.Type {
	case "user":
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
		_ = like.SetID(userID)
	case "company":
		fallthrough
	case "organization":
		isAdmin := s.checkAdminLevel(
			ctx,
			like.GetID(),
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return errors.New("not_allowed")
		}
	default:
		return errors.New("wrong_type")
	}

	if commentID != "" {
		err := s.repository.UnlikeComment(ctx, postID, commentID, id)
		if err != nil {
			return err
		}
	} else {
		err := s.repository.UnlikePost(ctx, postID, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLikedList ...
func (s Service) GetLikedList(ctx context.Context, postID, commentID string, emoji *post.EmojiType, first, after uint32) ([]*post.Like, error) {
	span := s.tracer.MakeSpan(ctx, "GetLikedList")
	defer span.Finish()

	var likes []*post.Like
	var err error

	if commentID == "" {
		likes, err = s.repository.GetLikedListInPost(ctx, postID, emoji, first, after)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
	} else {
		likes, err = s.repository.GetLikedListInComment(ctx, postID, commentID, emoji, first, after)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
	}

	for _, l := range likes {
		if l != nil {
			log.Println(*l)
		}
	}

	return likes, nil
}

// SearchAmongPosts ...
func (s Service) SearchAmongPosts(ctx context.Context, companyID, newsfeedID string, keyword string, first, after uint32) (*post.Newsfeed, error) {
	span := s.tracer.MakeSpan(ctx, "SearchAmongPosts")
	defer span.Finish()

	myID := ""

	if companyID == "" {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
		myID = userID
	} else {
		isAdmin := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
		)
		if !isAdmin {
			return nil, errors.New("not_allowed")
		}
		myID = companyID
	}

	k, h := splitKeywordsAndHashtags(keyword)

	feed, err := s.repository.Search(ctx, myID, newsfeedID, k, h, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// get amount of shares
	if feed != nil {
		postIDs := make([]string, 0, len(feed.Posts))
		for _, post := range feed.Posts {
			postIDs = append(postIDs, post.GetID())
		}

		shares, err := s.repository.GetAmountOfSharedPosts(ctx, postIDs)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		for _, post := range feed.Posts {
			post.SharesAmount = shares[post.GetID()]
		}
	}

	return feed, nil
}

func (s Service) retriveToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}

// checkAdminLevel return false if level doesn't much
func (s Service) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...companyadmin.AdminLevel) bool {
	span := s.tracer.MakeSpan(ctx, "checkAdminLevel")
	defer span.Finish()

	actualLevel, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		log.Println("Error: checkAdminLevel:", err)
		return false
	}

	for _, lvl := range requiredLevels {
		if lvl == actualLevel {
			return true
		}
	}

	return false
}

func splitKeywordsAndHashtags(keyword string) (string, []string) {
	words := strings.Fields(keyword)

	str := strings.Builder{}
	hashtag := make([]string, 0)

	for _, word := range words {
		if !strings.HasPrefix(word, "#") {
			str.WriteString(word + " ")
		} else {
			hashtag = append(hashtag, word)
		}
	}

	return str.String(), hashtag
}
