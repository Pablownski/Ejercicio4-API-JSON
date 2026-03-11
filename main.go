package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Player struct {
	ID 	 int    `json:"id"`
	Name string `json:"name"`
	Nationality string `json:"nationality"`
	Position string `json:"position"`
	CurrentTeam string `json:"current_team"`
	Age int `json:"age"`
	CareerGoals int `json:"career_goals"`
	Active bool `json:"active"`
}

type Message struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
    Code    int    `json:"code"`
    Error   string `json:"error"`
    Message string `json:"message"`
}



var players []Player

func main() {
	loadPlayers()

	http.HandleFunc("/api/players", playersHandler)
	http.HandleFunc("/api/players/", playerByIDHandler)
	log.Println("API running on port 24374")
	log.Fatal(http.ListenAndServe(":24374", nil))
}

func loadPlayers() {
	file, err := os.ReadFile("./data/players.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	err = json.Unmarshal(file, &players)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}
}

func savePlayers() {
	data, err := json.MarshalIndent(players, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("./data/players.json", data, 0644)
	if err != nil {
		log.Println("Error writing file:", err)
	}
}

func playersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetPlayers(w, r)
	case http.MethodPost:
		handleCreatePlayer(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetPlayers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	result := players

	
	if val := query.Get("nationality"); val != "" {
		filtered := []Player{}
		for _, p := range result {
			if strings.EqualFold(p.Nationality, val) {
				filtered = append(filtered, p)
			}
		}
		result = filtered
	}
	
	if val := query.Get("position"); val != "" {
		filtered := []Player{}
		for _, p := range result {
			if strings.EqualFold(p.Position, val) {
				filtered = append(filtered, p)
			}
		}
		result = filtered
	}

	if val := query.Get("active"); val != "" {
		activeBool, err := strconv.ParseBool(val)
		if err != nil {
			http.Error(w, "Invalid 'active' parameter: must be true or false", http.StatusBadRequest)
			return
		}
		filtered := []Player{}
		for _, p := range result {
			if p.Active == activeBool {
				filtered = append(filtered, p)
			}
		}
		result = filtered
	}

	if val := query.Get("min_goals"); val != "" {
		min, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Invalid 'min_goals' parameter: must be an integer", http.StatusBadRequest)
			return
		}
		filtered := []Player{}
		for _, p := range result {
			if p.CareerGoals >= min {
				filtered = append(filtered, p)
			}
		}
		result = filtered
	}

	if val := query.Get("max_goals"); val != "" {
		max, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Invalid 'max_goals' parameter: must be an integer", http.StatusBadRequest)
			return
		}
		filtered := []Player{}
		for _, p := range result {
			if p.CareerGoals <= max {
				filtered = append(filtered, p)
			}
		}
		result = filtered
	}

	if val := query.Get("search"); val != "" {
		filtered := []Player{}
		for _, p := range result {
			if strings.Contains(strings.ToLower(p.Name), strings.ToLower(val)) {
				filtered = append(filtered, p)
			}
		}
		result = filtered
	}

	writeJSON(w, http.StatusOK, result)
}


func handleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	var newPlayer Player

	err := json.NewDecoder(r.Body).Decode(&newPlayer)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if newPlayer.Name == "" {
		http.Error(w, "Field 'name' is required", http.StatusBadRequest)
		return
	}
	if newPlayer.Nationality == "" {
		http.Error(w, "Field 'nationality' is required", http.StatusBadRequest)
		return
	}
	if newPlayer.Position == "" {
		http.Error(w, "Field 'position' is required", http.StatusBadRequest)
		return
	}
	if newPlayer.CurrentTeam == "" {
		http.Error(w, "Field 'current_team' is required", http.StatusBadRequest)
		return
	}
	if newPlayer.Age <= 0 || newPlayer.Age > 100 {
		http.Error(w, "Field 'age' must be between 1 and 100", http.StatusBadRequest)
		return
	}
	if newPlayer.CareerGoals < 0 {
		http.Error(w, "Field 'career_goals' must be a non-negative number", http.StatusBadRequest)
		return
	}
		newPlayer.ID = generateNextID()
	players = append(players, newPlayer)
	savePlayers()

	writeJSON(w, http.StatusCreated, newPlayer)
}

func playerByIDHandler(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/players/")

	if idStr == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid player ID: must be an integer", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetPlayerByID(w, r, id)
	case http.MethodPut:
		handleUpdatePlayer(w, r, id)
	case http.MethodDelete:
		handleDeletePlayer(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetPlayerByID(w http.ResponseWriter, r *http.Request, id int) {
	for _, p := range players {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	http.Error(w, "Player not found", http.StatusNotFound)
}

func handleUpdatePlayer(w http.ResponseWriter, r *http.Request, id int) {
	var updated Player

	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

		if updated.Name == "" {
		http.Error(w, "Field 'name' is required", http.StatusBadRequest)
		return
	}
	if updated.Nationality == "" {
		http.Error(w, "Field 'nationality' is required", http.StatusBadRequest)
		return
	}
	if updated.Position == "" {
		http.Error(w, "Field 'position' is required", http.StatusBadRequest)
		return
	}
	if updated.CurrentTeam == "" {
		http.Error(w, "Field 'current_team' is required", http.StatusBadRequest)
		return
	}
	if updated.Age <= 0 || updated.Age > 100 {
		http.Error(w, "Field 'age' must be between 1 and 100", http.StatusBadRequest)
		return
	}
	if updated.CareerGoals < 0 {
		http.Error(w, "Field 'career_goals' must be a non-negative number", http.StatusBadRequest)
		return
	}

	for i, p := range players {
		if p.ID == id {
			updated.ID = id
			players[i] = updated
			savePlayers()
			writeJSON(w, http.StatusOK, updated)
			return
		}
	}

	http.Error(w, "Player not found", http.StatusNotFound)
}

func handleDeletePlayer(w http.ResponseWriter, r *http.Request, id int) {
	for i, p := range players {
		if p.ID == id {
			players = append(players[:i], players[i+1:]...)
			savePlayers()
			writeJSON(w, http.StatusOK, Message{Message: "Player deleted successfully"})
			return
		}
	}
	http.Error(w, "Player not found", http.StatusNotFound)
}



func generateNextID() int {
	maxID := 0
	for _, p := range players {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return maxID + 1
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorMessage{
        Code:    status,
        Error:   http.StatusText(status),
        Message: message,
    })
}