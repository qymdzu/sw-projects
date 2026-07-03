// 时间/数字格式化 — Phase B

/**
 * 把 ISO 时间字符串格式化为 "YYYY-MM-DD HH:mm"。
 * 输入非法时返回原字符串。
 */
export function formatDateTime(iso: string | Date | null | undefined): string {
  if (!iso) return '-'
  const d = typeof iso === 'string' ? new Date(iso) : iso
  if (Number.isNaN(d.getTime())) return String(iso)
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/**
 * 把秒数格式化为 "X 分钟" / "X 小时 Y 分"。
 */
export function formatDuration(seconds: number): string {
  if (!seconds || seconds < 0) return '0 分钟'
  if (seconds < 60) return `${seconds} 秒`
  const m = Math.floor(seconds / 60)
  if (m < 60) return `${m} 分钟`
  const h = Math.floor(m / 60)
  const remain = m % 60
  return remain > 0 ? `${h} 小时 ${remain} 分` : `${h} 小时`
}

/**
 * 把小数转百分比字符串，保留 0 位小数。
 */
export function formatPercent(n: number | null | undefined): string {
  if (n === null || n === undefined) return '-'
  return `${Math.round(n * 100)}%`
}

/**
 * 脱敏手机号：138****8000
 */
export function maskPhone(p: string | null | undefined): string {
  if (!p || p.length < 7) return p ?? ''
  return p.slice(0, 3) + '****' + p.slice(-4)
}

/**
 * 脱敏邮箱：a****@example.com
 */
export function maskEmail(e: string | null | undefined): string {
  if (!e) return ''
  const at = e.indexOf('@')
  if (at <= 1) return e
  return e.slice(0, 1) + '****' + e.slice(at - 1)
}