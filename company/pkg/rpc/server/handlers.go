package serverRPC

import (
	"context"
	"log"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	careercenter "gitlab.lan/Rightnao-site/microservices/company/pkg/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/status"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
)

// CheckIfURLForCompanyIsTaken ...
func (s Server) CheckIfURLForCompanyIsTaken(ctx context.Context, data *companyRPC.URL) (*companyRPC.BooleanValue, error) {
	b, err := s.service.CheckIfURLForCompanyIsTaken(ctx, data.GetURL())
	if err != nil {
		return nil, err
	}

	return &companyRPC.BooleanValue{
		Value: b,
	}, nil
}

// RegisterCompany ...
func (s Server) RegisterCompany(ctx context.Context, data *companyRPC.RegisterCompanyRequest) (*companyRPC.RegisterCompanyResponse, error) {
	id, url, err := s.service.CreateNewAccount(ctx, registerCompanyRequestRPCToCompanyAccount(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.RegisterCompanyResponse{
		Id:  id,
		URL: url,
	}, nil
}

// DeactivateCompany ...
func (s Server) DeactivateCompany(ctx context.Context, data *companyRPC.DeactivateCompanyRequest) (*companyRPC.Empty, error) {
	err := s.service.DeactivateCompany(ctx, data.GetId(), data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetCompanyAccount ...
func (s Server) GetCompanyAccount(ctx context.Context, data *companyRPC.ID) (*companyRPC.Account, error) {
	acc, err := s.service.GetCompanyAccount(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return accountToAccountRPC(acc), nil
}

// ChangeCompanyName ...
func (s Server) ChangeCompanyName(ctx context.Context, data *companyRPC.ChangeCompanyNameRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyName(ctx, data.GetId(), data.GetName())
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanyURL ...
func (s Server) ChangeCompanyURL(ctx context.Context, data *companyRPC.ChangeCompanyUrlRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyURL(ctx, data.GetId(), data.GetUrl())
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanyFoundationDate ...
func (s Server) ChangeCompanyFoundationDate(ctx context.Context, data *companyRPC.ChangeCompanyFoundationDateRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyFoundationDate(ctx, data.GetId(), stringDateToTime(data.GetFoundationDate()))
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanyIndustry ...
func (s Server) ChangeCompanyIndustry(ctx context.Context, data *companyRPC.ChangeCompanyIndustryRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyIndustry(ctx, data.GetID(), industryRPCToAccountIndustry(data.GetIndustry()))
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanyType ...
func (s Server) ChangeCompanyType(ctx context.Context, data *companyRPC.ChangeCompanyTypeRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyType(ctx, data.GetId(), companyTypeRPCToAccountType(data.GetType()))
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanySize ...
func (s Server) ChangeCompanySize(ctx context.Context, data *companyRPC.ChangeCompanySizeRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanySize(ctx, data.GetId(), sizeRPCToAccountSize(data.GetSize()))
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// AddCompanyEmail ...
func (s Server) AddCompanyEmail(ctx context.Context, data *companyRPC.AddCompanyEmailRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyEmail(ctx, data.GetID(), emailRPCToAccountEmail(data.GetEmail()))
	if err != nil {
		return nil, err
	}
	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyEmail ...
func (s Server) DeleteCompanyEmail(ctx context.Context, data *companyRPC.DeleteCompanyEmailRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyEmail(ctx, data.GetId(), data.GetEmailId())
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanyEmail makes email primary if it activated
func (s Server) ChangeCompanyEmail(ctx context.Context, data *companyRPC.ChangeCompanyEmailRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyEmail(ctx, data.GetId(), data.GetEmailId())
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// AddCompanyPhone ...
func (s Server) AddCompanyPhone(ctx context.Context, data *companyRPC.AddCompanyPhoneRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyPhone(ctx, data.GetId(), phoneRPCToAccountPhone(data.GetPhone()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyPhone ...
func (s Server) DeleteCompanyPhone(ctx context.Context, data *companyRPC.DeleteCompanyPhoneRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyPhone(ctx, data.GetId(), data.GetPhoneId())
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// ChangeCompanyPhone ...
func (s Server) ChangeCompanyPhone(ctx context.Context, data *companyRPC.ChangeCompanyPhoneRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyPhone(ctx, data.GetId(), data.GetPhoneId())
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// AddCompanyAddress ...
func (s Server) AddCompanyAddress(ctx context.Context, data *companyRPC.AddCompanyAddressRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyAddress(ctx, data.GetId(), addressRPCToAccountAddress(data.GetAddress()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyAddress ...
func (s Server) DeleteCompanyAddress(ctx context.Context, data *companyRPC.DeleteCompanyAddressRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyAddress(ctx, data.GetId(), data.GetAddressId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyAddress ...
func (s Server) ChangeCompanyAddress(ctx context.Context, data *companyRPC.ChangeCompanyAddressRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyAddress(ctx, data.GetId(), addressRPCToAccountAddress(data.GetAddress()))
	if err != nil {
		return nil, err
	}
	return &companyRPC.Empty{}, nil
}

// AddCompanyWebsite ...
func (s Server) AddCompanyWebsite(ctx context.Context, data *companyRPC.AddCompanyWebsiteRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyWebsite(ctx, data.GetId(), data.GetWebsite())
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyWebsite ...
func (s Server) DeleteCompanyWebsite(ctx context.Context, data *companyRPC.DeleteCompanyWebsiteRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyWebsite(ctx, data.GetId(), data.GetWebsiteId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyWebsite ...
func (s Server) ChangeCompanyWebsite(ctx context.Context, data *companyRPC.ChangeCompanyWebsiteRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyWebsite(ctx, data.GetId(), data.GetWebsiteId(), data.GetWebsite())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyParking ...
func (s Server) ChangeCompanyParking(ctx context.Context, data *companyRPC.ChangeCompanyParkingRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyParking(ctx, data.GetId(), parkingRPCToAccountParking(data.GetParking()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyBenefits ...
func (s Server) ChangeCompanyBenefits(ctx context.Context, data *companyRPC.ChangeCompanyBenefitsRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyBenefits(ctx, data.GetID(), benefitsRPCToProfileBenefits(data.GetCompanyBenefit()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddCompanyAdmin ...
func (s Server) AddCompanyAdmin(ctx context.Context, data *companyRPC.AddCompanyAdminRequest) (*companyRPC.Empty, error) {
	err := s.service.AddCompanyAdmin(ctx, data.GetId(), data.GetUserId(), adminRoleRPCToAccountAdminLevel(data.GetRole()), data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// DeleteCompanyAdmin ...
func (s Server) DeleteCompanyAdmin(ctx context.Context, data *companyRPC.DeleteCompanyAdminRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyAdmin(ctx, data.GetId(), data.GetUserId(), data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetCompanyProfile ...
func (s Server) GetCompanyProfile(ctx context.Context, data *companyRPC.GetCompanyProfileRequest) (*companyRPC.GetCompanyProfileResponse, error) {
	profile, lvl, err := s.service.GetCompanyProfile(ctx, data.GetUrl(), data.GetLang())
	if err != nil {
		return nil, err
	}

	prof := profileToProfileRPC(profile)
	prof.Role = accountAdminLevelToAdminRoleRPC(lvl)

	return &companyRPC.GetCompanyProfileResponse{
		Profile: prof,
	}, nil
}

// GetCompanyProfileByID ...
func (s Server) GetCompanyProfileByID(ctx context.Context, data *companyRPC.ID) (*companyRPC.GetCompanyProfileResponse, error) {
	profile, lvl, err := s.service.GetCompanyProfileByID(ctx, data.GetID(), "") // TODO: add lang
	if err != nil {
		return nil, err
	}

	prof := profileToProfileRPC(profile)
	prof.Role = accountAdminLevelToAdminRoleRPC(lvl)

	return &companyRPC.GetCompanyProfileResponse{
		Profile: prof,
	}, nil
}

// GetCompanyProfiles ...
func (s Server) GetCompanyProfiles(ctx context.Context, data *companyRPC.GetCompanyProfilesRequest) (*companyRPC.GetCompanyProfilesResponse, error) {
	profiles, err := s.service.GetCompanyProfiles(ctx, data.GetIds())
	if err != nil {
		return nil, err
	}

	profilesRPC := make([]*companyRPC.Profile, 0, len(profiles))
	for i := range profiles {
		profilesRPC = append(profilesRPC, profileToProfileRPC(profiles[i]))
	}

	return &companyRPC.GetCompanyProfilesResponse{
		Profiles: profilesRPC,
	}, nil
}

// ChangeCompanyAboutUs ...
func (s Server) ChangeCompanyAboutUs(ctx context.Context, data *companyRPC.ChangeCompanyAboutUsRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyAboutUs(ctx, data.GetId(), aboutUsRPCToProfileAboutUs(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetFounders ...
func (s Server) GetFounders(ctx context.Context, data *companyRPC.GetFoundersRequest) (*companyRPC.Founders, error) {

	founders, err := s.service.GetFounders(ctx, data.GetCompanyID(), data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	foundersRPC := make([]*companyRPC.Founder, 0, len(founders))

	for i := range founders {
		foundersRPC = append(foundersRPC, profileFounderToFounderRPC(founders[i]))
	}

	return &companyRPC.Founders{
		Founders: foundersRPC,
	}, nil
}

// AddCompanyFounder ...
func (s Server) AddCompanyFounder(ctx context.Context, data *companyRPC.AddCompanyFounderRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyFounder(ctx, data.GetID(), founderRPCToProfileFounder(data.GetFounder()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyFounder ...
func (s Server) DeleteCompanyFounder(ctx context.Context, data *companyRPC.DeleteCompanyFounderRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyFounder(ctx, data.GetId(), data.GetFounderId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyFounder ...
func (s Server) ChangeCompanyFounder(ctx context.Context, data *companyRPC.ChangeCompanyFounderRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyFounder(ctx, data.GetId(), founderRPCToProfileFounder(data.GetFounder()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyFounderAvatar ...
func (s Server) ChangeCompanyFounderAvatar(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	id, err := s.service.ChangeCompanyFounderAvatar(ctx, data.GetCompanyID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}
	return &companyRPC.ID{
		ID: id,
	}, nil
}

// RemoveCompanyFounderAvatar ...
func (s Server) RemoveCompanyFounderAvatar(ctx context.Context, data *companyRPC.File) (*companyRPC.Empty, error) {
	err := s.service.RemoveCompanyFounderAvatar(ctx, data.GetCompanyID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ApproveFounderRequest ...
func (s Server) ApproveFounderRequest(ctx context.Context, data *companyRPC.ApproveFounderRequestRequest) (*companyRPC.Empty, error) {
	err := s.service.ApproveFounderRequest(ctx, data.GetCompanyID(), data.GetRequestID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// RemoveFounderRequest ...
func (s Server) RemoveFounderRequest(ctx context.Context, data *companyRPC.RemoveFounderRequestRequest) (*companyRPC.Empty, error) {
	err := s.service.RemoveFounderRequest(ctx, data.GetCompanyID(), data.GetRequestID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddCompanyAward ...
func (s Server) AddCompanyAward(ctx context.Context, data *companyRPC.AddCompanyAwardRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyAward(ctx, data.GetID(), awardRPCToProfileAward(data.GetAward()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyAward ...
func (s Server) DeleteCompanyAward(ctx context.Context, data *companyRPC.DeleteCompanyAwardRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyAward(ctx, data.GetId(), data.GetAwardId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyAward ...
func (s Server) ChangeCompanyAward(ctx context.Context, data *companyRPC.ChangeCompanyAwardRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyAward(ctx, data.GetId(), awardRPCToProfileAward(data.GetAward()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddLinksInCompanyAward ...
func (s Server) AddLinksInCompanyAward(ctx context.Context, data *companyRPC.AddLinksRequest) (*companyRPC.Empty, error) {
	urls := make([]*profile.Link, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		urls = append(urls, linkRPCToProfileLink(data.GetLinks()[i]))
	}

	_ /*ids*/, err := s.service.AddLinksInCompanyAward(ctx, data.GetID(), data.GetAwardID(), urls)
	if err != nil {
		return nil, err
	}

	// TODO: return ids

	return &companyRPC.Empty{}, nil
}

// AddFileInCompanyAward ...
func (s Server) AddFileInCompanyAward(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	id, err := s.service.AddFileInCompanyAward(ctx, data.GetCompanyID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilesInCompanyAward ...
func (s Server) RemoveFilesInCompanyAward(ctx context.Context, data *companyRPC.RemoveFilesRequest) (*companyRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInCompanyAward(ctx, data.GetID(), data.GetAwardID(), ids)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeLinkInCompanyAward ...
func (s Server) ChangeLinkInCompanyAward(ctx context.Context, data *companyRPC.ChangeLinkRequest) (*companyRPC.Empty, error) {
	// ids := make([]string, 0, len(data.GetLinks()))
	// for i := range data.GetLinks() {
	// 	ids = append(ids, data.GetLinks()[i].GetID())
	// }

	err := s.service.ChangeLinkInCompanyAward(ctx, data.GetID(), data.GetAwardID(), data.GetLink().GetID(), data.GetLink().GetURL())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// RemoveLinksInCompanyAward ...
func (s Server) RemoveLinksInCompanyAward(ctx context.Context, data *companyRPC.RemoveLinksRequest) (*companyRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		ids = append(ids, data.GetLinks()[i].GetID())
	}

	err := s.service.RemoveLinksInCompanyAward(ctx, data.GetID(), data.GetAwardID(), ids)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetUploadedFilesInCompanyAward ...
func (s Server) GetUploadedFilesInCompanyAward(ctx context.Context, data *companyRPC.ID) (*companyRPC.Files, error) {
	files, err := s.service.GetUploadedFilesInCompanyAward(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	fileRPC := make([]*companyRPC.File, 0, len(files))
	for i := range files {
		fileRPC = append(fileRPC, profileFiletoFileRPC(files[i]))
	}

	return &companyRPC.Files{
		Files: fileRPC,
	}, nil
}

// GetCompanyGallery ...
func (s Server) GetCompanyGallery(ctx context.Context, data *companyRPC.RequestGallery) (*companyRPC.GalleryFiles, error) {
	files, err := s.service.GetCompanyGallery(ctx, data.GetCompanyID(), data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	galleriesRPC := make([]*companyRPC.GalleryFile, 0, len(files))

	for i := range files {
		galleriesRPC = append(galleriesRPC, profileFileToCompanyRPCGalleryFile(files[i]))
	}

	return &companyRPC.GalleryFiles{
		GalleryFiles: galleriesRPC,
	}, nil
}

// AddFileInCompanyGallery ...
func (s Server) AddFileInCompanyGallery(ctx context.Context, data *companyRPC.GalleryFile) (*companyRPC.ID, error) {
	id, err := s.service.AddFileInCompanyGallery(ctx, data.GetCompanyID(), galleryFileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilesInCompanyGallery ...
func (s Server) RemoveFilesInCompanyGallery(ctx context.Context, data *companyRPC.RemoveGalleryFileRequest) (*companyRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInCompanyGallery(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetUploadedFilesInCompanyGallery ...
func (s Server) GetUploadedFilesInCompanyGallery(ctx context.Context, data *companyRPC.ID) (*companyRPC.GalleryFiles, error) {
	files, err := s.service.GetUploadedFilesInCompanyGallery(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	fileRPC := make([]*companyRPC.GalleryFile, 0, len(files))
	for i := range files {
		fileRPC = append(fileRPC, profileGalleryFiletoFileRPC(files[i]))
	}

	return &companyRPC.GalleryFiles{
		GalleryFiles: fileRPC,
	}, nil
}

// AddCompanyMilestone ...
func (s Server) AddCompanyMilestone(ctx context.Context, data *companyRPC.AddCompanyMilestoneRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyMilestone(ctx, data.GetID(), milestoneRPCToProfileMilestone(data.GetMilestone()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyMilestone ...
func (s Server) DeleteCompanyMilestone(ctx context.Context, data *companyRPC.DeleteCompanyMilestoneRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyMilestone(ctx, data.GetId(), data.GetMilestoneId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyMilestone ...
func (s Server) ChangeCompanyMilestone(ctx context.Context, data *companyRPC.ChangeCompanyMilestoneRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyMilestone(ctx, data.GetId(), milestoneRPCToProfileMilestone(data.GetMilestone()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeImageMilestone ...
func (s *Server) ChangeImageMilestone(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	id, err := s.service.ChangeImageMilestone(ctx, data.GetCompanyID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}
	return &companyRPC.ID{
		ID: id,
	}, nil
}

// RemoveImageInMilestone ...
func (s *Server) RemoveImageInMilestone(ctx context.Context, data *companyRPC.RemoveImageInMilestoneRequest) (*companyRPC.Empty, error) {
	err := s.service.RemoveImageInMilestone(ctx, data.GetCompanyID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddCompanyProduct ...
func (s Server) AddCompanyProduct(ctx context.Context, data *companyRPC.AddCompanyProductRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyProduct(ctx, data.GetId(), productRPCToProfileProduct(data.GetProduct()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyProduct ...
func (s Server) DeleteCompanyProduct(ctx context.Context, data *companyRPC.DeleteCompanyProductRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyProduct(ctx, data.GetId(), data.GetProductId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyProduct ...
func (s Server) ChangeCompanyProduct(ctx context.Context, data *companyRPC.ChangeCompanyProductRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyProduct(ctx, data.GetId(), productRPCToProfileProduct(data.GetProduct()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeImageProduct ...
func (s *Server) ChangeImageProduct(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	id, err := s.service.ChangeImageProduct(ctx, data.GetCompanyID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}
	return &companyRPC.ID{
		ID: id,
	}, nil
}

// RemoveImageInProduct ...
func (s *Server) RemoveImageInProduct(ctx context.Context, data *companyRPC.RemoveImageInProductRequest) (*companyRPC.Empty, error) {
	err := s.service.RemoveImageInProduct(ctx, data.GetCompanyID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddCompanyService ...
func (s Server) AddCompanyService(ctx context.Context, data *companyRPC.AddCompanyServiceRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyService(ctx, data.GetID(), serviceRPCToProfileService(data.GetService()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyService ...
func (s Server) DeleteCompanyService(ctx context.Context, data *companyRPC.DeleteCompanyServiceRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyService(ctx, data.GetId(), data.GetServiceId())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeCompanyService ...
func (s Server) ChangeCompanyService(ctx context.Context, data *companyRPC.ChangeCompanyServiceRequest) (*companyRPC.Empty, error) {
	err := s.service.ChangeCompanyService(ctx, data.GetId(), serviceRPCToProfileService(data.GetService()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// ChangeImageService ...
func (s *Server) ChangeImageService(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	id, err := s.service.ChangeImageService(ctx, data.GetCompanyID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}
	return &companyRPC.ID{
		ID: id,
	}, nil
}

// RemoveImageInService ...
func (s *Server) RemoveImageInService(ctx context.Context, data *companyRPC.RemoveImageInServiceRequest) (*companyRPC.Empty, error) {
	err := s.service.RemoveImageInService(ctx, data.GetCompanyID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddCompanyReport ...
func (s Server) AddCompanyReport(ctx context.Context, data *companyRPC.AddCompanyReportRequest) (*companyRPC.Empty, error) {
	err := s.service.AddCompanyReport(ctx, data.GetID(), reportRPCToProfileReport(data.GetReport()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddCompanyReview ...
func (s Server) AddCompanyReview(ctx context.Context, data *companyRPC.AddCompanyReviewRequest) (*companyRPC.ID, error) {
	id, err := s.service.AddCompanyReview(ctx, data.GetId(), reviewRPCToProfileReview(data.GetReview()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{
		ID: id,
	}, nil
}

// DeleteCompanyReview ...
func (s Server) DeleteCompanyReview(ctx context.Context, data *companyRPC.DeleteCompanyReviewRequest) (*companyRPC.Empty, error) {
	err := s.service.DeleteCompanyReview(ctx, data.GetID(), data.GetReviewID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetCompanyReviews ...
func (s Server) GetCompanyReviews(ctx context.Context, data *companyRPC.GetCompanyReviewsRequest) (*companyRPC.GetCompanyReviewsResponse, error) {
	reviews, err := s.service.GetCompanyReviews(ctx, data.GetID(), data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	reviewRPC := make([]*companyRPC.Review, 0, len(reviews))
	for i := range reviews {
		reviewRPC = append(reviewRPC, profileReviewToReviewRPC(reviews[i]))
	}

	return &companyRPC.GetCompanyReviewsResponse{
		Reviews: reviewRPC,
	}, nil
}

// GetUsersReviews ...
func (s Server) GetUsersReviews(ctx context.Context, data *companyRPC.GetCompanyReviewsRequest) (*companyRPC.GetCompanyReviewsResponse, error) {
	reviews, err := s.service.GetUsersRevies(ctx, data.GetID(), data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	reviewRPC := make([]*companyRPC.Review, 0, len(reviews))
	for i := range reviews {
		reviewRPC = append(reviewRPC, profileReviewToReviewRPC(reviews[i]))
	}

	return &companyRPC.GetCompanyReviewsResponse{
		Reviews: reviewRPC,
	}, nil
}

// AddCompanyReviewReport ...
func (s Server) AddCompanyReviewReport(ctx context.Context, data *companyRPC.AddCompanyReviewReportRequest) (*companyRPC.Empty, error) {
	err := s.service.AddCompanyReviewReport(ctx, reviewReportRPCToProfileReviewReport(data.GetReviewReport()))
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetAvarageRateOfCompany ...
func (s Server) GetAvarageRateOfCompany(ctx context.Context, data *companyRPC.ID) (*companyRPC.Rate, error) {
	avg, amount, err := s.service.GetAvarageRateOfCompany(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Rate{
		AmountReviews: amount,
		AvarageRate:   avg,
	}, nil
}

// GetAmountOfEachRate ...
func (s Server) GetAmountOfEachRate(ctx context.Context, data *companyRPC.ID) (*companyRPC.AmountOfRates, error) {
	rates, err := s.service.GetAmountOfEachRate(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.AmountOfRates{
		Rate: rates,
	}, nil
}

// ChangeAvatar ...
func (s Server) ChangeAvatar(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	err := s.service.ChangeAvatar(ctx, data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{}, nil
}

// ChangeOriginAvatar ...
func (s Server) ChangeOriginAvatar(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	err := s.service.ChangeOriginAvatar(ctx, data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{}, nil
}

// RemoveAvatar ...
func (s Server) RemoveAvatar(ctx context.Context, data *companyRPC.ID) (*companyRPC.Empty, error) {
	err := s.service.RemoveAvatar(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetOriginAvatar ...
func (s Server) GetOriginAvatar(ctx context.Context, data *companyRPC.ID) (*companyRPC.File, error) {
	url, err := s.service.GetOriginAvatar(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.File{
		URL: url,
	}, nil
}

// ChangeCover ...
func (s Server) ChangeCover(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	err := s.service.ChangeCover(ctx, data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{}, nil
}

// ChangeOriginCover ...
func (s Server) ChangeOriginCover(ctx context.Context, data *companyRPC.File) (*companyRPC.ID, error) {
	err := s.service.ChangeOriginCover(ctx, data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &companyRPC.ID{}, nil
}

// RemoveCover ...
func (s Server) RemoveCover(ctx context.Context, data *companyRPC.ID) (*companyRPC.Empty, error) {
	err := s.service.RemoveCover(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetOriginCover ...
func (s Server) GetOriginCover(ctx context.Context, data *companyRPC.ID) (*companyRPC.File, error) {
	url, err := s.service.GetOriginCover(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &companyRPC.File{
		URL: url,
	}, nil
}

// SaveCompanyProfileTranslation ...
func (s Server) SaveCompanyProfileTranslation(ctx context.Context, data *companyRPC.ProfileTranslation) (*companyRPC.Empty, error) {
	err := s.service.SaveCompanyProfileTranslation(
		ctx,
		data.GetCompanyID(),
		data.GetLanguage(),
		&profile.Translation{
			Name:        data.GetName(),
			Mission:     data.GetMission(),
			Description: data.GetDescription(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// GetAmountOfReviewsOfUser ...
func (s Server) GetAmountOfReviewsOfUser(ctx context.Context, data *companyRPC.ID) (*companyRPC.Amount, error) {
	amount, err := s.service.GetAmountOfReviewsOfUser(
		ctx,
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Amount{
		Amount: amount,
	}, nil
}

// SaveCompanyMilestoneTranslation ...
func (s Server) SaveCompanyMilestoneTranslation(ctx context.Context, data *companyRPC.MilestoneTranslation) (*companyRPC.Empty, error) {
	err := s.service.SaveCompanyMilestoneTranslation(
		ctx,
		data.GetCompanyID(),
		data.GetMilestoneID(),
		data.GetLanguage(),
		&profile.MilestoneTranslation{
			Description: data.GetDesciption(),
			Title:       data.GetTitle(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// SaveCompanyAwardTranslation ...
func (s Server) SaveCompanyAwardTranslation(ctx context.Context, data *companyRPC.AwardTranslation) (*companyRPC.Empty, error) {
	err := s.service.SaveCompanyAwardTranslation(
		ctx,
		data.GetCompanyID(),
		data.GetAwardID(),
		data.GetLanguage(),
		&profile.AwardTranslation{
			Title:  data.GetTitle(),
			Issuer: data.GetIssuer(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// AddGoldCoinsToWallet ...
func (s Server) AddGoldCoinsToWallet(ctx context.Context, data *companyRPC.WalletAddGoldCoins) (*companyRPC.Empty, error) {
	err := s.service.AddGoldCoinsToWallet(
		ctx,
		data.GetUserID(),
		data.GetCoins(),
	)

	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// OpenCareerCenter ...
func (s Server) OpenCareerCenter(ctx context.Context, data *companyRPC.OpenCareerCenterRequest) (*companyRPC.Empty, error) {
	err := s.service.OpenCareerCenter(
		ctx,
		data.GetCompanyID(),
		&careercenter.CareerCenter{
			CVButtonEnabled:     data.GetCVButtonEnabled(),
			CustomButtonEnabled: data.GetCustomButtonEnabled(),
			CustomButtonTitle:   data.GetCustomButtontitle(),
			CustomButtonURL:     data.GetCustomButtonURL(),
			Description:         data.GetDescription(),
			Title:               data.GetTitle(),
		},
	)

	if err != nil {
		return nil, err
	}

	return &companyRPC.Empty{}, nil
}

// From RPC to struct

func registerCompanyRequestRPCToCompanyAccount(data *companyRPC.RegisterCompanyRequest) *account.Account {
	if data == nil {
		return nil
	}

	acc := account.Account{
		Name:           data.GetName(),
		URL:            data.GetURL(),
		Website:        make([]*account.Website, 0, len(data.GetWebsites())),
		Emails:         make([]*account.Email, 1),
		Phones:         make([]*account.Phone, 1),
		FoundationDate: stringDateToTime(data.GetFoundationDate()),
		Type:           companyTypeRPCToAccountType(data.GetType()),
		Addresses:      make([]*account.Address, 1),
	}

	if data.GetInvitedBy() != "" {
		acc.SetInvitedByID(data.GetInvitedBy())
	}

	acc.Industry = account.Industry{
		Main: data.GetIndustry().GetMain(),
		Sub:  make([]string, 0, len(data.GetIndustry().GetSubs())),
	}

	for i := range data.GetIndustry().GetSubs() {
		acc.Industry.Sub = append(acc.Industry.Sub, data.GetIndustry().GetSubs()[i])
	}

	acc.Emails[0] = &account.Email{
		Email: data.GetEmail(),
	}

	acc.Phones[0] = &account.Phone{
		Number: data.GetPhone().GetNumber(),
	}
	if cc := countryCodeRPCToAccountCountryCode(data.GetPhone().GetCountryCode()); cc != nil {
		acc.Phones[0].CountryCode = *cc
	}

	for i := range data.GetWebsites() {
		w := &account.Website{
			Site: data.Websites[i],
		}
		w.GenerateID()

		acc.Website = append(acc.Website, w)
	}

	acc.Addresses[0] = &account.Address{
		Street:    data.GetStreetAddress(),
		Apartment: data.GetApartment(),
		ZIPCode:   data.GetZipCode(),
		Location: location.Location{
			City: &location.City{
				ID: data.GetCityId(),
			},
			Country: &location.Country{},
		},
	}

	if data.GetVAT() != "" {
		vat := data.GetVAT()
		acc.VAT = &vat
	}

	return &acc
}

func companyTypeRPCToAccountType(data companyRPC.Type) account.Type {
	switch data {
	case companyRPC.Type_TYPE_PARTNERSHIP:
		return account.TypePartnership
	case companyRPC.Type_TYPE_SELF_EMPLOYED:
		return account.TypeSelfEmployed
	case companyRPC.Type_TYPE_PRIVATELY_HELD:
		return account.TypePrivatelyHeld
	case companyRPC.Type_TYPE_PUBLIC_COMPANY:
		return account.TypePublicCompany
	case companyRPC.Type_TYPE_GOVERNMENT_AGENCY:
		return account.TypeGovernmentAgency
	case companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP:
		return account.TypeSoleProprietorship
	case companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION:
		return account.TypeEducationalInstitution
	}

	return account.TypeUnknown
}

func sizeRPCToAccountSize(data companyRPC.Size) account.Size {
	size := account.SizeUnknown

	switch data {
	case companyRPC.Size_SIZE_SELF_EMPLOYED:
		size = account.SizeSelfEmployed
	case companyRPC.Size_SIZE_1_10_EMPLOYEES:
		size = account.SizeFrom1Till10Employees
	case companyRPC.Size_SIZE_11_50_EMPLOYEES:
		size = account.SizeFrom11Till50Employees
	case companyRPC.Size_SIZE_51_200_EMPLOYEES:
		size = account.SizeFrom51Till200Employees
	case companyRPC.Size_SIZE_201_500_EMPLOYEES:
		size = account.SizeFrom201Till500Employees
	case companyRPC.Size_SIZE_501_1000_EMPLOYEES:
		size = account.SizeFrom501Till1000Employees
	case companyRPC.Size_SIZE_1001_5000_EMPLOYEES:
		size = account.SizeFrom1001Till5000Employees
	case companyRPC.Size_SIZE_5001_10000_EMPLOYEES:
		size = account.SizeFrom5001Till10000Employees
	case companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES:
		size = account.SizeFrom10001AndMoreEmployees
	}

	return size
}

func emailRPCToAccountEmail(data *companyRPC.Email) *account.Email {
	if data == nil {
		return nil
	}

	email := account.Email{
		Email:   data.GetEmail(),
		Primary: data.GetIsPrimary(),
		// Activated: data.GetIsActivated(),
	}

	_ = email.SetID(email.GetID())

	return &email
}

func phoneRPCToAccountPhone(data *companyRPC.Phone) *account.Phone {
	if data == nil {
		return nil
	}

	phone := account.Phone{
		Number:  data.GetNumber(),
		Primary: data.GetIsPrimary(),
		// Activated
	}

	_ = phone.SetID(data.GetID())

	if cc := countryCodeRPCToAccountCountryCode(data.GetCountryCode()); cc != nil {
		phone.CountryCode = *cc
	}

	return &phone
}

func countryCodeRPCToAccountCountryCode(data *companyRPC.CountryCode) *account.CountryCode {
	if data == nil {
		return nil
	}

	cc := account.CountryCode{
		ID:   uint32(data.GetId()),
		Code: data.GetCode(),
		// CountryID:
	}

	return &cc
}

func addressRPCToAccountAddress(data *companyRPC.Address) *account.Address {
	if data == nil {
		return nil
	}

	address := account.Address{
		Name:          data.GetName(),
		Apartment:     data.GetApartment(),
		Street:        data.GetStreetAddress(),
		ZIPCode:       data.GetZipCode(),
		IsPrimary:     data.GetIsPrimary(),
		BusinessHours: make([]*account.BusinessHour, 0, len(data.GetBusinessHours())),
		Phones:        make([]*account.Phone, 0, len(data.GetPhones())),
	}

	_ = address.SetID(data.GetID())

	for i := range data.GetBusinessHours() {
		address.BusinessHours = append(address.BusinessHours, businessHourItemRPCToAccountBusinessHour(data.GetBusinessHours()[i]))
	}

	if gp := geoPosRPCToAccountGeoPos(data.GetGeoPos()); gp != nil {
		address.GeoPos = *gp
	}

	if loc := locationRPCToLocation(data.GetLocation()); loc != nil {
		address.Location = *loc
	}

	for i := range data.GetPhones() {
		address.Phones = append(address.Phones, phoneRPCToAccountPhone(data.GetPhones()[i]))
	}

	return &address
}

func businessHourItemRPCToAccountBusinessHour(data *companyRPC.BusinessHoursItem) *account.BusinessHour {
	if data == nil {
		return nil
	}

	bh := account.BusinessHour{
		Weekdays:   data.GetWeekDays(), // TODO: change
		FinishHour: data.GetHourTo(),
		StartHour:  data.GetHourFrom(),
	}

	_ = bh.SetID(data.GetID())

	return &bh
}

func geoPosRPCToAccountGeoPos(data *companyRPC.GeoPos) *account.GeoPos {
	if data == nil {
		return nil
	}

	return &account.GeoPos{
		Lantitude: data.GetLantitude(),
		Longitude: data.GetLongitude(),
	}
}

func locationRPCToLocation(data *companyRPC.Location) *location.Location {
	if data == nil {
		return nil
	}

	return &location.Location{
		City:    cityRPCToLocationCity(data.GetCity()),
		Country: countryRPCToLocationCountry(data.GetCountry()),
	}
}

func cityRPCToLocationCity(data *companyRPC.City) *location.City {
	if data == nil {
		return nil
	}

	city := location.City{
		Name:        data.GetTitle(),
		Subdivision: data.GetSubdivision(),
	}

	id, _ := strconv.Atoi(data.GetId())
	city.ID = int32(id)

	return &city
}

func countryRPCToLocationCountry(data *companyRPC.Country) *location.Country {
	if data == nil {
		return nil
	}

	return &location.Country{
		ID: data.GetId(),
	}
}

func websiteRPCToAccountWebsite(data *companyRPC.Website) *account.Website {
	if data == nil {
		return nil
	}

	ws := account.Website{
		Site: data.GetWebsite(),
	}

	_ = ws.SetID(data.GetId())

	return &ws
}

func parkingRPCToAccountParking(data companyRPC.Parking) account.Parking {
	parking := account.ParkingUnknown

	switch data {
	case companyRPC.Parking_PARKING_NO_PARKING:
		parking = account.ParkingNoParking
	case companyRPC.Parking_PARKING_PARKING_LOT:
		parking = account.ParkingParkingLot
	case companyRPC.Parking_PARKING_STREET_PARKING:
		parking = account.ParkingStreetParking
	}

	return parking
}

func benefitsRPCToProfileBenefits(data []companyRPC.BenefitEnum) []profile.Benefit {
	benefits := make([]profile.Benefit, 0, len(data))

	for _, t := range data {
		switch t {
		case companyRPC.BenefitEnum_labor_agreement:
			benefits = append(benefits, profile.LaborAgreement)
		case companyRPC.BenefitEnum_remote_working:
			benefits = append(benefits, profile.RemoteWorking)
		case companyRPC.BenefitEnum_floater:
			benefits = append(benefits, profile.Floater)
		case companyRPC.BenefitEnum_paid_timeoff:
			benefits = append(benefits, profile.PaidTimeoff)
		case companyRPC.BenefitEnum_flexible_working_hours:
			benefits = append(benefits, profile.FlexibleWorkingHours)
		case companyRPC.BenefitEnum_additional_timeoff:
			benefits = append(benefits, profile.AdditionalTimeoff)
		case companyRPC.BenefitEnum_additional_parental_leave:
			benefits = append(benefits, profile.AdditionalParentalLeave)
		case companyRPC.BenefitEnum_sick_leave_for_family_members:
			benefits = append(benefits, profile.SickLeaveForFamilyMembers)
		case companyRPC.BenefitEnum_childcare:
			benefits = append(benefits, profile.Childcare)
		case companyRPC.BenefitEnum_company_canteen:
			benefits = append(benefits, profile.CompanyCanteen)
		case companyRPC.BenefitEnum_sport_facilities:
			benefits = append(benefits, profile.SportFacilities)
		case companyRPC.BenefitEnum_access_for_handicapped_persons:
			benefits = append(benefits, profile.AccessForHandicappedPersons)
		case companyRPC.BenefitEnum_employee_parking:
			benefits = append(benefits, profile.EmployeeParking)
		case companyRPC.BenefitEnum_transportation:
			benefits = append(benefits, profile.Transportation)
		case companyRPC.BenefitEnum_multiple_work_spaces:
			benefits = append(benefits, profile.MultipleWorkSpaces)
		case companyRPC.BenefitEnum_corporate_events:
			benefits = append(benefits, profile.CorporateEvents)
		case companyRPC.BenefitEnum_trainig_and_development:
			benefits = append(benefits, profile.TrainigAndDevelopment)
		case companyRPC.BenefitEnum_pets_allowed:
			benefits = append(benefits, profile.PetsAllowed)
		case companyRPC.BenefitEnum_corporate_medical_staff:
			benefits = append(benefits, profile.CorporateMedicalStaff)
		case companyRPC.BenefitEnum_game_consoles:
			benefits = append(benefits, profile.GameConsoles)
		case companyRPC.BenefitEnum_snack_and_drink_selfservice:
			benefits = append(benefits, profile.SnackAndDrinkSelfservice)
		case companyRPC.BenefitEnum_private_pension_scheme:
			benefits = append(benefits, profile.PrivatePensionScheme)
		case companyRPC.BenefitEnum_health_insurance:
			benefits = append(benefits, profile.HealthInsurance)
		case companyRPC.BenefitEnum_dental_care:
			benefits = append(benefits, profile.DentalCare)
		case companyRPC.BenefitEnum_car_insurance:
			benefits = append(benefits, profile.CarInsurance)
		case companyRPC.BenefitEnum_tution_fees:
			benefits = append(benefits, profile.TutionFees)
		case companyRPC.BenefitEnum_permfomance_related_bonus:
			benefits = append(benefits, profile.PermfomanceRelatedBonus)
		case companyRPC.BenefitEnum_stock_options:
			benefits = append(benefits, profile.StockOptions)
		case companyRPC.BenefitEnum_profit_earning_bonus:
			benefits = append(benefits, profile.ProfitEarningBonus)
		case companyRPC.BenefitEnum_additional_months_salary:
			benefits = append(benefits, profile.AdditionalMonthsSalary)
		case companyRPC.BenefitEnum_employers_matching_contributions:
			benefits = append(benefits, profile.EmployersMatchingContributions)
		case companyRPC.BenefitEnum_parental_bonus:
			benefits = append(benefits, profile.ParentalBonus)
		case companyRPC.BenefitEnum_tax_deductions:
			benefits = append(benefits, profile.TaxDeductions)
		case companyRPC.BenefitEnum_language_courses:
			benefits = append(benefits, profile.LanguageCourses)
		case companyRPC.BenefitEnum_company_car:
			benefits = append(benefits, profile.CompanyCar)
		case companyRPC.BenefitEnum_laptop:
			benefits = append(benefits, profile.Laptop)
		case companyRPC.BenefitEnum_discounts_on_company_products_and_services:
			benefits = append(benefits, profile.DiscountsOnCompanyProductsAndServices)
		case companyRPC.BenefitEnum_holiday_vouchers:
			benefits = append(benefits, profile.HolidayVouchers)
		case companyRPC.BenefitEnum_restraunt_vouchers:
			benefits = append(benefits, profile.RestrauntVouchers)
		case companyRPC.BenefitEnum_corporate_housing:
			benefits = append(benefits, profile.CorporateHousing)
		case companyRPC.BenefitEnum_mobile_phone:
			benefits = append(benefits, profile.MobilePhone)
		case companyRPC.BenefitEnum_gift_vouchers:
			benefits = append(benefits, profile.GiftVouchers)
		case companyRPC.BenefitEnum_cultural_or_sporting_activites:
			benefits = append(benefits, profile.CulturalOrSportingActivites)
		case companyRPC.BenefitEnum_employee_service_vouchers:
			benefits = append(benefits, profile.EmployeeServiceVouchers)
		case companyRPC.BenefitEnum_corporate_credit_card:
			benefits = append(benefits, profile.CorporateCreditCard)
		case companyRPC.BenefitEnum_other:
			benefits = append(benefits, profile.Other)
		case companyRPC.BenefitEnum_relocation_package:
			benefits = append(benefits, profile.RelocationPackage)
		}
	}
	return benefits
}

func profileBenetisToCompanyrRPCBenefits(data []profile.Benefit) []companyRPC.BenefitEnum {
	benefits := make([]companyRPC.BenefitEnum, 0, len(data))

	for _, t := range data {
		switch t {
		case profile.LaborAgreement:
			benefits = append(benefits, companyRPC.BenefitEnum_labor_agreement)
		case profile.RemoteWorking:
			benefits = append(benefits, companyRPC.BenefitEnum_remote_working)
		case profile.Floater:
			benefits = append(benefits, companyRPC.BenefitEnum_floater)
		case profile.PaidTimeoff:
			benefits = append(benefits, companyRPC.BenefitEnum_paid_timeoff)
		case profile.FlexibleWorkingHours:
			benefits = append(benefits, companyRPC.BenefitEnum_flexible_working_hours)
		case profile.AdditionalTimeoff:
			benefits = append(benefits, companyRPC.BenefitEnum_additional_timeoff)
		case profile.AdditionalParentalLeave:
			benefits = append(benefits, companyRPC.BenefitEnum_additional_parental_leave)
		case profile.SickLeaveForFamilyMembers:
			benefits = append(benefits, companyRPC.BenefitEnum_sick_leave_for_family_members)
		case profile.Childcare:
			benefits = append(benefits, companyRPC.BenefitEnum_childcare)
		case profile.CompanyCanteen:
			benefits = append(benefits, companyRPC.BenefitEnum_company_canteen)
		case profile.SportFacilities:
			benefits = append(benefits, companyRPC.BenefitEnum_sport_facilities)
		case profile.AccessForHandicappedPersons:
			benefits = append(benefits, companyRPC.BenefitEnum_access_for_handicapped_persons)
		case profile.EmployeeParking:
			benefits = append(benefits, companyRPC.BenefitEnum_employee_parking)
		case profile.Transportation:
			benefits = append(benefits, companyRPC.BenefitEnum_transportation)
		case profile.MultipleWorkSpaces:
			benefits = append(benefits, companyRPC.BenefitEnum_multiple_work_spaces)
		case profile.CorporateEvents:
			benefits = append(benefits, companyRPC.BenefitEnum_corporate_events)
		case profile.TrainigAndDevelopment:
			benefits = append(benefits, companyRPC.BenefitEnum_trainig_and_development)
		case profile.PetsAllowed:
			benefits = append(benefits, companyRPC.BenefitEnum_pets_allowed)
		case profile.CorporateMedicalStaff:
			benefits = append(benefits, companyRPC.BenefitEnum_corporate_medical_staff)
		case profile.GameConsoles:
			benefits = append(benefits, companyRPC.BenefitEnum_game_consoles)
		case profile.SnackAndDrinkSelfservice:
			benefits = append(benefits, companyRPC.BenefitEnum_snack_and_drink_selfservice)
		case profile.PrivatePensionScheme:
			benefits = append(benefits, companyRPC.BenefitEnum_private_pension_scheme)
		case profile.HealthInsurance:
			benefits = append(benefits, companyRPC.BenefitEnum_health_insurance)
		case profile.DentalCare:
			benefits = append(benefits, companyRPC.BenefitEnum_dental_care)
		case profile.CarInsurance:
			benefits = append(benefits, companyRPC.BenefitEnum_car_insurance)
		case profile.TutionFees:
			benefits = append(benefits, companyRPC.BenefitEnum_tution_fees)
		case profile.PermfomanceRelatedBonus:
			benefits = append(benefits, companyRPC.BenefitEnum_permfomance_related_bonus)
		case profile.StockOptions:
			benefits = append(benefits, companyRPC.BenefitEnum_stock_options)
		case profile.ProfitEarningBonus:
			benefits = append(benefits, companyRPC.BenefitEnum_profit_earning_bonus)
		case profile.AdditionalMonthsSalary:
			benefits = append(benefits, companyRPC.BenefitEnum_additional_months_salary)
		case profile.EmployersMatchingContributions:
			benefits = append(benefits, companyRPC.BenefitEnum_employers_matching_contributions)
		case profile.ParentalBonus:
			benefits = append(benefits, companyRPC.BenefitEnum_parental_bonus)
		case profile.TaxDeductions:
			benefits = append(benefits, companyRPC.BenefitEnum_tax_deductions)
		case profile.LanguageCourses:
			benefits = append(benefits, companyRPC.BenefitEnum_language_courses)
		case profile.CompanyCar:
			benefits = append(benefits, companyRPC.BenefitEnum_company_car)
		case profile.Laptop:
			benefits = append(benefits, companyRPC.BenefitEnum_laptop)
		case profile.DiscountsOnCompanyProductsAndServices:
			benefits = append(benefits, companyRPC.BenefitEnum_discounts_on_company_products_and_services)
		case profile.HolidayVouchers:
			benefits = append(benefits, companyRPC.BenefitEnum_holiday_vouchers)
		case profile.RestrauntVouchers:
			benefits = append(benefits, companyRPC.BenefitEnum_restraunt_vouchers)
		case profile.CorporateHousing:
			benefits = append(benefits, companyRPC.BenefitEnum_corporate_housing)
		case profile.MobilePhone:
			benefits = append(benefits, companyRPC.BenefitEnum_mobile_phone)
		case profile.GiftVouchers:
			benefits = append(benefits, companyRPC.BenefitEnum_gift_vouchers)
		case profile.CulturalOrSportingActivites:
			benefits = append(benefits, companyRPC.BenefitEnum_cultural_or_sporting_activites)
		case profile.EmployeeServiceVouchers:
			benefits = append(benefits, companyRPC.BenefitEnum_employee_service_vouchers)
		case profile.CorporateCreditCard:
			benefits = append(benefits, companyRPC.BenefitEnum_corporate_credit_card)
		case profile.Other:
			benefits = append(benefits, companyRPC.BenefitEnum_other)
		case profile.RelocationPackage:
			benefits = append(benefits, companyRPC.BenefitEnum_relocation_package)
		}

	}
	return benefits
}

func adminRoleRPCToAccountAdminLevel(data companyRPC.AdminRole) account.AdminLevel {
	// lvl := account.AdminRole_ROLE_UNKNOWN
	lvl := account.AdminLevelAdmin

	switch data {
	case companyRPC.AdminRole_ROLE_UNKNOWN:
		// lvl = account.Admin
	case companyRPC.AdminRole_ROLE_ADMIN:
		lvl = account.AdminLevelAdmin
	case companyRPC.AdminRole_ROLE_COMMERCIAL_ADMIN:
		lvl = account.AdminLevelCommercial
	case companyRPC.AdminRole_ROLE_JOB_EDITOR:
		lvl = account.AdminLevelJob
	case companyRPC.AdminRole_ROLE_V_SHOP_ADMIN:
		lvl = account.AdminLevelVShop
	case companyRPC.AdminRole_ROLE_V_SERVICE_ADMIN:
		lvl = account.AdminLevelVService
	}

	return lvl
}

func awardRPCToProfileAward(data *companyRPC.Award) *profile.Award {
	if data == nil {
		return nil
	}

	award := profile.Award{
		Title:  data.GetTitle(),
		Issuer: data.GetIssuer(),
		Date:   stringYearToDate(strconv.Itoa(int(data.GetYear()))),
		Links:  make([]*profile.Link, 0, len(data.GetLinks())),
		Files:  make([]*profile.File, 0, len(data.GetFiles())),
	}

	_ = award.SetID(data.GetID())

	for i := range data.GetFiles() {
		award.Files = append(award.Files, fileRPCToProfileFile(data.Files[i]))
	}

	for i := range data.GetLinks() {
		award.Links = append(award.Links, linkRPCToProfileLink(data.Links[i]))
	}

	return &award
}

func milestoneRPCToProfileMilestone(data *companyRPC.Milestone) *profile.Milestone {
	if data == nil {
		return nil
	}

	ms := profile.Milestone{
		Title:       data.GetTitle(),
		Image:       data.GetImage(),
		Description: data.GetDescription(),
		Date:        stringYearToDate(strconv.Itoa(int(data.GetYear()))),
	}

	_ = ms.SetID(data.GetId())

	return &ms
}

func founderRPCToProfileFounder(data *companyRPC.Founder) *profile.Founder {
	if data == nil {
		return nil
	}

	founder := profile.Founder{
		Name:     data.GetName(),
		Position: data.GetPositionTitle(),
		Avatar:   data.GetAvatar(),
	}

	_ = founder.SetID(data.GetID())
	_ = founder.SetUserID(data.GetUserID())

	return &founder
}

func productRPCToProfileProduct(data *companyRPC.Product) *profile.Product {
	if data == nil {
		return nil
	}

	product := profile.Product{
		Image:   data.GetImage(),
		Name:    data.GetName(),
		Website: data.GetWebsite(),
	}

	_ = product.SetID(data.GetID())

	return &product
}

func serviceRPCToProfileService(data *companyRPC.Service) *profile.Service {
	if data == nil {
		return nil
	}

	service := profile.Service{
		Name:    data.GetName(),
		Website: data.GetWebsite(),
		Image:   data.GetImage(),
	}

	_ = service.SetID(data.GetID())

	return &service
}

func reportRPCToProfileReport(data *companyRPC.Report) *profile.Report {
	if data == nil {
		return nil
	}

	report := profile.Report{
		Description: data.GetExplanation(),
		Reason:      reportEnumRPCToProfileReasonType(data.GetReport()),
	}

	return &report
}

func reportEnumRPCToProfileReasonType(data companyRPC.ReportEnum) profile.ReasonType {
	reason := profile.ReasonOther

	switch data {
	case companyRPC.ReportEnum_REPORT_DUPLICATE:
		reason = profile.ReasonDuplicate
	case companyRPC.ReportEnum_REPORT_PICTURE_IS_NOT_LOGO:
		reason = profile.ReasonPictureNotLogo
	case companyRPC.ReportEnum_REPORT_MAY_HAVE_BEEN_HACKED:
		reason = profile.ReasonMayBeHacked
	case companyRPC.ReportEnum_REPORT_NOT_REAL_ORGANIZATION:
		reason = profile.ReasonNotRealOrganization
	case companyRPC.ReportEnum_REPORT_VIOLATES_TERMS_OF_USE:
		reason = profile.ReasonViolationTermOfUse
	}

	return reason
}

func reviewRPCToProfileReview(data *companyRPC.Review) *profile.Review {
	if data == nil {
		return nil
	}

	review := profile.Review{
		Rate:        uint8(data.Rate),
		Headline:    data.GetHeadline(),
		Description: data.GetDescription(),
	}

	_ = review.SetID(data.GetId())

	return &review
}

func reviewReportRPCToProfileReviewReport(data *companyRPC.ReviewReport) *profile.ReviewReport {
	if data == nil {
		return nil
	}

	rr := profile.ReviewReport{
		Explanation: data.GetExplanation(),
		Reason:      reviewReportEnumRPCToProfileReviewReportReason(data.GetReport()),
	}

	_ = rr.SetReviewID(data.GetReviewId())
	_ = rr.SetCompanyID(data.GetCompanyID())

	return &rr
}

func reviewReportEnumRPCToProfileReviewReportReason(data companyRPC.ReviewReportEnum) profile.ReviewReportReason {
	reason := profile.ReviewReportReasonOther

	switch data {
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_SCAM:
		reason = profile.ReviewReportReasonScam
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_SPAM:
		reason = profile.ReviewReportReasonSpam
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_UNKNOWN:
		reason = profile.ReviewReportReasonUnknown
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_OFF_TOPIC:
		reason = profile.ReviewReportReasonOffTopic
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_FALSE_FAKE:
		reason = profile.ReviewReportReasonFake
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_SOMTHING_ELSE:
		reason = profile.ReviewReportReasonOther
	case companyRPC.ReviewReportEnum_REVIEW_REPORT_INAPPROPRIATE_OFFENSIVE:
		reason = profile.ReviewReportReasonOffensive
	}

	return reason
}

func industryRPCToAccountIndustry(data *companyRPC.Industry) *account.Industry {
	if data == nil {
		return nil
	}

	industry := account.Industry{
		Main: data.GetMain(),
		Sub:  data.GetSubs(),
	}

	return &industry
}

func aboutUsRPCToProfileAboutUs(data *companyRPC.ChangeCompanyAboutUsRequest) *profile.AboutUs {
	if data == nil {
		return nil
	}
	aboutUs := profile.AboutUs{
		Description:       data.GetDescription(),
		FoundationDate:    stringDateToTime(data.GetFoundationDate()),
		Type:              companyTypeRPCToAccountType(data.GetType()),
		Mission:           data.GetMission(),
		Parking:           parkingRPCToAccountParking(data.GetParking()),
		Size:              sizeRPCToAccountSize(data.GetSize()),
		BusinessHours:     make([]*account.BusinessHour, 0, len(data.GetBusinessHours())),
		IsDescriptionNull: data.GetIsDescriptionNull(),
		IsMissionNull:     data.GetIsMissionNull(),
		IsTypeNull:        data.GetIsTypeNull(),
		IsSizeNull:        data.GetIsSizeNull(),
		IsParkingNull:     data.GetIsParkingNull(),
	}

	if data.GetIndustry() != nil {
		aboutUs.Industry = *industryRPCToAccountIndustry(data.GetIndustry())
	}

	for i := range data.GetBusinessHours() {
		aboutUs.BusinessHours = append(aboutUs.BusinessHours, businessHourItemRPCToAccountBusinessHour(data.GetBusinessHours()[i]))
	}

	return &aboutUs
}

func fileRPCToProfileFile(data *companyRPC.File) *profile.File {
	if data == nil {
		return nil
	}

	f := profile.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = f.SetID(data.GetID())

	return &f
}

func galleryFileRPCToProfileFile(data *companyRPC.GalleryFile) *profile.File {
	if data == nil {
		return nil
	}

	f := profile.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = f.SetID(data.GetID())

	return &f
}

func profileFileToCompanyRPCGalleryFile(data *profile.File) *companyRPC.GalleryFile {
	if data == nil {
		return nil
	}

	f := companyRPC.GalleryFile{
		ID:       data.GetID(),
		MimeType: data.MimeType,
		Name:     data.Name,
		URL:      data.URL,
	}

	return &f
}

// To RPC

// TODO:
func accountToAccountRPC(data *account.Account) *companyRPC.Account {
	if data == nil {
		return nil
	}
	acc := companyRPC.Account{
		ID:        data.GetID(),
		Name:      data.Name,
		URL:       data.URL,
		OwnerID:   data.GetOwnerID(),
		Status:    statusToStatusRPC(data.Status),
		Addresses: make([]*companyRPC.Address, 0, len(data.Addresses)),
		Emails:    make([]*companyRPC.Email, 0, len(data.Emails)),
		Phones:    make([]*companyRPC.Phone, 0, len(data.Phones)),
		Type:      accountTypeToTypeRPC(data.Type),
		Size:      accountSizeToSizeRPC(data.Size),
		Industry: &companyRPC.Industry{
			Main: data.Industry.Main,
			Subs: data.Industry.Sub,
		},
		Websites: make([]*companyRPC.Website, 0, len(data.Website)),
		// TODO:
		// CreatedAt
		FoundationDate: timeToStringDayMonthAndYear(data.FoundationDate),
		BusinessHours:  make([]*companyRPC.BusinessHoursItem, 0, len(data.BusinessHours)),
	}

	if data.Parking != nil {
		acc.Parking = accountParkingToParkingRPC(*data.Parking)
	}

	for i := range data.Addresses {
		acc.Addresses = append(acc.Addresses, accountAddressToAddressRPC(data.Addresses[i]))
	}

	for i := range data.Emails {
		acc.Emails = append(acc.Emails, accountEmailToEmailRPC(data.Emails[i]))
	}

	for i := range data.Phones {
		acc.Phones = append(acc.Phones, accountPhoneToPhoneRPC(data.Phones[i]))
	}

	for i := range data.Website {
		acc.Websites = append(acc.Websites, accountWebsiteToWebsiteRPC(data.Website[i]))
	}

	for i := range data.BusinessHours {
		acc.BusinessHours = append(acc.BusinessHours, accountBusinessHourToBusinessHourRPC(data.BusinessHours[i]))
	}

	return &acc
}

func statusToStatusRPC(data status.CompanyStatus) companyRPC.Status {
	var st companyRPC.Status

	switch data {
	case status.CompanyStatusActivated:
		st = companyRPC.Status_STATUS_ACTIVATED
	case status.CompanyStatusDeactivated:
		st = companyRPC.Status_STATUS_DEACTIVATED
	case status.CompanyStatusNotActivated:
		st = companyRPC.Status_STATUS_NOT_ACTIVATED
	}

	return st
}

func accountAddressToAddressRPC(data *account.Address) *companyRPC.Address {
	if data == nil {
		return nil
	}

	a := companyRPC.Address{
		ID:            data.GetID(),
		Name:          data.Name,
		Apartment:     data.Apartment,
		StreetAddress: data.Street,
		ZipCode:       data.ZIPCode,
		BusinessHours: make([]*companyRPC.BusinessHoursItem, 0, len(data.BusinessHours)),
		Location:      locationToRPC(&data.Location),
		GeoPos:        accountGeoPostToGeoPosRPC(&data.GeoPos),
		IsPrimary:     data.IsPrimary,
		Phones:        make([]*companyRPC.Phone, 0, len(data.Phones)),
	}

	for i := range data.BusinessHours {
		a.BusinessHours = append(a.BusinessHours, accountBusinessHourToBusinessHourRPC(data.BusinessHours[i]))
	}

	for i := range data.Phones {
		a.Phones = append(a.Phones, accountPhoneToPhoneRPC(data.Phones[i]))
	}

	return &a
}

func accountBusinessHourToBusinessHourRPC(data *account.BusinessHour) *companyRPC.BusinessHoursItem {
	if data == nil {
		return nil
	}

	return &companyRPC.BusinessHoursItem{
		ID:       data.GetID(),
		WeekDays: data.Weekdays,
		HourFrom: data.StartHour,
		HourTo:   data.FinishHour,
	}
}

func accountGeoPostToGeoPosRPC(data *account.GeoPos) *companyRPC.GeoPos {
	if data == nil {
		return nil
	}

	g := companyRPC.GeoPos{
		Lantitude: data.Lantitude,
		Longitude: data.Longitude,
	}

	return &g
}

func locationCityToRPC(data *location.City) *companyRPC.City {
	if data == nil {
		return nil
	}

	return &companyRPC.City{
		Id:          strconv.Itoa(int(data.ID)),
		Subdivision: data.Subdivision,
		Title:       data.Name,
	}
}

func locationCountryToRPC(data *location.Country) *companyRPC.Country {
	if data == nil {
		return nil
	}

	return &companyRPC.Country{
		Id: data.ID,
		// : data.Name,
	}
}

func locationToRPC(data *location.Location) *companyRPC.Location {
	if data == nil {
		return nil
	}

	return &companyRPC.Location{
		City:    locationCityToRPC(data.City),
		Country: locationCountryToRPC(data.Country),
		// Permission
	}
}

func accountPhoneToPhoneRPC(data *account.Phone) *companyRPC.Phone {
	if data == nil {
		return nil
	}

	return &companyRPC.Phone{
		CountryAbbreviation: data.CountryCode.CountryID,
		ID:                  data.GetID(),
		IsActivated:         data.Activated,
		IsPrimary:           data.Primary,
		Number:              data.Number,
		CountryCode:         accountCountryCodeToRPC(&data.CountryCode),
	}
}

func accountCountryCodeToRPC(data *account.CountryCode) *companyRPC.CountryCode {
	if data == nil {
		return nil
	}

	return &companyRPC.CountryCode{
		Code: data.Code,
		Id:   int32(data.ID),
	}
}

func accountEmailToEmailRPC(data *account.Email) *companyRPC.Email {
	if data == nil {
		return nil
	}

	return &companyRPC.Email{
		ID:          data.GetID(),
		Email:       data.Email,
		IsActivated: data.Activated,
		IsPrimary:   data.Primary,
	}
}

func accountTypeToTypeRPC(data account.Type) companyRPC.Type {
	companyType := companyRPC.Type_TYPE_UNKNOWN

	switch data {
	case account.TypePartnership:
		companyType = companyRPC.Type_TYPE_PARTNERSHIP
	case account.TypeSelfEmployed:
		companyType = companyRPC.Type_TYPE_SELF_EMPLOYED
	case account.TypePrivatelyHeld:
		companyType = companyRPC.Type_TYPE_PRIVATELY_HELD
	case account.TypePublicCompany:
		companyType = companyRPC.Type_TYPE_PUBLIC_COMPANY
	case account.TypeGovernmentAgency:
		companyType = companyRPC.Type_TYPE_GOVERNMENT_AGENCY
	case account.TypeSoleProprietorship:
		companyType = companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP
	case account.TypeEducationalInstitution:
		companyType = companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION
	}

	return companyType
}

func accountSizeToSizeRPC(data account.Size) companyRPC.Size {
	size := companyRPC.Size_SIZE_UNKNOWN

	switch data {
	case account.SizeSelfEmployed:
		size = companyRPC.Size_SIZE_SELF_EMPLOYED
	case account.SizeFrom1Till10Employees:
		size = companyRPC.Size_SIZE_1_10_EMPLOYEES
	case account.SizeFrom11Till50Employees:
		size = companyRPC.Size_SIZE_11_50_EMPLOYEES
	case account.SizeFrom51Till200Employees:
		size = companyRPC.Size_SIZE_51_200_EMPLOYEES
	case account.SizeFrom201Till500Employees:
		size = companyRPC.Size_SIZE_201_500_EMPLOYEES
	case account.SizeFrom501Till1000Employees:
		size = companyRPC.Size_SIZE_501_1000_EMPLOYEES
	case account.SizeFrom1001Till5000Employees:
		size = companyRPC.Size_SIZE_1001_5000_EMPLOYEES
	case account.SizeFrom5001Till10000Employees:
		size = companyRPC.Size_SIZE_5001_10000_EMPLOYEES
	case account.SizeFrom10001AndMoreEmployees:
		size = companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES
	}

	return size
}

func accountParkingToParkingRPC(data account.Parking) companyRPC.Parking {
	var parking companyRPC.Parking

	switch data {
	case account.ParkingNoParking:
		parking = companyRPC.Parking_PARKING_NO_PARKING
	case account.ParkingParkingLot:
		parking = companyRPC.Parking_PARKING_PARKING_LOT
	case account.ParkingStreetParking:
		parking = companyRPC.Parking_PARKING_STREET_PARKING
	}

	return parking
}

// TODO:
func profileToProfileRPC(data *profile.Profile) *companyRPC.Profile {
	if data == nil {
		return nil
	}

	profile := companyRPC.Profile{
		Id:          data.GetID(),
		Avatar:      data.Avatar,
		Name:        data.Name,
		URL:         data.URL,
		Cover:       data.Cover,
		Description: data.Description,
		Mission:     data.Mission,
		Size:        accountSizeToSizeRPC(data.Size),
		Type:        accountTypeToTypeRPC(data.Type),
		Addresses:   make([]*companyRPC.Address, 0, len(data.Addresses)),
		Awards:      make([]*companyRPC.Award, 0, len(data.Awards)),
		// Galleries:      make([]*companyRPC.Gallery, 0, len(data.Gallery)),
		BusinessHours:  make([]*companyRPC.BusinessHoursItem, 0),
		Milestones:     make([]*companyRPC.Milestone, 0, len(data.Milestones)),
		FoundationDate: timeToStringDayMonthAndYear(data.FoundationDate),
		// Founders:       make([]*companyRPC.Founder, 0, len(data.Founders)),
		Products:              make([]*companyRPC.Product, 0, len(data.Products)),
		Services:              make([]*companyRPC.Service, 0, len(data.Services)),
		Websites:              make([]string, 0, len(data.Website)),
		WasAboutUsSet:         data.IsAboutUsSet,
		IsFollow:              data.IsFollow,
		IsFavorite:            data.IsFavorite,
		IsOnline:              data.IsOnline,
		IsBlocked:             data.IsBlocked,
		Benefits:              profileBenetisToCompanyrRPCBenefits(data.Benefits),
		Emails:                make([]string, 0, len(data.Emails)),
		Phones:                make([]string, 0, len(data.Phones)),
		AmountOfJobs:          data.AmountOfJobs,
		AvailableTranslations: data.AvailableTranslations,
		CurrentTranslation:    data.CurrentTranslation,
		CareerCenter:          careerCenterToCareerCenterRPC(data.CareerCenter),
	}

	for i := range data.Addresses {
		profile.Addresses = append(profile.Addresses, accountAddressToAddressRPC(data.Addresses[i]))
	}

	if parking := data.Parking; parking != nil {
		profile.Parking = accountParkingToParkingRPC(*parking)
	}

	for i := range data.BusinessHours {
		profile.BusinessHours = append(profile.BusinessHours, profileBusinessHourToBusinessHourItem(data.BusinessHours[i]))
	}

	// for i := range data.Gallery {
	// 	profile.Galleries = append(profile.Galleries, profileGalleryToCompanyRPCGallery(data.Gallery[i]))
	// }

	for i := range data.Awards {
		profile.Awards = append(profile.Awards, profileAwardToAwardRPC(data.Awards[i]))
	}

	for i := range data.Milestones {
		profile.Milestones = append(profile.Milestones, profileMilestoneToMilestoneRPC(data.Milestones[i]))
	}

	profile.Industry = &companyRPC.Industry{
		Main: data.Industry.Main,
		Subs: data.Industry.Sub,
	}

	// for i := range data.Founders {
	// 	profile.Founders = append(profile.Founders, profileFounderToFounderRPC(data.Founders[i]))
	// }

	for i := range data.Products {
		profile.Products = append(profile.Products, profileProductToProductRPC(data.Products[i]))
	}

	for i := range data.Services {
		profile.Services = append(profile.Services, profileServiceToServiceRPC(data.Services[i]))
	}

	for i := range data.Website {
		profile.Websites = append(profile.Websites, data.Website[i].Site)
	}

	for _, e := range data.Emails {
		if e.Activated {
			if e.Primary {
				profile.Email = &companyRPC.EmailProfile{
					Email: e.Email,
				}
			}
			profile.Emails = append(profile.Emails, e.Email)
		}
	}

	for _, p := range data.Phones {
		// if p.Activated {
		if p.Primary {
			profile.Phone = &companyRPC.PhoneProfile{
				CountryAbbreviation: p.CountryCode.CountryID,
				Number:              p.Number,
				CountryCode: &companyRPC.CountryCode{
					Code: p.CountryCode.Code,
				},
			}
		}

		profile.Phones = append(profile.Phones, p.CountryCode.Code+p.Number)
		// }
	}

	profile.AmountOfEmployees = data.AmountOfEmployees
	profile.AmountOfFollowers = data.AmountOfFollowers
	profile.AmountOfFollowings = data.AmountOfFollowings
	profile.AvarageRating = data.AvarageRating

	return &profile
}

func profileBusinessHourToBusinessHourItem(data *account.BusinessHour) *companyRPC.BusinessHoursItem {
	if data == nil {
		return nil
	}

	bh := companyRPC.BusinessHoursItem{
		ID:       data.GetID(),
		WeekDays: data.Weekdays, // TOOD: change
		HourFrom: data.StartHour,
		HourTo:   data.FinishHour,
	}

	return &bh
}

func profileAwardToAwardRPC(data *profile.Award) *companyRPC.Award {
	if data == nil {
		return nil
	}

	award := companyRPC.Award{
		ID:     data.GetID(),
		Issuer: data.Issuer,
		Title:  data.Title,
		Year:   int32(data.Date.Year()),
		Links:  make([]*companyRPC.Link, 0, len(data.Links)),
		Files:  make([]*companyRPC.File, 0, len(data.Files)),
	}

	for i := range data.Links {
		if data.Links[i] != nil { // TODO: findout why it why is saves nulls
			award.Links = append(award.Links, profileLinktoLinkRPC(data.Links[i]))
		}
	}

	for i := range data.Files {
		award.Files = append(award.Files, profileFiletoFileRPC(data.Files[i]))
	}

	return &award
}

func profileMilestoneToMilestoneRPC(data *profile.Milestone) *companyRPC.Milestone {
	if data == nil {
		return nil
	}

	m := companyRPC.Milestone{
		Id:          data.GetID(),
		Description: data.Description,
		Image:       data.Image,
		Title:       data.Title,
		Year:        int32(data.Date.Year()),
	}

	return &m
}

func profileReviewToReviewRPC(data *profile.Review) *companyRPC.Review {
	if data == nil {
		return nil
	}

	review := companyRPC.Review{
		Id:          data.GetID(),
		Rate:        uint32(data.Rate),
		Headline:    data.Headline,
		Description: data.Description,
		AuthorID:    data.GetAuthorID(),
		Company:     profileToProfileRPC(&data.Company),
		Date:        timeToStringDayMonthAndYear(data.CreatedAt),
	}

	return &review
}

func accountWebsiteToWebsiteRPC(data *account.Website) *companyRPC.Website {
	if data == nil {
		return nil
	}

	website := companyRPC.Website{
		Id:      data.GetID(),
		Website: data.Site,
	}

	return &website
}

func profileFounderToFounderRPC(data *profile.Founder) *companyRPC.Founder {
	if data == nil {
		return nil
	}

	founder := companyRPC.Founder{
		ID:            data.GetID(),
		Avatar:        data.Avatar,
		Name:          data.Name,
		PositionTitle: data.Position,
		UserID:        data.GetUserID(),
		IsApproved:    data.Approved,
	}

	return &founder
}

func profileGalleryToGalleryRPC(data *profile.Gallery) *companyRPC.Gallery {
	if data == nil {
		return nil
	}

	gallery := companyRPC.Gallery{
		ID:   data.GetID(),
		File: make([]*companyRPC.GalleryFile, 0, len(data.Files)),
	}

	for _, f := range data.Files {
		gallery.File = append(gallery.File, profileFileToCompanyRPCGalleryFile(f))
	}

	return &gallery
}

func profileProductToProductRPC(data *profile.Product) *companyRPC.Product {
	if data == nil {
		return nil
	}
	product := companyRPC.Product{
		ID:      data.GetID(),
		Image:   data.Image,
		Name:    data.Name,
		Website: data.Website,
	}

	return &product
}

func profileServiceToServiceRPC(data *profile.Service) *companyRPC.Service {
	if data == nil {
		return nil
	}

	service := companyRPC.Service{
		ID:      data.GetID(),
		Image:   data.Image,
		Name:    data.Name,
		Website: data.Website,
	}

	return &service
}

func accountAdminLevelToAdminRoleRPC(data account.AdminLevel) companyRPC.AdminRole {
	lvl := companyRPC.AdminRole_ROLE_UNKNOWN

	switch data {
	case account.AdminLevelAdmin:
		lvl = companyRPC.AdminRole_ROLE_ADMIN
	case account.AdminLevelCommercial:
		lvl = companyRPC.AdminRole_ROLE_COMMERCIAL_ADMIN
	case account.AdminLevelJob:
		lvl = companyRPC.AdminRole_ROLE_JOB_EDITOR
	case account.AdminLevelVShop:
		lvl = companyRPC.AdminRole_ROLE_V_SHOP_ADMIN
	case account.AdminLevelVService:
		lvl = companyRPC.AdminRole_ROLE_V_SERVICE_ADMIN
	}

	return lvl
}

// ---------

func stringDateToTime(s string) time.Time {
	if date, err := time.Parse("2-1-2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func stringDayMonthAndYearToTime(s string) time.Time {
	var err error
	if date, err := time.Parse("1-2006", s); err == nil {
		return date
	}

	log.Println("error: date:", err)

	return time.Time{}
}

func stringYearToDate(s string) time.Time {
	if date, err := time.Parse("2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func timeToStringMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, _ := t.UTC().Date()
	return strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func timeToStringDayMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, d := t.UTC().Date()
	return strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func linkRPCToProfileLink(data *companyRPC.Link) *profile.Link {
	if data == nil {
		return nil
	}

	l := profile.Link{
		URL: data.GetURL(),
	}

	_ = l.SetID(data.GetID())

	return &l
}

func profileFiletoFileRPC(data *profile.File) *companyRPC.File {
	if data == nil {
		return nil
	}

	return &companyRPC.File{
		ID:       data.GetID(),
		URL:      data.URL,
		MimeType: data.MimeType,
		Name:     data.Name,
	}
}

func profileGalleryFiletoFileRPC(data *profile.File) *companyRPC.GalleryFile {
	if data == nil {
		return nil
	}

	return &companyRPC.GalleryFile{
		ID:       data.GetID(),
		URL:      data.URL,
		MimeType: data.MimeType,
		Name:     data.Name,
	}
}

func profileLinktoLinkRPC(data *profile.Link) *companyRPC.Link {
	if data == nil {
		return nil
	}

	return &companyRPC.Link{
		ID:  data.GetID(),
		URL: data.URL,
	}
}

func profileGalleryToCompanyRPCGallery(data *profile.Gallery) *companyRPC.Gallery {
	if data == nil {
		return nil
	}

	return &companyRPC.Gallery{
		ID: data.GetID(),
	}
}

func careerCenterToCareerCenterRPC(data *careercenter.CareerCenter) *companyRPC.CareerCenter {
	if data == nil {
		return nil
	}

	return &companyRPC.CareerCenter{
		CVButtonEnabled:     data.CVButtonEnabled,
		CustomButtonEnabled: data.CustomButtonEnabled,
		CustomButtontitle:   data.CustomButtonTitle,
		CustomButtonURL:     data.CustomButtonURL,
		Description:         data.Description,
		Title:               data.Title,
	}
}
