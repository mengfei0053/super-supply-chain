package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL = "http://localhost:8081"
	defaultTimeout = 10 * time.Second
	configDirName  = "super-supply-chain"
	configFileName = "cli.json"
)

type config struct {
	BaseURL   string    `json:"baseUrl"`
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
	SavedAt   time.Time `json:"savedAt"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Avatar   string `json:"avatar"`
}

type claimsPayload struct {
	Username string `json:"username"`
	Expires  int64  `json:"exp"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "login":
		return runLogin(args[1:])
	case "status":
		return runStatus(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runLogin(args []string) error {
	flags := flag.NewFlagSet("login", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	baseURL := flags.String("base-url", envOrDefault("SSC_BASE_URL", defaultBaseURL), "backend base URL")
	username := flags.String("username", os.Getenv("SSC_USERNAME"), "login username")
	password := flags.String("password", os.Getenv("SSC_PASSWORD"), "login password")
	configPath := flags.String("config", defaultConfigPath(), "config file path")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if *username == "" {
		return errors.New("username is required; pass --username or set SSC_USERNAME")
	}
	if *password == "" {
		return errors.New("password is required; pass --password or set SSC_PASSWORD")
	}

	backendURL, err := normalizeBaseURL(*baseURL)
	if err != nil {
		return err
	}

	loginResult, err := login(backendURL, *username, *password)
	if err != nil {
		return err
	}

	expiresAt, _ := tokenExpiresAt(loginResult.Token)
	savedConfig := config{
		BaseURL:   backendURL,
		Token:     loginResult.Token,
		Username:  loginResult.Username,
		ExpiresAt: expiresAt,
		SavedAt:   time.Now(),
	}
	if err := saveConfig(*configPath, savedConfig); err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", loginResult.Username)
	fmt.Printf("Config saved to %s\n", *configPath)
	if !expiresAt.IsZero() {
		fmt.Printf("Token expires at %s\n", expiresAt.Format(time.RFC3339))
	}
	return nil
}

func runStatus(args []string) error {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	configPath := flags.String("config", defaultConfigPath(), "config file path")
	checkRemote := flags.Bool("remote", true, "call an authenticated backend endpoint")

	if err := flags.Parse(args); err != nil {
		return err
	}

	storedConfig, err := loadConfig(*configPath)
	if err != nil {
		return err
	}
	if storedConfig.Token == "" {
		return errors.New("not logged in; run `ssc login` first")
	}

	fmt.Printf("Config: %s\n", *configPath)
	fmt.Printf("Base URL: %s\n", storedConfig.BaseURL)
	if storedConfig.Username != "" {
		fmt.Printf("User: %s\n", storedConfig.Username)
	}

	expiresAt := storedConfig.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt, _ = tokenExpiresAt(storedConfig.Token)
	}
	if !expiresAt.IsZero() {
		fmt.Printf("Token expires at: %s\n", expiresAt.Format(time.RFC3339))
		if time.Now().After(expiresAt) {
			return errors.New("token is expired; run `ssc login` again")
		}
	}

	if *checkRemote {
		if err := checkAuthenticatedEndpoint(storedConfig.BaseURL, storedConfig.Token); err != nil {
			return err
		}
		fmt.Println("Backend auth check: OK")
	}

	fmt.Println("Status: logged in")
	return nil
}

func login(baseURL, username, password string) (loginResponse, error) {
	payload, err := json.Marshal(loginRequest{Username: username, Password: password})
	if err != nil {
		return loginResponse{}, err
	}

	endpoint := strings.TrimRight(baseURL, "/") + "/api/login"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return loginResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient().Do(req)
	if err != nil {
		return loginResponse{}, fmt.Errorf("login request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return loginResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return loginResponse{}, fmt.Errorf("login failed with HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var loginResult loginResponse
	if err := json.Unmarshal(body, &loginResult); err != nil {
		return loginResponse{}, fmt.Errorf("decode login response: %w", err)
	}
	if loginResult.Token == "" {
		return loginResponse{}, errors.New("login response did not include token")
	}
	return loginResult, nil
}

func checkAuthenticatedEndpoint(baseURL, token string) error {
	endpoint := strings.TrimRight(baseURL, "/") + "/api/admin/menus"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("backend status request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(res.Body)
	return fmt.Errorf("backend auth check failed with HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
}

func saveConfig(path string, savedConfig config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(savedConfig, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	return os.WriteFile(path, payload, 0o600)
}

func loadConfig(path string) (config, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config{}, errors.New("not logged in; run `ssc login` first")
		}
		return config{}, err
	}

	var storedConfig config
	if err := json.Unmarshal(payload, &storedConfig); err != nil {
		return config{}, fmt.Errorf("decode config %s: %w", path, err)
	}
	return storedConfig, nil
}

func tokenExpiresAt(token string) (time.Time, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, errors.New("token is not a JWT")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, err
	}

	var claims claimsPayload
	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, err
	}
	if claims.Expires == 0 {
		return time.Time{}, nil
	}
	return time.Unix(claims.Expires, 0), nil
}

func normalizeBaseURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("base URL is required")
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", errors.New("base URL must start with http:// or https://")
	}
	if parsedURL.Host == "" {
		return "", errors.New("base URL must include a host")
	}
	return strings.TrimRight(parsedURL.String(), "/"), nil
}

func defaultConfigPath() string {
	if value := os.Getenv("SSC_CONFIG"); value != "" {
		return value
	}
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", configFileName)
	}
	return filepath.Join(configRoot, configDirName, configFileName)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func httpClient() *http.Client {
	return &http.Client{Timeout: defaultTimeout}
}

func printUsage() {
	fmt.Print(`Super Supply Chain CLI

Usage:
  ssc login --base-url http://localhost:8081 --username USER --password PASS
  ssc status [--remote=true]

Commands:
  login   Call POST /api/login and save the returned JWT locally
  status  Show saved login state and verify auth via GET /api/admin/menus

Environment:
  SSC_BASE_URL   Backend base URL, default http://localhost:8081
  SSC_USERNAME   Login username
  SSC_PASSWORD   Login password
  SSC_CONFIG     Config path, default OS user config dir
`)
}
