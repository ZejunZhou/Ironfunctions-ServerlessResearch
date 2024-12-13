import matplotlib.pyplot as plt
import numpy as np

# Data from the two test results
iron_latency = [21.58, 32.51, 58.53, 256.89, 323.58, 323.58]
night_latency = [11.74,  13.15,  15.34,  63.39, 252.03, 252.03]
percentiles = [50, 75, 90, 99, 99.9, 99.99]
percentile_labels = ['p50', 'p75', 'p90', 'p99', 'p99.9', 'p99.99']

# Create the plot
bar_width = 0.35
index = np.arange(len(percentiles))

plt.figure(figsize=(10, 6))
plt.bar(index, night_latency, bar_width, label='Nightcore')
plt.bar(index  + bar_width, iron_latency, bar_width, label='Iron Function')

# Add data labels
# for i, (iron, night) in enumerate(zip(iron_latency, night_latency)):
#     plt.text(i, night + 10, f'{night}ms', fontsize=9, ha='center')
#     plt.text(i  + bar_width, iron + 10, f'{iron}ms', fontsize=9, ha='center')

# Graph details
plt.title("Latency Distribution Comparison(RPS 50)", fontsize=14)
plt.xlabel("Percentiles", fontsize=12)
plt.ylabel("Latency (ms)", fontsize=12)
plt.yscale('log')  # Log scale for better visualization of high variance
# plt.grid(True, which="both", linestyle="--", linewidth=0.5)
plt.legend(fontsize=12)
plt.xticks(index + bar_width / 2, labels=percentile_labels)

# Save the graph to a file
plt.savefig('latency_distribution_comparison.png')
