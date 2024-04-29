<template>
  <div style="display: flex; justify-content: space-between; align-items: center">
    <el-text size="large">
      {{ currentDir }}
    </el-text>
    <div>
      <el-button type="primary" @click="back">
        <el-icon>
          <Back/>
        </el-icon>
      </el-button>
      <el-button type="primary" @click="isMkdirDialogShow = true">
        新建目录
      </el-button>
      <el-button type="primary" @click="uploadFiles">
        上传文件
      </el-button>
      <el-button type="primary" @click="uploadDir">
        上传目录
      </el-button>
    </div>
  </div>
  <el-table v-loading="isLoading" :data="fileList" @row-dblclick="openDir">
    <el-table-column width="50px">
      <template #default="scope">
        <el-icon v-if="scope.row.isDir">
          <Folder/>
        </el-icon>
        <el-icon v-else>
          <Document/>
        </el-icon>
      </template>
    </el-table-column>
    <el-table-column label="名称">
      <template #default="scope">
        {{ getBaseName(scope.row.path) }}
      </template>
    </el-table-column>
    <el-table-column label="大小">
      <template #default="scope">
        {{ formatBytes(scope.row.size) }}
      </template>
    </el-table-column>
    <el-table-column label="创建时间">
      <template #default="scope">
        {{ dayjs(scope.row.createTime).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
    </el-table-column>
    <el-table-column label="操作">
      <template #default="scope">
        <el-button type="primary" @click="download(scope.row.id)">
          <el-icon>
            <Download/>
          </el-icon>
        </el-button>

        <el-button type="primary" @click="fileIdToShare = scope.row.id; isShareDialogShow = true">
          <el-icon>
            <Share/>
          </el-icon>
        </el-button>

        <el-button type="danger" @click="fileToDelete = scope.row; isDeleteDialogShow = true">
          <el-icon>
            <Delete/>
          </el-icon>
        </el-button>
      </template>
    </el-table-column>
  </el-table>

  <el-dialog v-model="isMkdirDialogShow">
    <template #header>
      创建新目录
    </template>

    <el-form-item label="新目录名">
      <el-input v-model="newDir">

      </el-input>
    </el-form-item>

    <template #footer>
      <el-button type="primary" @click="mkdir">
        创建
      </el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="isDeleteDialogShow">
    <template #header>
      确认删除？
    </template>
    <div v-if="fileToDelete.isDir">
      确认要删除{{ getBaseName(fileToDelete.path) }}及其子文件吗？
    </div>
    <div v-else>
      确认要删除{{ getBaseName(fileToDelete.path) }}吗？
    </div>

    <template #footer>
      <el-button type="primary" @click="isDeleteDialogShow = false">
        取消
      </el-button>
      <el-button type="danger" @click="deleteFile">
        删除
      </el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="isShareDialogShow">
    <template #header>
      分享
    </template>

    <el-form-item label="有效期">
      <el-radio-group v-model="expireInSecond">
        <el-radio :value="60 * 60 * 24" label="1 天"/>
        <el-radio :value="60 * 60 * 24 * 3" label="3 天"/>
        <el-radio :value="60 * 60 * 24 * 7" label="7 天"/>
      </el-radio-group>
    </el-form-item>
    <template #footer>
      <el-button type="primary" @click="share">创建分享链接</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {DeleteFile, DownloadFile, GetFileList, GetShareUrl, Mkdir, UploadFiles} from "../../wailsjs/go/main/App";
import {ElMessage} from "element-plus";
import dayjs from "dayjs";

const fileList = ref([])
const isLoading = ref(false)
const currentDir = ref('/')

function formatBytes(bytes: number): string {
  if (!bytes) {
    return '-'
  }
  if (bytes < 1024) {
    return bytes + " Bytes";
  }
  if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(2) + " KB";
  }
  if (bytes < 1024 * 1024 * 1024) {
    return (bytes / (1024 * 1024)).toFixed(2) + " MB";
  }
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
}


async function loadFileList() {
  isLoading.value = true
  const result = await GetFileList(currentDir.value)
  isLoading.value = false
  if (result.code != 2000) {
    ElMessage.error(result.msg)
    return
  }
  fileList.value = result.data
}

onMounted(loadFileList)

// TODO 进度展示
async function uploadFiles() {
  const result = await UploadFiles(currentDir.value)
  if (result.code != 2000) {
    ElMessage.error(result.msg)
  }
  if (result.data > 0) {
    ElMessage.success(`成功上传${result.data}个文件`)
  }
  await loadFileList()
}

onMounted(loadFileList)

async function download(id: number) {
  let result = await DownloadFile(id.toString());
  if (result.code != 2000) {
    ElMessage.error(result.msg)
    return
  }
  ElMessage.success('下载成功')
}

const isMkdirDialogShow = ref(false)
const newDir = ref('')

async function mkdir() {
  const result = await Mkdir(currentDir.value, newDir.value)
  if (result.code != 2000) {
    ElMessage.error(result.msg)
    return
  }
  ElMessage.success('创建成功')
  isMkdirDialogShow.value = false
  newDir.value = ''
  await loadFileList()
}

const openDir = async (row: any, column: any, event: Event): Promise<void> => {
  if (row.isDir) {
    currentDir.value = row.path + '/'
    await loadFileList()
  }
}

async function back() {
  const withoutLastSlash = currentDir.value.substring(0, currentDir.value.length - 1)
  currentDir.value = withoutLastSlash.substring(0, withoutLastSlash.lastIndexOf('/') + 1)
  await loadFileList()
}

async function uploadDir() {
  // TODO 上传目录
}

const isShareDialogShow = ref(false)
const expireInSecond = ref(60 * 60 * 24 * 7)
const fileIdToShare = ref(0)

async function share() {
  const result = await GetShareUrl(String(fileIdToShare.value), expireInSecond.value)
  if (result.code != 2000) {
    ElMessage.error(result.msg)
    return
  }
  ElMessage.success('分享链接已复制')
  isShareDialogShow.value = false
}

const isDeleteDialogShow = ref(false)
const fileToDelete: any = ref({})

async function deleteFile() {
  const result = await DeleteFile(String(fileToDelete.value.id))
  if (result.code != 2000) {
    ElMessage.error(result.msg)
    return
  }
  ElMessage.success('删除成功')
  isDeleteDialogShow.value = false
  await loadFileList()
}

function getBaseName(path: string): string {
  return path.substring(path.lastIndexOf('/') + 1)
}
</script>

<style scoped>

</style>