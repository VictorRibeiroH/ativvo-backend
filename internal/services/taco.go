package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TACOFood struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Energy      float64 `json:"energy_kcal"`
	Protein     float64 `json:"protein_g"`
	Lipid       float64 `json:"lipid_g"`
	Carb        float64 `json:"carbohydrate_g"`
	Fiber       float64 `json:"fiber_g"`
	Base        float64 `json:"base_qty"`
}

type TACOResponse struct {
	Foods []TACOFood `json:"foods"`
	Total int        `json:"total"`
}

const (
	TACO_BASE_URL = "https://api-taco.herokuapp.com/api"
)

func SearchTACOFoods(query string) ([]TACOFood, error) {
	if query == "" {
		return []TACOFood{}, nil
	}

	// Encode query para URL
	encodedQuery := url.QueryEscape(strings.ToLower(query))
	apiURL := fmt.Sprintf("%s/food/search?q=%s", TACO_BASE_URL, encodedQuery)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar na API TACO: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API TACO retornou status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da API: %v", err)
	}

	var tacoResp TACOResponse
	if err := json.Unmarshal(body, &tacoResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %v", err)
	}

	return tacoResp.Foods, nil
}

func GetTACOFoodByID(id int) (*TACOFood, error) {
	apiURL := fmt.Sprintf("%s/food/%d", TACO_BASE_URL, id)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar alimento TACO: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("alimento não encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API TACO retornou status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var food TACOFood
	if err := json.Unmarshal(body, &food); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse: %v", err)
	}

	return &food, nil
}

func ConvertTACOFoodToLocal(tacoFood TACOFood) map[string]interface{} {
	return map[string]interface{}{
		"name":         tacoFood.Description,
		"calories":     tacoFood.Energy,
		"protein":      tacoFood.Protein,
		"carbs":        tacoFood.Carb,
		"fat":          tacoFood.Lipid,
		"fiber":        tacoFood.Fiber,
		"serving_size": 100,
		"category":     tacoFood.Category,
		"taco_id":      tacoFood.ID,
		"source":       "TACO",
	}
}
