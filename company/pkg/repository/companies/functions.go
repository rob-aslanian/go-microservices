package companiesRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	careercenter "gitlab.lan/Rightnao-site/microservices/company/pkg/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/status"
)

const (
	companiesCollection = "companies"
	reportsCollection   = "reports"
)

// SaveNewCompanyAccount ...
func (r Repository) SaveNewCompanyAccount(ctx context.Context, data *account.Account) error {
	err := r.collections[companiesCollection].Insert(data)
	if err != nil {
		return err
	}

	return nil
}

// ChangeStatusOfCompany ...
func (r Repository) ChangeStatusOfCompany(ctx context.Context, companyID string, stat status.CompanyStatus) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"status": stat,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyAccount ...
func (r Repository) GetCompanyAccount(ctx context.Context, companyID string) (*account.Account, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	acc := account.Account{}

	err := r.collections[companiesCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
	).One(&acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// ChangeCompanyName ...
func (r Repository) ChangeCompanyName(ctx context.Context, companyID string, name string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"name": name,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsURLBusy checks if url is busy already
func (r Repository) IsURLBusy(ctx context.Context, url string) (bool, error) {
	count, err := r.collections[companiesCollection].Find(
		bson.M{
			// "status": status.CompanyStatusActivated,
			"url": url,
		}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// ChangeCompanyURL ...
func (r Repository) ChangeCompanyURL(ctx context.Context, companyID string, url string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"url": url,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyFoundationDate ...
func (r Repository) ChangeCompanyFoundationDate(ctx context.Context, companyID string, foundationDate time.Time) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"foundation_date": foundationDate,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyIndustry ...
func (r Repository) ChangeCompanyIndustry(ctx context.Context, companyID string, industry *account.Industry) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"industry": industry,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyType ...
func (r Repository) ChangeCompanyType(ctx context.Context, companyID string, companyType account.Type) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"type": companyType,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanySize ...
func (r Repository) ChangeCompanySize(ctx context.Context, companyID string, size account.Size) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"size": size,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyEmail ...
func (r Repository) AddCompanyEmail(ctx context.Context, companyID string, email *account.Email) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"emails": email,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyEmail ...
func (r Repository) DeleteCompanyEmail(ctx context.Context, companyID string, emailID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(emailID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"emails": bson.M{
					"id": bson.ObjectIdHex(emailID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsEmailExists returns true if such activated email was added in DB
func (r Repository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.collections[companiesCollection].Find(
		bson.M{
			"status": status.CompanyStatusActivated,
			"emails": bson.M{
				"$elemMatch": bson.M{
					"email":     email,
					"activated": true,
				},
			},
		}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// MakeEmailPrimary ...
func (r Repository) MakeEmailPrimary(ctx context.Context, companyID string, emailID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(emailID) {
		return errors.New("wrong_id")
	}

	// set primary false for all emails
	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"emails.id": bson.ObjectIdHex(emailID),
		},
		bson.M{
			"$set": bson.M{
				"emails.$[].primary": false,
			},
		},
	)
	if err != nil {
		return err
	}

	// set primary true given email
	err = r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"emails.id": bson.ObjectIdHex(emailID),
		},
		bson.M{
			"$set": bson.M{
				"emails.$.primary": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsEmailActivated returns true if email is activated
func (r Repository) IsEmailActivated(ctx context.Context, companyID string, emailID string) (bool, error) {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(emailID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(companyID),
		"emails": bson.M{
			"$elemMatch": bson.M{
				"id":        bson.ObjectIdHex(emailID),
				"activated": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IsEmailPrimary returns true if email is primary
func (r Repository) IsEmailPrimary(ctx context.Context, companyID string, emailID string) (bool, error) {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(emailID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(companyID),
		"emails": bson.M{
			"$elemMatch": bson.M{
				"id":      bson.ObjectIdHex(emailID),
				"primary": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IsEmailAdded checks if email already added for in this company
func (r Repository) IsEmailAdded(ctx context.Context, companyID string, email string) (bool, error) {
	if !bson.IsObjectIdHex(companyID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id":          bson.ObjectIdHex(companyID),
		"emails.email": email,
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// ActivateEmail ...
func (r Repository) ActivateEmail(ctx context.Context, companyID string, email string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(companyID),
			"emails.email": email,
		},
		bson.M{
			"$set": bson.M{
				"emails.$.activated": true,
			},
		},
	)

	return err
}

// AddCompanyPhone ...
func (r Repository) AddCompanyPhone(ctx context.Context, companyID string, phone *account.Phone) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"phones": phone,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyPhone ...
func (r Repository) DeleteCompanyPhone(ctx context.Context, companyID string, phoneID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(phoneID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"phones": bson.M{
					"id": bson.ObjectIdHex(phoneID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// MakePhonePrimary ...
func (r Repository) MakePhonePrimary(ctx context.Context, companyID string, phoneID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(phoneID) {
		return errors.New("wrong_id")
	}

	// set primary false for all phones
	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"phones.id": bson.ObjectIdHex(phoneID),
		},
		bson.M{
			"$set": bson.M{
				"phones.$[].primary": false,
			},
		},
	)
	if err != nil {
		return err
	}

	// set primary true given phone
	err = r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"phones.id": bson.ObjectIdHex(phoneID),
		},
		bson.M{
			"$set": bson.M{
				"phones.$.primary": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsPhoneExists returns true if phone exists
func (r Repository) IsPhoneExists(ctx context.Context, companyID string, phone *account.Phone) (bool, error) {
	count, err := r.collections[companiesCollection].Find(bson.M{
		"phones.id":              phone.Number,
		"phones.country_code.id": phone.CountryCode.ID,
		"activated":              true,
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IsPhonesActivated returns true if phone is activated
func (r Repository) IsPhonesActivated(ctx context.Context, companyID string, phoneID string) (bool, error) {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(phoneID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(companyID),
		"phones": bson.M{
			"$elemMatch": bson.M{
				"id":        bson.ObjectIdHex(phoneID),
				"activated": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IsPhonePrimary returns true if phone is primary
func (r Repository) IsPhonePrimary(ctx context.Context, companyID string, phoneID string) (bool, error) {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(phoneID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(companyID),
		"phones": bson.M{
			"$elemMatch": bson.M{
				"id":      bson.ObjectIdHex(phoneID),
				"primary": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IsPhoneAdded checks if phone already added for in this company
func (r Repository) IsPhoneAdded(ctx context.Context, companyID string, phone *account.Phone) (bool, error) {
	if !bson.IsObjectIdHex(companyID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id":                    bson.ObjectIdHex(companyID),
		"phones.number":          phone.Number,
		"phones.country_code.id": phone.CountryCode.ID,
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// IsPhoneActivated returns true if phone is activated
func (r Repository) IsPhoneActivated(ctx context.Context, companyID string, phoneID string) (bool, error) {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(phoneID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[companiesCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(companyID),
		"phones": bson.M{
			"$elemMatch": bson.M{
				"id":        bson.ObjectIdHex(phoneID),
				"activated": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// AddCompanyAddress ...
func (r Repository) AddCompanyAddress(ctx context.Context, companyID string, address *account.Address) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	// remove current primary fields
	if address.IsPrimary {
		err := r.collections[companiesCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
			bson.M{
				"$set": bson.M{
					"addresses.$[].is_primary": false,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	// adding new address
	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"addresses": address,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyAddress ...
func (r Repository) ChangeCompanyAddress(ctx context.Context, companyID string, address *account.Address) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(address.GetID()) {
		return errors.New("wrong_id")
	}

	// remove current primary fields
	if address.IsPrimary {
		err := r.collections[companiesCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
			bson.M{
				"$set": bson.M{
					"addresses.$[].is_primary": false,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(companyID),
			"addresses.id": bson.ObjectIdHex(address.GetID()),
		},
		bson.M{
			"$set": bson.M{
				"addresses.$": address,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyAddress ...
func (r Repository) DeleteCompanyAddress(ctx context.Context, companyID string, addressID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(addressID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"addresses": bson.M{
					"id": bson.ObjectIdHex(addressID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyWebsite ...
func (r Repository) AddCompanyWebsite(ctx context.Context, companyID string, website *account.Website) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"websites": website,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyWebsite ...
func (r Repository) DeleteCompanyWebsite(ctx context.Context, companyID string, websiteID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(websiteID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"websites": bson.M{
					"id": bson.ObjectIdHex(websiteID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyWebsite ...
func (r Repository) ChangeCompanyWebsite(ctx context.Context, companyID string, websiteID string, website string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(websiteID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"websites.id": bson.ObjectIdHex(websiteID),
		},
		bson.M{
			"$set": bson.M{
				"websites.$.website": website,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyParking ...
func (r Repository) ChangeCompanyParking(ctx context.Context, companyID string, parking account.Parking) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"parking": parking,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyBenefits ...
func (r Repository) ChangeCompanyBenefits(ctx context.Context, companyID string, benefits []profile.Benefit) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"benefits": benefits,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyProfile ...
func (r Repository) GetCompanyProfile(ctx context.Context, url string) (*profile.Profile, error) {
	prof := profile.Profile{}

	err := r.collections[companiesCollection].Find(
		bson.M{
			"url": url,
		},
	).One(&prof)
	if err != nil {
		return nil, err
	}
	return &prof, nil
}

// GetCompanyProfileByID ...
func (r Repository) GetCompanyProfileByID(ctx context.Context, companyID string) (*profile.Profile, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	prof := profile.Profile{}

	err := r.collections[companiesCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
	).One(&prof)
	if err != nil {
		return nil, err
	}
	return &prof, nil
}

// GetCompanyProfiles ...
func (r Repository) GetCompanyProfiles(ctx context.Context, ids []string) ([]*profile.Profile, error) {
	prof := []*profile.Profile{}

	idsBSON := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsBSON = append(idsBSON, bson.ObjectIdHex(ids[i]))
		} else {
			log.Println("wrong company_id:", ids[i])
		}
	}

	err := r.collections[companiesCollection].Find(
		bson.M{
			"_id": bson.M{
				"$in": idsBSON,
			},
		},
	).All(&prof)
	if err != nil {
		return nil, err
	}
	return prof, nil
}

// ChangeCompanyAboutUs ...
func (r Repository) ChangeCompanyAboutUs(ctx context.Context, companyID string, aboutUs *profile.AboutUs) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	bs := bson.M{}

	if !aboutUs.IsDescriptionNull {
		bs["description"] = aboutUs.Description
	}

	if !aboutUs.IsMissionNull {
		bs["mission"] = aboutUs.Mission
	}

	if !aboutUs.IsTypeNull {
		bs["type"] = aboutUs.Type
	}

	if !aboutUs.IsSizeNull {
		bs["size"] = aboutUs.Size
	}

	if !aboutUs.IsParkingNull {
		bs["parking"] = aboutUs.Parking
	}

	if aboutUs.Industry.Main != "" {
		bs["industry.main"] = aboutUs.Industry.Main
	}

	if !aboutUs.IsSubindustryNull {
		bs["industry.sub"] = aboutUs.Industry.Sub
	}

	if !aboutUs.FoundationDate.IsZero() {
		bs["foundation_date"] = aboutUs.FoundationDate
	}

	bs["business_hours"] = aboutUs.BusinessHours

	bs["changed_about_us"] = true

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bs,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetFounders ...
func (r Repository) GetFounders(ctx context.Context, companyID string, first int32, afterNumber int) ([]*profile.Founder, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	result := struct {
		Founders []*profile.Founder `bson:"founders"`
	}{}

	err := r.collections[companiesCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$project": bson.M{
				"founders": 1,
			},
		},
		{
			"$skip": afterNumber,
		},
		{
			"$limit": first,
		},
	},
	).One(&result)
	if err != nil {
		return nil, err
	}

	return result.Founders, nil
}

// AddCompanyFounder ...
func (r Repository) AddCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}
	if founder.GetUserID() != "" && !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	// bs := bson.M{}
	//
	// bs["id"] = founder.ID
	// bs["position"] = founder.Position
	// bs["created_at"] = founder.CreatedAt
	//
	// if founder.UserID != nil {
	// 	bs["user_id"] = *founder.UserID
	// } else {
	// 	bs["name"] = founder.Name
	// 	bs["avatar"] = founder.Avatar
	// }

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"founders": founder,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyFounder ...
func (r Repository) DeleteCompanyFounder(ctx context.Context, companyID string, founderID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(founderID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"founders": bson.M{
					"id": bson.ObjectIdHex(founderID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyFounder ...
func (r Repository) ChangeCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(founder.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}
	if founder.Name != "" {
		updateQuery["founders.$.name"] = founder.Name
	}
	if founder.Position != "" {
		updateQuery["founders.$.position"] = founder.Position
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"founders.id": bson.ObjectIdHex(founder.GetID()),
		},
		bson.M{
			"$set": updateQuery,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyFounderAvatar ...
func (r Repository) ChangeCompanyFounderAvatar(ctx context.Context, companyID string, founderID string, image *profile.File) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}
	if founderID != "" && !bson.IsObjectIdHex(founderID) {
		return errors.New("wrong_id")
	}

	var err error

	if founderID != "" {
		// add to the existing milestone
		err = r.collections[companiesCollection].Update(
			bson.M{
				"_id":         bson.ObjectIdHex(companyID),
				"founders.id": bson.ObjectIdHex(founderID),
			},
			bson.M{
				"$set": bson.M{
					"founders.$.avatar": image.URL,
				},
			},
		)
	} else {
		// in case founder it is nesaccery to upload image in certaing item
		return errors.New("wrong_id_of_founder")
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveCompanyFounderAvatar ...
func (r Repository) RemoveCompanyFounderAvatar(ctx context.Context, companyID string, founderID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(founderID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"founders.id": bson.ObjectIdHex(founderID),
		},
		bson.M{
			"$set": bson.M{
				"founders.$.avatar": nil,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ApproveFounderRequest ...
func (r Repository) ApproveFounderRequest(ctx context.Context, companyID, requestID, userID string, byCompanyAdmin bool) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(requestID) || !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	if byCompanyAdmin {
		err := r.collections[companiesCollection].Update(
			bson.M{
				"_id":         bson.ObjectIdHex(companyID),
				"founders.id": bson.ObjectIdHex(requestID),
			},
			bson.M{
				"$set": bson.M{
					"founders.$.approved": true,
				},
			})
		if err != nil {
			return err
		}
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":              bson.ObjectIdHex(companyID),
			"founders.user_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"founders.$.approved": true,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveFounderRequest ...
func (r Repository) RemoveFounderRequest(ctx context.Context, companyID, requestID, userID string, byCompanyAdmin bool) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(requestID) || !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	if byCompanyAdmin {
		err := r.collections[companiesCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
			bson.M{
				"$pull": bson.M{
					"founders": bson.M{
						"id": bson.ObjectIdHex(requestID),
					},
				},
			})
		if err != nil {
			return err
		}
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"founders": bson.M{
					"user_id": bson.ObjectIdHex(userID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyAward ...
func (r Repository) AddCompanyAward(ctx context.Context, companyID string, award *profile.Award) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	fileIDs := make([]bson.ObjectId, 0, len(award.Files))
	for i := range award.Files {
		fileIDs = append(fileIDs, award.Files[i].ID)
	}

	// append files
	m := []struct {
		File *profile.File `bson:"unused_files"`
	}{}
	err := r.collections[companiesCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(companyID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "CompanyAward",
					"unused_files.id": bson.M{
						"$in": fileIDs,
					},
				},
			},
		},
	).All(&m)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	award.Files = make([]*profile.File, len(m))
	for i := range m {
		award.Files[i] = m[i].File
	}

	// adding experience
	err = r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"awards": award,
			},
		})

	if err != nil {
		return err
	}

	// delete old files
	err = r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Awards",
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyAward ...
func (r Repository) DeleteCompanyAward(ctx context.Context, companyID string, awardID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"awards": bson.M{
					"id": bson.ObjectIdHex(awardID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyAward ...
func (r Repository) ChangeCompanyAward(ctx context.Context, companyID string, award *profile.Award) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(award.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}
	if award.Issuer != "" {
		updateQuery["awards.$.issuer"] = award.Issuer
	}
	if award.Title != "" {
		updateQuery["awards.$.title"] = award.Title
	}
	if !award.Date.IsZero() {
		updateQuery["awards.$.date"] = award.Date
	}

	if len(award.Links) > 0 {
		updateQuery["awards.$.links"] = award.Links
	}

	if len(updateQuery) == 0 {
		return errors.New("nothing_to_change")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"awards.id": bson.ObjectIdHex(award.GetID()),
		},
		bson.M{
			"$set": updateQuery,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveCompanyAward ...
func (r Repository) RemoveCompanyAward(ctx context.Context, companyID, awardID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"awards": bson.M{
					"id": bson.ObjectIdHex(awardID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddFileInCompanyAward ...
func (r Repository) AddFileInCompanyAward(ctx context.Context, companyID, awardID string, file *profile.File) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_company_id")
	}
	if awardID != "" && !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_award_id")
	}

	var err error

	f := struct {
		File *profile.File `bson:",inline"`
		Type string
	}{
		File: file,
		Type: "CompanyAward",
	}

	if awardID != "" {
		err = r.collections[companiesCollection].Update(
			bson.M{
				"_id":       bson.ObjectIdHex(companyID),
				"awards.id": bson.ObjectIdHex(awardID),
			},
			bson.M{
				"$push": bson.M{
					"awards.$.files": f,
				},
			},
		)
	} else {
		err = r.collections[companiesCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
			bson.M{
				"$push": bson.M{
					"unused_files": f,
				},
			},
		)
	}

	if err != nil {
		return err
	}
	return nil
}

// RemoveFilesInCompanyAward ...
func (r Repository) RemoveFilesInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"awards.id": bson.ObjectIdHex(awardID),
		},
		bson.M{
			"$pull": bson.M{
				"awards.$.files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddLinksInCompanyAward ...
func (r Repository) AddLinksInCompanyAward(ctx context.Context, companyID, awardID string, links []*profile.Link) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"awards.id": bson.ObjectIdHex(awardID),
		},
		bson.M{
			"$push": bson.M{
				"awards.$.links": bson.M{
					"$each": links,
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeLinkInCompanyAward ...
func (r Repository) ChangeLinkInCompanyAward(ctx context.Context, companyID, awardID string, linkID string, url string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) || !bson.IsObjectIdHex(linkID) {
		return errors.New("wrong_id")
	}

	result := struct {
		Index int `bson:"index"`
	}{}

	err := r.collections[companiesCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$project": bson.M{
				"awards": 1,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$awards",
				"includeArrayIndex":          "awards.index",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"awards.id": bson.ObjectIdHex(awardID),
			},
		},
		{
			"$project": bson.M{
				"index": "$awards.index",
			},
		},
	}).One(&result)
	if err != nil {
		return err
	}

	fmt.Println("Index:", result.Index)

	err = r.collections[companiesCollection].Update(
		bson.M{
			"_id":             bson.ObjectIdHex(companyID),
			"awards.links.id": bson.ObjectIdHex(linkID),
		},
		bson.M{
			"$set": bson.M{
				"awards." + strconv.Itoa(result.Index) + ".links.$.url": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveLinksInCompanyAward ...
func (r Repository) RemoveLinksInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"awards.id": bson.ObjectIdHex(awardID),
		},
		bson.M{
			"$pull": bson.M{
				"awards.$.links": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetUploadedFilesInCompanyAward ...
func (r Repository) GetUploadedFilesInCompanyAward(ctx context.Context, companyID string) ([]*profile.File, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	f := []struct {
		File *profile.File `bson:"unused_files"`
	}{}

	result := r.collections[companiesCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(companyID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "CompanyAward",
				},
			},
		},
	)
	err := result.All(&f)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		return nil, err
	}

	files := make([]*profile.File, 0, len(f))

	for i := range f {
		files = append(files, f[i].File)
	}

	return files, nil
}

// GetCompanyGallery ...
func (r Repository) GetCompanyGallery(ctx context.Context, companyID string, first uint32, afterNumber uint32) ([]*profile.File, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	// result := struct {
	// 	gallery []*profile.Gallery `bson:"gallery"`
	// }{}
	result := struct {
		Files []*profile.File `bson:"gallery"`
	}{}

	err := r.collections[companiesCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$project": bson.M{
				"gallery": 1,
			},
		},
		{
			"$skip": afterNumber,
		},
		{
			"$limit": first,
		},
	},
	).One(&result)
	if err != nil {
		return nil, err
	}

	return result.Files, nil
}

// AddFileInCompanyGallery ...
func (r Repository) AddFileInCompanyGallery(ctx context.Context, companyID string, file *profile.File) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_company_id")
	}

	var err error

	err = r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"gallery": file,
			},
		},
	)

	if err != nil {
		return err
	}
	return nil
}

// RemoveFilesInCompanyGallery ...
func (r Repository) RemoveFilesInCompanyGallery(ctx context.Context, companyID string, ids []string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"gallery": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetUploadedFilesInCompanyGallery ...
func (r Repository) GetUploadedFilesInCompanyGallery(ctx context.Context, companyID string) ([]*profile.File, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	f := []struct {
		File *profile.File `bson:"unused_files"`
	}{}

	result := r.collections[companiesCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(companyID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "CompanyGallery",
				},
			},
		},
	)
	err := result.All(&f)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		return nil, err
	}

	files := make([]*profile.File, 0, len(f))

	for i := range f {
		files = append(files, f[i].File)
	}

	return files, nil
}

// AddCompanyMilestone ...
func (r Repository) AddCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"milestones": milestone,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyMilestone ...
func (r Repository) DeleteCompanyMilestone(ctx context.Context, companyID string, milestoneID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(milestoneID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"milestones": bson.M{
					"id": bson.ObjectIdHex(milestoneID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyMilestone ...
func (r Repository) ChangeCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(milestone.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}
	if milestone.Title != "" {
		updateQuery["milestones.$.title"] = milestone.Title
	}
	// if milestone.Image!=
	if milestone.Description != "" {
		updateQuery["milestones.$.description"] = milestone.Description
	}
	if !milestone.Date.IsZero() {
		updateQuery["milestones.$.date"] = milestone.Date
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(companyID),
			"milestones.id": bson.ObjectIdHex(milestone.GetID()),
		},
		bson.M{
			"$set": updateQuery,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeImageMilestone ...
func (r Repository) ChangeImageMilestone(ctx context.Context, companyID string, milestoneID string, image *profile.File) error {

	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}
	if milestoneID != "" && !bson.IsObjectIdHex(milestoneID) {
		return errors.New("wrong_id")
	}

	var err error

	if milestoneID != "" {
		// add to the existing milestone
		err = r.collections[companiesCollection].Update(
			bson.M{
				"_id":           bson.ObjectIdHex(companyID),
				"milestones.id": bson.ObjectIdHex(milestoneID),
			},
			bson.M{
				"$set": bson.M{
					"milestones.$.image": image.URL,
				},
			},
		)
	} else {
		// log.Println("unused image for milestone") // delete
		//
		// // add for new milestone
		// f := struct {
		// 	File *profile.File `bson:",inline"`
		// 	Type string
		// }{
		// 	File: image,
		// 	Type: "Milestone",
		// }
		//
		// // retrive uploaded file
		// err = r.collections[companiesCollection].Pipe(
		// 	[]bson.M{
		// 		bson.M{
		// 			"$match": bson.M{
		// 				"_id": bson.ObjectIdHex(companyID),
		// 			},
		// 		},
		// 		bson.M{
		// 			"$project": bson.M{
		// 				"_id":          0,
		// 				"unused_files": 1,
		// 			},
		// 		},
		// 		bson.M{
		// 			"$unwind": bson.M{
		// 				"path": "$unused_files",
		// 			},
		// 		},
		// 		bson.M{
		// 			"$match": bson.M{
		// 				"unused_files.type": "Milestone",
		// 			},
		// 		},
		// 	},
		// ).One(&f)
		// if err != nil && err != mgo.ErrNotFound {
		// 	return err
		// }
		//
		// // if file has been uploaded
		// if f.File != nil {
		// 	log.Println("file has been uploaded")
		// 	err = r.collections[companiesCollection].Update(
		// 		bson.M{
		// 			"_id": bson.ObjectIdHex(companyID),
		// 		},
		// 		bson.M{
		// 			"$push": bson.M{
		// 				"unused_files": image,
		// 			},
		// 		},
		// 	)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {

		// in case milestone it is nesaccery to upload image in certaing item
		return errors.New("wrong_id_of_milestone")

		// log.Println("file has not been uploaded")
		// err = r.collections[companiesCollection].Update(
		// 	bson.M{
		// 		"_id":             bson.ObjectIdHex(companyID),
		// 		"unused_files.id": f.File.ID,
		// 	},
		// 	bson.M{
		// 		"$set": bson.M{
		// 			"unused_files.$.url": image.URL,
		// 		},
		// 	},
		// )
		// if err != nil {
		// 	return err
		// }
		// }
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveImageInMilestone ...
func (r Repository) RemoveImageInMilestone(ctx context.Context, companyID string, milestoneID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(milestoneID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(companyID),
			"milestones.id": bson.ObjectIdHex(milestoneID),
		},
		bson.M{
			"$set": bson.M{
				"milestones.$.image": nil,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyProduct ...
func (r Repository) AddCompanyProduct(ctx context.Context, companyID string, product *profile.Product) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"products": product,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyProduct ...
func (r Repository) DeleteCompanyProduct(ctx context.Context, companyID string, productID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(productID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"products": bson.M{
					"id": bson.ObjectIdHex(productID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyProduct ...
func (r Repository) ChangeCompanyProduct(ctx context.Context, companyID string, product *profile.Product) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(product.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}
	// if product.Image != "" {
	// 	updateQuery["products.$.image"] = product.Image
	// }
	if product.Name != "" {
		updateQuery["products.$.name"] = product.Name
	}
	if product.Website != "" {
		updateQuery["products.$.website"] = product.Website
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"products.id": bson.ObjectIdHex(product.GetID()),
		},
		bson.M{
			"$set": updateQuery,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeImageProduct ...
func (r Repository) ChangeImageProduct(ctx context.Context, companyID string, productID string, image *profile.File) error {

	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}
	if productID != "" && !bson.IsObjectIdHex(productID) {
		return errors.New("wrong_id")
	}

	var err error

	if productID != "" {
		// add to the existing milestone
		err = r.collections[companiesCollection].Update(
			bson.M{
				"_id":         bson.ObjectIdHex(companyID),
				"products.id": bson.ObjectIdHex(productID),
			},
			bson.M{
				"$set": bson.M{
					"products.$.image": image.URL,
				},
			},
		)
	} else {
		// log.Println("unused image for milestone") // delete
		//
		// // add for new milestone
		// f := struct {
		// 	File *profile.File `bson:",inline"`
		// 	Type string
		// }{
		// 	File: image,
		// 	Type: "Milestone",
		// }
		//
		// // retrive uploaded file
		// err = r.collections[companiesCollection].Pipe(
		// 	[]bson.M{
		// 		bson.M{
		// 			"$match": bson.M{
		// 				"_id": bson.ObjectIdHex(companyID),
		// 			},
		// 		},
		// 		bson.M{
		// 			"$project": bson.M{
		// 				"_id":          0,
		// 				"unused_files": 1,
		// 			},
		// 		},
		// 		bson.M{
		// 			"$unwind": bson.M{
		// 				"path": "$unused_files",
		// 			},
		// 		},
		// 		bson.M{
		// 			"$match": bson.M{
		// 				"unused_files.type": "Milestone",
		// 			},
		// 		},
		// 	},
		// ).One(&f)
		// if err != nil && err != mgo.ErrNotFound {
		// 	return err
		// }
		//
		// // if file has been uploaded
		// if f.File != nil {
		// 	log.Println("file has been uploaded")
		// 	err = r.collections[companiesCollection].Update(
		// 		bson.M{
		// 			"_id": bson.ObjectIdHex(companyID),
		// 		},
		// 		bson.M{
		// 			"$push": bson.M{
		// 				"unused_files": image,
		// 			},
		// 		},
		// 	)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {

		// in case milestone it is nesaccery to upload image in certaing item
		return errors.New("wrong_id_of_product")

		// log.Println("file has not been uploaded")
		// err = r.collections[companiesCollection].Update(
		// 	bson.M{
		// 		"_id":             bson.ObjectIdHex(companyID),
		// 		"unused_files.id": f.File.ID,
		// 	},
		// 	bson.M{
		// 		"$set": bson.M{
		// 			"unused_files.$.url": image.URL,
		// 		},
		// 	},
		// )
		// if err != nil {
		// 	return err
		// }
		// }
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveImageInProduct ...
func (r Repository) RemoveImageInProduct(ctx context.Context, companyID string, productID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(productID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"products.id": bson.ObjectIdHex(productID),
		},
		bson.M{
			"$set": bson.M{
				"products.$.image": nil,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyService ...
func (r Repository) AddCompanyService(ctx context.Context, companyID string, service *profile.Service) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$push": bson.M{
				"services": service,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyService ...
func (r Repository) DeleteCompanyService(ctx context.Context, companyID string, serviceID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(serviceID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$pull": bson.M{
				"services": bson.M{
					"id": bson.ObjectIdHex(serviceID),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyService ...
func (r Repository) ChangeCompanyService(ctx context.Context, companyID string, service *profile.Service) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(service.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}

	if service.Name != "" {
		updateQuery["services.$.name"] = service.Name
	}
	if service.Website != "" {
		updateQuery["services.$.website"] = service.Website
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"services.id": bson.ObjectIdHex(service.GetID()),
		},
		bson.M{
			"$set": updateQuery,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeImageService ...
func (r Repository) ChangeImageService(ctx context.Context, companyID string, serviceID string, image *profile.File) error {

	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}
	if serviceID != "" && !bson.IsObjectIdHex(serviceID) {
		return errors.New("wrong_id")
	}

	var err error

	if serviceID != "" {
		// add to the existing milestone
		err = r.collections[companiesCollection].Update(
			bson.M{
				"_id":         bson.ObjectIdHex(companyID),
				"services.id": bson.ObjectIdHex(serviceID),
			},
			bson.M{
				"$set": bson.M{
					"services.$.image": image.URL,
				},
			},
		)
	} else {
		// log.Println("unused image for milestone") // delete
		//
		// // add for new milestone
		// f := struct {
		// 	File *profile.File `bson:",inline"`
		// 	Type string
		// }{
		// 	File: image,
		// 	Type: "Milestone",
		// }
		//
		// // retrive uploaded file
		// err = r.collections[companiesCollection].Pipe(
		// 	[]bson.M{
		// 		bson.M{
		// 			"$match": bson.M{
		// 				"_id": bson.ObjectIdHex(companyID),
		// 			},
		// 		},
		// 		bson.M{
		// 			"$project": bson.M{
		// 				"_id":          0,
		// 				"unused_files": 1,
		// 			},
		// 		},
		// 		bson.M{
		// 			"$unwind": bson.M{
		// 				"path": "$unused_files",
		// 			},
		// 		},
		// 		bson.M{
		// 			"$match": bson.M{
		// 				"unused_files.type": "Milestone",
		// 			},
		// 		},
		// 	},
		// ).One(&f)
		// if err != nil && err != mgo.ErrNotFound {
		// 	return err
		// }
		//
		// // if file has been uploaded
		// if f.File != nil {
		// 	log.Println("file has been uploaded")
		// 	err = r.collections[companiesCollection].Update(
		// 		bson.M{
		// 			"_id": bson.ObjectIdHex(companyID),
		// 		},
		// 		bson.M{
		// 			"$push": bson.M{
		// 				"unused_files": image,
		// 			},
		// 		},
		// 	)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {

		// in case milestone it is nesaccery to upload image in certaing item
		return errors.New("wrong_id_of_service")

		// log.Println("file has not been uploaded")
		// err = r.collections[companiesCollection].Update(
		// 	bson.M{
		// 		"_id":             bson.ObjectIdHex(companyID),
		// 		"unused_files.id": f.File.ID,
		// 	},
		// 	bson.M{
		// 		"$set": bson.M{
		// 			"unused_files.$.url": image.URL,
		// 		},
		// 	},
		// )
		// if err != nil {
		// 	return err
		// }
		// }
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveImageInService ...
func (r Repository) RemoveImageInService(ctx context.Context, companyID string, serviceID string) error {
	if !bson.IsObjectIdHex(serviceID) || !bson.IsObjectIdHex(serviceID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":         bson.ObjectIdHex(companyID),
			"services.id": bson.ObjectIdHex(serviceID),
		},
		bson.M{
			"$set": bson.M{
				"services.$.image": nil,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyReport ...
func (r Repository) AddCompanyReport(ctx context.Context, companyID string, report *profile.Report) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(report.GetID()) {
		return errors.New("wrong_id")
	}

	if err := r.collections[reportsCollection].Insert(report); err != nil {
		return err
	}

	return nil
}

// ChangeAvatar ...
func (r Repository) ChangeAvatar(ctx context.Context, companyID string, image string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"avatar": image,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeOriginAvatar ...
func (r Repository) ChangeOriginAvatar(ctx context.Context, companyID string, image string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"avatar_origin": image,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAvatar ...
func (r Repository) RemoveAvatar(ctx context.Context, companyID string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"avatar_origin": nil,
				"avatar":        nil,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetOriginAvatar ...
func (r Repository) GetOriginAvatar(ctx context.Context, companyID string) (string, error) {
	if !bson.IsObjectIdHex(companyID) {
		return "", errors.New("wrong_id")
	}

	result := struct {
		Avatar string `bson:"avatar_origin"`
	}{}

	err := r.collections[companiesCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$project": bson.M{
				"avatar_origin": 1,
			},
		},
	}).One(&result)
	if err != nil {
		return "", err
	}

	return result.Avatar, nil
}

// ChangeCover ...
func (r Repository) ChangeCover(ctx context.Context, companyID string, image string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"cover": image,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeOriginCover ...
func (r Repository) ChangeOriginCover(ctx context.Context, companyID string, image string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"cover_origin": image,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveCover ...
func (r Repository) RemoveCover(ctx context.Context, companyID string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"cover_origin": nil,
				"cover":        nil,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetOriginCover ...
func (r Repository) GetOriginCover(ctx context.Context, companyID string) (string, error) {
	if !bson.IsObjectIdHex(companyID) {
		return "", errors.New("wrong_id")
	}

	result := struct {
		Cover string `bson:"avatar_origin"`
	}{}

	err := r.collections[companiesCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$project": bson.M{
				"cover_origin": 1,
			},
		},
	}).One(&result)
	if err != nil {
		return "", err
	}

	return result.Cover, nil
}

// SaveCompanyProfileTranslation ...
func (r Repository) SaveCompanyProfileTranslation(ctx context.Context, companyID string, lang string, tr *profile.Translation) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"translation." + lang: tr,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveCompanyMilestoneTranslation ...
func (r Repository) SaveCompanyMilestoneTranslation(ctx context.Context, companyID, milestoneID, language string, translation *profile.MilestoneTranslation) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(milestoneID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(companyID),
			"milestones.id": bson.ObjectIdHex(milestoneID),
		},
		bson.M{
			"$set": bson.M{
				"milestones.$.translations." + language: translation,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveCompanyAwardTranslation ...
func (r Repository) SaveCompanyAwardTranslation(ctx context.Context, companyID, awardID, language string, translation *profile.AwardTranslation) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(awardID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(companyID),
			"awards.id": bson.ObjectIdHex(awardID),
		},
		bson.M{
			"$set": bson.M{
				"awards.$.translations." + language: translation,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// OpenCareerCenter ...
func (r Repository) OpenCareerCenter(ctx context.Context, companyID string, cc *careercenter.CareerCenter) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companiesCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				"career_center": cc,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}
