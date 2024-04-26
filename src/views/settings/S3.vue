<template>
  <h1>s3 配置</h1>
  <el-form label-position="left" label-width="auto">
    <el-form-item label="endpoint">
      <el-input v-model="s3Config.endpoint">
      </el-input>
    </el-form-item>
    <el-form-item label="accessKey">
      <el-input v-model="s3Config.accessKey">
      </el-input>
    </el-form-item>
    <el-form-item label="secretKey">
      <el-input v-model="s3Config.secretKey">
      </el-input>
    </el-form-item>
    <el-form-item label="bucket">
      <el-input v-model="s3Config.bucket">
      </el-input>
    </el-form-item>
    <el-form-item>
      <div style="width: 100%;display: flex;justify-content: center;align-items: center">
        <el-button type="primary" @click="updateS3Config">保存</el-button>
      </div>
    </el-form-item>
  </el-form>


</template>

<script lang="ts" setup>
import {onMounted, reactive} from "vue";
import {invoke} from "@tauri-apps/api/tauri";
import {ElMessage} from "element-plus";

class S3Config {
  endpoint: string
  accessKey: string
  secretKey: string
  bucket: string

  constructor() {
    this.endpoint = ''
    this.accessKey = ''
    this.secretKey = ''
    this.bucket = ''
  }
}

const s3Config = reactive(new S3Config())
onMounted(async () => {
  let fromRust = await invoke("get_s3_config");
  Object.assign(s3Config, fromRust)
})

async function updateS3Config() {
  await invoke('update_s3_config', {s3Config: s3Config})
  ElMessage.success('保存成功')
}
</script>

<style scoped>

</style>