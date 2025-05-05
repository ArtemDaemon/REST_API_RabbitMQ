package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"
)

type Item struct {
	IndicatorId     string `json:"indicator_id"`
	IndicatorValue  string `json:"indicator_value"`
	CountryId       string `json:"country_id"`
	CountryValue    string `json:"country_value"`
	CountryISO3Code string `json:"country_iso3_code"`
	Date            string `json:"date"`
	Value           uint   `json:"value"`
	Unit            string `json:"unit"`
	ObsStatus       string `json:"obs_status"`
	Decimal         uint   `json:"decimal"`
	CreatedAt       string `json:"created_at"`
}

// CountItemsHandler возвращает количество элементов в таблице items
// @Summary Получить количество элементов
// @Description Возвращает общее число записей в базе данных
// @Tags items
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {string} string "Ошибка при запросе количества"
// @Security ApiKeyAuth
// @Router /api/count [get]
func CountItemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
		if err != nil {
			http.Error(w, "Ошибка при запросе количества", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]int{"count": count})
	}
}

// LastCreatedAtHandler возвращает дату последнего добавления записи
// @Summary Получить дату последнего элемента
// @Description Возвращает максимальное значение поля created_at
// @Tags items
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Ошибка при получении даты"
// @Security ApiKeyAuth
// @Router /api/last_created_at [get]
func LastCreatedAtHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ts sql.NullString
		err := db.QueryRow("SELECT MAX(created_at) FROM items").Scan(&ts)
		if err != nil {
			http.Error(w, "Ошибка при получении даты", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"last_created_at": ts.String,
		})
	}
}

// GetItemByDateHandler получает элемент по дате
// @Summary Получить элемент по дате
// @Description Возвращает первую запись с указанной датой
// @Tags items
// @Produce json
// @Param year query string true "Год (например, 2022)"
// @Success 200 {object} handlers.Item
// @Failure 400 {string} string "Параметр 'year' обязателен"
// @Failure 404 {string} string "Элемент не найден"
// @Failure 500 {string} string "Ошибка при запросе"
// @Security ApiKeyAuth
// @Router /api/get_item [get]
func GetItemByDateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year := r.URL.Query().Get("year")
		if year == "" {
			http.Error(w, "Параметр 'year' обязателен", http.StatusBadRequest)
			return
		}

		var item Item
		err := db.QueryRow(
			`SELECT indicator_id, indicator_value, country_id, country_value, country_iso3_code, date, value, unit, obs_status, decimal, created_at FROM items WHERE date = ?`, year,
		).Scan(&item.IndicatorId, &item.IndicatorValue, &item.CountryId, &item.CountryValue,
			&item.CountryISO3Code, &item.Date, &item.Value, &item.Unit, &item.ObsStatus, &item.Decimal, &item.CreatedAt)

		if err == sql.ErrNoRows {
			http.Error(w, "Элемент не найден", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Ошибка при запросе: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(item)
	}
}

// AddItemHandler добавляет новый элемент в базу данных
// @Summary Добавить элемент
// @Description Принимает JSON-структуру и добавляет её в базу
// @Tags items
// @Accept json
// @Produce json
// @Param item body handlers.Item true "Элемент данных"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Неверный JSON"
// @Failure 500 {string} string "Ошибка при добавлении записи"
// @Security ApiKeyAuth
// @Router /api/add_item [post]
func AddItemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO items (
				indicator_id, indicator_value, country_id, country_value,
				country_iso3_code, date, value, unit, obs_status, decimal
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = db.Exec(query,
			item.IndicatorId, item.IndicatorValue, item.CountryId, item.CountryValue,
			item.CountryISO3Code, item.Date, item.Value, item.Unit, item.ObsStatus,
			item.Decimal,
		)
		if err != nil {
			http.Error(w, "Ошибка при добавлении записи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Элемент добавлен"})
	}
}
