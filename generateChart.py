# generate_chart.py

import sys
import requests
import matplotlib.pyplot as plt
from datetime import datetime, timedelta

def fetch_crypto_data(symbol):
    # Fetch historical data for a specific cryptocurrency from CoinGecko
    url = f"https://api.coingecko.com/api/v3/coins/{symbol}/market_chart"
    params = {
        'vs_currency': 'usd',
        'days': '180'
    }

    response = requests.get(url, params=params)
    if response.status_code != 200:
        print(f"Error fetching data: {response.status_code} - {response.text}")
        return [], []

    data = response.json().get("prices", [])
    prices = [price[1] for price in data]
    dates = [datetime.fromtimestamp(price[0] / 1000).strftime('%Y-%m-%d') for price in data]

    return prices, dates

def create_chart(title, dates, values, filename):
    plt.figure(figsize=(12, 6))
    plt.plot(dates, values, marker='o', color='blue')
    plt.title(title)
    plt.xlabel('Date')
    plt.ylabel('Price (USD)')
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.savefig(filename)
    plt.close()

if __name__ == "__main__":
    if len(sys.argv) > 1:
        crypto_symbol = sys.argv[1]
        prices, dates = fetch_crypto_data(crypto_symbol)
        if prices and dates:
            create_chart(f"{crypto_symbol.capitalize()} Price Over Last 6 Months", dates, prices, f"{crypto_symbol}_price_chart.png")
        else:
            print(f"Failed to fetch data for {crypto_symbol}")
    else:
        print("Please specify a cryptocurrency symbol as an argument.")
