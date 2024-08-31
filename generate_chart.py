# generate_chart.py

import matplotlib.pyplot as plt
import pandas as pd
from datetime import datetime

def create_chart(data, dates, title, filename):
    plt.figure(figsize=(10, 5))
    plt.plot(dates, data, marker='o', color='b')
    plt.title(title)
    plt.xlabel('Date')
    plt.ylabel('Value')
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.savefig(filename)
    plt.close()

# Example data (use real data instead)
btc_dominance = [60, 58, 62, 63, 61, 60]
dates = ["2024-02-01", "2024-03-01", "2024-04-01", "2024-05-01", "2024-06-01", "2024-07-01"]
formatted_dates = [datetime.strptime(date, "%Y-%m-%d") for date in dates]

create_chart(btc_dominance, formatted_dates, "BTC Dominance Over Time", "btc_dominance_chart.png")
