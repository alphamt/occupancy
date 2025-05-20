package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Device struct {
	Index  int     `json:"index"`
	UUID   string  `json:"uuid"`
	Memory float64 `json:"memory"` // used
	Total  float64 `json:"total"`
}

type Summary struct {
	TotalGPUs      int     `json:"total_gpus"`
	TotalUsedMem   float64 `json:"total_used_mem"`
	TotalFreeMem   float64 `json:"total_free_mem"`
	AverageFreeMem float64 `json:"avg_free_mem_per_gpu"`
}

var (
	mu    sync.RWMutex
	store = make(map[string][]Device)
)

func handleDump(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data map[string][]Device
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	mu.Lock()
	for host, devs := range data {
		store[host] = devs
	}
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

func handleSummary(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	out := make(map[string]Summary)
	for host, devs := range store {
		var used, total float64
		for _, d := range devs {
			used += d.Memory
			total += d.Total
		}
		free := total - used
		count := len(devs)
		avgFree := 0.0
		if count > 0 {
			avgFree = free / float64(count)
		}
		out[host] = Summary{
			TotalGPUs:      count,
			TotalUsedMem:   used,
			TotalFreeMem:   free,
			AverageFreeMem: avgFree,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func main() {
	http.HandleFunc("/dump", handleDump)       // POST JSON dump here
	http.HandleFunc("/devices", handleGetAll)  // GET full map
	http.HandleFunc("/summary", handleSummary) // GET per-host summary
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
