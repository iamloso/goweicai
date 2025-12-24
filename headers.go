package gowencai

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// UserAgents contains a list of common user agents
var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
}

// GetRandomUserAgent returns a random user agent string
func GetRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return UserAgents[rand.Intn(len(UserAgents))]
}

// GetToken executes the JavaScript using Node.js to generate the token
func GetToken() (string, error) {
	// Get the current executable directory or working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Look for hexin-v.js
	scriptPaths := []string{
		filepath.Join(dir, "pywencai", "hexin-v.js"),
		filepath.Join(dir, "hexin-v.js"),
		filepath.Join(dir, "..", "pywencai", "hexin-v.js"),
		filepath.Join(dir, "..", "..", "pywencai", "hexin-v.js"),
		"/home/administrator/workplace/gowencai/pywencai/hexin-v.js", // 绝对路径
	}

	var scriptPath string
	for _, path := range scriptPaths {
		if _, err := os.Stat(path); err == nil {
			scriptPath = path
			break
		}
	}

	if scriptPath == "" {
		return "", fmt.Errorf("hexin-v.js not found in any of the expected locations")
	}

	// Use Node.js to execute the JavaScript and get the token
	// The hexin-v.js file defines v() function globally at the end
	cmd := exec.Command("node", scriptPath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute JS with Node.js: %w, output: %s", err, string(output))
	}

	token := strings.TrimSpace(string(output))
	return token, nil
}

// Headers creates HTTP headers map with the token
func Headers(cookie, userAgent string) (map[string]string, error) {
	if userAgent == "" {
		userAgent = GetRandomUserAgent()
	}

	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"hexin-v":    token,
		"User-Agent": userAgent,
	}

	if cookie != "" {
		headers["Cookie"] = cookie
	}

	return headers, nil
}
