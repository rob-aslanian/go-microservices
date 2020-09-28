package clientRPC

import (
	"context"
	"log"
	"strconv"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/utils"
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

// // GetCountryIDAndCountryCode returns country ID and phone country code by country code's ID
// func (i Info) GetCountryIDAndCountryCode(ctx context.Context, countryCodeID int32) (countryCode string, countryID string, err error) {
// 	result, err := i.infoClient.GetCountryCodeByID(ctx, &infoRPC.CountryCode{
// 		Id: countryCodeID,
// 	})
//
// 	// TODO: handle error
//
// 	if err != nil {
// 		return "", "", err
// 	}
//
// 	return result.GetCountryCode(), result.GetCountry(), err
// }

// GetCityInformationByID returns names of city, subdivision and country
func (i Info) GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error) {
	var l string

	if lang != nil {
		l = *lang
	}

	result, err := i.infoClient.GetCityInfoByID(ctx, &infoRPC.IDWithLang{
		ID:   strconv.Itoa(int(cityID)),
		Lang: l,
	})

	// TOOD: handle error

	if err != nil {
		return
	}

	return result.GetTitle(), result.GetSubdivision(), result.GetCountry(), nil
}

// GetUserCountry returns country's ISO code
func (i Info) GetUserCountry(ctx context.Context) (string, error) {
	city, err := i.infoClient.GetCityByIP(utils.ToOutContext(ctx), &infoRPC.Empty{})
	if err != nil {
		return "", err
	}

	return city.CountryIso, nil
}
