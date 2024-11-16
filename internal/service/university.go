package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/constants"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
)

type universityService struct {
	cfg            config.Config
	universityRepo repository.UniversityRepository
	userRepo       repository.UserRepository
	db             *bun.DB
}

func NewUniversityService(cfg config.Config, universityRepo repository.UniversityRepository, userRepo repository.UserRepository, db *bun.DB) UniversityService {

	return &universityService{
		cfg:            cfg,
		universityRepo: universityRepo,
		userRepo:       userRepo,
		db:             db,
	}
}

func (s *universityService) CreateUniversityRating(ctx context.Context, req request.RateUniversityReq) error {
	//ctx, endFunc := trace.Start(ctx, "UniversityService.CreateUniversityRating", "service")
	//defer endFunc()

	var (
		overallRating   float64
		uniRatingPoints []model.UniversityRatingPoints
	)

	user, err := s.userRepo.GetByEmail(req.UserEmail)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] User not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.NotFound).Errorf("User not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to get user by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user data")
	}

	if user.UniversityID == nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] user does not belong to any university", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User does not belong to any university")
	}

	if *user.UniversityID != req.UniversityID {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] user does not belong to this university", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User does not belong to this university")
	}

	existingUniRating, err := s.universityRepo.GetUniversityRatingByUserIDAndUniversityID(req.UserID, req.UniversityID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to get university rating by user id and university id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get existing university rating")
	}

	if existingUniRating != nil && existingUniRating.ID != "" {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] University rating already exists")
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("University rating already exists")
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	totalRating := req.FacilityRating +
		req.PriceToValueRating +
		req.EducationQualityRating +
		req.StudentOrganizationRating +
		req.SocialEnvironmentRating

	overallRating = float64(totalRating) / float64(5)

	uniRating := &model.UniversityRating{
		ID:                        uuid.NewString(),
		UserID:                    req.UserID,
		UniversityID:              req.UniversityID,
		Title:                     req.Title,
		Content:                   req.Content,
		UniversityMajor:           req.UniversityMajor,
		FacilityRating:            req.FacilityRating,
		StudentOrganizationRating: req.StudentOrganizationRating,
		SocialEnvironmentRating:   req.SocialEnvironmentRating,
		EducationQualityRating:    req.EducationQualityRating,
		PriceToValueRating:        req.PriceToValueRating,
		OverallRating:             overallRating,
		CreatedBy:                 req.UserEmail,
		UpdatedBy:                 req.UserEmail,
	}

	err = s.universityRepo.SaveUniversityRatingTx(uniRating, tx)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to insert university rating to database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to create university rating")
	}

	for _, pro := range req.Pros {
		uniRatingPoints = append(uniRatingPoints, model.UniversityRatingPoints{
			ID:                 uuid.NewString(),
			UniversityRatingID: uniRating.ID,
			Type:               constants.UNI_RATING_PRO,
			Content:            pro,
			CreatedBy:          req.UserEmail,
			UpdatedBy:          req.UserEmail,
		})
	}

	for _, con := range req.Cons {
		uniRatingPoints = append(uniRatingPoints, model.UniversityRatingPoints{
			ID:                 uuid.NewString(),
			UniversityRatingID: uniRating.ID,
			Type:               constants.UNI_RATING_CON,
			Content:            con,
			CreatedBy:          req.UserEmail,
			UpdatedBy:          req.UserEmail,
		})
	}

	err = s.universityRepo.BulkSaveUniversityRatingPointsTx(uniRatingPoints, tx)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to insert university rating points to database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to create university rating pros and cons")
	}

	err = s.userRepo.SetUserHasRateUniversityTx(req.UserID, true, tx)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to update user data has rate university", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update user data")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUniversityRating] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *universityService) GetUniversityRatingList(ctx context.Context, req request.GetUniversityRatingListReq) (response.GetUniversityRatingListResponse, error) {
	//ctx, endFunc := trace.Start(ctx, "UniversityService.GetUniversityRatingList", "service")
	//defer endFunc()

	var resp response.GetUniversityRatingListResponse

	uniRatings, pagination, err := s.universityRepo.GetList(req)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetUniversityRatingList] Failed to get university rating list", zap.Error(err))

		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get university rating list")
	}

	resp.Meta = response.PaginationMeta{
		CurrentCursor: pagination.CurrentCursor,
		NextCursor:    pagination.NextCursor,
	}

	resp.Data = s.mapUniversityRatingListData(uniRatings)

	return resp, nil
}

func (s *universityService) mapUniversityRatingListData(uniRatings []model.UniversityRating) []response.UniversityRatingListData {

	uniRatingData := []response.UniversityRatingListData{}

	for _, ur := range uniRatings {

		uniRatingPros := []string{}
		uniRatingCons := []string{}

		urd := response.UniversityRatingListData{
			ID:                        ur.ID,
			UserID:                    ur.UserID,
			UserName:                  ur.User.Username,
			UniversityID:              ur.UniversityID,
			UniversityAbbreviatedName: ur.University.AbbreviatedName,
			UniversityImageURL:        ur.University.ImageURL,
			UniversityMajor:           ur.UniversityMajor,
			Title:                     ur.Title,
			Content:                   ur.Content,
			FacilityRating:            ur.FacilityRating,
			StudentOrganizationRating: ur.StudentOrganizationRating,
			SocialEnvironmentRating:   ur.SocialEnvironmentRating,
			EducationQualityRating:    ur.EducationQualityRating,
			PriceToValueRating:        ur.PriceToValueRating,
			OverallRating:             ur.OverallRating,
			CreatedBy:                 ur.CreatedBy,
			CreatedAt:                 ur.CreatedAt,
			UpdatedBy:                 ur.UpdatedBy,
			UpdatedAt:                 ur.UpdatedAt,
		}

		for _, urp := range ur.UniversityRatingPoints {
			if urp.Type == constants.UNI_RATING_PRO {
				uniRatingPros = append(uniRatingPros, urp.Content)
				continue
			}

			if urp.Type == constants.UNI_RATING_CON {
				uniRatingCons = append(uniRatingCons, urp.Content)
				continue
			}
		}

		urd.Pros = uniRatingPros
		urd.Cons = uniRatingCons

		uniRatingData = append(uniRatingData, urd)
	}

	return uniRatingData
}
