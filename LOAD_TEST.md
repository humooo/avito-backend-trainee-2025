# Результаты нагрузочного тестирования

Проведено тестирование эндпоинта `GET /stats` для проверки производительности сервиса и соответствия SLI (300 мс).

**Инструменты:**
- Utility: ApacheBench (ab)
- Requests: 1000
- Concurrency: 100
- Environment: Docker container

## Результаты

```text
Document Path:          /stats
Document Length:        114 bytes

Concurrency Level:      100
Time taken for tests:   0.471 seconds
Complete requests:      1000
Failed requests:        0
Requests per second:    2121.83 [#/sec] (mean)
Time per request:       47.129 [ms] (mean)
Time per request:       0.471 [ms] (mean, across all concurrent requests)

Percentage of the requests served within a certain time (ms)
  50%     42
  66%     47
  75%     50
  80%     53
  90%     61
  95%     71
  98%     82
  99%     88
 100%    108 (longest request)
