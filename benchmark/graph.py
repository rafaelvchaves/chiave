import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from matplotlib import rc

rc('font',**{'family':'sans-serif','sans-serif':['DejaVu Sans'],'size':10})
rc('mathtext',**{'default':'regular'})

df = pd.read_csv('results.csv')
df.sort_values(by='TP', inplace=True)
plt.xlabel('throughput (operations per second)')
plt.ylabel('95th percentile latency (microseconds)')
plt.title('Throughput vs. Latency')
plt.plot(df['TP'], df['L95'] / 1000)
plt.savefig('test.png', dpi=300)