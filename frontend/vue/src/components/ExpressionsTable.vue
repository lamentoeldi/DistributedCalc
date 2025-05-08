<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  NTag,
  NPagination,
  NButton,
  NIcon,
  useMessage
} from 'naive-ui'
import {
  RefreshOutline as RefreshIcon,
  AddOutline as PlusIcon
} from '@vicons/ionicons5'
import { treaty } from '@elysiajs/eden'
import type { App } from 'bff/src'

const msg = useMessage()
const defaultFetchLimit = 30
const defaultExpsShown = 5

const chunkArray = <T>(array: T[], size: number): T[][] =>
    array.reduce((acc, val, i) => {
      if (i % size === 0) acc.push([])
      acc[acc.length - 1].push(val)
      return acc
    }, [] as T[][])

const expressions = ref<{ id: string, result: number, status: string }[]>([])
const expsToRender = ref<{ id: string, result: number, status: string }[][]>([])
const page = ref(1)
const pages = ref(0)

const fetchExpressions = async ({
                                  cursor,
                                  limit
                                }: { cursor: string; limit: number }): Promise<[data: any | null, status: number]> => {
  const app = treaty<App>(window.location.origin)

  const { data, status } = await app.bff.api.v1.expressions.get({
    query: { cursor, limit }
  })
  if (!data || (status < 200 || status > 299)) {
    return [null, status]
  }

  return [data, status]
}

const refreshExpressions = async () => {
  const loading = msg.loading("Refreshing expressions...")
  const currentCount = expressions.value.length || defaultFetchLimit

  const [data, status] = await fetchExpressions({ cursor: "", limit: currentCount })

  if ((status > 299 || status < 200) || !data) {
    loading.destroy()

    switch(status) {
      case 401:
        msg.error("Unauthorized")
        break
      case 404:
        msg.error("No expressions were found")
        break
      default:
        msg.error("Unknown error occurred")
    }
    return
  }

  expressions.value = data.expressions
  expsToRender.value = chunkArray(expressions.value, defaultExpsShown)
  pages.value = expsToRender.value.length
  page.value = 1

  msg.success("Expressions refreshed.")
  loading.destroy()
}

const getExpressions = async () => {
  const loading = msg.loading("Fetching expressions...")

  const lastId = expressions.value.length > 0
      ? expressions.value[expressions.value.length - 1].id
      : ""

  const [data, status] = await fetchExpressions({ cursor: lastId, limit: defaultFetchLimit })

  if ((status > 299 || status < 200) || !data) {
    loading.destroy()

    switch(status) {
      case 401:
        msg.error("Unauthorized")
        break
      case 404:
        msg.error("No expressions were found")
        break
      default:
        msg.error("Unknown error occurred")
    }
    return
  }

  expressions.value = expressions.value.concat(data.expressions)
  expsToRender.value = chunkArray(expressions.value, defaultExpsShown)
  pages.value = expsToRender.value.length

  msg.success("Expressions fetched.")
  loading.destroy()
}

onMounted(async () => {
  const lastId = expressions.value.length > 0
      ? expressions.value[expressions.value.length - 1].id
      : ""

  const [data, status] = await fetchExpressions({ cursor: lastId, limit: defaultFetchLimit })

  if ((status > 299 || status < 200) || !data) {
    return
  }

  expressions.value = expressions.value.concat(data.expressions)
  expsToRender.value = chunkArray(expressions.value, defaultExpsShown)
  pages.value = expsToRender.value.length
})
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
      <n-button @click="refreshExpressions" class="rb">
        <n-icon>
          <RefreshIcon/>
        </n-icon>
      </n-button>
      <n-pagination v-model:page="page" :page-count="pages"/>
      <n-button @click="getExpressions" class="gb">
        <n-icon>
          <PlusIcon/>
        </n-icon>
      </n-button>
    </div>
  </div>
</template>

<style scoped>
  .rb, .gb {
    margin: auto 0.7em auto 0.7em;
  }
  .rb {
    font-size: 1.3em;
  }
  .gb {
    font-size: 1.6em;
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