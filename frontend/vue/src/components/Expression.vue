<script setup lang="ts">
import {
  defineProps
} from 'vue'
import {
  NTag
} from 'naive-ui'

defineProps<{
  id: string,
  result: number,
  status: string
}>()
</script>

<template>
  <div class="flex justify-center flex-col d-wrap">
    <table class="t">
      <tbody>
        <tr>
          <th>ID</th>
          <td>{{id}}</td>
        </tr>
        <tr>
          <th>Result</th>
          <td>{{result}}</td>
        </tr>
        <tr>
          <th>Status</th>
          <td>
            <n-tag
                v-if="status
          .trim()
          .toLowerCase()==='completed'"
                type="success"
            >
              Completed
            </n-tag>
            <n-tag
                v-else-if="status
          .trim()
          .toLowerCase()==='failed'"
                type="error"
            >
              Failed
            </n-tag>
            <n-tag
                v-else-if="status
          .trim()
          .toLowerCase()==='pending'"
                type="warning"
            >
              Pending
            </n-tag>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
  .t {
    border-collapse: collapse;
    display: table;
    table-layout: fixed;
    text-align: center;
  }
  .t td, .t th {
    padding: 0.3em 1em 0.3em 1em;
  }
  .t th {
    border-right: 0.1em solid #2A2A2E;
  }
  .d-wrap {
    height: fit-content;
    width: fit-content;
    padding: 0.5em 0.5em 0.5em 0.5em;
    margin: 0.5em 0.5em 0.5em 0.5em;
    border: 0.1em solid #2A2A2E;
    border-radius: 0.7em;
    background-color: #19191D;
  }
</style>