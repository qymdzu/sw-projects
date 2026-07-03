// AppButton 单元测试 — Phase B

import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AppButton from '@/components/common/AppButton.vue'

describe('AppButton', () => {
  it('渲染默认插槽内容', () => {
    const wrapper = mount(AppButton, {
      slots: { default: '点我' }
    })
    expect(wrapper.text()).toContain('点我')
  })

  it('click 触发 emit', async () => {
    const wrapper = mount(AppButton, {
      slots: { default: '提交' }
    })
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
    expect(wrapper.emitted('click')?.length).toBe(1)
  })

  it('loading 时禁用点击', async () => {
    const wrapper = mount(AppButton, {
      props: { loading: true },
      slots: { default: '加载中' }
    })
    // el-button 在 loading 时点击会被吞掉
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('disabled 时禁用点击', async () => {
    const wrapper = mount(AppButton, {
      props: { disabled: true },
      slots: { default: '禁用' }
    })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })
})