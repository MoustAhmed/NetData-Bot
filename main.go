package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gopkg.in/resty.v1"
)

// Struct to handle API response for Fear and Greed Index
type FearGreedIndexResponse struct {
	Data []struct {
		Value     string `json:"value"`
		ValueText string `json:"value_classification"`
		Timestamp string `json:"timestamp"`
	} `json:"data"`
}

// Struct to handle CoinMarketCap API response for global metrics (including BTC dominance)
type CoinMarketCapGlobalMetricsResponse struct {
	Status struct {
		ErrorCode int `json:"error_code"`
	} `json:"status"`
	Data struct {
		BTCPercentage float64 `json:"btc_dominance"`
	} `json:"data"`
}

// Struct to handle CoinMarketCap API response for top cryptocurrencies
type CoinMarketCapResponse struct {
	Status struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_message"`
	} `json:"status"`
	Data []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  struct {
			USD struct {
				Price            float64 `json:"price"`
				PercentChange1h  float64 `json:"percent_change_1h"`
				PercentChange24h float64 `json:"percent_change_24h"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

// Struct to handle CoinMarketCap API response for specific coin historical data
type CoinHistoricalDataResponse struct {
	Status struct {
		ErrorCode int `json:"error_code"`
	} `json:"status"`
	Data struct {
		Quotes []struct {
			TimeOpen string  `json:"time_open"`
			Close    float64 `json:"close"`
		} `json:"quotes"`
	} `json:"data"`
}

// Function to fetch historical Fear and Greed Index data
func getHistoricalFearGreedIndex() (FearGreedIndexResponse, error) {
	var result FearGreedIndexResponse

	resp, err := resty.New().R().
		SetHeader("Accept", "application/json").
		Get("https://api.alternative.me/fng/?limit=180") // Fetch data for the last 6 months (approx 180 days)

	if err != nil {
		fmt.Println("Error making API request:", err)
		return result, err
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Println("API request failed with status:", resp.StatusCode(), "response:", resp.String())
		return result, fmt.Errorf("API request failed with status: %d", resp.StatusCode())
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		fmt.Println("Error parsing API response:", err)
	}
	return result, err
}

// Function to fetch historical BTC dominance (mock function, replace with real API call)
func getHistoricalBTCDominance() ([]float64, []string, error) {
	var btcDominance []float64
	var dates []string

	// Mock data for demonstration purposes
	btcDominance = []float64{60, 58, 62, 63, 61, 60} // Example data
	dates = []string{"2024-02-01", "2024-03-01", "2024-04-01", "2024-05-01", "2024-06-01", "2024-07-01"}

	return btcDominance, dates, nil
}

// Function to fetch historical data for a specific coin (mock function, replace with real API call)
func getHistoricalCoinData(coinSymbol string) ([]float64, []string, error) {
	var prices []float64
	var dates []string

	// Mock data for demonstration purposes
	prices = []float64{30000, 32000, 31000, 33000, 34000, 35000} // Example data for BTC
	dates = []string{"2024-02-01", "2024-03-01", "2024-04-01", "2024-05-01", "2024-06-01", "2024-07-01"}

	return prices, dates, nil
}

// Function to fetch Fear and Greed Index
func getFearGreedIndex() (FearGreedIndexResponse, error) {
	var result FearGreedIndexResponse

	resp, err := resty.New().R().
		SetHeader("Accept", "application/json").
		Get("https://api.alternative.me/fng/?limit=1") // Latest data

	if err != nil {
		fmt.Println("Error making API request:", err)
		return result, err
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Println("API request failed with status:", resp.StatusCode(), "response:", resp.String())
		return result, fmt.Errorf("API request failed with status: %d", resp.StatusCode())
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		fmt.Println("Error parsing API response:", err)
	}
	return result, err
}

// Function to fetch global metrics including BTC dominance from CoinMarketCap API
func getGlobalMetrics() (CoinMarketCapGlobalMetricsResponse, error) {
	var result CoinMarketCapGlobalMetricsResponse
	apiKey := os.Getenv("COINMARKETCAP_API_KEY")

	resp, err := resty.New().R().
		SetHeader("X-CMC_PRO_API_KEY", apiKey).
		SetHeader("Accept", "application/json").
		Get("https://pro-api.coinmarketcap.com/v1/global-metrics/quotes/latest")

	if err != nil {
		fmt.Println("Error making API request:", err)
		return result, err
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Println("API request failed with status:", resp.StatusCode(), "response:", resp.String())
		return result, fmt.Errorf("API request failed with status: %d", resp.StatusCode())
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		fmt.Println("Error parsing API response:", err)
	}
	return result, err
}

func generateChart(symbol string) error {
	var cmd *exec.Cmd
	pythonPath := "python" // Adjust to the path of Python if needed

	if symbol == "" {
		cmd = exec.Command(pythonPath, "generateChart.py")
	} else {
		cmd = exec.Command(pythonPath, "generateChart.py", symbol)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error generating chart:", err)
		return err
	}
	return nil
}

// Function to fetch top 5 cryptocurrencies from CoinMarketCap API
func getTop5Cryptos() (CoinMarketCapResponse, error) {
	var result CoinMarketCapResponse
	apiKey := os.Getenv("COINMARKETCAP_API_KEY")

	resp, err := resty.New().R().
		SetHeader("X-CMC_PRO_API_KEY", apiKey).
		SetHeader("Accept", "application/json").
		Get("https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=5")

	if err != nil {
		fmt.Println("Error making API request:", err)
		return result, err
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Println("API request failed with status:", resp.StatusCode(), "response:", resp.String())
		return result, fmt.Errorf("API request failed with status: %d", resp.StatusCode())
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		fmt.Println("Error parsing API response:", err)
	}
	return result, err
}

// Function to format percentage change using red and green circles
func formatPercentageChange(change float64) string {
	if change < 0 {
		return fmt.Sprintf("🔴 %.2f%%", change) // Red circle for negative changes
	}
	return fmt.Sprintf("🟢 %.2f%%", change) // Green circle for positive changes
}

// Function to create a line chart from data
func createLineChart(filename string, title string, xValues []string, yValues []float64) error {
	p := plot.New() // No error return here now

	p.Title.Text = title
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Value"

	pts := make(plotter.XYs, len(yValues))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = yValues[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red line
	p.Add(line)

	p.NominalX(xValues...)

	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		return err
	}

	return nil
}

// Function to handle new messages
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Respond to "!HELLO" command
	if m.Content == "!HELLO" {
		s.ChannelMessageSend(m.ChannelID, "JELLOOOOOO!")
	}

	// Respond to "!Market" command for market overview
	if m.Content == "!Market" {
		fearGreedResponse, err := getFearGreedIndex()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching Fear and Greed Index.")
			return
		}

		cryptoResponse, err := getTop5Cryptos()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching top 5 cryptocurrencies.")
			return
		}

		globalMetricsResponse, err := getGlobalMetrics()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching global market metrics.")
			return
		}

		// Prepare response message with improved formatting
		message := fmt.Sprintf("**📊 Fear and Greed Index:**\n**Value:** %s\n**Description:** %s\n\n", fearGreedResponse.Data[0].Value, fearGreedResponse.Data[0].ValueText)
		message += "**🚀 Top 5 Cryptocurrencies:**\n"
		for _, crypto := range cryptoResponse.Data {
			message += fmt.Sprintf("**%s (%s)**\nPrice: `$%.2f`\n24h Change: %s\n1h Change: %s\n\n",
				crypto.Name, crypto.Symbol, crypto.Quote.USD.Price, formatPercentageChange(crypto.Quote.USD.PercentChange24h), formatPercentageChange(crypto.Quote.USD.PercentChange1h))
		}

		// Add Overview section with Bitcoin Dominance and Market Sentiment Analysis
		btcDominance := globalMetricsResponse.Data.BTCPercentage
		message += fmt.Sprintf("**Overview:**\n**🌐 Bitcoin Dominance:** %.2f%%\n", btcDominance)

		fearValue := 0
		fmt.Sscanf(fearGreedResponse.Data[0].Value, "%d", &fearValue) // Convert value to integer for analysis

		if fearValue < 50 {
			message += "😱 The market is fearful. This might be a good time to buy.\n"
			if btcDominance > 50 {
				message += "💡 Consider buying Bitcoin due to its dominance.\n"
			} else {
				message += "💡 Consider looking into Altcoins due to their current dominance.\n"
			}
		} else {
			message += "😄 The market is optimistic. Be cautious of potential corrections.\n"
		}

		if btcDominance > 50 {
			message += "📈 It's currently **Bitcoin season** based on BTC dominance.\n"
		} else {
			message += "📉 It's currently **Altcoin season** based on BTC dominance.\n"
		}

		// Send the complete message to the Discord channel
		s.ChannelMessageSend(m.ChannelID, message)
	}

	// Respond to "!chart" command to show a list of available cryptocurrencies
	if m.Content == "!chart" {
		cryptoList := "**Available Cryptocurrencies for Charting:**\n"
		cryptoList += "- Bitcoin (`bitcoin`)\n"
		cryptoList += "- Ethereum (`ethereum`)\n"
		cryptoList += "- Polkadot (`polkadot`)\n"
		cryptoList += "- Ripple (`ripple`)\n"
		cryptoList += "- Cardano (`cardano`)\n"
		cryptoList += "- Dogecoin (`dogecoin`)\n"
		cryptoList += "- Solana (`solana`)\n"
		cryptoList += "- Chainlink (`chainlink`)\n"
		cryptoList += "- Binance Coin (`binancecoin`)\n"
		cryptoList += "- Litecoin (`litecoin`)\n"
		cryptoList += "\nUse `!chart [crypto name]` to view a chart for a specific cryptocurrency."

		s.ChannelMessageSend(m.ChannelID, cryptoList)
	}

	// Respond to "!chart [crypto symbol]" command
	if strings.HasPrefix(m.Content, "!chart ") {
		symbol := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(m.Content, "!chart ")))
		err := generateChart(symbol)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error generating chart for %s.", symbol))
			return
		}

		filename := fmt.Sprintf("%s_price_chart.png", symbol)
		file, err := os.Open(filename)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error opening chart file for %s.", symbol))
			return
		}
		defer file.Close()

		_, err = s.ChannelFileSend(m.ChannelID, filename, file)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error sending chart file for %s.", symbol))
			return
		}
	}

	// Respond to "!Help" command
	if m.Content == "!Help" {
		helpMessage := "**ℹ️ Available Commands:**\n"
		helpMessage += "`!Market` - Display Fear and Greed Index, BTC dominance, top 5 cryptocurrencies, and investment advice.\n"
		helpMessage += "`!chart` - Display a list of available cryptocurrencies for charting.\n"
		helpMessage += "`!chart [crypto name]` - Display a chart for the specified cryptocurrency over the last 6 months.\n"
		helpMessage += "`!Help` - Display this help message."
		s.ChannelMessageSend(m.ChannelID, helpMessage)
	}
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Get the bot token from environment variables
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	// Wait here until CTRL+C or other term signal is received
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	// Cleanly close down the Discord session
	dg.Close()
}
