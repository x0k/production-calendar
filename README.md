# Производственный календарь

Граббер и парсер производственного календаря из [SuperJob](https://superjob.ru/proizvodstvennyj_kalendar)

## Выходной формат данных

```yaml
type: object
patternProperties:
  ^[0-9]{4}-[0-9]{2}-[0-9]{2}$:
    enum:
      - 1 # Weekend
      - 2 # Holiday
      - 3 # PreHoliday
```

## Запуск:

```bash
curl https://www.superjob.ru/proizvodstvennyj_kalendar/2024/ -L | go run . > calendar.json
```
