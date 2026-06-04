<template>
  <div>
    <div class="page-header">
      <h2>技能浏览器</h2>
      <div class="page-header-actions">
        <el-button size="small" @click="loadSkillTree()" :icon="Refresh">刷新</el-button>
      </div>
    </div>

    <el-container style="height: calc(100vh - 160px)">
      <!-- 左侧技能目录树 -->
      <div class="skill-tree-panel">
        <div class="tree-panel-header">
          <span style="font-size: 13px; font-weight: 500">已安装技能</span>
        </div>
        <div class="tree-panel-body" v-loading="treeLoading">
          <FileTree
            :nodes="skillTree"
            :lazy="true"
            :load-node="lazyLoadSkillNode"
            @node-click="handleSkillSelect"
          />
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="skill-content-panel">
        <!-- 未选择 -->
        <div v-if="!selectedNode" style="display: flex; justify-content: center; align-items: center; height: 100%">
          <el-empty description="选择一个技能查看详情" />
        </div>

        <template v-else>
          <!-- 加载中 -->
          <div v-if="contentLoading" v-loading="true" style="height: 100%"></div>

          <!-- 内容 -->
          <template v-else-if="skillContent">
            <div class="file-info-bar">
              <div style="display: flex; align-items: center; gap: 8px">
                <span style="font-weight: 500">{{ skillContent.name }}</span>
                <el-tag size="small" type="info">{{ skillContent.language }}</el-tag>
                <span style="font-size: 12px; color: var(--text-secondary)">
                  {{ skillContent.size ? (skillContent.size / 1024).toFixed(1) + ' KB' : '' }}
                </span>
              </div>
            </div>
            <!-- Markdown 渲染 -->
            <div v-if="skillContent.language === 'markdown'" class="markdown-preview" v-html="renderedMarkdown" />
            <!-- 其他源码只读展示 -->
            <CodeEditor
              v-else
              :content="skillContent.content"
              :language="skillContent.language || 'text'"
              :readonly="true"
              height="calc(100vh - 220px)"
            />
          </template>
        </template>
      </div>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { getTree, getFile } from '@/api/skills'
import FileTree from '@/components/FileTree.vue'
import type { FileNode } from '@/components/FileTree.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import { Refresh } from '@element-plus/icons-vue'
import { marked } from 'marked'

const skillTree = ref<FileNode[]>([])
const treeLoading = ref(false)
const selectedNode = ref<FileNode | null>(null)
const skillContent = ref<any>(null)
const contentLoading = ref(false)

const renderedMarkdown = computed(() => {
  if (!skillContent.value?.content) return ''
  try {
    return marked.parse(skillContent.value.content) as string
  } catch {
    return skillContent.value.content
  }
})

async function loadSkillTree(path = '/') {
  treeLoading.value = true
  try {
    const res = await getTree(path)
    skillTree.value = res.data?.data || res.data || []
  } catch {
    skillTree.value = []
  } finally {
    treeLoading.value = false
  }
}

function lazyLoadSkillNode(node: any, resolve: (data: FileNode[]) => void) {
  if (node.data.type !== 'directory') {
    resolve([])
    return
  }
  getTree(node.data.path)
    .then(res => resolve(res.data?.data || res.data || []))
    .catch(() => resolve([]))
}

async function handleSkillSelect(node: FileNode) {
  if (node.type !== 'file') return

  selectedNode.value = node
  contentLoading.value = true

  try {
    const res = await getFile(node.path)
    skillContent.value = res.data?.data || res.data
  } catch {
    skillContent.value = null
  } finally {
    contentLoading.value = false
  }
}

loadSkillTree()
</script>

<style scoped>
.skill-tree-panel {
  width: 300px;
  min-width: 200px;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
}
.tree-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-color);
}
.tree-panel-body {
  flex: 1;
  overflow: auto;
  padding: 8px;
}
.skill-content-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: auto;
  padding-left: 12px;
}
.file-info-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 8px;
}
</style>
