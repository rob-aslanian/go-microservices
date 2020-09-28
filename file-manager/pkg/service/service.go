package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	fileResponse "gitlab.lan/Rightnao-site/microservices/file-manager/pkg/file_response"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"

	"google.golang.org/grpc/metadata"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type fileService struct {
	rpcClient   RPCClient
	fileList    FileListRepository
	fileStorage FileStorageRepository
	// validator
}

// NewFileService ...
func NewFileService(rpcClient RPCClient, fileList FileListRepository, fileStorage FileStorageRepository) (*fileService, error) {
	return &fileService{
		rpcClient:   rpcClient,
		fileList:    fileList,
		fileStorage: fileStorage,
	}, nil
}

func (f *fileService) Upload(ctx context.Context, token, companyID, target, targetID string, itemID string, request *http.Request, isBase64 bool) ([]fileResponse.FileResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span = span.Tracer().StartSpan("Upload", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	passThroughContext(&ctx)

	userID, err := f.getUserID(ctx, token)
	if err != nil {
		return []fileResponse.FileResponse{}, err
	}

	path := "./data/"

	var fileHeaders []multipart.FileHeader
	var names []string

	if isBase64 {
		if request.ContentLength > (10485760) { // 10 MB
			log.Println("ContentLength is ", request.ContentLength)
			return []fileResponse.FileResponse{}, errors.New("file more than 10 MB")
		}

		// log.Println("uploading base64 file")
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println(err)
		}

		var js map[string]interface{}

		var fileName string

		err = json.Unmarshal(body, &js)
		if err != nil {
			log.Println(err)
		}

		// log.Println(js)

		if _, isExists := js["file"]; !isExists {
			return []fileResponse.FileResponse{}, errors.New("file_is_absent")
		}

		if n, isExists := js["name"]; isExists {
			if s, ok := n.(string); ok {
				fileName = s
			}
		}

		r := regexp.MustCompile(`data:(?P<mime>.*);base64,(?P<image>.*)`)
		fil, ok := js["file"].(string)
		if !ok {
			return []fileResponse.FileResponse{}, errors.New("wrong request: file field is absent")
		}
		out := r.FindStringSubmatch(fil)

		dec, err := base64.StdEncoding.DecodeString(out[2])
		if err != nil {
			panic(err)
		}

		if len(out) < 2 {
			log.Println("wrong request", "len == ", len(out))
			return []fileResponse.FileResponse{}, errors.New("wrong request")
		}

		// n, _ := generateName()
		names = make([]string, 1)
		names[0], _ = generateName()

		target, err := os.OpenFile(path+names[0], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600) // move it
		if err != nil {
			log.Println(err)
		}
		defer target.Close()

		if _, err = target.Write(dec); err != nil {
			log.Println(err)
			return []fileResponse.FileResponse{}, err
		}
		if err = target.Sync(); err != nil {
			log.Println(err)
			return []fileResponse.FileResponse{}, err
		}

		fileHeaders = make([]multipart.FileHeader, 1)

		fileHeaders[0].Filename = fileName
		h := textproto.MIMEHeader{}
		h.Add("Content-Type", out[1])
		fileHeaders[0].Header = h

	} else {
		fileHeaders, names, err = f.fileStorage.Upload(ctx, request)

		if err != nil {
			return []fileResponse.FileResponse{}, err
		}
	}

	// log.Println("file uploaded:", names)

	fileInfoResponse := make([]fileResponse.FileResponse, 0, len(fileHeaders))

	for i := range fileHeaders {
		url := generateURL()
		err := f.fileList.SaveFileInfo(fileHeaders[i], userID, url, names[i])
		if err == nil {

			file := userRPC.File{
				UserID:   userID,
				Type:     userRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
				TargetID: targetID,
				Name:     fileHeaders[i].Filename,
				MimeType: fileHeaders[i].Header.Get("Content-Type"),
				URL:      url,
				// ID:
			}

			// galleryFile := companyRPC.GalleryFile{
			// 	UserID:    userID,
			// 	CompanyID: companyID,
			// 	Name:      fileHeaders[i].Filename,
			// 	MimeType:  fileHeaders[i].Header.Get("Content-Type"),
			// 	URL:       url,
			// }

			id := &userRPC.ID{}

			switch target {

			// User

			case "experience":
				id, err = f.rpcClient.User().AddFileInExperience(ctx, &file)

			case "accomplishment":
				id, err = f.rpcClient.User().AddFileInAccomplishment(ctx, &file)

			case "education":
				id, err = f.rpcClient.User().AddFileInEducation(ctx, &file)

			case "portfolio":
				id, err = f.rpcClient.User().AddFileInPortfolio(ctx, &file)

			case "company_award":
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().AddFileInCompanyAward(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "gallery":
				idComp := &companyRPC.ID{}
				// if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
				// 	return []fileResponse.FileResponse{}, errors.New("not_image_gallery")
				// }
				idComp, err = f.rpcClient.Company().AddFileInCompanyGallery(ctx, &companyRPC.GalleryFile{
					UserID:    userID,
					CompanyID: companyID,
					GalleryID: targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "interest":
				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				id, err = f.rpcClient.User().ChangeImageInterest(ctx, &file)

			case "interest_origin":
				_, err = f.rpcClient.User().ChangeOriginImageInInterest(ctx, &file)

			case "avatar_origin":
				_, err = f.rpcClient.User().ChangeOriginAvatar(ctx, &file)

			case "avatar":
				if isBase64 == false && !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				_, err = f.rpcClient.User().ChangeAvatar(ctx, &file)

				// TODO: make gif not animated
				// makeNotAnimated(ctx, path, names[i])

				// Company

			case "founder":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeCompanyFounderAvatar(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "company_avatar":
				// if targetID == "" {
				// 	return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				// }

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeAvatar(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "company_avatar_origin":
				// if targetID == "" {
				// 	return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				// }

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeOriginAvatar(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "company_cover":
				// if targetID == "" {
				// 	return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				// }

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeCover(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "company_cover_origin":
				// if targetID == "" {
				// 	return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				// }

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeOriginCover(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "milestone":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeImageMilestone(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "product":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeImageProduct(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

			case "service":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &companyRPC.ID{}
				idComp, err = f.rpcClient.Company().ChangeImageService(ctx, &companyRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idComp.GetID()

				// Jobs

			case "job_applicant":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				idJob := &jobsRPC.ID{}
				idJob, err = f.rpcClient.Jobs().UploadFileForApplication(ctx, &jobsRPC.File{
					UserID: userID,
					// CompanyID: companyID,
					// Type:     companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID: targetID,
					Name:     fileHeaders[i].Filename,
					MimeType: fileHeaders[i].Header.Get("Content-Type"),
					URL:      url,
				})
				id.ID = idJob.GetId()
			case "job_post":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				idJob := &jobsRPC.ID{}
				idJob, err = f.rpcClient.Jobs().UploadFileForJob(ctx, &jobsRPC.File{
					CompanyID: companyID,
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idJob.GetId()

			case "advert_gallery":
				if !isImage(fileHeaders[i].Header.Get("Content-Type")) {
					return []fileResponse.FileResponse{}, errors.New("not_image")
				}
				idComp := &advertRPC.ID{}
				idComp, err = f.rpcClient.Advert().AddImageToGallery(ctx, &advertRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					// Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID: targetID,
					Name:     fileHeaders[i].Filename,
					MimeType: fileHeaders[i].Header.Get("Content-Type"),
					URL:      url,
				})
				id.ID = idComp.GetID()

			case "service_request":
				idService := &servicesRPC.ID{}
				idService, err = f.rpcClient.Services().AddFileInServiceRequest(ctx, &servicesRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID, // service id
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idService.GetID()

			case "service_order":
				idService := &servicesRPC.ID{}
				idService, err = f.rpcClient.Services().AddFileInOrderService(ctx, &servicesRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID, // order id
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idService.GetID()

			/// V-office
			case "v_office":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				idService := &servicesRPC.ID{}
				idService, err = f.rpcClient.Services().ChangeVofficeCover(ctx, &servicesRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID,
					URL:       url,
					// Type:     companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
				})

				id.ID = idService.GetID()

			case "v_office_origin":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				idService := &servicesRPC.ID{}
				idService, err = f.rpcClient.Services().ChangeVofficeOriginCover(ctx, &servicesRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					// Type:     companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID: targetID,
					URL:      url,
				})

				// id.ID = idComp.GetID()
				id.ID = idService.GetID()
			case "v_service":
				idService := &servicesRPC.ID{}
				idService, err = f.rpcClient.Services().AddFileInVofficeService(ctx, &servicesRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID, // office id
					ItemID:    itemID,   // service id
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idService.GetID()

			case "v_portfolio":

				idPortfolio := &servicesRPC.ID{}
				idPortfolio, err = f.rpcClient.Services().AddFileInVOfficePortfolio(ctx, &servicesRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID, // office id
					ItemID:    itemID,   // portfolio id
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})
				id.ID = idPortfolio.GetID()

			case "feedback_bugs":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				idFeedbackBug := &stuffRPC.ID{}
				idFeedbackBug, err = f.rpcClient.Stuff().AddFileToFeedBackBug(ctx, &stuffRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})

				id.ID = idFeedbackBug.GetID()

			case "feedback_suggestion":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}

				idFeedbackSuggestion := &stuffRPC.ID{}
				idFeedbackSuggestion, err = f.rpcClient.Stuff().AddFileToFeedBackSuggestion(ctx, &stuffRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					TargetID:  targetID,
					Name:      fileHeaders[i].Filename,
					MimeType:  fileHeaders[i].Header.Get("Content-Type"),
					URL:       url,
				})

				id.ID = idFeedbackSuggestion.GetID()
			case "newsfeed_post":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}
				idComp := &newsfeedRPC.ID{}
				idComp, err = f.rpcClient.Newsfeed().AddFileInPost(ctx, &newsfeedRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					// Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID: targetID,
					Name:     fileHeaders[i].Filename,
					MimeType: fileHeaders[i].Header.Get("Content-Type"),
					URL:      url,
				})
				id.ID = idComp.GetID()

			case "newsfeed_post_comment":
				if targetID == "" {
					return []fileResponse.FileResponse{}, errors.New("target_id_is_empty")
				}
				idComp := &newsfeedRPC.ID{}
				idComp, err = f.rpcClient.Newsfeed().AddFileInPost(ctx, &newsfeedRPC.File{
					UserID:    userID,
					CompanyID: companyID,
					// Type:      companyRPC.File_TargetType(userRPC.File_TargetType_value[strings.Title(target)]),
					TargetID: targetID,
					ItemID:   itemID,
					Name:     fileHeaders[i].Filename,
					MimeType: fileHeaders[i].Header.Get("Content-Type"),
					URL:      url,
				})
				id.ID = idComp.GetID()

			default:
				log.Println("error! Unknown target: ", target) // TODO:
			}
			if err != nil {
				log.Println(err)
			}

			if !isBase64 && isImage(file.GetMimeType()) {
				err = postProcessImage(path, names[i])
				if err != nil {
					return []fileResponse.FileResponse{}, err
				}
			}

			fileInfoResponse = append(fileInfoResponse,
				fileResponse.FileResponse{
					URL:      file.GetURL(),
					ID:       id.GetID(),
					MimeType: file.GetMimeType(),
				},
			)

		} else {
			return []fileResponse.FileResponse{}, err
		}
	}
	return fileInfoResponse, nil
}

func (f *fileService) GetFile(ctx context.Context, response http.ResponseWriter, request *http.Request, fileURL string) error {
	span := opentracing.SpanFromContext(ctx)
	span = span.Tracer().StartSpan("GetFile", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	fileInfo, err := f.fileList.GetFileInfo(ctx, fileURL)
	if err != nil {
		http.Error(response, "file not found", http.StatusInternalServerError)
		return err
	}

	file, err := os.Open("./data/" + fileInfo.GetInternalName()) // should be fileID
	defer file.Close()
	if err != nil {
		http.Error(response, "file not found", http.StatusNotFound)
		return errors.New("file not found")
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		http.Error(response, "file not found", http.StatusNotFound)
		return errors.New("file not found")
	}

	http.ServeContent(response, request, fileURL, fileStat.ModTime(), file)

	return nil
}

func generateURL() string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	b := make([]rune, 60)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func passThroughContext(ctx *context.Context) error {
	span := opentracing.SpanFromContext(*ctx)
	span = span.Tracer().StartSpan("passThroughContext", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		return errors.New("token is empty")
	}
	return nil
}

func (f *fileService) getUserID(ctx context.Context, token string) (string, error) {
	span := opentracing.SpanFromContext(ctx)
	span = span.Tracer().StartSpan("getUserID", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	user, err := f.rpcClient.Auth().GetUser(ctx, &authRPC.Session{Token: token})
	if err != nil {
		return "", err
	}

	return user.GetId(), nil
}

// ---

func getResolution(width, height, resolution uint) (uint, uint) {
	w := width
	h := height

	if height > width {
		if height < resolution {
			return width, height
		}

		factor := float32(resolution) / float32(height)

		w = uint(float32(width) * factor)
		h = uint(float32(height) * factor)
	} else {
		if width < resolution {
			return width, height
		}

		factor := float32(resolution) / float32(width)

		w = uint(float32(width) * factor)
		h = uint(float32(height) * factor)
	}

	return w, h
}

func postProcessImage(path, image string) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(path + image)
	if err != nil {
		return err
	}

	cols := mw.GetImageWidth()
	rows := mw.GetImageHeight()

	size := func() uint {
		if cols > rows {
			return cols
		}
		return rows
	}()

	switch {
	case size > 1080:
		mw.ThumbnailImage(getResolution(cols, rows, 1080))
		mw.WriteImage(path + image + "-1080")
		fallthrough

	case size > 720:
		mw.ThumbnailImage(getResolution(cols, rows, 720))
		mw.WriteImage(path + image + "-720")
		fallthrough

	case size > 480:
		mw.ThumbnailImage(getResolution(cols, rows, 480))
		mw.WriteImage(path + image + "-480")
		fallthrough

	case size > 320:
		mw.ThumbnailImage(getResolution(cols, rows, 320))
		mw.WriteImage(path + image + "-320")
		fallthrough

	case size > 240:
		mw.ThumbnailImage(getResolution(cols, rows, 240))
		mw.WriteImage(path + image + "-240")
	}
	return nil
}

func isImage(m string) bool {
	if m == "image/jpeg" ||
		m == "image/png" /* || m == "image/gif" */ {
		return true
	}
	return false
}

func generateName() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}
