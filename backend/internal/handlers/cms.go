package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CMSHandler struct{ repo *repository.CMSRepo }

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func toSlug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return strings.Trim(slugRe.ReplaceAllString(b.String(), "-"), "-")
}

func paginate(r *http.Request) (limit, offset int) {
	p := &models.PaginationParams{}
	if v := r.URL.Query().Get("page"); v != "" {
		for _, c := range v {
			if c >= '0' && c <= '9' {
				p.Page = p.Page*10 + int(c-'0')
			}
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		for _, c := range v {
			if c >= '0' && c <= '9' {
				p.PerPage = p.PerPage*10 + int(c-'0')
			}
		}
	}
	p.Offset()
	if p.PerPage == 0 {
		p.PerPage = 20
	}
	return p.PerPage, p.Offset()
}

// ── Posts ─────────────────────────────────────────────────────────

// GET /api/v1/cms/posts?category=blog&status=published&page=1
func (h *CMSHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	limit, offset := paginate(r)
	posts, total, err := h.repo.ListPosts(r.Context(), category, status, limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list posts")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": posts})
}

// GET /api/v1/cms/posts/{slug}
func (h *CMSHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := h.repo.GetPostBySlug(r.Context(), slug)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get post")
		return
	}
	response.JSON(w, http.StatusOK, post)
}

// POST /api/v1/admin/cms/posts
func (h *CMSHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	authorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Title           string  `json:"title"`
		Content         string  `json:"content"`
		Excerpt         string  `json:"excerpt"`
		Category        string  `json:"category"`
		Status          string  `json:"status"`
		CoverImageURL   string  `json:"cover_image_url"`
		Slug            *string `json:"slug"`
		MetaTitle       string  `json:"meta_title"`
		MetaDescription string  `json:"meta_description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	slug := req.Title
	if req.Slug != nil && *req.Slug != "" {
		slug = *req.Slug
	}
	if req.Category == "" {
		req.Category = "blog"
	}
	if req.Status == "" {
		req.Status = "draft"
	}

	now := time.Now()
	var publishedAt *time.Time
	if req.Status == "published" {
		publishedAt = &now
	}

	post := &models.CMSPost{
		ID:              uuid.NewString(),
		Title:           req.Title,
		Slug:            toSlug(slug),
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Category:        req.Category,
		Status:          req.Status,
		CoverImageURL:   req.CoverImageURL,
		AuthorID:        &authorID,
		PublishedAt:     publishedAt,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}
	if err := h.repo.CreatePost(r.Context(), post); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A post with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create post")
		return
	}
	response.JSON(w, http.StatusCreated, post)
}

// PUT /api/v1/admin/cms/posts/{id}
func (h *CMSHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.repo.GetPostByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get post")
		return
	}

	var req struct {
		Title           string  `json:"title"`
		Content         string  `json:"content"`
		Excerpt         string  `json:"excerpt"`
		Category        string  `json:"category"`
		Status          string  `json:"status"`
		CoverImageURL   string  `json:"cover_image_url"`
		Slug            *string `json:"slug"`
		MetaTitle       string  `json:"meta_title"`
		MetaDescription string  `json:"meta_description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.Category != "" {
		post.Category = req.Category
	}
	post.CoverImageURL = req.CoverImageURL
	if req.Slug != nil && *req.Slug != "" {
		post.Slug = toSlug(*req.Slug)
	}
	post.MetaTitle = req.MetaTitle
	post.MetaDescription = req.MetaDescription
	if req.Status != "" {
		prevStatus := post.Status
		post.Status = req.Status
		if prevStatus != "published" && req.Status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	if err := h.repo.UpdatePost(r.Context(), post); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update post")
		return
	}
	response.JSON(w, http.StatusOK, post)
}

// DELETE /api/v1/admin/cms/posts/{id}
func (h *CMSHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeletePost(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete post")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Post deleted")
}

// ── Events ────────────────────────────────────────────────────────

// GET /api/v1/cms/events?status=published
func (h *CMSHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, offset := paginate(r)
	events, total, err := h.repo.ListEvents(r.Context(), status, limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list events")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": events})
}

// GET /api/v1/cms/events/{slug}
func (h *CMSHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	event, err := h.repo.GetEventBySlug(r.Context(), slug)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Event not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get event")
		return
	}
	response.JSON(w, http.StatusOK, event)
}

