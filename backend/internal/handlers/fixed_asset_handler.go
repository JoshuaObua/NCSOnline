package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"ncsintranet/internal/models"
	"ncsintranet/internal/repository"
	"ncsintranet/internal/response"
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
		response.Error(w, http.StatusInternalServerError, "Failed to fetch fixed assets: "+err.Error())
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
		response.Error(w, http.StatusInternalServerError, "Failed to fetch asset summary: "+err.Error())
		return
	}
	response.JSON(w, http.StatusOK, summary)
}

func (h *FixedAssetHandler) GetFixedAssetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	asset, err := h.repo.GetFixedAssetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Asset not found")
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
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if asset.AssetNumber == "" || asset.AssetDescription == "" {
		response.Error(w, http.StatusBadRequest, "Asset number and description are required")
		return
	}

	if asset.AdjustedCost == 0 && asset.FBCost > 0 {
		asset.AdjustedCost = asset.FBCost
	}

	err := h.repo.CreateFixedAsset(r.Context(), &asset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create fixed asset: "+err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, asset)
}

func (h *FixedAssetHandler) RevalueFixedAsset(w http.ResponseWriter, r *http.Request) {
	var req models.AssetRevaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AssetID == "" || req.NewCost < 0 {
		response.Error(w, http.StatusBadRequest, "Asset ID and valid new cost are required")
		return
	}

	userID := ""
	if claims, ok := r.Context().Value("user_claims").(map[string]interface{}); ok {
		if sub, ok := claims["sub"].(string); ok {
			userID = sub
		}
	}

	err := h.repo.RevalueFixedAsset(r.Context(), req.AssetID, req.NewCost, req.Notes, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to revalue asset: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Asset revalued successfully",
	})
}

func (h *FixedAssetHandler) RunDepreciation(w http.ResponseWriter, r *http.Request) {
	var req models.AssetDepreciationRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	userID := ""
	if claims, ok := r.Context().Value("user_claims").(map[string]interface{}); ok {
		if sub, ok := claims["sub"].(string); ok {
			userID = sub
		}
	}

	count, totalAmount, err := h.repo.RunDepreciation(r.Context(), req.Period, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to run depreciation: "+err.Error())
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
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := ""
	if claims, ok := r.Context().Value("user_claims").(map[string]interface{}); ok {
		if sub, ok := claims["sub"].(string); ok {
			userID = sub
		}
	}

	err := h.repo.VerifyFixedAsset(r.Context(), req.AssetID, req.Status, req.Notes, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to verify asset: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Asset verification recorded successfully",
	})
}
