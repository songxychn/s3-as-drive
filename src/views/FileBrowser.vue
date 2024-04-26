<template>
  <el-table v-loading="isLoading" :data="fileList">
    <el-table-column label="名称">
      <template #default="scope">
        {{ scope.row.path.substring(scope.row.path.lastIndexOf('/') + 1)}}
      </template>
    </el-table-column>
    <el-table-column label="大小">
      <template #default="scope">
        {{ formatBytes(scope.row.size)}}
      </template>
    </el-table-column>
  </el-table>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from "vue";

const fileList: any[] = reactive([])
const isLoading = ref(false)

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

}

onMounted(loadFileList)

</script>

<style scoped>

</style>