// POST /api/v1/admin/cms/events
func (h *CMSHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	authorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Title         string     `json:"title"`
		Description   string     `json:"description"`
		Location      string     `json:"location"`
		EventDate     *time.Time `json:"event_date"`
		EndDate       *time.Time `json:"end_date"`
		CoverImageURL string     `json:"cover_image_url"`
		Status        string     `json:"status"`
		Slug          *string    `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	slug := req.Title
	if req.Slug != nil && *req.Slug != "" {
		slug = *req.Slug
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	event := &models.CMSEvent{
		ID:            uuid.NewString(),
		Title:         req.Title,
		Slug:          toSlug(slug),
		Description:   req.Description,
		Location:      req.Location,
		EventDate:     req.EventDate,
		EndDate:       req.EndDate,
		CoverImageURL: req.CoverImageURL,
		Status:        req.Status,
		AuthorID:      &authorID,
	}
	if err := h.repo.CreateEvent(r.Context(), event); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "An event with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create event")
		return
	}
	response.JSON(w, http.StatusCreated, event)
}

// PUT /api/v1/admin/cms/events/{id}
func (h *CMSHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	event, err := h.repo.GetEventByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Event not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get event")
		return
	}
	var req struct {
		Title         string     `json:"title"`
		Description   string     `json:"description"`
		Location      string     `json:"location"`
		EventDate     *time.Time `json:"event_date"`
		EndDate       *time.Time `json:"end_date"`
		CoverImageURL string     `json:"cover_image_url"`
		Status        string     `json:"status"`
		Slug          *string    `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.EventDate != nil {
		event.EventDate = req.EventDate
	}
	if req.EndDate != nil {
		event.EndDate = req.EndDate
	}
	if req.CoverImageURL != "" {
		event.CoverImageURL = req.CoverImageURL
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if req.Slug != nil && *req.Slug != "" {
		event.Slug = toSlug(*req.Slug)
	}
	if err := h.repo.UpdateEvent(r.Context(), event); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update event")
		return
	}
	response.JSON(w, http.StatusOK, event)
}

// DELETE /api/v1/admin/cms/events/{id}
func (h *CMSHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteEvent(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete event")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Event deleted")
}

// ── Careers ───────────────────────────────────────────────────────

// GET /api/v1/cms/careers?status=published&category=jobs
func (h *CMSHandler) ListCareers(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	category := r.URL.Query().Get("category")
	limit, offset := paginate(r)
	careers, total, err := h.repo.ListCareers(r.Context(), status, category, limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list careers")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": careers})
}

// GET /api/v1/cms/careers/{id}
func (h *CMSHandler) GetCareer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	career, err := h.repo.GetCareerByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Job posting not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get career")
		return
	}
	response.JSON(w, http.StatusOK, career)
}

