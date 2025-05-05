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

  const { data, error } = await app.bff.api.v1.calculate.post({
    expression: exp.value
  })

  loading.destroy()

  if (!data || error) {
    msg.error("Failed to send expression")
    return
  }

  msg.success("Expression sent successfully")
  expID.value = data.id
}
</script>

<template>
  <div class="container flex items-center justify-center flex-col">
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
    <span v-if="expID" class="text">Expression ID: {{expID}}</span>
  </div>
</template>

<style scoped>
  .text {
    text-align: center;
    border: 0.1em solid transparent;
    background-color: #141418;
    border-radius: 1em;
    padding: 0.3em 1em 0.3em 1em;
  }
  .n-button {
    font-size: 3em;
  }
  .n-input {
    height: 1.7em;
    font-size: 1.5em;
  }
  .container {
    height: 100%;
  }
</style>