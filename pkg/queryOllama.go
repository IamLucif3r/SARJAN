package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kaptinlin/jsonrepair"
)

func QueryOllamaScoreMap(prompt string) (map[string]float64, error) {

	log.Println("[Debug] Sending prompt here: ", prompt)
	OllamaAPIURL := os.Getenv("OLLAMA_URL") + "/api/generate"
	payload := map[string]any{
		"model":  os.Getenv("LLM_MODEL"),
		"prompt": prompt,
		"stream": false,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(OllamaAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData map[string]any
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var completion string
	if c, ok := responseData["completion"].(string); ok {
		completion = c
	} else if c, ok := responseData["response"].(string); ok {
		completion = c
	} else {
		return nil, fmt.Errorf("could not find completion text in response")
	}

	repaired, err := jsonrepair.JSONRepair(completion)
	if err != nil {
		log.Fatalf("Failed to repair JSON: %v", err)
	}
	var result map[string]float64
	if err := json.Unmarshal([]byte(repaired), &result); err != nil {
		return nil, fmt.Errorf("failed to parse extracted JSON: %v", err)
	}

	return result, nil
}
