package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func almostEqual(a, b, eps float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= eps
}

func TestConvertLength(t *testing.T) {
	tcs := []struct {
		val      float64
		from, to string
		want     float64
	}{
		{1, "meter", "kilometer", 0.001},
		{1000, "millimeter", "meter", 1},
		{2.54, "centimeter", "inch", 1},
		{1, "mile", "kilometer", 1.609344},
		{3, "yard", "foot", 9},
	}
	for _, tc := range tcs {
		got, err := convertLength(tc.val, tc.from, tc.to)
		if err != "" {
			t.Fatalf("convertLength(%v,%q,%q) err=%q", tc.val, tc.from, tc.to, err)
		}
		if !almostEqual(got, tc.want, 1e-9) {
			t.Fatalf("convertLength(%v,%q,%q)=%v, want %v", tc.val, tc.from, tc.to, got, tc.want)
		}
	}
}

func TestConvertWeight(t *testing.T) {
	tcs := []struct {
		val      float64
		from, to string
		want     float64
	}{
		{1000, "gram", "kilogram", 1},
		{1, "kilogram", "gram", 1000},
		{16, "ounce", "pound", 1},
		{2.2, "pound", "kilogram", 0.9979032},
	}
	for _, tc := range tcs {
		got, err := convertWeight(tc.val, tc.from, tc.to)
		if err != "" {
			t.Fatalf("convertWeight(%v,%q,%q) err=%q", tc.val, tc.from, tc.to, err)
		}
		if !almostEqual(got, tc.want, 1e-6) {
			t.Fatalf("convertWeight(%v,%q,%q)=%v, want %v", tc.val, tc.from, tc.to, got, tc.want)
		}
	}
}

func TestConvertTemp(t *testing.T) {
	tcs := []struct {
		val      float64
		from, to string
		want     float64
	}{
		{0, "Celsius", "Kelvin", 273.15},
		{100, "Celsius", "Fahrenheit", 212},
		{32, "Fahrenheit", "Celsius", 0},
		{300, "Kelvin", "Celsius", 26.85},
	}
	for _, tc := range tcs {
		got, err := convertTemp(tc.val, tc.from, tc.to)
		if err != "" {
			t.Fatalf("convertTemp(%v,%q,%q) err=%q", tc.val, tc.from, tc.to, err)
		}
		if !almostEqual(got, tc.want, 1e-6) {
			t.Fatalf("convertTemp(%v,%q,%q)=%v, want %v", tc.val, tc.from, tc.to, got, tc.want)
		}
	}
}

func TestLengthHandler_PostShowsResult(t *testing.T) {
	form := url.Values{}
	form.Set("value", "100")
	form.Set("from", "meter")
	form.Set("to", "kilometer")

	req := httptest.NewRequest(http.MethodPost, "/length", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	lengthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<strong>Result:</strong> 0.100000") {
		t.Fatalf("body missing formatted result, got:\n\n%s", body)
	}
	if !strings.Contains(body, `href="/length" class="active"`) {
		t.Fatalf("active nav not set for length page")
	}
}

func TestWeightHandler_PostShowsResult(t *testing.T) {
	form := url.Values{}
	form.Set("value", "2500")
	form.Set("from", "gram")
	form.Set("to", "kilogram")

	req := httptest.NewRequest(http.MethodPost, "/weight", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	weightHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<strong>Result:</strong> 2.500000") {
		t.Fatalf("body missing formatted result, got:\n\n%s", body)
	}
	if !strings.Contains(body, `href="/weight" class="active"`) {
		t.Fatalf("active nav not set for weight page")
	}
}

func TestTempHandler_PostShowsResult(t *testing.T) {
	form := url.Values{}
	form.Set("value", "37")
	form.Set("from", "Celsius")
	form.Set("to", "Fahrenheit")

	req := httptest.NewRequest(http.MethodPost, "/temperature", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	tempHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<strong>Result:</strong> 98.600000") {
		t.Fatalf("body missing formatted result, got:\n\n%s", body)
	}
	if !strings.Contains(body, `href="/temperature" class="active"`) {
		t.Fatalf("active nav not set for temperature page")
	}
}

func TestHandlers_InvalidValueShowsError(t *testing.T) {
	for _, path := range []string{"/length", "/weight", "/temperature"} {
		form := url.Values{}
		form.Set("value", "abc") // invalid
		switch path {
		case "/length":
			form.Set("from", "meter")
			form.Set("to", "kilometer")
		case "/weight":
			form.Set("from", "gram")
			form.Set("to", "kilogram")
		case "/temperature":
			form.Set("from", "Celsius")
			form.Set("to", "Kelvin")
		}

		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		switch path {
		case "/length":
			lengthHandler(rr, req)
		case "/weight":
			weightHandler(rr, req)
		case "/temperature":
			tempHandler(rr, req)
		}

		if rr.Code != http.StatusOK {
			t.Fatalf("%s status=%d, want 200", path, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Invalid number.") {
			t.Fatalf("%s should show error message, got:\n%s", path, rr.Body.String())
		}
	}
}
