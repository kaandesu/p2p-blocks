package network

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"p2pBlocks/blockchain"
)

func (s *Server) addBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if requestBody.Data == "" {
		http.Error(w, "Data field is required", http.StatusBadRequest)
		return
	}

	s.Blockchain.AddBlock(requestBody.Data)

	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Block added successfully"))
}

func (s *Server) getAllBlocks(w http.ResponseWriter, r *http.Request) {
	blocks := []blockchain.Block{}
	iter := s.Blockchain.Iterator()
	for {
		block := iter.Next()
		blocks = append(blocks, *block)

		if len(block.PrevHash) == 0 {
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blocks)
}

func (s *Server) StartHttpServer() {
	slog.Info("Starting the HTTP server on", "addr", "0.0.0.0:80")
	http.HandleFunc("/blocks/add", s.addBlock)
	http.HandleFunc("/blocks", s.getAllBlocks)

	slog.Info("HTTP server started on:", "ADDR", "0.0.0.0:80")
	http.ListenAndServe("0.0.0.0:80", nil)
}
