<template>
  <el-container>
    <el-header>
      <div><h1>s3 as drive</h1></div>
    </el-header>
    <el-container>
      <el-aside width="200px">
        <el-menu :router="true">
          <el-menu-item index="/home/file-browser">
            文件
          </el-menu-item>
          <el-menu-item index="/home/config">
            设置
          </el-menu-item>
          <el-menu-item index="/home/about">
            关于
          </el-menu-item>
          <el-menu-item index="/home/sync-dir">
            同步目录
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <RouterView/>
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>


</script>

<style scoped>

</style>