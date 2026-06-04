import { ElMessageBox } from 'element-plus'

export interface ConfirmOptions {
  title?: string
  message: string
  type?: 'warning' | 'info' | 'success' | 'error'
  confirmText?: string
  cancelText?: string
}

export function showConfirm(options: ConfirmOptions): Promise<boolean> {
  return ElMessageBox.confirm(
    options.message,
    options.title || '确认操作',
    {
      confirmButtonText: options.confirmText || '确认',
      cancelButtonText: options.cancelText || '取消',
      type: options.type || 'warning',
      draggable: true
    }
  ).then(() => true).catch(() => false)
}

export function showPrompt(title: string, message: string, inputValue?: string): Promise<string | false> {
  return ElMessageBox.prompt(message, title, {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    inputValue,
    inputPattern: /.*\S.*/,
    inputErrorMessage: '内容不能为空'
  }).then(({ value }) => value).catch(() => false)
}
