A Discord bot that provides real-time market sentiment analysis using the Fear and Greed Index, Bitcoin dominance, and top cryptocurrencies data. The bot also generates visual charts to depict historical trends.

## Features

- Fetches and displays the latest Fear and Greed Index.
- Provides market sentiment analysis based on Bitcoin dominance.
- Displays top 5 cryptocurrencies and their price changes.
- Generates charts for BTC dominance and Fear and Greed Index using Python's Matplotlib.

## Requirements

- Go (Golang) 1.18 or later
- Python 3.x
- Python packages: `requests`, `matplotlib`
- A Discord Bot Token (create a bot in the Discord Developer Portal)
- CoinMarketCap API Key

## Setup

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/your-username/DiscordBot-MarketAnalysis.git
    cd DiscordBot-MarketAnalysis
    ```

2. **Install Go and Python Dependencies**:

    - Install Go packages:

        ```bash
        go get github.com/bwmarrin/discordgo
        go get github.com/joho/godotenv
        ```

    - Install Python packages:

        ```bash
        pip install requests matplotlib
        ```

3. **Configure Environment Variables**:

    Create a `.env` file in the project root and add your Discord Bot Token and CoinMarketCap API Key:

    ```plaintext
    DISCORD_BOT_TOKEN=your_discord_bot_token
    COINMARKETCAP_API_KEY=your_coinmarketcap_api_key
    ```

4. **Run the Bot**:

    ```bash
    go run main.go
    Commands: !Market, !Help, !chart
    ```
