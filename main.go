package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// вхідні дані
type Data struct {
	Power       float64 `json:"power"`
	Electricity float64 `json:"electricity"`
	Deviation1  float64 `json:"deviation1"`
	Deviation2  float64 `json:"deviation2"`
}

// результати розрахунків
type CalculationResults struct {
	ProfitBefore float64 `json:"profitBefore"`
	ProfitAfter  float64 `json:"profitAfter"`
}

// функція нормального закону розподілу потужності (формула 9.1)
func normalDistribution(x, power, sigma float64) float64 {
	return (1 / (sigma * math.Sqrt(2*math.Pi))) *
		math.Exp(-(math.Pow(x-power, 2)) / (2*math.Pow(sigma, 2)))
}

// інтегрування
func integrate(
	a float64, // нижня межа
	b float64, // верхня межа
	n int, // кількість точок для інтегрування
	power float64,
	sigma float64,
) float64 {
	h := (b - a) / float64(n)
	sum := (normalDistribution(a, power, sigma) +
		normalDistribution(b, power, sigma)) / 2

	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		sum += normalDistribution(x, power, sigma)
	}

	return h * sum
}

func calculateEnergyWithoutImbalance(
	power float64,
	sigma float64,
	lowerBound float64,
	upperBound float64,
) float64 {
	return integrate(
		lowerBound,
		upperBound,
		100000, // кількість точок для інтегрування
		power,
		sigma,
	) * 100 // переводимо у відсотки
}

func calculateResults(data Data) CalculationResults {
	// діапазони
	lowerBound := 4.75
	upperBound := 5.25

	// розрахунок частки енергії без небалансів до покращення (δW1)
	energyWithoutImbalance1 := math.Round(calculateEnergyWithoutImbalance(
		data.Power,
		data.Deviation1,
		lowerBound,
		upperBound,
	))

	// розрахунок частки енергії без небалансів після покращення (δW2)
	energyWithoutImbalance2 := math.Round(calculateEnergyWithoutImbalance(
		data.Power,
		data.Deviation2,
		lowerBound,
		upperBound,
	))

	// енергія W1
	energy1 := data.Power * 24 * energyWithoutImbalance1 / 100

	// прибуток П1
	profit1 := energy1 * data.Electricity

	// енергія W2
	energy2 := data.Power * 24 * (1 - energyWithoutImbalance1/100)

	// штраф Ш1
	fine1 := energy2 * data.Electricity

	// загальний прибуток перед покращенням
	profitBefore := profit1 - fine1

	// енергія W3
	energy3 := data.Power * 24 * energyWithoutImbalance2 / 100

	// прибуток П2
	profit2 := energy3 * data.Electricity

	// енергія W4
	energy4 := data.Power * 24 * (1 - energyWithoutImbalance2/100)

	// штраф Ш2
	fine2 := energy4 * data.Electricity

	// загальний прибуток після покращення
	profitAfter := profit2 - fine2

	return CalculationResults{
		ProfitBefore: profitBefore,
		ProfitAfter:  profitAfter,
	}
}

func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results := calculateResults(data)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// статичні файли
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// маршрути для головної сторінки
	http.HandleFunc("/calculator", calculatorHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}