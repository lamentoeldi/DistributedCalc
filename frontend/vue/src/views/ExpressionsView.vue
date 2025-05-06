<script setup lang="ts">
import {
  NSpace
} from 'naive-ui'
import ExpressionsTable from "@/components/ExpressionsTable.vue";
</script>

<template>
  <div class="flex justify-center items-center container">
    <ExpressionsTable/>
  </div>
</template>

<style scoped>
  .container {
    height: 100%;
  }
</style>