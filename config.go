package htmx

import (
	"os"
	"strconv"
	"strings"
)

func parseConfig() *Config {
	return &Config{
		DefaultTemplates:   envAsSlice("DEFAULT_TEMPLATES", []string{"index.gohtml"}, ","),
		DefaultTemplatesHx: envAsSlice("DEFAULT_TEMPLATES_HX", []string{"hx/index.gohtml"}, ","),
		ServerAddress:      env("SERVER_ADDR", "localhost:8888"),
		TemplateDir:        env("TEMPLATE_DIR", "/templates"),
	}
}

func env(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}

func envAsInt(name string, def int) int {
	value := env(name, "")
	if val, err := strconv.Atoi(value); err == nil {
		return val
	}

	return def
}

func envAsBool(name string, def bool) bool {
	value := env(name, "")
	if val, err := strconv.ParseBool(value); err == nil {
		return val
	}

	return def
}

func envAsSlice(name string, def []string, sep string) []string {
	value := env(name, "")
	if value == "" {
		return def
	}

	return strings.Split(value, sep)
}
