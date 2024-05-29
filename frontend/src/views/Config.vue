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
import {onMounted, reactive} from "vue";
import {ElMessage} from "element-plus";
import {GetConfig, UpdateConfig} from "../../wailsjs/go/services/ConfigService";
import {types} from "../../wailsjs/go/models";
import Config = types.Config;

const config: Config = reactive(new Config({
  s3Config: {},
  downloadConfig: {}
}))

async function loadConfig() {
  try {
    const result = await GetConfig();
    Object.assign(config, result)
  } catch (e: any) {
    ElMessage.error(e)
  }
}

onMounted(loadConfig)

async function updateConfig() {
  try {
    await UpdateConfig(config);
    ElMessage.success('保存成功')
  } catch (e: any) {
    ElMessage.error(e)
  }
}
</script>

<style scoped>

</style>