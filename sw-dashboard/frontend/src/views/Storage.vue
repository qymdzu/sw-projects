<template>
  <div>
    <div class="page-header">
      <h2>存储结构浏览器</h2>
      <div class="page-header-actions">
        <el-radio-group v-model="activeRepo" @change="switchRepo" size="small">
          <el-radio-button value="all">全部</el-radio-button>
          <el-radio-button value="workspace">Workspace</el-radio-button>
          <el-radio-button value="skills-library">Skills Library</el-radio-button>
          <el-radio-button value="projects">Projects</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <el-container style="height: calc(100vh - 160px)">
      <!-- 左侧目录树 -->
      <div class="tree-panel">
        <div class="tree-panel-header">
          <span style="font-size: 13px; font-weight: 500">文件目录</span>
          <el-button text size="small" @click="loadTree(activeRepo)">刷新</el-button>
        </div>
        <div class="tree-panel-body" v-loading="treeLoading">
          <FileTree
            :nodes="treeData"
            :lazy="false"
            @node-click="handleFileSelect"
          />
        </div>
      </div>

      <!-- 右侧文件内容 -->
      <div class="content-panel">
        <!-- 未选择文件 -->
        <div v-if="!selectedFile" style="display: flex; justify-content: center; align-items: center; height: 100%">
          <el-empty description="选择文件以预览" />
        </div>

        <!-- 文件内容 -->
        <template v-else>
          <div class="file-info-bar">
            <div style="display: flex; align-items: center; gap: 8px">
              <span style="font-weight: 500">{{ fileContent?.name }}</span>
              <el-tag size="small" type="info">{{ fileContent?.language }}</el-tag>
              <span style="font-size: 12px; color: var(--text-secondary)">
                {{ fileContent?.size ? (fileContent.size / 1024).toFixed(1) + ' KB' : '' }}
              </span>
            </div>
            <div style="display: flex; gap: 8px">
              <el-button
                v-if="!isEditing && fileContent?.editable"
                size="small"
                type="primary"
                @click="toggleEdit"
              >编辑</el-button>
              <template v-if="isEditing">
                <el-button size="small" type="primary" :loading="saving" @click="handleSaveFile">保存</el-button>
                <el-button size="small" @click="handleCancelEdit">取消</el-button>
              </template>
              <el-button v-else size="small" @click="refreshFile">刷新</el-button>
            </div>
          </div>
          <div v-loading="fileLoading">
            <CodeEditor
              v-if="fileContent"
              :content="isEditing ? editedContent : fileContent.content"
              :language="fileContent.language || 'text'"
              :readonly="!isEditing"
              @update:content="editedContent = $event"
            />
          </div>
        </template>
      </div>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getTree, getFile, saveFile } from '@/api/storage'
import FileTree from '@/components/FileTree.vue'
import type { FileNode } from '@/components/FileTree.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import { showConfirm } from '@/components/ConfirmDialog'
import { ElMessage } from 'element-plus'

const activeRepo = ref('all')
const treeData = ref<FileNode[]>([])
const treeLoading = ref(false)
const selectedFile = ref<string | null>(null)
const fileContent = ref<any>(null)
const fileLoading = ref(false)
const isEditing = ref(false)
const editedContent = ref('')
const saving = ref(false)

async function loadTree(repo: string, path = '/') {
  treeLoading.value = true
  try {
    const res = await getTree(repo, path)
    treeData.value = res.data?.data || res.data || []
  } catch {
    treeData.value = []
  } finally {
    treeLoading.value = false
  }
}

function switchRepo(repo: string) {
  activeRepo.value = repo
  selectedFile.value = null
  fileContent.value = null
  isEditing.value = false
  loadTree(repo)
}

// lazy 模式已禁用（后端一次返回完整树），保留空函数以防 FileTree 调用

async function handleFileSelect(node: FileNode) {
  if (node.type !== 'file') return

  selectedFile.value = node.path
  fileLoading.value = true
  isEditing.value = false

  try {
    const res = await getFile(node.path)
    fileContent.value = res.data?.data || res.data
  } catch {
    fileContent.value = null
  } finally {
    fileLoading.value = false
  }
}

function refreshFile() {
  if (selectedFile.value) {
    handleFileSelect({ path: selectedFile.value, type: 'file' } as FileNode)
  }
}

function toggleEdit() {
  isEditing.value = true
  editedContent.value = fileContent.value?.content || ''
}

async function handleSaveFile() {
  const confirmed = await showConfirm({
    title: '确认保存',
    message: `确认保存对 ${fileContent.value?.name} 的修改？`
  })
  if (!confirmed) return

  saving.value = true
  try {
    await saveFile(selectedFile.value!, editedContent.value)
    ElMessage.success('保存成功')
    isEditing.value = false
    await refreshFile()
  } catch {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function handleCancelEdit() {
  const confirmed = await showConfirm({
    title: '放弃修改',
    message: '放弃未保存的修改？'
  })
  if (confirmed) {
    isEditing.value = false
    editedContent.value = fileContent.value?.content || ''
  }
}

// 初始化加载
loadTree('all')
</script>

<style scoped>
.tree-panel {
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
.content-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
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
