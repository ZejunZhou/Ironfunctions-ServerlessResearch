import matplotlib.pyplot as plt
import numpy as np

# Data from the two test results
iron_latency = [21.58, 32.51, 58.53, 256.89, 323.58, 323.58]
night_latency = [11.74,  13.15,  15.34,  63.39, 252.03, 252.03]
percentiles = [50, 75, 90, 99, 99.9, 99.99]
percentile_labels = ['p50', 'p75', 'p90', 'p99', 'p99.9', 'p99.99']

# Create the plot
plt.figure(figsize=(10, 6))
plt.plot(range(len(percentiles)), iron_latency, label='Iron Function', marker='o')
plt.plot(range(len(percentiles)), night_latency, label='Nightcore', marker='o')

# Add data labels
for i, (iron, night) in enumerate(zip(iron_latency, night_latency)):
    plt.text(i, iron, f'{iron}ms', fontsize=9, ha='right')
    plt.text(i, night, f'{night}ms', fontsize=9, ha='left')

# Graph details
plt.title("Latency Distribution Comparison", fontsize=14)
plt.xlabel("Percentiles", fontsize=12)
plt.ylabel("Latency (ms)", fontsize=12)
plt.yscale('log')  # Log scale for better visualization of high variance
plt.grid(True, which="both", linestyle="--", linewidth=0.5)
plt.legend(fontsize=12)
plt.xticks(range(len(percentiles)), labels=percentile_labels)

# Save the graph to a file
plt.savefig('latency_distribution_comparison.png')
