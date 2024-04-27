<template>
  <h1>s3 配置</h1>
  <el-form label-position="left" label-width="auto">
    <el-form-item label="endpoint">
      <el-input v-model="config.s3Config.endpoint">
      </el-input>
    </el-form-item>
    <el-form-item label="accessKey">
      <el-input v-model="config.s3Config.accessKey">
      </el-input>
    </el-form-item>
    <el-form-item label="secretKey">
      <el-input v-model="config.s3Config.secretKey">
      </el-input>
    </el-form-item>
    <el-form-item label="bucket">
      <el-input v-model="config.s3Config.bucket">
      </el-input>
    </el-form-item>
  </el-form>
  <el-divider/>

  <h1>下载设置</h1>
  <el-form label-position="left" label-width="auto">
    <el-form-item label="下载目录">
      <el-input v-model="config.downloadConfig.dir">
      </el-input>
    </el-form-item>
  </el-form>

  <div style="width: 100%;display: flex;justify-content: center;align-items: center">
    <el-button type="primary" @click="updateConfig">保存</el-button>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {ElMessage} from "element-plus";
import {GetConfig, UpdateConfig} from "../../wailsjs/go/main/App";
import {types} from "../../wailsjs/go/models";
import Config = types.Config;
import S3Config = types.S3Config;
import DownloadConfig = types.DownloadConfig;

const config: any = ref(new Config({
  s3Config: new S3Config({
    endpoint: '',
    accessKey: '',
    secretKey: '',
    bucket: '',
  }),
  downloadConfig: new DownloadConfig({
    dir: ''
  })
}))

async function loadConfig() {
  const result = await GetConfig();
  if (result.code != 2000) {
    ElMessage.error(result.msg)
    return
  }
  config.value = result.data
}

onMounted(loadConfig)

async function updateConfig() {
  const resp = await UpdateConfig(config.value);
  if (resp.code == 2000) {
    ElMessage.success('保存成功')
  } else {
    ElMessage.error(resp.msg)
  }
}
</script>

<style scoped>

</style>