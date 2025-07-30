package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CEPRequest struct {
	Cep string `json:"cep"`
}

func validateCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	if !strings.HasPrefix(cep, "0") && !strings.HasPrefix(cep, "1") && !strings.HasPrefix(cep, "2") && !strings.HasPrefix(cep, "3") && !strings.HasPrefix(cep, "4") && !strings.HasPrefix(cep, "5") && !strings.HasPrefix(cep, "6") && !strings.HasPrefix(cep, "7") && !strings.HasPrefix(cep, "8") && !strings.HasPrefix(cep, "9") {
		return false
	}
	for _, c := range cep {
		if !('0' <= c && c <= '9') {
			return false
		}
	}
	return true
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	var req CEPRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if !validateCEP(req.Cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `"invalid zipcode"`)
		return
	}

	// Forward to Servico B (example URL: http://localhost:8080/api/cep)
	resp, err := http.Post("http://localhost:8080/api/cep", "application/json", strings.NewReader(fmt.Sprintf(`{"cep": "%s"}`, req.Cep)))
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyResp, _ := ioutil.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, string(bodyResp))
}

func main() {
	http.HandleFunc("/api/cep", handleCEP)
	fmt.Println("Servico A is running on :8081")
	http.ListenAndServe(":8081", nil)
}