// POST /api/v1/admin/cms/careers
func (h *CMSHandler) CreateCareer(w http.ResponseWriter, r *http.Request) {
	authorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Title        string     `json:"title"`
		Department   string     `json:"department"`
		Location     string     `json:"location"`
		JobType      string     `json:"job_type"`
		Category     string     `json:"category"`
		Description  string     `json:"description"`
		Requirements string     `json:"requirements"`
		SalaryRange  string     `json:"salary_range"`
		Status       string     `json:"status"`
		DeadlineAt   *time.Time `json:"deadline_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title == "" || req.Description == "" {
		errs := map[string]string{}
		if req.Title == "" {
			errs["title"] = "required"
		}
		if req.Description == "" {
			errs["description"] = "required"
		}
		response.ValidationErr(w, errs)
		return
	}
	if req.JobType == "" {
		req.JobType = "full_time"
	}
	if req.Category == "" {
		req.Category = "jobs"
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	career := &models.CMSCareer{
		ID:           uuid.NewString(),
		Title:        req.Title,
		Department:   req.Department,
		Location:     req.Location,
		JobType:      req.JobType,
		Category:     req.Category,
		Description:  req.Description,
		Requirements: req.Requirements,
		SalaryRange:  req.SalaryRange,
		Status:       req.Status,
		DeadlineAt:   req.DeadlineAt,
		AuthorID:     &authorID,
	}
	if err := h.repo.CreateCareer(r.Context(), career); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create career")
		return
	}
	response.JSON(w, http.StatusCreated, career)
}

// PUT /api/v1/admin/cms/careers/{id}
func (h *CMSHandler) UpdateCareer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	career, err := h.repo.GetCareerByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Career not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get career")
		return
	}
	var req struct {
		Title        string     `json:"title"`
		Department   string     `json:"department"`
		Location     string     `json:"location"`
		JobType      string     `json:"job_type"`
		Category     string     `json:"category"`
		Description  string     `json:"description"`
		Requirements string     `json:"requirements"`
		SalaryRange  string     `json:"salary_range"`
		Status       string     `json:"status"`
		DeadlineAt   *time.Time `json:"deadline_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title != "" {
		career.Title = req.Title
	}
	if req.Department != "" {
		career.Department = req.Department
	}
	if req.Location != "" {
		career.Location = req.Location
	}
	if req.JobType != "" {
		career.JobType = req.JobType
	}
	if req.Category != "" {
		career.Category = req.Category
	}
	if req.Description != "" {
		career.Description = req.Description
	}
	if req.Requirements != "" {
		career.Requirements = req.Requirements
	}
	if req.SalaryRange != "" {
		career.SalaryRange = req.SalaryRange
	}
	if req.Status != "" {
		career.Status = req.Status
	}
	if req.DeadlineAt != nil {
		career.DeadlineAt = req.DeadlineAt
	}
	if err := h.repo.UpdateCareer(r.Context(), career); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update career")
		return
	}
	response.JSON(w, http.StatusOK, career)
}

// DELETE /api/v1/admin/cms/careers/{id}
func (h *CMSHandler) DeleteCareer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteCareer(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete career")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Career posting deleted")
}

// ── Slides ────────────────────────────────────────────────────────

// GET /api/v1/cms/slides
func (h *CMSHandler) ListSlides(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	slides, err := h.repo.ListSlides(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list slides")
		return
	}
	response.JSON(w, http.StatusOK, slides)
}

