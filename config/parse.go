package config

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func lookupEnv(env string) (string, bool, error) {
	if raw, ok := os.LookupEnv(env); ok {
		return raw, true, nil
	}
	path, ok := os.LookupEnv(env + "_FILE")
	if !ok {
		return "", false, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, fmt.Errorf("read file for %s_FILE (%s): %w", env, path, err)
	}
	return strings.TrimRight(string(data), "\r\n"), true, nil
}

func parseString(target *string, env string) error {
	raw, ok, err := lookupEnv(env)
	if err != nil {
		return err
	}
	if ok {
		*target = raw
	}
	return nil
}

func parseInt(target *int, env string) error {
	raw, ok, err := lookupEnv(env)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("invalid int for %s (%q): %w", env, raw, err)
	}
	*target = n
	return nil
}

func parseBool(target *bool, env string) error {
	raw, ok, err := lookupEnv(env)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	b, err := strconv.ParseBool(raw)
	if err != nil {
		return fmt.Errorf("invalid bool for %s (%q): %w", env, raw, err)
	}
	*target = b
	return nil
}

func parseList(target *[]string, env string) error {
	raw, ok, err := lookupEnv(env)
	if err != nil {
		return err
	}
	if !ok || raw == "" {
		return nil
	}
	reader := csv.NewReader(strings.NewReader(raw))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	record, err := reader.Read()
	if err != nil {
		return fmt.Errorf("invalid CSV for %s (%q): %w", env, raw, err)
	}
	*target = record
	return nil
}

func parseMap(target *map[string]string, env string) error {
	raw, ok, err := lookupEnv(env)
	if err != nil {
		return err
	}
	if !ok || raw == "" {
		return nil
	}
	out := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return fmt.Errorf("invalid JSON for %s: %w", env, err)
	}
	*target = out
	return nil
}

func parseLogLevel(target *LogLevel, env string) error {
	raw, ok, err := lookupEnv(env)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return target.Decode(raw)
}
