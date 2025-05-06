<script setup lang="ts">
import {
  defineProps,
  ref
} from 'vue'

import {
  NTag,
  NPagination,
  NButton,
  NIcon,
  useMessage
} from 'naive-ui'

import {
  RefreshOutline as RefreshIcon
} from '@vicons/ionicons5'

let expressions = (() =>
    Array.from({ length: 15 }, (_, i) => ({
      id: `${i + 1}`,
      expression: '1+3+2',
      result: 12,
      status: ['completed', 'failed', 'pending'][(i % 3)],
    }))
)()

const chunkArray = <T>(array: T[], size: number): T[][] =>
    array.reduce((acc, val, i) => {
      if (i % size === 0) acc.push([])
      acc[acc.length - 1].push(val)
      return acc;
    }, [] as T[][])

const msg = useMessage()

const fetchExps = async () => {
  let loading = msg.loading("Fetching more expressions...")

  const newExps = (() =>
      Array.from({ length: 15 }, (_, i) => ({
        id: `${i + 1}`,
        expression: '1+3+2',
        result: 12,
        status: ['completed', 'failed', 'pending'][(i % 3)],
      }))
  )()

  loading.destroy()
  msg.success("New expressions fetched")

  expressions = expressions.concat(newExps)

  expsToRender.value = chunkArray(expressions, 7)
  pages.value = expsToRender.value.length
}

let expsToRender = ref(chunkArray(
    expressions, 7
))

const page = ref(1)
const pages = ref(expsToRender.value.length)
</script>

<template>
  <div class="flex items-center justify-center flex-col">
    <div id="table-div" class="flex justify-center items-center">
      <table class="exp-table">
        <thead>
          <tr>
            <th>
              ID
            </th>
            <th>
              Expression
            </th>
            <th>
              Result
            </th>
            <th>
              Status
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="expression in expsToRender[page-1]" :key="expression.id">
            <td>{{ expression.id }}</td>
            <td>{{ expression.expression }}</td>
            <td>{{ expression.result }}</td>
            <td>
              <n-tag v-if="expression.status.toLowerCase()==='completed'" type="success">{{expression.status}}</n-tag>
              <n-tag v-else-if="expression.status.toLowerCase()==='failed'" type="error">{{expression.status}}</n-tag>
              <n-tag v-else-if="expression.status.toLowerCase()==='pending'" type="warning">{{expression.status}}</n-tag>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="flex justify-center items-center">
      <n-pagination v-model:page="page" :page-count="pages"/>
      <n-button text @click="fetchExps">
        <n-icon>
          <RefreshIcon/>
        </n-icon>
      </n-button>
    </div>
  </div>
</template>

<style scoped>
  .n-button {
    margin-left: 1em;
  }
  #table-div {
    width: 70%;
    max-height: 70%;
    border: 0.05em solid #2A2A2E;
    margin-bottom: 1em;
  }
  table {
    display: table;
    border-collapse: collapse;
    font-size: 1.5em;
    text-align: center;
    table-layout: fixed;
    width: 100%;
  }
  thead {
    background-color: #19191D;
  }
  th, td {
    padding: 0.5em 0.5em 0.5em 0.5em;
    border-bottom: 0.05em solid #2A2A2E;
  }
</style>