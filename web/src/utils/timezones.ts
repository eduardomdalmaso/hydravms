// HydraVMS - Native IANA Timezone Discovery (< 40 lines)

export function getAvailableTimezones(): string[] {
  try {
    if (typeof Intl !== 'undefined' && typeof Intl.supportedValuesOf === 'function') {
      const all = Intl.supportedValuesOf('timeZone')
      const priority = [
        'America/Sao_Paulo', 'America/Manaus', 'America/Cuiaba', 'America/Belem',
        'America/Fortaleza', 'America/Recife', 'America/Noronha', 'UTC',
        'America/New_York', 'Europe/London', 'Europe/Lisbon', 'Europe/Madrid'
      ]
      const remaining = all.filter(tz => !priority.includes(tz))
      return [...priority, ...remaining]
    }
  } catch {}
  return ['America/Sao_Paulo', 'America/Manaus', 'America/Cuiaba', 'UTC', 'America/New_York', 'Europe/London']
}

export function detectLocalTimezone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/Sao_Paulo'
  } catch {
    return 'America/Sao_Paulo'
  }
}
