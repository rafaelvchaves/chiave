import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from matplotlib import rc

rc('font', **{'family':'sans-serif','sans-serif':['DejaVu Sans'],'size':10})
rc('mathtext', **{'default':'regular'})

graph = 'conv'

df_op = pd.read_csv(f'set/results_op_{graph}.csv')
df_op.sort_values(by='size', inplace=True)
plt.plot(df_op['size'], df_op['time'] / 1e9, label='op')

df_state = pd.read_csv(f'set/results_state_{graph}.csv')
df_state.sort_values(by='size', inplace=True)
plt.plot(df_state['size'], df_state['time'] / 1e9, label='state')

df_delta = pd.read_csv(f'set/results_delta_{graph}.csv')
df_delta.sort_values(by='size', inplace=True)
plt.plot(df_delta['size'], df_delta['time'] / 1e9, label='delta')

plt.xlabel('throughput (operations per second)')
plt.ylabel('95th percentile convergence time (seconds)')
plt.title('Convergence Time vs. Throughput')
plt.legend()
# plt.plot(df['throughput'], df['latency'] / 1000000)
plt.savefig(f'{graph}.png', dpi=300)