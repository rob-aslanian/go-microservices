package chat

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/net/context"
	"io"
)

func (r *ChatRepo) SaveFile(ctx context.Context, input io.Reader, name string, md interface{}) (*mgo.GridFile, error) {
	if gridFile, err := r.fs.Create(name); err != nil {
		return nil, err
	} else {
		gridFile.SetName(name)
		gridFile.SetMeta(md)
		size, err := io.Copy(gridFile, input)
		if err != nil {
			return nil, err
		}
		fmt.Println("file size: ", size)
		gridFile.Close()
		return gridFile, nil
	}
}

func (r *ChatRepo) ReadFile(ctx context.Context, fileId string) (*mgo.GridFile, error) {
	return r.fs.OpenId(bson.ObjectIdHex(fileId))
}
