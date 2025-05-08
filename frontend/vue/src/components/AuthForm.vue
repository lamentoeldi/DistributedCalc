<script setup lang="ts">
import { ref } from 'vue'
import {
  NSpace,
  NInput,
  NButton,
  useMessage
} from 'naive-ui'
import { PersonCircleOutline as UserIcon } from "@vicons/ionicons5"
import { treaty } from '@elysiajs/eden'
import type { App } from 'bff/src'
import {useUserStore} from "@/stores/user.ts";

const username = ref('')
const password = ref('')
const isRegisterForm = ref(true)

const msg = useMessage()

const sendForm = async () => {
  const app = treaty<App>(window.location.origin)

  if (isRegisterForm.value) {
    let loading = msg.loading("Processing registration...")

    const { status } = await app.bff.api.v1.register.post({
      login: username.value,
      password: password.value
    })

    loading.destroy()

    if (status < 200 || status > 299) {
      switch (status) {
        case 400:
          msg.error("Invalid credentials")
          break
        case 409:
          msg.error("This username was already registered")
          break
        case 500:
          msg.error("Unknown error occurred")
          break
      }
      return
    }

    msg.success("Registration succeeded")
  } else {
    let loading = msg.loading("Processing authorization...")

    const { status } = await app.bff.api.v1.login.post({
      login: username.value,
      password: password.value
    })

    const userStore = useUserStore()
    await userStore.fetchUser()

    loading.destroy()

    if (status < 200 || status > 299) {
      switch (status) {
        case 400:
        case 401:
        case 404:
          msg.error("Invalid credentials")
          break
        case 500:
          msg.error("Unknown error occurred")
      }
      return
    }

    msg.success("Authorization succeeded")
  }
}

const changeForm = () => {
  username.value = ''
  password.value = ''
  isRegisterForm.value = !isRegisterForm.value
}
</script>

<template>
  <n-space vertical align="center">
    <div class="container flex justify-center items-center flex-col">
      <UserIcon class="icon"/>
      <h1 v-if="isRegisterForm">Register</h1>
      <h1 v-else>Authorize</h1>
    </div>
    <n-input
        size="large"
        round
        placeholder="Enter username"
        v-model:value="username"
    />
    <n-input
      size="large"
      round
      type="password"
      show-password-on="mousedown"
      placeholder="Enter password"
      v-model:value="password"
    />
    <n-space>
      <n-button
        size="large"
        round
        @click="sendForm"
      >
        <span v-if="isRegisterForm">Sign Up</span>
        <span v-else>Sign In</span>
      </n-button>
      <n-button
        size="large"
        round
        secondary
        @click="changeForm"
      >
        <span v-if="isRegisterForm">Sign In</span>
        <span v-else>Sign Up</span>
      </n-button>
    </n-space>
  </n-space>
</template>

<style scoped>
  h1 {
    font-size: 2em;
  }
  .icon {
    width: 4em;
    height: 4em;
    font-size: 1em;
  }
</style>