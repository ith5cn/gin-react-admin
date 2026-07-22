package system

import "testing"

func TestLoginIPLocation(t *testing.T) {
	cases := []struct {
		ip   string
		want string
	}{
		{"127.0.0.1", "本机"},
		{"::1", "本机"},
		{"192.168.1.20", "内网"},
		{"10.0.0.8", "内网"},
		{"8.8.8.8", ""},
		{"invalid", ""},
	}
	for _, tc := range cases {
		if got := loginIPLocation(tc.ip); got != tc.want {
			t.Fatalf("loginIPLocation(%q) = %q, want %q", tc.ip, got, tc.want)
		}
	}
}

func TestParseUserAgent(t *testing.T) {
	osName, browser := parseUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126.0 Safari/537.36 Edg/126.0")
	if osName != "Windows" || browser != "Edge" {
		t.Fatalf("parseUserAgent edge = %q/%q", osName, browser)
	}

	osName, browser = parseUserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 Version/17.0 Mobile/15E148 Safari/604.1")
	if osName != "iOS" || browser != "Safari" {
		t.Fatalf("parseUserAgent safari = %q/%q", osName, browser)
	}
}
