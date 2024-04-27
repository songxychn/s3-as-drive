<template>
  <div style="display: flex; justify-content: space-between; align-items: center">
    <el-text size="large">
      {{ dir }}
    </el-text>
    <el-button type="primary" @click="uploadFiles">
      <el-icon>
        <Upload/>
      </el-icon>
    </el-button>
  </div>
  <el-table v-loading="isLoading" :data="fileList">
    <el-table-column label="名称">
      <template #default="scope">
        {{ scope.row.path.substring(scope.row.path.lastIndexOf('/') + 1) }}
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
</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {DownloadFile, GetFileList, SelectFiles} from "../../wailsjs/go/main/App";
import {ElMessage} from "element-plus";

const fileList = ref([])
const isLoading = ref(false)
const dir = ref('/')

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
  const result = await GetFileList()
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
  const result = await SelectFiles(dir.value)
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
</script>

<style scoped>

</style>