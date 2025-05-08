<script setup lang="ts">
import {
  NInput,
  NIcon,
  NButton,
  NSpace,
  useMessage
} from 'naive-ui'

import {
  ArrowForwardCircleOutline as SendIcon
} from "@vicons/ionicons5"

import { treaty } from '@elysiajs/eden'
import type { App } from 'bff/src'

import { ref } from 'vue'

const expID = ref<string | undefined>(undefined)

const exp = ref('')

const msg = useMessage()

const sendExp = async () => {
  const loading = msg.loading("Sending expression")

  const app = treaty<App>(window.location.origin)

  const { data, status } = await app.bff.api.v1.calculate.post({
    expression: exp.value
  })

  loading.destroy()

  if (!data || (status < 200 || status > 299)) {
    switch (status) {
      case 400:
      case 422:
        msg.error("Invalid expression")
        break
      case 401:
        msg.error("Unauthorized")
        break
      case 500:
        msg.error("Unknown error occurred")
        break
    }
    return
  }

  msg.success("Expression sent successfully")
  expID.value = data.id
}

const copyExpID = async () => {
  if (!expID) {
    return
  }

  await navigator.clipboard.writeText(`${expID.value}`)
}
</script>

<template>
  <div class="container flex items-center justify-center flex-col">
    <div class="d-wrap" v-if="expID">
      <table class="t">
        <tbody>
        <tr>
          <th>ID</th>
          <td>
            <span
                @click="copyExpID"
                class="btn"
            >{{expID}}</span>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <n-space align="center" justify="center">
      <n-input
        size="large"
        round
        placeholder="Enter your expression"
        v-model:value="exp"
      />
      <n-button
        text
        @click="sendExp"
      >
        <n-icon>
          <SendIcon/>
        </n-icon>
      </n-button>
    </n-space>
    <div>
    </div>
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
  .n-button {
    font-size: 2.5em;
  }
  .n-input {
    height: 1.7em;
    font-size: 1.5em;
  }
  .container {
    height: 100%;
  }
  .d-wrap {
    height: fit-content;
    width: fit-content;
    padding: 0.5em 0.5em 0.5em 0.5em;
    margin: 0.5em 0.5em 0.5em 0.5em;
    border: 0.1em solid #2A2A2E;
    border-radius: 1em;
    background-color: #19191D;
  }
  .btn:hover {
    color: rgb(255, 255, 255, 0.7);
    border-bottom: 0.07em dashed rgb(255, 255, 255, 0.7);
  }
  .btn:active {
    color: white;
  }
</style>