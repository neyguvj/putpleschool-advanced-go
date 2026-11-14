package product

import (
	"errors"
	"net/http"
	"strconv"

	"advancedgo/pkg/request"
	"advancedgo/pkg/response"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo   Repository
	logger logrus.FieldLogger
}

func NewHandler(repo Repository, logger logrus.FieldLogger) *Handler {
	return &Handler{repo: repo, logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("GET /products/{id}", h.Get)
	mux.HandleFunc("PUT /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)
}

type upsertRequest struct {
	Name        string          `json:"name" validate:"required,max=255"`
	Description string          `json:"description" validate:"required"`
	Price       decimal.Decimal `json:"price"`
	Images      []string        `json:"images" validate:"omitempty,dive,url"`
}

func (r upsertRequest) toProduct() *Product {
	images := r.Images
	if images == nil {
		images = []string{}
	}
	return &Product{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Images:      pq.StringArray(images),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := request.DecodeAndValidate[upsertRequest](w, r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	p := req.toProduct()
	if err := p.Validate(); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.Create(r.Context(), p); err != nil {
		h.writeError(w, err)
		return
	}

	response.OkResponse(w, http.StatusCreated, p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.OkResponse(w, http.StatusOK, p)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	req, err := request.DecodeAndValidate[upsertRequest](w, r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	p := req.toProduct()
	p.ID = id
	if err := p.Validate(); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.Update(r.Context(), p); err != nil {
		h.writeError(w, err)
		return
	}

	p, err = h.repo.GetByID(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.OkResponse(w, http.StatusOK, p)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		h.writeError(w, err)
		return
	}

	response.OkResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.ErrorResponse(w, http.StatusNotFound, err)
	default:
		h.logger.WithError(err).Error("unexpected error")
		response.ErrorResponse(w, http.StatusInternalServerError, nil)
	}
}

func idParam(r *http.Request) (uint, error) {
	raw := r.PathValue("id")
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || v == 0 {
		return 0, errors.New("invalid product id")
	}
	return uint(v), nil
}
