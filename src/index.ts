import { HTMLElement, parse } from 'node-html-parser'

enum DayType {
  Weekend = 1,
  Holiday = 2,
  PreHoliday = 3,
}

function assert<T>(
  value: T | null | undefined,
  msg = 'value not present'
): asserts value is T {
  if (value === undefined || value === null) {
    throw new TypeError(msg)
  }
}

function getDayType(dayText: string): DayType {
  const lowered = dayText.toLocaleLowerCase('ru')
  switch (true) {
    case lowered.startsWith('праздничный'):
      return DayType.Holiday
    case lowered.startsWith('выходной'):
      return DayType.Weekend
    case lowered.includes('сокращен'):
      return DayType.PreHoliday
    default:
      throw new TypeError(`Unknown day type with text: ${dayText}`)
  }
}

function makeProductionCalendar(html: HTMLElement): Record<string, DayType> {
  const calendarTitleElement = html.querySelector('h1')
  assert(calendarTitleElement, 'Calendar title')
  const yearFullNumber = calendarTitleElement.text.match(/\d+/)?.[0]
  assert(yearFullNumber, 'Calendar year')
  const calendarRootElement = calendarTitleElement.parentNode
  assert(calendarRootElement, 'Calendar root element')

  const months = [
    'yanvar',
    'fevral',
    'mart',
    'aprel',
    'maj',
    'iyun',
    'iyul',
    'avgust',
    'sentyabr',
    'oktbyar',
    'noyabr',
    'dekabr',
  ]

  const monthsData = months.flatMap((month, monthIndex) => {
    const monthRootElement: HTMLElement | undefined =
      calendarRootElement.querySelector(`[href$="/${month}/"]`)?.parentNode
        ?.parentNode
    assert(monthRootElement, `Month element "${month}"`)

    const days = monthRootElement.querySelectorAll('._1c_LS,._1YS-8')
    const monthNumber = String(monthIndex + 1).padStart(2, '0')
    return days.map((dayElement): [string, DayType] => {
      const dayNumberText = dayElement.querySelector('[role="button"]')?.text
      const dayText =
        dayElement.querySelector('[role="tooltip"]')?.firstChild.text
      assert(dayNumberText, 'Day number')
      assert(dayText, 'Day text')
      const dayNumber = dayNumberText.padStart(2, '0')
      return [
        `${yearFullNumber}-${monthNumber}-${dayNumber}`,
        getDayType(dayText),
      ]
    })
  })

  return Object.fromEntries(monthsData)
}

const inputLines: string[] = []

process.stdin.on('data', (data) => {
  inputLines.push(data.toString())
})

process.stdin.once('end', () => {
  const calendar = makeProductionCalendar(parse(inputLines.join('')))
  process.stdout.write(JSON.stringify(calendar, null, 2))
})
