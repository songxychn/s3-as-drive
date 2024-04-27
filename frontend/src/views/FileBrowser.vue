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
        <el-icon>
          <Plus/>
        </el-icon>
      </el-button>
      <el-button type="primary" @click="uploadFiles">
        <el-icon>
          <Upload/>
        </el-icon>
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
        {{ scope.row.path.substring(scope.row.path.lastIndexOf('/') + 1, scope.row.path.length) }}
      </template>
    </el-table-column>
    <!--    <el-table-column label="大小">-->
    <!--      <template #default="scope">-->
    <!--        {{ formatBytes(scope.row.size)}}-->
    <!--      </template>-->
    <!--    </el-table-column>-->
    <el-table-column label="操作">
      <template #default="scope">
        <el-button type="primary" @click="download(scope.row.id)">
          <el-icon>
            <Download/>
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
</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {DownloadFile, GetFileList, Mkdir, UploadFiles} from "../../wailsjs/go/main/App";
import {ElMessage} from "element-plus";

const fileList = ref([])
const isLoading = ref(false)
const currentDir = ref('/')

function formatBytes(bytes: number): string {
  if (bytes < 1024) {
    return bytes + " Bytes";
  } else if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(2) + " KB";
  } else if (bytes < 1024 * 1024 * 1024) {
    return (bytes / (1024 * 1024)).toFixed(2) + " MB";
  } else {
    return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
  }
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
</script>

<style scoped>

</style>