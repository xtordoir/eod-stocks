package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/gin-gonic/gin"
)

// Ticker represents a stock ticker
type Ticker struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// TickerStore holds our in-memory ticker data
type TickerStore struct {
	tickers []Ticker
}

var store *TickerStore

func main() {
	// Initialize data from DuckDB
	if err := loadTickersFromDuckDB(); err != nil {
		log.Fatal("Failed to load tickers:", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Serve HTML page
	r.GET("/", servePage)

	// API endpoint for searching tickers
	r.GET("/api/tickers/search", searchTickers)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

// loadTickersFromDuckDB loads ticker data from DuckDB into memory
func loadTickersFromDuckDB() error {
	// Connect to DuckDB
	db, err := sql.Open("duckdb", "tickers.duckdb")
	if err != nil {
		return err
	}
	defer db.Close()

	// Query all tickers
	rows, err := db.Query("SELECT symbol, name FROM tickers ORDER BY symbol")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Load into memory
	var tickers []Ticker
	for rows.Next() {
		var t Ticker
		if err := rows.Scan(&t.Symbol, &t.Name); err != nil {
			return err
		}
		tickers = append(tickers, t)
	}

	store = &TickerStore{tickers: tickers}
	log.Printf("Loaded %d tickers into memory", len(tickers))
	return nil
}

// searchTickers handles the search API endpoint
func searchTickers(c *gin.Context) {
	query := strings.ToUpper(strings.TrimSpace(c.Query("q")))

	if query == "" {
		c.JSON(http.StatusOK, []Ticker{})
		return
	}

	// Search in memory
	var results []Ticker
	for _, ticker := range store.tickers {
		// Match symbol or name
		if strings.Contains(ticker.Symbol, query) ||
			strings.Contains(strings.ToUpper(ticker.Name), query) {
			results = append(results, ticker)
			// Limit results to 50
			if len(results) >= 50 {
				break
			}
		}
	}

	c.JSON(http.StatusOK, results)
}

// servePage serves the HTML page
func servePage(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, pageHTML)
}

const pageHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ticker Search</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        
        h1 {
            color: white;
            text-align: center;
            margin-bottom: 40px;
            font-size: 2.5rem;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }
        
        .search-box {
            background: white;
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            margin-bottom: 30px;
        }
        
        input[type="text"] {
            width: 100%;
            padding: 15px;
            font-size: 18px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            transition: border-color 0.3s;
        }
        
        input[type="text"]:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .results {
            background: white;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            max-height: 500px;
            overflow-y: auto;
        }
        
        .result-item {
            padding: 15px 20px;
            border-bottom: 1px solid #f0f0f0;
            display: flex;
            justify-content: space-between;
            align-items: center;
            transition: background-color 0.2s;
        }
        
        .result-item:hover {
            background-color: #f8f9ff;
        }
        
        .result-item:last-child {
            border-bottom: none;
        }
        
        .symbol {
            font-weight: bold;
            color: #667eea;
            font-size: 16px;
            min-width: 80px;
        }
        
        .name {
            color: #666;
            flex: 1;
        }
        
        .empty-state {
            padding: 60px 20px;
            text-align: center;
            color: #999;
        }
        
        .loading {
            padding: 40px;
            text-align: center;
            color: #667eea;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔍 Ticker Search</h1>
        
        <div class="search-box">
            <input 
                type="text" 
                id="searchInput" 
                placeholder="Search by symbol or company name..."
                autocomplete="off"
            />
        </div>
        
        <div class="results" id="results">
            <div class="empty-state">
                Start typing to search for tickers...
            </div>
        </div>
    </div>

    <script>
        const searchInput = document.getElementById('searchInput');
        const resultsDiv = document.getElementById('results');
        let debounceTimer;

        searchInput.addEventListener('input', function(e) {
            clearTimeout(debounceTimer);
            const query = e.target.value.trim();
            
            if (!query) {
                resultsDiv.innerHTML = '<div class="empty-state">Start typing to search for tickers...</div>';
                return;
            }
            
            resultsDiv.innerHTML = '<div class="loading">Searching...</div>';
            
            debounceTimer = setTimeout(() => {
                searchTickers(query);
            }, 300);
        });

        async function searchTickers(query) {
            try {
                const response = await fetch('/api/tickers/search?q=' + encodeURIComponent(query));
                const tickers = await response.json();
                
                if (tickers.length === 0) {
                    resultsDiv.innerHTML = '<div class="empty-state">No tickers found</div>';
                    return;
                }
                
                resultsDiv.innerHTML = tickers.map(ticker =>
                    '<div class="result-item">' +
                        '<span class="symbol">' + ticker.symbol + '</span>' +
                        '<span class="name">' + ticker.name + '</span>' +
                    '</div>'
                ).join('');
            } catch (error) {
                resultsDiv.innerHTML = '<div class="empty-state">Error loading results</div>';
                console.error('Search error:', error);
            }
        }
    </script>
</body>
</html>
`
