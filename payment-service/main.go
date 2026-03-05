package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Payment struct {
	ID        int       `json:"id"`
	Amount    float32   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	payments = make(map[int]Payment)
	mutex    sync.Mutex
	nextID   = 1
)

func createPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p Payment
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	p.ID = nextID
	nextID++
	p.Status = "pending"
	p.CreatedAt = time.Now()
	payments[p.ID] = p

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func getPayment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	p, exists := payments[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	p, exists := payments[id]
	if !exists {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	p.Status = status
	p.UpdatedAt = time.Now()

	payments[id] = p
	json.NewEncoder(w).Encode(p)

}

func deletePayment(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	_, exists := payments[id]
	if !exists {
		http.Error(w, "Paymnet not found", http.StatusNotFound)
		return
	}

	delete(payments, id)

	w.Write([]byte("Payment Deleted"))
}

func listPayments(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var list []Payment
	for _, p := range payments {
		list = append(list, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/payments/create", createPayment)
	mux.HandleFunc("/payments/get", getPayment)
	mux.HandleFunc("/payments/list", listPayments)
	mux.HandleFunc("/payments/delete", deletePayment)
	mux.HandleFunc("/payments/update-status", updatePayment)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to payment - service"))
	})

	loggedMux := loggingMiddleware(mux)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
