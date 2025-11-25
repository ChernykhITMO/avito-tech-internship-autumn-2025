package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const url = "http://localhost:8080"

func TestE2E_AddTeam(t *testing.T) {
	t.Helper()

	body := map[string]any{
		"team_name": "team_e2e_001",
		"members": []map[string]any{
			{"user_id": "123", "username": "Alice", "is_active": true},
			{"user_id": "124", "username": "Bob", "is_active": true},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url+"/team/add", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("http request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestE2E_AddTeam_AlreadyExists(t *testing.T) {
	t.Helper()

	body := map[string]any{
		"team_name": "team_e2e_001",
		"members": []map[string]any{
			{"user_id": "123", "username": "Alice", "is_active": true},
			{"user_id": "124", "username": "Bob", "is_active": true},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}

	client := &http.Client{}

	req1, _ := http.NewRequest(http.MethodPost, url+"/team/add", bytes.NewReader(data))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := client.Do(req1)

	if err != nil {
		t.Fatalf("first request: %v", err)
	}

	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 on first creation, got %d", resp1.StatusCode)
	}

	req2, _ := http.NewRequest(http.MethodPost, url+"/team/add", bytes.NewReader(data))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("second request: %v", err)
	}

	if resp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 on duplicate team, got %d", resp2.StatusCode)
	}

	var errResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp2.Body).Decode(&errResp); err != nil {
		t.Fatalf("json decode: %v", err)
	}

	if errResp.Error.Code != "TEAM_EXISTS" {
		t.Fatalf("expected error code TEAM_EXISTS, got %s", errResp.Error.Code)
	}
}
