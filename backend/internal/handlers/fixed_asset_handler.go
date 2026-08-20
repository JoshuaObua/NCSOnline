package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
)

type FixedAssetHandler struct {
	repo repository.FixedAssetRepository
}

func NewFixedAssetHandler(repo repository.FixedAssetRepository) *FixedAssetHandler {
	return &FixedAssetHandler{repo: repo}
}

func (h *FixedAssetHandler) ListFixedAssets(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	assets, total, err := h.repo.ListFixedAssets(r.Context(), category, search, limit, offset)
	if err != nil {
		fixedAssetError(w, http.StatusInternalServerError, "Failed to fetch fixed assets: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"assets": assets,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *FixedAssetHandler) GetFixedAssetSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.repo.GetFixedAssetSummary(r.Context())
	if err != nil {
		fixedAssetError(w, http.StatusInternalServerError, "Failed to fetch asset summary: "+err.Error())
		return
	}
	response.JSON(w, http.StatusOK, summary)
}

func (h *FixedAssetHandler) GetFixedAssetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	asset, err := h.repo.GetFixedAssetByID(r.Context(), id)
	if err != nil {
		fixedAssetError(w, http.StatusNotFound, "Asset not found")
		return
	}
	logs, _ := h.repo.GetTransactionLogs(r.Context(), asset.ID)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"asset": asset,
		"logs":  logs,
	})
}

func (h *FixedAssetHandler) CreateFixedAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.FixedAsset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		fixedAssetError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if asset.AssetNumber == "" || asset.AssetDescription == "" {
		fixedAssetError(w, http.StatusBadRequest, "Asset number and description are required")
		return
	}
	if strings.TrimSpace(asset.TagNumber) == "" {
		fixedAssetError(w, http.StatusBadRequest, "Tag number is required")
		return
	}
	applyFixedAssetDefaults(&asset)

	if asset.AdjustedCost == 0 && asset.FBCost > 0 {
		asset.AdjustedCost = asset.FBCost
	}

	err := h.repo.CreateFixedAsset(r.Context(), &asset)
	if err != nil {
		fixedAssetError(w, http.StatusInternalServerError, "Failed to create fixed asset: "+err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, asset)
}

func (h *FixedAssetHandler) RevalueFixedAsset(w http.ResponseWriter, r *http.Request) {
	var req models.AssetRevaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fixedAssetError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AssetID == "" || req.NewCost < 0 {
		fixedAssetError(w, http.StatusBadRequest, "Asset ID and valid new cost are required")
		return
	}

	userID := currentUserID(r)

	err := h.repo.RevalueFixedAsset(r.Context(), req.AssetID, req.NewCost, req.Notes, userID)
	if err != nil {
		fixedAssetError(w, http.StatusInternalServerError, "Failed to revalue asset: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Asset revalued successfully",
	})
}

func (h *FixedAssetHandler) RunDepreciation(w http.ResponseWriter, r *http.Request) {
	var req models.AssetDepreciationRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	userID := currentUserID(r)

	count, totalAmount, err := h.repo.RunDepreciation(r.Context(), req.Period, userID)
	if err != nil {
		fixedAssetError(w, http.StatusInternalServerError, "Failed to run depreciation: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":          "Depreciation run executed successfully",
		"assets_processed": count,
		"total_deprec_run": totalAmount,
	})
}

func (h *FixedAssetHandler) VerifyFixedAsset(w http.ResponseWriter, r *http.Request) {
	var req models.AssetVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fixedAssetError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := currentUserID(r)
	if req.Status == "" {
		req.Status = "VERIFIED"
	}

	err := h.repo.VerifyFixedAsset(r.Context(), req.AssetID, req.Status, req.Notes, userID)
	if err != nil {
		fixedAssetError(w, http.StatusInternalServerError, "Failed to verify asset: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Asset verification recorded successfully",
	})
}

func fixedAssetError(w http.ResponseWriter, status int, message string) {
	response.Err(w, status, "FIXED_ASSET_ERROR", message)
}

func currentUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(models.CtxUserID).(string); ok {
		return userID
	}
	if claims, ok := r.Context().Value("user_claims").(map[string]interface{}); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub
		}
	}
	return ""
}

func applyFixedAssetDefaults(asset *models.FixedAsset) {
	asset.AssetNumber = strings.TrimSpace(asset.AssetNumber)
	asset.TagNumber = strings.TrimSpace(asset.TagNumber)
	asset.AssetDescription = strings.TrimSpace(asset.AssetDescription)
	if strings.TrimSpace(asset.InterfaceLineNumber) == "" {
		asset.InterfaceLineNumber = "MANUAL-" + strings.ToUpper(asset.AssetNumber)
	}
	if strings.TrimSpace(asset.AssetBook) == "" {
		asset.AssetBook = "NCS FA BOOK"
	}
	if strings.TrimSpace(asset.CategorySegment1) == "" {
		asset.CategorySegment1 = "MACHINERY AND EQUIPMENT"
	}
	if strings.TrimSpace(asset.CategorySegment3) == "" {
		asset.CategorySegment3 = "LIGHT ICT HARDWARE"
	}
	if strings.TrimSpace(asset.CategorySegment4) == "" {
		asset.CategorySegment4 = "General Asset"
	}
	if asset.AssetUnits <= 0 {
		asset.AssetUnits = 1
	}
	if strings.TrimSpace(asset.DatePlacedInService) == "" {
		asset.DatePlacedInService = time.Now().Format("2006-01-02")
	}
	if strings.TrimSpace(asset.CustodianDepartment) == "" {
		asset.CustodianDepartment = "General Administration"
	}
	if strings.TrimSpace(asset.LocationBuilding) == "" {
		asset.LocationBuilding = "NCS Lugogo Head Office"
	}
	if strings.TrimSpace(asset.LocationRoom) == "" {
		asset.LocationRoom = "Main Facility"
	}
	if strings.TrimSpace(asset.DepreciationMethod) == "" {
		asset.DepreciationMethod = "STRAIGHT_LINE"
	}
	if asset.UsefulLifeYears <= 0 {
		asset.UsefulLifeYears = 5
	}
	if strings.TrimSpace(asset.Status) == "" {
		asset.Status = "ACTIVE"
	}
	if strings.TrimSpace(asset.WorksheetSource) == "" {
		asset.WorksheetSource = asset.CategorySegment3
	}
}
