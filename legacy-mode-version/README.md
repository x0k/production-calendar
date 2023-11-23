# Производственный календарь

Граббер и парсер производственного календаря из [SuperJob](https://superjob.ru/proizvodstvennyj_kalendar)

## Формат данных

```typescript
enum DayType {
  Weekend = 1,
  Holiday = 2,
  PreHoliday = 3,
}

type ProductionCalendar = { [date: string]: DayType }
```

## Запуск:

```bash
curl https://www.superjob.ru/proizvodstvennyj_kalendar/2023/ -L | npm start > calendar.json
```
