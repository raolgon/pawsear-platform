package httpapi

import (
	"database/sql"
	"net/http"
	"time"
)

func NewRouter(database *sql.DB) http.Handler {
	mux := http.NewServeMux()
	households := newHouseholdHandler(database)
	mvp := newMVPHandler(database)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := database.PingContext(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /api/meta", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"name":      "pawsear-api",
			"mode":      "local-first",
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		})
	})

	mux.HandleFunc("GET /api/households", households.list)
	mux.HandleFunc("POST /api/households", households.create)
	mux.HandleFunc("GET /api/households/{id}", households.get)
	mux.HandleFunc("PATCH /api/households/{id}", households.update)
	mux.HandleFunc("GET /api/households/{id}/contacts", mvp.listHouseholdContacts)
	mux.HandleFunc("POST /api/households/{id}/contacts", mvp.linkHouseholdContact)

	mux.HandleFunc("GET /api/contacts", mvp.listContacts)
	mux.HandleFunc("POST /api/contacts", mvp.createContact)
	mux.HandleFunc("GET /api/contacts/{id}", mvp.getContact)
	mux.HandleFunc("PATCH /api/contacts/{id}", mvp.updateContact)

	mux.HandleFunc("GET /api/pets", mvp.listPets)
	mux.HandleFunc("POST /api/pets", mvp.createPet)
	mux.HandleFunc("GET /api/pets/{id}", mvp.getPet)
	mux.HandleFunc("PATCH /api/pets/{id}", mvp.updatePet)

	mux.HandleFunc("GET /api/staff", mvp.listStaff)
	mux.HandleFunc("POST /api/staff", mvp.createStaff)
	mux.HandleFunc("GET /api/staff/{id}", mvp.getStaff)
	mux.HandleFunc("PATCH /api/staff/{id}", mvp.updateStaff)

	mux.HandleFunc("GET /api/bookings", mvp.listBookings)
	mux.HandleFunc("POST /api/bookings", mvp.createBooking)
	mux.HandleFunc("GET /api/bookings/{id}", mvp.getBooking)
	mux.HandleFunc("PATCH /api/bookings/{id}", mvp.updateBooking)

	mux.HandleFunc("GET /api/care-tasks", mvp.listCareTasks)
	mux.HandleFunc("POST /api/care-tasks", mvp.createCareTask)
	mux.HandleFunc("GET /api/care-tasks/{id}", mvp.getCareTask)
	mux.HandleFunc("PATCH /api/care-tasks/{id}", mvp.updateCareTask)

	mux.HandleFunc("GET /api/charges", mvp.listCharges)
	mux.HandleFunc("POST /api/charges", mvp.createCharge)
	mux.HandleFunc("GET /api/charges/{id}", mvp.getCharge)
	mux.HandleFunc("PATCH /api/charges/{id}", mvp.updateCharge)

	mux.HandleFunc("GET /api/payments", mvp.listPayments)
	mux.HandleFunc("POST /api/payments", mvp.createPayment)
	mux.HandleFunc("GET /api/payments/{id}", mvp.getPayment)

	mux.HandleFunc("GET /api/dashboard/today", mvp.dashboardToday)

	return mux
}
