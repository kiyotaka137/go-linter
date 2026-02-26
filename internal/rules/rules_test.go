package rules

import "testing"

func TestIsLowercase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "simple lowercase", in: "server started", want: true},
		{name: "leading spaces", in: "   server started", want: true},
		{name: "starts uppercase", in: "Server started", want: false},
		{name: "starts digit", in: "8080 started", want: false},
		{name: "empty", in: "", want: false},
		{name: "spaces only", in: "   ", want: false},
		{name: "russian lowercase", in: "ошибка подключения", want: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsLowercase(tt.in)
			if got != tt.want {
				t.Fatalf("IsLowercase(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestIsWithoutSymbols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "letters and spaces", in: "server started", want: true},
		{name: "letters numbers spaces", in: "server started 2", want: true},
		{name: "leading and trailing spaces", in: "  server started  ", want: true},
		{name: "empty", in: "", want: false},
		{name: "spaces only", in: "   ", want: false},
		{name: "exclamation marks", in: "connection failed!!!", want: false},
		{name: "colon", in: "warning: failed", want: false},
		{name: "underscore", in: "api_key leaked", want: false},
		{name: "dash", in: "token-expired", want: false},
		{name: "emoji", in: "server started 🚀", want: false},
		{name: "russian letters", in: "ошибка подключения", want: false},
		{name: "tab", in: "server\tstarted", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsWithoutSymbols(tt.in)
			if got != tt.want {
				t.Fatalf("IsWithoutSymbols(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestFindSensitiveKeyword(t *testing.T) {
	t.Parallel()

	keywords := []string{
		"password",
		" token ",
		"",
		"client_secret",
	}

	tests := []struct {
		name        string
		msg         string
		wantFound   bool
		wantKeyword string
	}{
		{name: "password", msg: "user password leaked", wantFound: true, wantKeyword: "password"},
		{name: "case insensitive token", msg: "TOKEN exposed", wantFound: true, wantKeyword: "token"},
		{name: "client secret", msg: "client_secret exposed", wantFound: true, wantKeyword: "client_secret"},
		{name: "no match", msg: "request completed", wantFound: false, wantKeyword: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotKeyword, gotFound := FindSensitiveKeyword(tt.msg, keywords)
			if gotFound != tt.wantFound {
				t.Fatalf("FindSensitiveKeyword(%q) found = %v, want %v", tt.msg, gotFound, tt.wantFound)
			}
			if gotKeyword != tt.wantKeyword {
				t.Fatalf("FindSensitiveKeyword(%q) keyword = %q, want %q", tt.msg, gotKeyword, tt.wantKeyword)
			}
		})
	}
}
