package watcher

import "testing"

func TestGetWeiboVisitorCookies(t *testing.T) {
	cookies, err := GetWeiboVisitorCookies()
	if err != nil {
		t.Fatalf("Failed to get Weibo visitor cookies: %v", err)
	}
	t.Logf("Successfully fetched Weibo visitor cookies: SUB=%s, SUBP=%s", cookies.SUB, cookies.SUBP)
}
