package hardware

import (
	"encoding/json"
	"log"
	"net/http"
)

// AllServicesUp checks if all services returned a 200 response
func AllServicesUp() (bool, error) {
	var err error

	req, err := http.NewRequest("GET", "https://monit.depado.eu/api/status", nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not fetch monit status")
		return false, err
	}
	defer resp.Body.Close()
	var rsp map[string]int
	if err = json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		log.Println("Could not decode response")
		return false, err
	}
	for _, s := range rsp {
		if s != 200 && s != 0 {
			return false, nil
		}
	}
	return true, nil
}
