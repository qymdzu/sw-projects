<template>
  <div class="code-editor-wrapper" v-loading="loading">
    <Codemirror
      v-if="!loading"
      :model-value="props.content"
      :extensions="extensions"
      :disabled="props.readonly"
      :style="{ height: height }"
      @update:model-value="onUpdate"
      :placeholder="'加载中...'"
      :autofocus="!props.readonly"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { EditorView } from '@codemirror/view'
import { oneDark } from '@codemirror/theme-one-dark'
import { yaml } from '@codemirror/lang-yaml'
import { python } from '@codemirror/lang-python'
import { javascript } from '@codemirror/lang-javascript'
import { markdown } from '@codemirror/lang-markdown'
import { EditorState } from '@codemirror/state'

const props = withDefaults(defineProps<{
  content?: string
  language?: string
  readonly?: boolean
  height?: string
  loading?: boolean
}>(), {
  content: '',
  language: 'text',
  readonly: true,
  height: 'calc(100vh - 300px)',
  loading: false
})

const emit = defineEmits<{
  (e: 'update:content', content: string): void
}>()

const extensions = computed(() => {
  const exts: any[] = [oneDark]

  switch (props.language) {
    case 'yaml':
      exts.push(yaml())
      break
    case 'python':
      exts.push(python())
      break
    case 'javascript':
    case 'typescript':
      exts.push(javascript())
      break
    case 'markdown':
      exts.push(markdown())
      break
    case 'env':
    case 'text':
    default:
      break
  }

  if (props.readonly) {
    exts.push(EditorView.editable.of(false))
    exts.push(EditorState.readOnly.of(true))
  }

  return exts
})

function onUpdate(value: string) {
  emit('update:content', value)
}
</script>

<style scoped>
.code-editor-wrapper {
  border: 1px solid #333;
  border-radius: 4px;
  overflow: hidden;
}
</style>
