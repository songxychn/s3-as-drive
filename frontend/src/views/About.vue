<template>
  <div style="width: 100%; text-align: center">
    <img src="@/assets/logo.svg"/>
    <h1>S3AsDrive</h1>
    <el-text>v1.0.0</el-text>
    <div style="padding: 10px 0;color: #606266">
      <el-link @click="BrowserOpenURL('https://github.com/songxychn/s3-as-drive')">源码地址</el-link>
      <el-divider direction="vertical"/>
      <!--TODO-->
      <el-link>官方网站</el-link>
    </div>
    <el-text>Created by songxychina@gmail.com</el-text>
  </div>
</template>

<script lang="ts" setup>
import {BrowserOpenURL} from "../../wailsjs/runtime";
</script>

<style scoped>

</style>