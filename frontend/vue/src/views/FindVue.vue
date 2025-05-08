<script setup lang="ts">
import {
  SearchOutline as SearchIcon,
} from '@vicons/ionicons5'
import {
  NButton,
  NIcon,
  NInput,
  useMessage
} from "naive-ui";
import {
  ref
} from 'vue'
import Expression from "@/components/Expression.vue";

import { treaty } from '@elysiajs/eden'
import type { App } from 'bff/src'

const expID = ref<string>('')

interface Expression {
  id: string,
  result: number,
  status: string
}

const exp = ref<Expression | null>(null)

const msg = useMessage()

const findExp = async () => {
  const loading = msg.loading("Searching expression")

  const app = treaty<App>(window.location.origin)
  const { data, status } = await app.bff.api.v1.expressions({
    id: expID.value
  }).get()

  loading.destroy()
  if (!data || (status < 200 || status > 299)) {
    exp.value = null

    switch (status) {
      case 400:
      case 422:
        msg.error("Invalid expression id")
        break
      case 401:
        msg.error("Unauthorized")
        break
      case 404:
        msg.error("Expression not found")
        break
      case 500:
        msg.error("Unknown error occurred")
        break
    }

    return
  }

  msg.success("Expression found")
  exp.value = data.expression
}
</script>

<template>
  <div class="container flex justify-center items-center">
    <div class="flex items-center flex-col">
      <Expression v-if="exp !== null" :id="exp.id" :result="exp.result" :status="exp.status"/>
      <div class="flex items-center justify-center">
        <n-input
            size="large"
            round
            placeholder="Enter expression id"
            v-model:value="expID"
        />
        <n-button
            class="b-1"
            text
            @click="findExp"
        >
          <n-icon>
            <SearchIcon/>
          </n-icon>
        </n-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .b-1 {
    font-size: 1.5em;
    margin: 0.3em 0.5em 0.3em 0.5em;
  }
  .container {
    height: 100%;
  }
  .n-button {
    font-size: 2em;
  }
  .n-input {
    height: 1.7em;
    font-size: 1.5em;
  }
</style>