type slideReq struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ButtonText  string `json:"button_text"`
	ButtonURL   string `json:"button_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// POST /api/v1/admin/cms/slides
func (h *CMSHandler) CreateSlide(w http.ResponseWriter, r *http.Request) {
	var req slideReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}
	if req.Title == "" {
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "Title is required")
		return
	}
	s := &models.CMSSlide{
		ID: uuid.NewString(), Title: req.Title, Subtitle: req.Subtitle,
		Description: req.Description, ImageURL: req.ImageURL,
		ButtonText: req.ButtonText, ButtonURL: req.ButtonURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateSlide(r.Context(), s); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create slide")
		return
	}
	response.JSON(w, http.StatusCreated, s)
}

// PUT /api/v1/admin/cms/slides/{id}
func (h *CMSHandler) UpdateSlide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req slideReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}
	s := &models.CMSSlide{
		ID: id, Title: req.Title, Subtitle: req.Subtitle,
		Description: req.Description, ImageURL: req.ImageURL,
		ButtonText: req.ButtonText, ButtonURL: req.ButtonURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.UpdateSlide(r.Context(), s); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update slide")
		return
	}
	response.JSON(w, http.StatusOK, s)
}

// DELETE /api/v1/admin/cms/slides/{id}
func (h *CMSHandler) DeleteSlide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteSlide(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete slide")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Slide deleted")
}

// ── Menus ─────────────────────────────────────────────────────────

// GET /api/v1/cms/menus/{name}
func (h *CMSHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	menu, err := h.repo.GetMenu(r.Context(), name)
	if errors.Is(err, repository.ErrNotFound) {
		response.JSON(w, http.StatusOK, &models.CMSMenu{Name: name, Items: []models.CMSMenuItem{}})
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get menu")
		return
	}
	response.JSON(w, http.StatusOK, menu)
}

// PUT /api/v1/admin/cms/menus/{name}
func (h *CMSHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var items []models.CMSMenuItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid menu items")
		return
	}
	if err := h.repo.UpdateMenu(r.Context(), name, items); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update menu")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Menu updated")
}

// ── Settings (generic JSON key/value) ─────────────────────────────

// GET /api/v1/cms/settings/{key}
func (h *CMSHandler) GetSetting(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	s, err := h.repo.GetSetting(r.Context(), key)
	if errors.Is(err, repository.ErrNotFound) {
		// Treat missing as empty value so the public site can render
		// defaults without per-request error handling.
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": map[string]interface{}{},
		})
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get setting")
		return
	}
	response.JSON(w, http.StatusOK, s)
}

// PUT /api/v1/admin/cms/settings/{key}
func (h *CMSHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Body required")
		return
	}
	// Validate the body is well-formed JSON; store raw bytes.
	var probe interface{}
	if err := json.Unmarshal(body, &probe); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Body must be valid JSON")
		return
	}
	if err := h.repo.UpdateSetting(r.Context(), key, body); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save setting")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Setting saved")
}

// ── Media Upload ──────────────────────────────────────────────────

// POST /api/v1/admin/media/upload
func (h *CMSHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Could not parse form (max 32MB)")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "No file uploaded (field: 'file')")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".svg": true, ".pdf": true,
	}
	if !allowed[ext] {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "File type not allowed. Allowed: jpg, jpeg, png, gif, webp, svg, pdf")
		return
	}

	subDir := "images"
	if ext == ".pdf" {
		subDir = "documents"
	}

	uploadDir := fmt.Sprintf("/app/uploads/%s", subDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create upload directory")
		return
	}

	filename := uuid.NewString() + ext
	filePath := filepath.Join(uploadDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not write file")
		return
	}

	url := fmt.Sprintf("/uploads/%s/%s", subDir, filename)
	response.JSON(w, http.StatusOK, map[string]string{"url": url, "filename": filename})
}

// ── Fun Facts ─────────────────────────────────────────────────────

func (h *CMSHandler) ListFunFacts(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListFunFacts(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list fun facts")
		return
	}
	if items == nil {
		items = []*models.CMSFunFact{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateFunFact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Label     string `json:"label"`
		Value     string `json:"value"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Label == "" || req.Value == "" {
		response.ValidationErr(w, map[string]string{"label": "required", "value": "required"})
		return
	}
	f := &models.CMSFunFact{
		ID: uuid.NewString(), Label: req.Label, Value: req.Value,
		Icon: req.Icon, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateFunFact(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create fun fact")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

func (h *CMSHandler) UpdateFunFact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Label     string `json:"label"`
		Value     string `json:"value"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	f := &models.CMSFunFact{
		ID: id, Label: req.Label, Value: req.Value,
		Icon: req.Icon, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.UpdateFunFact(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update fun fact")
		return
	}
	response.JSON(w, http.StatusOK, f)
}

func (h *CMSHandler) DeleteFunFact(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteFunFact(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete fun fact")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── FAQs ──────────────────────────────────────────────────────────

func (h *CMSHandler) ListFAQs(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	items, err := h.repo.ListFAQs(r.Context(), category)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list FAQs")
		return
	}
	if items == nil {
		items = []*models.CMSFAQ{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		Category  string `json:"category"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Question == "" {
		response.ValidationErr(w, map[string]string{"question": "required"})
		return
	}
	if req.Category == "" {
		req.Category = "general"
	}
	f := &models.CMSFAQ{
		ID: uuid.NewString(), Question: req.Question, Answer: req.Answer,
		Category: req.Category, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateFAQ(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create FAQ")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

func (h *CMSHandler) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		Category  string `json:"category"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	f := &models.CMSFAQ{
		ID: id, Question: req.Question, Answer: req.Answer,
		Category: req.Category, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.UpdateFAQ(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update FAQ")
		return
	}
	response.JSON(w, http.StatusOK, f)
}

func (h *CMSHandler) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteFAQ(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete FAQ")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Resources ─────────────────────────────────────────────────────

func (h *CMSHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	items, err := h.repo.ListResources(r.Context(), category)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list resources")
		return
	}
	if items == nil {
		items = []*models.CMSResource{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		FileURL     string `json:"file_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	if req.Category == "" {
		req.Category = "general"
	}
	res := &models.CMSResource{
		ID: uuid.NewString(), Title: req.Title, Category: req.Category,
		FileURL: req.FileURL, Description: req.Description,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateResource(r.Context(), res); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create resource")
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *CMSHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetResourceByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get resource")
		return
	}
	var req struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		FileURL     string `json:"file_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.FileURL != "" {
		existing.FileURL = req.FileURL
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateResource(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update resource")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteResource(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete resource")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Facilities ────────────────────────────────────────────────────

func (h *CMSHandler) ListFacilities(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListFacilities(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list facilities")
		return
	}
	if items == nil {
		items = []*models.CMSFacility{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateFacility(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = toSlug(req.Name)
	}
	f := &models.CMSFacility{
		ID: uuid.NewString(), Name: req.Name, Slug: slug,
		Description: req.Description, ImageURL: req.ImageURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateFacility(r.Context(), f); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A facility with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create facility")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

func (h *CMSHandler) UpdateFacility(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetFacilityByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Facility not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get facility")
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateFacility(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update facility")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteFacility(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteFacility(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete facility")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Associations ──────────────────────────────────────────────────

func (h *CMSHandler) ListAssociations(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListAssociations(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list associations")
		return
	}
	if items == nil {
		items = []*models.CMSAssociation{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateAssociation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		LogoURL     string `json:"logo_url"`
		WebsiteURL  string `json:"website_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = toSlug(req.Name)
	}
	a := &models.CMSAssociation{
		ID: uuid.NewString(), Name: req.Name, Slug: slug,
		Description: req.Description, LogoURL: req.LogoURL, WebsiteURL: req.WebsiteURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateAssociation(r.Context(), a); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "An association with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create association")
		return
	}
	response.JSON(w, http.StatusCreated, a)
}

