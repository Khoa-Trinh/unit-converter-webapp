package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var tmpl = template.Must(
	template.New("base").Parse(baseHTML),
)

func init() {
	// Parse subtemplates into the SAME template object (tmpl)
	template.Must(tmpl.New("length").Parse(lengthHTML))
	template.Must(tmpl.New("weight").Parse(weightHTML))
	template.Must(tmpl.New("temperature").Parse(tempHTML))
}

// ---------- Conversion tables ----------

// Length → via meters
var lengthToMeters = map[string]float64{
	"millimeter": 0.001,
	"centimeter": 0.01,
	"meter":      1,
	"kilometer":  1000,
	"inch":       0.0254,
	"foot":       0.3048,
	"yard":       0.9144,
	"mile":       1609.344,
}
var lengthUnits = []string{"millimeter", "centimeter", "meter", "kilometer", "inch", "foot", "yard", "mile"}

// Weight → via kilograms
var weightToKg = map[string]float64{
	"milligram": 1e-6,
	"gram":      1e-3,
	"kilogram":  1,
	"ounce":     0.028349523125,
	"pound":     0.45359237,
}
var weightUnits = []string{"milligram", "gram", "kilogram", "ounce", "pound"}

// Temperature units
var tempUnits = []string{"Celsius", "Fahrenheit", "Kelvin"}

// ---------- Helpers ----------

func parseValue(r *http.Request, name string) (float64, string) {
	raw := strings.TrimSpace(r.FormValue(name))
	if raw == "" {
		return 0, "Please enter a value."
	}
	v, err := strconv.ParseFloat(strings.ReplaceAll(raw, ",", ""), 64)
	if err != nil {
		return 0, "Invalid number."
	}
	return v, ""
}

func selectedOrDefault(r *http.Request, name, def string) string {
	if v := r.FormValue(name); v != "" {
		return v
	}
	return def
}

func round(n float64) float64 {
	f, _ := strconv.ParseFloat(strconv.FormatFloat(n, 'f', 6, 64), 64)
	return f
}

// ---------- Converters ----------

func convertLength(val float64, from, to string) (float64, string) {
	fk, okF := lengthToMeters[from]
	tk, okT := lengthToMeters[to]
	if !okF || !okT {
		return 0, "Unsupported length unit."
	}
	meters := val * fk
	return meters / tk, ""
}

func convertWeight(val float64, from, to string) (float64, string) {
	fk, okF := weightToKg[from]
	tk, okT := weightToKg[to]
	if !okF || !okT {
		return 0, "Unsupported weight unit."
	}
	kg := val * fk
	return kg / tk, ""
}

func toCelsius(v float64, from string) (float64, string) {
	switch from {
	case "Celsius":
		return v, ""
	case "Fahrenheit":
		return (v - 32) * 5 / 9, ""
	case "Kelvin":
		return v - 273.15, ""
	default:
		return 0, "Unsupported temperature unit."
	}
}

func fromCelsius(c float64, to string) (float64, string) {
	switch to {
	case "Celsius":
		return c, ""
	case "Fahrenheit":
		return c*9/5 + 32, ""
	case "Kelvin":
		return c + 273.15, ""
	default:
		return 0, "Unsupported temperature unit."
	}
}

func convertTemp(val float64, from, to string) (float64, string) {
	c, err := toCelsius(val, from)
	if err != "" {
		return 0, err
	}
	out, err2 := fromCelsius(c, to)
	if err2 != "" {
		return 0, err2
	}
	return out, ""
}

// ---------- Handlers ----------

type viewData struct {
	Units     []string
	Value     string
	From, To  string
	Result    float64
	HasResult bool
	Error     string
	Active    string
}

