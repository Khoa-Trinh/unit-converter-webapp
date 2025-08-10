# 🌐 Unit Converter Web App (Go)

A simple server-rendered web application that converts between different units of **Length**, **Weight**, and **Temperature**.
Built in Go using `net/http` and `html/template`.

Project idea from: [https://roadmap.sh/projects/unit-converter](https://roadmap.sh/projects/unit-converter)

---

## ✨ Features

* Convert between common **Length** units:

  * millimeter, centimeter, meter, kilometer, inch, foot, yard, mile.
* Convert between **Weight** units:

  * milligram, gram, kilogram, ounce, pound.
* Convert between **Temperature** units:

  * Celsius, Fahrenheit, Kelvin.
* Clean responsive UI (pure HTML/CSS, no JS required).
* Server-rendered — form submits to itself and displays the result.
* Active navigation highlighting for current page.
* Fully tested conversion functions and HTTP handlers.

---

## 📂 Project Structure

```
.
├── main.go          # Web server, handlers, converters, templates
├── main_test.go     # Unit and handler tests
└── README.md
```

---

## 🚀 Getting Started

### 1. **Install Go**

Make sure you have Go 1.20+ installed:

```bash
go version
```

### 2. **Run the server**

```bash
go run main.go
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## 🧪 Running Tests

The project includes unit tests for:

* Conversion logic
* HTML output of handlers
* Error handling for invalid inputs

Run:

```bash
go test -v
```

---

## 📷 Screenshots

**Length Converter**

```
+----------------------------+
| Value: [ 100     ]         |
| From:  meter               |
| To:    kilometer           |
| [ Convert ]                |
------------------------------
Result: 0.100000
```

**Weight Converter**

```
Result: 2.500000
```

**Temperature Converter**

```
Result: 98.600000
```

---

## 🛠 Tech Stack

* **Language:** Go
* **Web Server:** `net/http`
* **Templating:** `html/template`
* **Testing:** Go’s `testing` package

---

## 📄 License

MIT License — feel free to use and modify. -- see [LICENSE](LICENSE) for details.
