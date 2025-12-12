package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type I18n struct {
	translations map[string]map[string]string
	fallback     string
}

func New(localesPath string, fallback string) (*I18n, error) {
	i18n := &I18n{
		translations: make(map[string]map[string]string),
		fallback:     fallback,
	}

	err := i18n.loadTranslations(localesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	return i18n, nil
}

func (i *I18n) loadTranslations(localesPath string) error {
	files, err := os.ReadDir(localesPath)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		lang := strings.TrimSuffix(file.Name(), ".json")
		filePath := filepath.Join(localesPath, file.Name())

		translations, err := i.loadTranslationFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to load translation file %s: %w", filePath, err)
		}

		i.translations[lang] = translations
	}

	return nil
}

func (i *I18n) loadTranslationFile(filePath string) (map[string]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var translations map[string]string
	err = json.Unmarshal(data, &translations)
	if err != nil {
		return nil, err
	}

	return translations, nil
}

func (i *I18n) Translate(lang, key string) string {
	// Try to get translation for the requested language
	if translations, exists := i.translations[lang]; exists {
		if translation, exists := translations[key]; exists {
			return translation
		}
	}

	// Fallback to default language
	if lang != i.fallback {
		if translations, exists := i.translations[i.fallback]; exists {
			if translation, exists := translations[key]; exists {
				return translation
			}
		}
	}

	// Return the key itself if no translation found
	return key
}

func (i *I18n) GetSupportedLanguages() []string {
	var languages []string
	for lang := range i.translations {
		languages = append(languages, lang)
	}
	return languages
}

func (i *I18n) IsLanguageSupported(lang string) bool {
	_, exists := i.translations[lang]
	return exists
}

// Helper function to extract language from Accept-Language header
func ExtractLanguageFromHeader(acceptLanguage string) string {
	if acceptLanguage == "" {
		return "en"
	}

	// Parse Accept-Language header (simplified)
	languages := strings.Split(acceptLanguage, ",")
	if len(languages) > 0 {
		lang := strings.TrimSpace(languages[0])
		// Remove quality factor if present (e.g., "en-US;q=0.9" -> "en-US")
		if idx := strings.Index(lang, ";"); idx != -1 {
			lang = lang[:idx]
		}
		// Extract primary language (e.g., "en-US" -> "en")
		if idx := strings.Index(lang, "-"); idx != -1 {
			lang = lang[:idx]
		}
		return strings.ToLower(lang)
	}

	return "en"
}
