package service

//
// import (
// 	"testing"
// 	"time"
//
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
// )
//
// func TestAccomplishmentCertificateValidator(t *testing.T) {
//
// 	time1, _ := time.Parse(time.RFC1123, "Thu, 31 Dec 1959 00:00:00 MST")
// 	time2 := time.Now()
// 	tru := true
// 	fals := false
// 	validNumber := "12"
// 	zero := uint32(0)
// 	nonZero := uint32(5)
// 	link := "www.google.com"
// 	inValid := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
//
// 	tables := []struct {
// 		Cert   *profile.Accomplishment
// 		Result error
// 	}{
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &validNumber,
// 				IsExpire:      &tru,
// 				FinishDate:    &time2,
// 				URL:           &link,
// 				Type:          profile.AccomplishmentTypeCertificate,
// 			},
// 			Result: usersErrors.SpecificRequired,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "dfd",
// 				LicenseNumber: nil,
// 				Issuer:        &validNumber,
// 				IsExpire:      &tru,
// 				FinishDate:    &time2,
// 				URL:           &link,
// 				Type:          profile.AccomplishmentTypeCertificate,
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232xx",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &validNumber,
// 				IsExpire:      &tru,
// 				FinishDate:    &time1,
// 				URL:           &link,
// 				Type:          profile.AccomplishmentTypeCertificate,
// 			},
// 			Result: usersErrors.InValidTime,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232assss",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &inValid,
// 				IsExpire:      &tru,
// 				FinishDate:    &time1,
// 				URL:           &link,
// 				Type:          profile.AccomplishmentTypeCertificate,
// 			},
// 			Result: usersErrors.Max100,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232df",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &inValid,
// 				IsExpire:      &fals,
// 				FinishDate:    &time1,
// 				URL:           &link,
// 				Type:          profile.AccomplishmentTypeCertificate,
// 			},
// 			Result: usersErrors.Max100,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232ff",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &inValid,
// 				IsExpire:      &fals,
// 				FinishDate:    &time1,
// 				URL:           &link,
// 				Description:   &validNumber,
// 				Type:          profile.AccomplishmentTypeAward,
// 			},
// 			Result: usersErrors.Max100,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232fgh",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &inValid,
// 				IsExpire:      &fals,
// 				StartDate:     &time1,
// 				FinishDate:    &time2,
// 				URL:           &link,
// 				Description:   &validNumber,
// 				Type:          profile.AccomplishmentTypeProject,
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232 cd",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &validNumber,
// 				IsExpire:      &fals,
// 				StartDate:     &time2,
// 				FinishDate:    &time1,
// 				URL:           &link,
// 				Description:   &validNumber,
// 				Type:          profile.AccomplishmentTypeProject,
// 			},
// 			Result: usersErrors.InValidTime,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232 sd",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &inValid,
// 				IsExpire:      &fals,
// 				StartDate:     &time1,
// 				FinishDate:    &time2,
// 				URL:           &link,
// 				Description:   &validNumber,
// 				Type:          profile.AccomplishmentTypePublication,
// 			},
// 			Result: usersErrors.Max100,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232 dvb",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &validNumber,
// 				IsExpire:      &fals,
// 				StartDate:     &time1,
// 				FinishDate:    &time2,
// 				Score:         &zero,
// 				URL:           &link,
// 				Description:   &validNumber,
// 				Type:          profile.AccomplishmentTypeTest,
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Cert: &profile.Accomplishment{
// 				Name:          "23232aa",
// 				LicenseNumber: &validNumber,
// 				Issuer:        &validNumber,
// 				IsExpire:      &fals,
// 				StartDate:     &time1,
// 				FinishDate:    &time2,
// 				Score:         &nonZero,
// 				URL:           &link,
// 				Description:   &validNumber,
// 				Type:          profile.AccomplishmentTypeTest,
// 			},
// 			Result: nil,
// 		},
// 	}
//
// 	for _, ta := range tables {
// 		err := accomplishmentValidator(ta.Cert)
// 		if err != ta.Result {
// 			t.Errorf("Error: accomplishmentValidator(%+v): got %v, expected: %v", ta.Cert, err, nil)
// 		}
// 	}
// }
