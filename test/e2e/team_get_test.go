package e2e

import (
	"net/http"
	"testing"
)

func TestE2E_GetTeam(t *testing.T) {
	t.Helper()

	const teamName = "avito"

	resp, err := http.Get(url + "/team/get?team_name=" + teamName)
	if err != nil {
		t.Errorf("get team name: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("not get team name")
	}
}
