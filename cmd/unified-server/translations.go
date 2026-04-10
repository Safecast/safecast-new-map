package main

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

// =====================
// Транслейт
// =====================
var translations map[string]map[string]string
var translationsMu sync.RWMutex

// loadTranslationsFromDB loads all translations from the database into the in-memory map.
// Returns true if successful, false if the table doesn't exist or DB is unavailable.
func loadTranslationsFromDB() bool {
	if db == nil || db.DB == nil {
		return false
	}

	rows, err := db.DB.Query("SELECT language_code, key, value FROM translations")
	if err != nil {
		log.Printf("[i18n] Cannot load translations from DB: %v", err)
		return false
	}
	defer rows.Close()

	newTranslations := make(map[string]map[string]string)
	count := 0
	for rows.Next() {
		var lang, key, value string
		if err := rows.Scan(&lang, &key, &value); err != nil {
			log.Printf("[i18n] Error scanning translation row: %v", err)
			continue
		}
		if newTranslations[lang] == nil {
			newTranslations[lang] = make(map[string]string)
		}
		newTranslations[lang][key] = value
		count++
	}

	if count == 0 {
		return false
	}

	translationsMu.Lock()
	translations = newTranslations
	translationsMu.Unlock()

	log.Printf("[i18n] Loaded %d translations from database", count)
	return true
}

// loadTranslationsFromFile loads translations from the embedded JSON file (fallback).
func loadTranslationsFromFile(fs embed.FS, filename string) {
	file, err := fs.Open(filename)
	if err != nil {
		log.Fatalf("Error opening translation file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading translation file: %v", err)
	}

	err = json.Unmarshal(data, &translations)
	if err != nil {
		log.Fatalf("Error parsing translations: %v", err)
	}
	log.Printf("[i18n] Loaded translations from embedded file (fallback)")
}

// seedTranslationsDB seeds the database from the embedded JSON file if the table is empty.
func seedTranslationsDB(fs embed.FS, filename string) {
	if db == nil || db.DB == nil {
		return
	}

	// Load from file — always check for missing keys (not just empty table)
	file, err := fs.Open(filename)
	if err != nil {
		log.Printf("[i18n] Cannot open translation file for seeding: %v", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[i18n] Cannot read translation file for seeding: %v", err)
		return
	}

	var fileTranslations map[string]map[string]string
	if err := json.Unmarshal(data, &fileTranslations); err != nil {
		log.Printf("[i18n] Cannot parse translation file for seeding: %v", err)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("[i18n] Cannot begin transaction for seeding: %v", err)
		return
	}

	inserted := 0
	for lang, keys := range fileTranslations {
		for key, value := range keys {
			_, err := tx.Exec(
				"INSERT INTO translations (language_code, key, value) VALUES ($1, $2, $3) ON CONFLICT (language_code, key) DO NOTHING",
				lang, key, value,
			)
			if err != nil {
				log.Printf("[i18n] Error seeding %s.%s: %v", lang, key, err)
				continue
			}
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[i18n] Error committing seed transaction: %v", err)
		return
	}

	if inserted > 0 {
		log.Printf("[i18n] Seeded %d new translations into database from embedded file", inserted)
	}
}

func getPreferredLanguage(r *http.Request) string {
	// Поддерживаемые языки (добавлены: da, fa)
	supported := map[string]struct{}{
		"en": {}, "zh": {}, "es": {}, "hi": {}, "ar": {}, "fr": {}, "ru": {}, "pt": {}, "de": {}, "ja": {}, "tr": {}, "it": {},
		"ko": {}, "pl": {}, "uk": {}, "mn": {}, "no": {}, "fi": {}, "ka": {}, "sv": {}, "he": {}, "nl": {}, "el": {}, "hu": {},
		"cs": {}, "ro": {}, "th": {}, "vi": {}, "id": {}, "ms": {}, "bg": {}, "lt": {}, "et": {}, "lv": {}, "sl": {},
		"da": {}, "fa": {},
	}

	// Check ?lang= query parameter first (explicit override)
	if langParam := r.URL.Query().Get("lang"); langParam != "" {
		code := strings.ToLower(strings.TrimSpace(langParam))
		if _, ok := supported[code]; ok {
			return code
		}
	}

	// Нормализация/синонимы: приводим варианты к поддерживаемым базовым кодам
	aliases := map[string]string{
		// Устаревшие коды
		"iw": "he", // he (Hebrew)
		"in": "id", // id (Indonesian)

		// Норвежский: часто приходит nb-NO/nn-NO
		"nb": "no",
		"nn": "no",

		// Китайский: сводим к "zh"
		"zh-cn":   "zh",
		"zh-sg":   "zh",
		"zh-hans": "zh",
		"zh-tw":   "zh",
		"zh-hk":   "zh",
		"zh-hant": "zh",

		// Португальский варианты → "pt"
		"pt-br": "pt",
		"pt-pt": "pt",
	}

	// Fall back to Accept-Language header
	langHeader := r.Header.Get("Accept-Language")
	if langHeader == "" {
		return "en"
	}

	langs := strings.Split(langHeader, ",")
	for _, raw := range langs {
		code := strings.TrimSpace(strings.SplitN(raw, ";", 2)[0])
		code = strings.ToLower(strings.ReplaceAll(code, "_", "-"))

		// Берём базовую часть до дефиса (например, "de" из "de-DE")
		base := code
		if i := strings.Index(code, "-"); i != -1 {
			base = code[:i]
		}

		// Применяем алиасы (и к полному коду, и к базе)
		if a, ok := aliases[code]; ok {
			base = a
		} else if a, ok := aliases[base]; ok {
			base = a
		}

		// Проверяем поддержку
		if _, ok := supported[base]; ok {
			return base
		}
	}

	return "en"
}

func translationsForLang(all map[string]map[string]string, lang string) map[string]map[string]string {
	filtered := make(map[string]map[string]string, 2)
	if en, ok := all["en"]; ok {
		filtered["en"] = en
	}
	if lang != "en" {
		if langMap, ok := all[lang]; ok {
			filtered[lang] = langMap
		}
	}
	return filtered
}
