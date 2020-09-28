package server

import (
	"context"
	"errors"
	"log"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"

	"google.golang.org/grpc/metadata"
)

// TODO: finish: find out language, refactor, optimize...
func (s *Server) GetListOfCities(ctx context.Context, in *infoRPC.GetCitiesRequest) (*infoRPC.ListOfCities, error) {
	after, err := strconv.Atoi(in.GetAfter())
	if err != nil {
		after = 0 // FIXME: find the better way
	}

	lang := s.retriveUILang(ctx)
	if lang == "" {
		lang = "en"
	}

	result, err := s.db.GetListOfCitiesNew(ctx, lang, in.GetCountryIso(), in.GetFindCity(), int(in.GetFirst()), after)
	if err != nil {
		return nil, err
	}

	cities := make([]*infoRPC.City, 0)
	for i := 0; i < len(*result); i++ {
		cities = append(cities, &infoRPC.City{
			Id: (*result)[i].Id,
			Title: func() string {
				if (*result)[i].City != nil {
					return *(*result)[i].City
				}
				return ""
			}(),
			Subdivision: func() string {
				if (*result)[i].Subdivision != nil {
					return *(*result)[i].Subdivision
				}
				return ""
			}(),
			Country: func() string {
				if (*result)[i].Country != nil {
					return *(*result)[i].Country
				}
				return ""
			}(),
		},
		)
	}

	return &infoRPC.ListOfCities{Cities: cities}, nil
}

func (s *Server) GetListOfAllCities(ctx context.Context, in *infoRPC.GetAllCitiesRequest) (*infoRPC.ListOfCities, error) {
	after, err := strconv.Atoi(in.GetAfter())
	if err != nil {
		after = 0 // FIXME: find the better way
	}
	// result, err := s.Db.GetListOfAllCities(ctx, "en", in.GetCountryIso(), in.GetFindCity(), int(in.GetFirst()), after)

	lang := s.retriveUILang(ctx)
	if lang == "" {
		lang = "en"
	}

	result, err := s.db.GetListOfAllCitiesNew(ctx, lang, in.GetFindCity(), int(in.GetFirst()), after)
	if err != nil {
		return nil, err
	}

	cities := make([]*infoRPC.City, 0)
	for i := 0; i < len(*result); i++ {
		cities = append(cities, &infoRPC.City{
			Id: (*result)[i].Id,
			Title: func() string {
				if (*result)[i].City != nil {
					return *(*result)[i].City
				}
				return ""
			}(),
			Subdivision: func() string {
				if (*result)[i].Subdivision != nil {
					return *(*result)[i].Subdivision
				}
				return ""
			}(),
			Country: func() string {
				if (*result)[i].Country != nil {
					return *(*result)[i].Country
				}
				return ""
			}(),
		},
		)
	}

	return &infoRPC.ListOfCities{Cities: cities}, nil
}

func (s *Server) GetListOfCountryCodes(ctx context.Context, in *infoRPC.Empty) (*infoRPC.ListOfCountryCodes, error) {
	result, err := s.db.GetListOfCountryCodesNew(ctx)
	if err != nil {
		return nil, err
	}
	countryCodes := make([]*infoRPC.CountryCode, 0)
	for i := 0; i < len(*result); i++ {
		countryCodes = append(countryCodes, &infoRPC.CountryCode{
			Id:          (*result)[i].Id,
			Country:     (*result)[i].Country,
			CountryCode: (*result)[i].CountryCode,
		})
	}
	return &infoRPC.ListOfCountryCodes{CountryCodes: countryCodes}, nil
}

func (s *Server) GetCountryCodeByID(ctx context.Context, in *infoRPC.CountryCode) (*infoRPC.CountryCode, error) {
	result, err := s.db.GetCountryCodeByIDNew(ctx, uint(in.GetId()))
	if err != nil {
		return nil, err
	}

	log.Println("result", result)

	return &infoRPC.CountryCode{
		Id:          in.GetId(),
		Country:     (*result).Country,
		CountryCode: (*result).CountryCode,
	}, nil
}

func (s *Server) GetCityInfoByID(ctx context.Context, in *infoRPC.IDWithLang) (*infoRPC.City, error) {
	id, err := strconv.Atoi(in.GetID())
	if err != nil {
		return nil, errors.New("wrong_id")
	}

	lang := s.retriveUILang(ctx)
	if lang == "" {
		lang = "en"
	}

	result, err := s.db.GetCityInfoByIDNew(ctx, lang, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	var city, subdivision, country string

	if result.City != nil {
		t := result.City
		city = *t
	}

	if result.Subdivision != nil {
		t := result.Subdivision
		subdivision = *t
	}

	if result.Country != nil {
		t := result.Country
		country = *t
	}

	return &infoRPC.City{
		Title:       city,
		Subdivision: subdivision,
		Country:     country,
	}, nil
}

func (s *Server) GetCityByIP(ctx context.Context, e *infoRPC.Empty) (*infoRPC.City, error) {
	res := make(map[string]string, 1)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("ip")
		if len(arr) > 0 {
			res["ip"] = arr[0]
		}
	}

	ip := res["ip"]

	city, err := s.mDB.GetCityByIp(ip, "")
	if err != nil {
		return nil, err
	}

	country, err := s.mDB.GetCountryByIp(ip, "")
	if err != nil {
		return nil, err
	}

	return &infoRPC.City{
		Id:          city.Id,
		Title:       *city.City,
		Subdivision: *city.Subdivision,
		Country:     *city.Country,
		CountryId:   country.Id,
		CountryIso:  country.CountryCode,
	}, nil

}

// func (s *Server) GetListOfUiLanguages(ctx context.Context, e *infoRPC.Empty) (*infoRPC.ListOfUiLanguages, error) {
// 	// TODO: remove hardcoded data
// 	languages := infoRPC.ListOfUiLanguages{Languages: []string{"en", "ru", "ka"}}
// 	return &languages, nil
// }

//
// func (s *Server) GetListOfCountries(ctx context.Context, in *infoRPC.Empty) (*infoRPC.ListOfCountries, error) {
// 	result, err := s.db.GetListOfCountries(ctx, "en")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &infoRPC.ListOfCountries{Country: *result}, nil
// }

func (s *Server) retriveUILang(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("ui_lang")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}
