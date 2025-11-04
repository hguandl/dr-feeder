package watcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)" +
	" AppleWebKit/605.1.15 (KHTML, like Gecko)" +
	" Version/26.1 Safari/605.1.15"

// WeiboVisitorCookies holds the SUB and SUBP cookie values
type WeiboVisitorCookies struct {
	SUB  string
	SUBP string
}

// GetWeiboVisitorCookies fetches visitor cookies from Weibo
func GetWeiboVisitorCookies() (*WeiboVisitorCookies, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar:           jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return nil },
	}

	req, err := http.NewRequest("GET", "https://m.weibo.cn", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to visit m.weibo.cn: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	requestID, err := extractRequestID(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to extract request_id: %w", err)
	}

	finalURL := resp.Request.URL.String()
	cookies, err := sendGenVisitorRequest(client, requestID, finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get visitor cookies: %w", err)
	}

	return cookies, nil
}

// extractRequestID extracts the request_id variable from JavaScript code
func extractRequestID(html string) (string, error) {
	re := regexp.MustCompile(`var\s+request_id\s*=\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		log.Printf("Got request_id = %v", matches[1])
		return matches[1], nil
	}

	return "", fmt.Errorf("request_id not found in page content")
}

// sendGenVisitorRequest sends POST request to genvisitor2 endpoint
func sendGenVisitorRequest(client *http.Client, requestID string, refererURL string) (*WeiboVisitorCookies, error) {
	rid := strconv.FormatInt(time.Now().UnixMilli(), 10)
	formData := url.Values{
		"cb":         {"visitor_gray_callback"},
		"ver":        {"20250916"},
		"request_id": {requestID},
		"tid":        {""},
		"from":       {"weibo"},
		"webdriver":  {"false"},
		"rid":        {rid},
		"return_url": {"https://m.weibo.cn"},
	}

	req, err := http.NewRequest("POST", "https://visitor.passport.weibo.cn/visitor/genvisitor2",
		strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", refererURL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	cookies := &WeiboVisitorCookies{}

	parsedURL, _ := url.Parse("https://visitor.passport.weibo.cn")
	for _, cookie := range client.Jar.Cookies(parsedURL) {
		switch cookie.Name {
		case "SUB":
			cookies.SUB = cookie.Value
		case "SUBP":
			cookies.SUBP = cookie.Value
		}
	}

	if cookies.SUB == "" || cookies.SUBP == "" {
		return nil, fmt.Errorf("failed to extract SUB or SUBP cookies")
	}

	return cookies, nil
}
