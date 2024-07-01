import matplotlib.pyplot as plt

# Данные
users = [1, 10, 100, 1000]
latency_before = [221, 867, 10879, 6159]  # Latency до индекса
throughput_before = [4.5, 7.6, 8.5, 55.6]  # Throughput до индекса
latency_after = [205, 1081, 1076, 6313]  # Latency после индекса
throughput_after = [4.9, 7.4, 23.6, 58.1]  # Throughput после индекса

# Создание графиков
fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(15, 5))

# График latency
ax1.plot(users, latency_before, marker='o', label='Latency до индекса')
ax1.plot(users, latency_after, marker='o', label='Latency после индекса')
ax1.set_title('Latency до и после индекса')
ax1.set_xlabel('Количество пользователей')
ax1.set_ylabel('Среднее время отклика (мс)')
ax1.set_xscale('log')
ax1.set_yscale('log')
ax1.legend()
ax1.grid(True)

# График throughput
ax2.plot(users, throughput_before, marker='o', label='Throughput до индекса')
ax2.plot(users, throughput_after, marker='o', label='Throughput после индекса')
ax2.set_title('Throughput до и после индекса')
ax2.set_xlabel('Количество пользователей')
ax2.set_ylabel('Пропускная способность (запросов/сек)')
ax2.set_xscale('log')
ax2.set_yscale('log')
ax2.legend()
ax2.grid(True)

# Сохранение графиков
plt.tight_layout()
plt.savefig('latency_throughput_comparison.png')
plt.show()

