package clientRPC

import (
	"context"
	"log"
	"strconv"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"google.golang.org/grpc"
)

// Info represents client of Info
type Info struct {
	infoClient infoRPC.InfoServiceClient
}

// NewInfoClient crates new gRPC client of Info
func NewInfoClient(settings Settings) Info {
	i := Info{}
	i.connect(settings.Address)
	return i
}

func (i *Info) connect(address string) {
	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	i.infoClient = infoRPC.NewInfoServiceClient(conn)
}

// GetCountryIDAndCountryCode ...
func (i Info) GetCountryIDAndCountryCode(ctx context.Context, countryCodeID int32) (countryCode string, countryID string, err error) {
	result, err := i.infoClient.GetCountryCodeByID(ctx, &infoRPC.CountryCode{
		Id: countryCodeID,
	})

	err = handleError(err)
	if err != nil {
		return "", "", err
	}

	return result.GetCountryCode(), result.GetCountry(), err
}

// GetCityInformationByID ...
func (i Info) GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error) {
	var l string

	if lang != nil {
		l = *lang
	}

	result, err := i.infoClient.GetCityInfoByID(ctx, &infoRPC.IDWithLang{
		ID:   strconv.Itoa(int(cityID)),
		Lang: l,
	})
	err = handleError(err)
	if err != nil {
		return
	}

	return result.GetTitle(), result.GetSubdivision(), result.GetCountry(), nil
}

// // GetUserID returns user id
// func (a Info) GetUserID(ctx context.Context, token string) (string, error) {
// 	u, err := a.infoClient.GetUser(ctx, &infoRPC.Session{
// 		Token: token,
// 	})
//
// 	handleError(err)
//
// 	// ---------------
//
// 	return u.GetId(), nil
// }