func (h *CMSHandler) UpdateAssociation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetAssociationByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Association not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get association")
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		LogoURL     string `json:"logo_url"`
		WebsiteURL  string `json:"website_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.LogoURL != "" {
		existing.LogoURL = req.LogoURL
	}
	if req.WebsiteURL != "" {
		existing.WebsiteURL = req.WebsiteURL
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateAssociation(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update association")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteAssociation(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteAssociation(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete association")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Invest with Us ────────────────────────────────────────────────

func (h *CMSHandler) ListInvest(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.ListInvest(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list invest items")
		return
	}
	if items == nil {
		items = []*models.CMSInvest{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateInvest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		Content   string `json:"content"`
		ImageURL  string `json:"image_url"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	inv := &models.CMSInvest{
		ID: uuid.NewString(), Title: req.Title, Subtitle: req.Subtitle,
		Content: req.Content, ImageURL: req.ImageURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateInvest(r.Context(), inv); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create invest item")
		return
	}
	response.JSON(w, http.StatusCreated, inv)
}

func (h *CMSHandler) UpdateInvest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetInvestByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Invest item not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get invest item")
		return
	}
	var req struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		Content   string `json:"content"`
		ImageURL  string `json:"image_url"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Subtitle != "" {
		existing.Subtitle = req.Subtitle
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateInvest(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update invest item")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteInvest(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteInvest(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete invest item")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Team Members ──────────────────────────────────────────────────

func (h *CMSHandler) ListTeam(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListTeamMembers(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list team members")
		return
	}
	if items == nil {
		items = []*models.CMSTeamMember{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateTeamMember(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName    string `json:"full_name"`
		Designation string `json:"designation"`
		ImageURL    string `json:"image_url"`
		Bio         string `json:"bio"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.FullName == "" {
		response.ValidationErr(w, map[string]string{"full_name": "required"})
		return
	}
	m := &models.CMSTeamMember{
		ID: uuid.NewString(), FullName: req.FullName, Designation: req.Designation,
		ImageURL: req.ImageURL, Bio: req.Bio, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateTeamMember(r.Context(), m); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create team member")
		return
	}
	response.JSON(w, http.StatusCreated, m)
}

func (h *CMSHandler) UpdateTeamMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetTeamMemberByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Team member not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get team member")
		return
	}
	var req struct {
		FullName    string `json:"full_name"`
		Designation string `json:"designation"`
		ImageURL    string `json:"image_url"`
		Bio         string `json:"bio"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	existing.Designation = req.Designation
	existing.ImageURL = req.ImageURL
	existing.Bio = req.Bio
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateTeamMember(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update team member")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteTeamMember(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteTeamMember(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete team member")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}
