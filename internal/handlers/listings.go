package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/Rajit-Dutta/GolangMonolith/internal/httpx"
	"github.com/Rajit-Dutta/GolangMonolith/internal/middleware"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func ListingConstructor(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx, `
			SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100
		`)
	if err != nil {
		log.Printf("db.query: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.Error_internal_error)
		return
	}
	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			log.Printf("rows.scan: %v", err)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.Error_internal_error)
			return
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.err: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.Error_internal_error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request_id := middleware.RetrieveCTXIDFromContext(ctx)
	id := r.PathValue("id")
	_, err := lh.db.ExecContext(ctx, `DELETE FROM LISTING WHERE id = $1`, id)
	if err != nil {
		log.Printf("db.query: %v", err)
		lh.logger.Error("delete failed", "listing_id", id, "request_id", request_id, "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.Error_internal_error)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req listing
	ctx := r.Context()
	request_id := middleware.RetrieveCTXIDFromContext((ctx))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("r.body: %v", err)
		lh.logger.Error("fetching failed", "request_id", request_id, "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid body", httpx.Error_malformed_json)
		return
	}
	row := lh.db.QueryRowContext(ctx, `
	INSERT INTO listings (title,description,price,city) VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Title, req.Description, req.Price, req.City)

	var id string
	if err := row.Scan(&id); err != nil {
		lh.logger.Error("creation failed", "request_id", request_id, "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.Error_internal_error)
		return
	}
	lh.logger.Info("listing created", "listing_id", id, "request_id", request_id)
}
