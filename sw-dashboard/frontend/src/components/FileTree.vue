<template>
  <el-tree
    :data="nodes"
    :props="treeProps"
    :lazy="lazy"
    :load="loadNode"
    :highlight-current="true"
    node-key="path"
    @node-click="handleNodeClick"
    :expand-on-click-node="false"
    :default-expand-all="!lazy"
  >
    <template #default="{ node, data }">
      <span style="display: flex; align-items: center; gap: 4px; font-size: 13px">
        <el-icon v-if="data.type === 'directory' || data.repo" style="color: #e6a23c">
          <FolderOpened />
        </el-icon>
        <el-icon v-else style="color: #409eff">
          <Document />
        </el-icon>
        <span>{{ data.name ?? data.repo ?? '(未命名)' }}</span>
      </span>
    </template>
  </el-tree>
</template>

<script setup lang="ts">
import { FolderOpened, Document } from '@element-plus/icons-vue'
import type { ElTree } from 'element-plus'

export interface FileNode {
  name: string
  path: string
  type: 'file' | 'directory'
  isLeaf: boolean
  size?: number
  mtime?: string
  children?: FileNode[]
  [key: string]: any
}

const props = withDefaults(defineProps<{
  nodes?: FileNode[]
  lazy?: boolean
  loadNode?: (node: any, resolve: (data: FileNode[]) => void) => void
}>(), {
  nodes: () => [],
  lazy: true
})

const emit = defineEmits<{
  (e: 'node-click', node: FileNode): void
}>()

const treeProps = {
  children: 'children',
  label: 'name',
  isLeaf: 'isLeaf'
}

function handleNodeClick(data: FileNode) {
  if (data.type === 'file') {
    emit('node-click', data)
  }
}
</script>
