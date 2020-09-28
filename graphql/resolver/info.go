package resolver

import (
	"context"
	"log"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
)

// // TODO: finish
// func (_ *Resolver) GetListOfCountries(ctx context.Context) (*[]string, error) {
// 	countries, err := info.GetListOfCountries(ctx, &infoRPC.Empty{})
// 	if err != nil {
// 		// TODO: handle error
// 		log.Println(err)
// 		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
// 		if !b {
// 			return nil, err
// 		}
// 		return nil, e
// 	}
// 	c := countries.GetCountry()
// 	return &c, nil
// }

// TODO: finish
func (_ *Resolver) GetListOfCities(ctx context.Context, input GetListOfCitiesRequest) (*[]CityResolver, error) {
	cities, err := info.GetListOfCities(
		ctx,
		&infoRPC.GetCitiesRequest{
			CountryIso: input.Search_city.Country_id,
			FindCity: func() string {
				if input.Search_city.Find_city != nil {
					return *input.Search_city.Find_city
				}
				return ""
			}(),
			First: func() int32 {
				if input.Pagination.First != nil {
					return *input.Pagination.First
				}
				return 0
			}(),
			After: func() string {
				if input.Pagination.After != nil {
					return *input.Pagination.After
				}
				return ""
			}(),
		},
	)
	if err != nil {
		// TODO: handle error
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return nil, err
		}
		return nil, e
	}

	cityResolver := make([]CityResolver, 0, len(cities.GetCities()))
	for _, c := range cities.GetCities() {
		cityResolver = append(cityResolver, CityResolver{
			R: &City{
				ID:          strconv.Itoa(int(c.GetId())),
				City:        c.GetTitle(),
				Subdivision: c.GetSubdivision(),
			},
		})
	}

	return &cityResolver, nil
}

func (_ *Resolver) GetListOfAllCities(ctx context.Context, input GetListOfCitiesRequest) (*[]CityResolver, error) {
	cities, err := info.GetListOfAllCities(
		ctx,
		&infoRPC.GetAllCitiesRequest{
			// CountryIso: input.Search_city.Country,
			FindCity: func() string {
				if input.Search_city.Find_city != nil {
					return *input.Search_city.Find_city
				}
				return ""
			}(),
			// First: func() int32 {
			// 	if input.Pagination.First != nil {
			// 		return *input.Pagination.First
			// 	}
			// 	return 0
			// }(),
			// After: func() string {
			// 	if input.Pagination.After != nil {
			// 		return *input.Pagination.After
			// 	}
			// 	return ""
			// }(),
		},
	)
	if err != nil {
		// TODO: handle error
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return nil, err
		}
		return nil, e
	}

	cityResolver := make([]CityResolver, 0, len(cities.GetCities()))
	for _, c := range cities.GetCities() {
		cityResolver = append(cityResolver, CityResolver{
			R: &City{
				ID:          strconv.Itoa(int(c.GetId())),
				City:        c.GetTitle(),
				Subdivision: c.GetSubdivision(),
				Country:     c.GetCountry(),
			},
		})
	}

	return &cityResolver, nil
}

func (_ *Resolver) GetListOfCountryCodes(ctx context.Context) (*[]CountryCodeResolver, error) {
	countryCodes, err := info.GetListOfCountryCodes(ctx, &infoRPC.Empty{})
	if err != nil {
		// TODO: handle error
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return nil, err
		}
		return nil, e
	}

	countryResolver := make([]CountryCodeResolver, 0, len(countryCodes.GetCountryCodes()))
	for _, c := range countryCodes.GetCountryCodes() {
		countryResolver = append(countryResolver, CountryCodeResolver{
			R: &CountryCode{
				ID:           strconv.Itoa(int(c.GetId())),
				Country_code: c.GetCountryCode(),
				Country:      c.GetCountry(),
			},
		})
	}

	return &countryResolver, nil
}

func (_ *Resolver) GetCityInfo(ctx context.Context, input GetCityInfoRequest) (*CityResolver, error) {
	res, err := info.GetCityInfoByID(ctx, &infoRPC.IDWithLang{
		ID:   input.City_id,
		Lang: "en",
	})
	if err != nil {
		return nil, err
	}

	// n, err := strconv.Atoi()
	// if err != nil {
	// 	return nil, err
	// }

	return &CityResolver{
		R: &City{
			City:        res.GetTitle(),
			Country:     res.GetCountry(),
			ID:          strconv.Itoa(int(res.GetId())),
			Subdivision: res.GetSubdivision(),
		},
	}, nil
}

// func (_ *Resolver) GetListOfUiLanguages(ctx context.Context) (*[]string, error) {
// 	lang, err := info.GetListOfUiLanguages(ctx, &infoRPC.Empty{})
// 	if err != nil {
// 		// TODO: handle error
// 		log.Println(err)
// 		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
// 		if !b {
// 			return nil, err
// 		}
// 		return nil, e
// 	}
// 	p := lang.GetLanguages()
// 	return &p, nil
// }

// func (_ *Resolver) GetListOfAllIndustries(ctx context.Context) (*[]IndustryResolver, error) {
// 	result, err := info.GetListOfAllIndustries(ctx, &infoRPC.Language{})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	industryResolver := make([]IndustryResolver, 0, len(result.GetIndustry()))
// 	for _, c := range result.GetIndustry() {
// 		industryResolver = append(industryResolver, IndustryResolver{
// 			R: &Industry{
// 				ID:       c.GetID(),
// 				Industry: c.GetIndustry(),
// 			},
// 		})
// 	}
//
// 	return &industryResolver, nil
// }
//
// func (_ *Resolver) GetListOfAllSubindustries(ctx context.Context, in GetListOfAllSubindustriesRequest) (*[]SubindustryResolver, error) {
// 	result, err := info.GetListOfAllSubindustries(ctx, &infoRPC.IDWithLang{
// 		ID: in.ID,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	industryResolver := make([]SubindustryResolver, 0, len(result.GetSubindustry()))
// 	for _, c := range result.GetSubindustry() {
// 		industryResolver = append(industryResolver, SubindustryResolver{
// 			R: &Subindustry{
// 				ID:          c.GetID(),
// 				Subindustry: c.GetSubindustry(),
// 			},
// 		})
// 	}
//
// 	return &industryResolver, nil
// }
