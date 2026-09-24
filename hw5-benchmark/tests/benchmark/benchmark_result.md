# Бенчмарк-результаты

**Платформа:**

* OS: darwin
* Arch: arm64
* CPU: Apple M4

## Результаты

| Benchmark                                | Итерации (N) | Время (ns/op) | Память (B/op) | Аллокации (allocs/op) |
| ---------------------------------------- | -----------: | ------------: | ------------: | --------------------: |
| BenchmarkStudentName_Direct-10           |  643,561,899 |         1.638 |             0 |                     0 |
| BenchmarkStudentName_Interface-10        |  713,676,289 |         1.697 |             0 |                     0 |
| BenchmarkStudentName_MethodValue-10      |  749,384,098 |         1.602 |             0 |                     0 |
| BenchmarkStudentName_MethodExpression-10 |  747,289,132 |         1.605 |             0 |                     0 |
| BenchmarkStudentName_Reflect-10          |   12,060,310 |         99.62 |            48 |                     3 |


## Выводы

* Direct / Interface / MethodValue / MethodExpression имеют близкую производительность без аллокаций.
* Reflect заметно медленнее и даёт дополнительные аллокации.
