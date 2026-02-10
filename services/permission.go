package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type permissionResponse struct {
	Permissions []struct {
		Module      string   `json:"module"`
		Permissions []string `json:"permissions"`
	} `json:"permissions"`
}

func CheckImportPermission(token string, companyID int) error {
	apiURL := os.Getenv("AMIGO_API_URL")
	if apiURL == "" {
		return fmt.Errorf("AMIGO_API_URL not configured")
	}

	if apiURL == "IGNORE" {
		return nil
	}

	req, err := http.NewRequest("GET", apiURL+"/api/user/info", nil)
	if err != nil {
		return fmt.Errorf("failed to create permission request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("company-id", strconv.Itoa(companyID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check permissions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UNABLE_TO_IMPORT_LEADS")
	}

	var body permissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("failed to parse permission response: %w", err)
	}

	for _, p := range body.Permissions {
		if p.Module == "LEADS" {
			for _, perm := range p.Permissions {
				if perm == "IMPORT_LEADS" {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("UNABLE_TO_IMPORT_LEADS")
}
