package resolver

import (
	"context"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

type addedPostSubscriber struct {
	stop       <-chan struct{}
	events     chan<- *NewsfeedPostResolverCustom
	NewsfeedID string
}

func (r *Resolver) broadcastAddedPost() {
	subscribers := map[string]*addedPostSubscriber{}
	unsubscribe := make(chan string)

	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
			log.Println("unsubscribe, id:", id)

		case s := <-r.addedPostSubscriber:
			id := generateID()
			log.Println("subscribe, id:", id)
			subscribers[id] = s

		case e := <-r.AddedPostEvents:
			for id, s := range subscribers {

				go func(id string, s *addedPostSubscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}

					select {
					case <-s.stop:
						unsubscribe <- id
					default:
						if s.NewsfeedID == e.R.NewsFeedUserID ||
							s.NewsfeedID == e.R.NewsFeedCompanyID {
							s.events <- e
						}
					case <-time.After(time.Second):
					}
				}(id, s)

			}
		}
	}
}

type addedCommentPostSubscriber struct {
	stop   <-chan struct{}
	events chan<- *NewsfeedPostCommentResolverCustom
	PostID string
	Ctx    context.Context
}

func (r *Resolver) broadcastAddedCommentPost() {
	subscribers := map[string]*addedCommentPostSubscriber{}
	unsubscribe := make(chan string)

	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)

		case s := <-r.addedCommentPostSubscriber:
			id := generateID()
			subscribers[id] = s

		case e := <-r.AddedCommentPostEvents:
			for id, s := range subscribers {

				go func(id string, s *addedCommentPostSubscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}

					select {
					case <-s.stop:
						unsubscribe <- id
					default:
						if e.R != nil && s.PostID == e.R.PostID {
							if e.R.UserID != "" {
								userResp, err := user.GetMapProfilesByID(s.Ctx, &userRPC.UserIDs{
									ID: []string{e.R.UserID},
								})
								if err != nil {
									return
								}
								profile := ToProfile(s.Ctx, userResp.GetProfiles()[e.R.UserID])
								e.R.User_profile = &profile
							} else if e.R.CompanyID != "" {
								companyResp, err := company.GetCompanyProfiles(s.Ctx, &companyRPC.GetCompanyProfilesRequest{
									Ids: []string{e.R.CompanyID},
								})
								if err != nil {
									return
								}

								var profile *companyRPC.Profile

								for _, company := range companyResp.GetProfiles() {
									if company.GetId() == e.R.CompanyID {
										profile = company
										break
									}
								}

								if profile != nil {
									pr := toCompanyProfile(s.Ctx, *profile)
									e.R.Company_profile = &pr
								}
							}
							s.events <- e
						}
					case <-time.After(time.Second):
					}
				}(id, s)

			}
		}
	}
}

type addedLikePostSubscriber struct {
	stop      <-chan struct{}
	events    chan<- *LikeResolverCustom
	PostID    string
	CommentID *string
}

func (r *Resolver) broadcastAddedLikePost() {
	subscribers := map[string]*addedLikePostSubscriber{}
	unsubscribe := make(chan string)

	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)

		case s := <-r.addedLikePostSubscriber:
			id := generateID()
			subscribers[id] = s

		case e := <-r.AddedLikePostEvents:
			for id, s := range subscribers {

				go func(id string, s *addedLikePostSubscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}

					select {
					case <-s.stop:
						unsubscribe <- id
					default:
						if (s.CommentID != nil && e.R.CommentID != nil) &&
							s.PostID == e.R.PostID &&
							*s.CommentID == *e.R.CommentID {
							fmt.Println("*s.CommentID:", *s.CommentID)
							fmt.Println("*e.R.CommentID:", *e.R.CommentID)
							s.events <- e
						} else {
							if (s.CommentID == nil && e.R.CommentID == nil) &&
								s.PostID == e.R.PostID {
								s.events <- e
							}
						}
					case <-time.After(time.Second):
					}
				}(id, s)

			}
		}
	}
}

// AddedPost ...
func (r *Resolver) AddedPost(ctx context.Context, input AddedPostRequest) (<-chan *NewsfeedPostResolverCustom, error) {

	c := make(chan *NewsfeedPostResolverCustom)
	r.addedPostSubscriber <- &addedPostSubscriber{
		events:     c,
		stop:       ctx.Done(),
		NewsfeedID: input.ID,
	}

	return c, nil
}

// AddedPostComment ...
func (r *Resolver) AddedPostComment(ctx context.Context, input AddedPostCommentRequest) (<-chan *NewsfeedPostCommentResolverCustom, error) {
	c := make(chan *NewsfeedPostCommentResolverCustom)
	r.addedCommentPostSubscriber <- &addedCommentPostSubscriber{
		events: c,
		stop:   ctx.Done(),
		PostID: input.Post_id,
		Ctx:    ctx,
	}

	return c, nil
}

// AddedPostLike ...
func (r *Resolver) AddedPostLike(ctx context.Context, input AddedPostLikeRequest) (<-chan *LikeResolverCustom, error) {
	c := make(chan *LikeResolverCustom)
	r.addedLikePostSubscriber <- &addedLikePostSubscriber{
		events:    c,
		stop:      ctx.Done(),
		PostID:    input.Post_id,
		CommentID: input.Comment_id,
	}

	return c, nil
}

func generateID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
