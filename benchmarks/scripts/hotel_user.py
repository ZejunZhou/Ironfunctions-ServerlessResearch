import matplotlib.pyplot as plt
import pandas as pd

# 数据
data = {
    'Nightcore RPS 50': {
        'p50': 13.74,
        'p99': 43.29
    },
    'Nightcore RPS 100': {
        'p50': 13.99,
        'p99': 32.43
    },
    'Iron Function RPS 50': {
        'p50': 15.72,
        'p99': 57.29
    },
    'Iron Function RPS 100': {
        'p50': 18.18,
        'p99': 64.34
    }
}

p50_rps50_nightcore = data['Nightcore RPS 50']['p50']
p99_rps50_nightcore = data['Nightcore RPS 50']['p99']
p50_rps100_nightcore = data['Nightcore RPS 100']['p50']
p99_rps100_nightcore = data['Nightcore RPS 100']['p99']

p50_rps50_iron = data['Iron Function RPS 50']['p50']
p99_rps50_iron = data['Iron Function RPS 50']['p99']
p50_rps100_iron = data['Iron Function RPS 100']['p50']
p99_rps100_iron = data['Iron Function RPS 100']['p99']

df = pd.DataFrame({
    'p50(RPS50)': [p50_rps50_nightcore, p50_rps50_iron],
    'p99(RPS50)': [p99_rps50_nightcore, p99_rps50_iron],
    'p50(RPS100)': [p50_rps100_nightcore, p50_rps100_iron],
    'p99(RPS100)': [p99_rps100_nightcore, p99_rps100_iron]
}, index=['Nightcore(RPS 50 and RPS 100)', 'Iron Function(RPS 50 and RPS 100)'])

ax = df.plot(kind='bar', figsize=(10, 6))
plt.ylabel('Latency (ms)')
plt.xlabel('Median and Tail Latency Comparison')
plt.title('Latency Comparison: Nightcore vs Iron Function')
plt.xticks(rotation=0)
plt.legend(title='Metrics')

# Adding value labels on top of the bars
for p in ax.patches:
    ax.annotate(f'{p.get_height():.2f}', (p.get_x() + p.get_width() / 2., p.get_height()),
                ha='center', va='center', xytext=(0, 10), textcoords='offset points')

plt.tight_layout()
plt.savefig('latency_comparison.png')