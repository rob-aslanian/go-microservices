package files

import (
	"context"
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

type files struct {
	path string
}

func (f files) GetPath() string {
	return f.path
}

type Configuration interface {
	GetPath() string
}

func NewFileStorageRepository(config Configuration) (*files, error) {
	//create folder if not exists
	_, err := os.Stat(config.GetPath())
	if os.IsNotExist(err) {
		os.Mkdir(config.GetPath(), 0760) // TODO: find out which permissions should be
	}
	//---

	return &files{
		path: config.GetPath(),
	}, nil
}

func generateName() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}

func getBoundary(request *http.Request) (string, error) {
	_, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil {
		return "", errors.New("getting boundary fail")
	}

	return params["boundary"], nil
}

// Upload ...
func (f files) Upload(ctx context.Context, request *http.Request) ([]multipart.FileHeader, []string, error) {

	err := request.ParseMultipartForm(32 << 20) // 32 GB
	if err != nil {

		log.Println("cant parse multipart form:")
		var b []byte
		request.Body.Read(b)
		log.Println(string(b))

		return []multipart.FileHeader{}, []string{}, err
	}

	fileHeaders := make([]multipart.FileHeader, 0, len(request.MultipartForm.File))
	internalNames := make([]string, 0, len(request.MultipartForm.File))

	for _, value := range request.MultipartForm.File {
		for i := range value {
			// file, header, _ := request.FormFile(key)
			// defer file.Close()

			file, err := value[i].Open()
			if err != nil {
				return nil, []string{}, err
			}
			defer file.Close()

			header := &multipart.FileHeader{
				Filename: value[i].Filename,
				Header:   value[i].Header,
				Size:     value[i].Size,
			}

			name, err := generateName()
			if err != nil {
				return nil, []string{}, err
			}

			f, err := os.OpenFile(f.GetPath()+name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // TODO: find out which permissions should be

			if err != nil {
				log.Println("error: opening file:", err)
			}

			io.Copy(f, file)

			fileHeaders = append(fileHeaders, *header)
			internalNames = append(internalNames, name)
		}

	}

	return fileHeaders, internalNames, nil
}
