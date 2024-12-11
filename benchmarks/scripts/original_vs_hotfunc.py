import matplotlib.pyplot as plt
import numpy as np

# Data extracted from the WRK2 results
latencies = {
    "OptimizedIronFunc": {
        "P50": 12.21,
        "P90": 13.85,
        "P99": 15.48,
        "P999": 15.48,
    },
    "OriginalIronFunc": {
        "P50": 4600,
        "P90": 7380,
        "P99": 9130,
        "P999": 9130,
    },
}

# Extracting data for plotting
categories = ["P50", "P90", "P99", "P999"]
x = np.arange(len(categories))  # x-axis locations

optimized_data = [latencies["OptimizedIronFunc"][key] for key in categories]
original_data = [latencies["OriginalIronFunc"][key] for key in categories]

# Plotting the histogram
bar_width = 0.35
fig, ax = plt.subplots(figsize=(10, 6))

# Bars for OptimizedIronFunc
bars1 = ax.bar(x - bar_width / 2, optimized_data, bar_width, label="Optimized_IronFunc", color="skyblue")
# Bars for OriginalIronFunc
bars2 = ax.bar(x + bar_width / 2, original_data, bar_width, label="Original_IronFunc", color="orange")

# Adding labels and title
ax.set_xlabel("Latency Percentiles", fontsize=12)
ax.set_ylabel("Latency (ms) [log scale]", fontsize=12)
ax.set_title("Latency Comparison Echo Function: OptimizedIronFunc vs OriginalIronFunc", fontsize=14)
ax.set_xticks(x)
ax.set_xticklabels(categories, fontsize=10)
ax.legend(fontsize=10)

# Set y-axis to log scale
ax.set_yscale('log')

# Adding value annotations
for bars in [bars1, bars2]:
    for bar in bars:
        height = bar.get_height()
        ax.annotate(f'{height:.2f}',
                    xy=(bar.get_x() + bar.get_width() / 2, height),
                    xytext=(0, 3),  # Offset text above the bar
                    textcoords="offset points",
                    ha='center', va='bottom', fontsize=9)

# Show the plot
plt.tight_layout()
plt.savefig('./images/latency_comparison.png')
