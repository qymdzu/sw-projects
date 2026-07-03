// LocalStorage 封装 — Phase B
//
// 提供统一的读写删接口，便于：
//   - SSR / 测试环境 fallback 到内存 Map
//   - 统一错误处理

const memoryStore = new Map<string, string>()

function hasLocalStorage(): boolean {
  try {
    return typeof window !== 'undefined' && !!window.localStorage
  } catch {
    return false
  }
}

export const storage = {
  get(key: string): string | null {
    if (hasLocalStorage()) {
      return window.localStorage.getItem(key)
    }
    return memoryStore.get(key) ?? null
  },

  set(key: string, value: string): void {
    if (hasLocalStorage()) {
      window.localStorage.setItem(key, value)
    } else {
      memoryStore.set(key, value)
    }
  },

  remove(key: string): void {
    if (hasLocalStorage()) {
      window.localStorage.removeItem(key)
    } else {
      memoryStore.delete(key)
    }
  }
}