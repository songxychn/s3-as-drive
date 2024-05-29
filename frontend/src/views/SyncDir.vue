<template>
  <div style="display: flex; justify-content: flex-end;width: 100%">
    <el-button type="primary" @click="isCreateSyncDirDialogShow = true">
      <el-icon>
        <Plus/>
      </el-icon>
    </el-button>
  </div>
  <el-table>
    <el-table-column label="名称"></el-table-column>
    <el-table-column label="路径"></el-table-column>
  </el-table>

  <el-dialog v-model="isCreateSyncDirDialogShow">
    <template #header>
      创建同步目录
    </template>
    <el-form ref="createSyncDirFormRef" :model="createSyncDirReq" :rules="rules" label-position="left" label-width="auto"
             status-icon>
      <el-form-item label="名称" prop="name">
        <el-input v-model="createSyncDirReq.name"></el-input>
      </el-form-item>
      <el-form-item label="路径" prop="path">
        <el-input v-model="createSyncDirReq.path" readonly>
          <template #append>
            <el-button @click="choseDir">
              选择目录
            </el-button>
          </template>
        </el-input>
      </el-form-item>
    </el-form>
    <div style="display: flex; justify-content: center;align-items: center">
      <el-button type="primary" @click="createSyncDir">创建</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts" setup>
import {ChoseDir, CreateSyncDir,} from '../../wailsjs/go/services/SyncDirService'
import {reactive, ref} from "vue";
import {ElMessage, FormInstance, FormRules} from "element-plus";

const isCreateSyncDirDialogShow = ref(false)

const rules = reactive<FormRules<CreateSyncDirReq>>({
  name: [
    {required: true, message: '请输入名称', trigger: 'blur'},
  ],
  path: [
    {required: true, message: '请选择目录', trigger: 'blur'},
  ],
})
const createSyncDirFormRef = ref<FormInstance>()

class CreateSyncDirReq {
  name: string
  path: string

  constructor() {
    this.name = ""
    this.path = ""
  }
}

const createSyncDirReq: CreateSyncDirReq = reactive(new CreateSyncDirReq())

async function createSyncDir() {
  try {
    await CreateSyncDir(createSyncDirReq.name, createSyncDirReq.path);
    ElMessage.success('创建成功')
    isCreateSyncDirDialogShow.value = false
  } catch (e: any) {
    ElMessage.error(e)
  }
}

async function choseDir() {
  try {
    const dir = await ChoseDir()
    createSyncDirReq.path = dir
  } catch (e: any) {
    ElMessage.error(e)
  }

}
</script>

<style scoped>

</style>