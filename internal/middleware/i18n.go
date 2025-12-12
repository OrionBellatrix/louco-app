package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/i18n"
)

func I18n(i18nService *i18n.I18n) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract language from Accept-Language header
		acceptLanguage := c.GetHeader("Accept-Language")
		lang := i18n.ExtractLanguageFromHeader(acceptLanguage)

		// Check if language is supported, fallback to English if not
		if !i18nService.IsLanguageSupported(lang) {
			lang = "en"
		}

		// Set language in context
		c.Set("language", lang)
		c.Set("i18n", i18nService)

		c.Next()
	}
}

// GetLanguage extracts language from context
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get("language"); exists {
		if language, ok := lang.(string); ok {
			return language
		}
	}
	return "en"
}

// Translate translates a key using the language from context
func Translate(c *gin.Context, key string) string {
	lang := GetLanguage(c)
	if i18nService, exists := c.Get("i18n"); exists {
		if service, ok := i18nService.(*i18n.I18n); ok {
			return service.Translate(lang, key)
		}
	}
	return key
}
