<template>
  <div style="display: flex; justify-content: space-between; align-items: center">
    <el-text size="large">
      {{ currentDir }}
    </el-text>
    <div>
      <el-button :disabled="currentDir === '/'" type="primary" @click="back">
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

        <el-button :disabled="scope.row.isDir" type="primary"
                   @click="fileIdToShare = scope.row.id; isShareDialogShow = true">
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
import {onMounted, reactive, ref} from "vue";
import {
  DeleteFile,
  DownloadFile,
  GetFileList,
  GetShareUrl,
  Mkdir,
  UploadDir,
  UploadFiles
} from "../../wailsjs/go/service/FileService";
import {ElMessage} from "element-plus";
import dayjs from "dayjs";
import {types} from "../../wailsjs/go/models";

const fileList: types.File[] = reactive([])
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
  try {
    const result = await GetFileList(currentDir.value)
    fileList.splice(0, fileList.length)
    result.forEach((item) => {
      fileList.push(item)
    })
    isLoading.value = false
  } catch (e: any) {
    ElMessage.error(e)
  }
}

onMounted(loadFileList)

// TODO 进度展示
async function uploadFiles() {
  try {
    const result = await UploadFiles(currentDir.value)
    if (result > 0) {
      ElMessage.success(`成功上传${result}个文件`)
    }
    await loadFileList()
  } catch (e: any) {
    ElMessage.error(e)
  }
}

onMounted(loadFileList)

async function download(id: number) {
  try {
    await DownloadFile(id);
    ElMessage.success('下载成功')
  } catch (e: any) {
    ElMessage.error(e)
  }
}


const isMkdirDialogShow = ref(false)
const newDir = ref('')

async function mkdir() {
  try {
    await Mkdir(currentDir.value, newDir.value);
    ElMessage.success('创建成功')
    isMkdirDialogShow.value = false
    newDir.value = ''
    await loadFileList()
  } catch (e: any) {
    ElMessage.error(e)
  }
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
  try {
    await UploadDir(currentDir.value)
    ElMessage.success('上传成功')
    await loadFileList()
  } catch (e: any) {
    ElMessage.error(e)
  }
}

const isShareDialogShow = ref(false)
const expireInSecond = ref(60 * 60 * 24 * 7)
const fileIdToShare = ref(0)

async function share() {
  try {
    const result = await GetShareUrl(fileIdToShare.value, expireInSecond.value)
    ElMessage.success('分享链接已复制')
    isShareDialogShow.value = false
  } catch (e: any) {
    ElMessage.error(e)
  }
}

const isDeleteDialogShow = ref(false)
const fileToDelete: any = ref({})

async function deleteFile() {
  try {
    await DeleteFile(fileToDelete.value.id)
    ElMessage.success('删除成功')
    isDeleteDialogShow.value = false
    await loadFileList()
  } catch (e: any) {
    ElMessage.error(e)
  }
}

function getBaseName(path: string): string {
  return path.substring(path.lastIndexOf('/') + 1)
}
</script>

<style scoped>

</style>