func lengthHandler(w http.ResponseWriter, r *http.Request) {
	data := viewData{
		Units:  lengthUnits,
		From:   selectedOrDefault(r, "from", "meter"),
		To:     selectedOrDefault(r, "to", "kilometer"),
		Value:  r.FormValue("value"),
		Active: "length",
	}
	if r.Method == http.MethodPost {
		val, msg := parseValue(r, "value")
		if msg != "" {
			data.Error = msg
		} else {
			out, msg := convertLength(val, data.From, data.To)
			if msg != "" {
				data.Error = msg
			} else {
				data.Result = round(out)
				data.HasResult = true
			}
		}
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func weightHandler(w http.ResponseWriter, r *http.Request) {
	data := viewData{
		Units:  weightUnits,
		From:   selectedOrDefault(r, "from", "gram"),
		To:     selectedOrDefault(r, "to", "kilogram"),
		Value:  r.FormValue("value"),
		Active: "weight",
	}
	if r.Method == http.MethodPost {
		val, msg := parseValue(r, "value")
		if msg != "" {
			data.Error = msg
		} else {
			out, msg := convertWeight(val, data.From, data.To)
			if msg != "" {
				data.Error = msg
			} else {
				data.Result = round(out)
				data.HasResult = true
			}
		}
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	data := viewData{
		Units:  tempUnits,
		From:   selectedOrDefault(r, "from", "Celsius"),
		To:     selectedOrDefault(r, "to", "Fahrenheit"),
		Value:  r.FormValue("value"),
		Active: "temperature",
	}
	if r.Method == http.MethodPost {
		val, msg := parseValue(r, "value")
		if msg != "" {
			data.Error = msg
		} else {
			out, msg := convertTemp(val, data.From, data.To)
			if msg != "" {
				data.Error = msg
			} else {
				data.Result = round(out)
				data.HasResult = true
			}
		}
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// ---------- Router / main ----------

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/length", http.StatusFound)
	})
	http.HandleFunc("/length", lengthHandler)
	http.HandleFunc("/weight", weightHandler)
	http.HandleFunc("/temperature", tempHandler)

	log.Println("Unit Converter running at http://localhost:8080 …")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---------- Templates ----------

const baseHTML = `
{{define "base"}}
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Unit Converter</title>
<style>
  :root{--bg:#0b0c10;--card:#15171c;--text:#e8eef2;--muted:#aab4bf;--accent:#60a5fa;--error:#ef4444;}
  *{box-sizing:border-box} body{margin:0;font-family:system-ui,-apple-system,Segoe UI,Roboto,Inter,Arial,sans-serif;background:var(--bg);color:var(--text);}
  header{padding:18px 20px;border-bottom:1px solid #23262d;background:#0f1116;position:sticky;top:0}
  nav a{color:var(--muted);text-decoration:none;margin-right:14px;padding:8px 10px;border-radius:10px}
  nav a.active{background:var(--card);color:var(--text);border:1px solid #262a33}
  main{max-width:760px;margin:30px auto;padding:0 16px}
  .card{background:var(--card);border:1px solid #23262d;border-radius:16px;padding:20px;box-shadow:0 4px 20px rgba(0,0,0,.25)}
  .row{display:grid;grid-template-columns:1fr 1fr;gap:12px}
  label{font-size:14px;color:var(--muted)}
  input, select{width:100%;margin-top:6px;background:#0f1116;border:1px solid #23262d;color:var(--text);padding:12px 10px;border-radius:12px;outline:none}
  input:focus, select:focus{border-color:var(--accent)}
  button{margin-top:14px;background:var(--accent);border:none;color:#04121f;padding:12px 14px;border-radius:12px;font-weight:600;cursor:pointer}
  .result{margin-top:16px;padding:14px;border-radius:12px;background:#0f1116;border:1px dashed #2a2f39}
  .error{margin-top:12px;color:var(--error)}
  footer{margin:30px auto 40px;max-width:760px;color:var(--muted);font-size:14px;padding:0 16px}
</style>
</head>
<body>
  <header>
    <nav>
      <a href="/length" {{if eq .Active "length"}}class="active"{{end}}>Length</a>
      <a href="/weight" {{if eq .Active "weight"}}class="active"{{end}}>Weight</a>
      <a href="/temperature" {{if eq .Active "temperature"}}class="active"{{end}}>Temperature</a>
    </nav>
  </header>
  <main>
    {{if eq .Active "length"}}
      {{template "length" .}}
    {{else if eq .Active "weight"}}
      {{template "weight" .}}
    {{else if eq .Active "temperature"}}
      {{template "temperature" .}}
    {{else}}
      {{template "length" .}}
    {{end}}
  </main>
  <footer>
    <div>Server-rendered, no database. Submits to self and displays the result.</div>
  </footer>
</body>
</html>
{{end}}
`

const lengthHTML = `
{{define "length"}}
  <div class="card">
    <h2>Length Converter</h2>
    <form method="post" action="/length" target="_self" autocomplete="off" novalidate>
      <div class="row">
        <div>
          <label for="value">Value</label>
          <input id="value" name="value" type="text" placeholder="e.g. 123.45" value="{{.Value}}">
        </div>
        <div>
          <label>&nbsp;</label>
          <button type="submit">Convert</button>
        </div>
      </div>
      <div class="row">
        <div>
          <label for="from">From</label>
          <select id="from" name="from">
            {{range .Units}}
              <option value="{{.}}" {{if eq $.From .}}selected{{end}}>{{.}}</option>
            {{end}}
          </select>
        </div>
        <div>
          <label for="to">To</label>
          <select id="to" name="to">
            {{range .Units}}
              <option value="{{.}}" {{if eq $.To .}}selected{{end}}>{{.}}</option>
            {{end}}
          </select>
        </div>
      </div>
      {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
      {{if .HasResult}}
        <div class="result"><strong>Result:</strong> {{printf "%.6f" .Result}}</div>
      {{end}}
    </form>
  </div>
{{end}}
`

const weightHTML = `
{{define "weight"}}
  <div class="card">
    <h2>Weight Converter</h2>
    <form method="post" action="/weight" target="_self" autocomplete="off" novalidate>
      <div class="row">
        <div>
          <label for="value">Value</label>
          <input id="value" name="value" type="text" placeholder="e.g. 2500" value="{{.Value}}">
        </div>
        <div>
          <label>&nbsp;</label>
          <button type="submit">Convert</button>
        </div>
      </div>
      <div class="row">
        <div>
          <label for="from">From</label>
          <select id="from" name="from">
            {{range .Units}}
              <option value="{{.}}" {{if eq $.From .}}selected{{end}}>{{.}}</option>
            {{end}}
          </select>
        </div>
        <div>
          <label for="to">To</label>
          <select id="to" name="to">
            {{range .Units}}
              <option value="{{.}}" {{if eq $.To .}}selected{{end}}>{{.}}</option>
            {{end}}
          </select>
        </div>
      </div>
      {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
      {{if .HasResult}}
        <div class="result"><strong>Result:</strong> {{printf "%.6f" .Result}}</div>
      {{end}}
    </form>
  </div>
{{end}}
`

const tempHTML = `
{{define "temperature"}}
  <div class="card">
    <h2>Temperature Converter</h2>
    <form method="post" action="/temperature" target="_self" autocomplete="off" novalidate>
      <div class="row">
        <div>
          <label for="value">Value</label>
          <input id="value" name="value" type="text" placeholder="e.g. 37" value="{{.Value}}">
        </div>
        <div>
          <label>&nbsp;</label>
          <button type="submit">Convert</button>
        </div>
      </div>
      <div class="row">
        <div>
          <label for="from">From</label>
          <select id="from" name="from">
            {{range .Units}}
              <option value="{{.}}" {{if eq $.From .}}selected{{end}}>{{.}}</option>
            {{end}}
          </select>
        </div>
        <div>
          <label for="to">To</label>
          <select id="to" name="to">
            {{range .Units}}
              <option value="{{.}}" {{if eq $.To .}}selected{{end}}>{{.}}</option>
            {{end}}
          </select>
        </div>
      </div>
      {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
      {{if .HasResult}}
        <div class="result"><strong>Result:</strong> {{printf "%.6f" .Result}}</div>
      {{end}}
    </form>
  </div>
{{end}}
`
