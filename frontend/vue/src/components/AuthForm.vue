<script setup lang="ts">
import { ref } from 'vue'
import {
  NSpace,
  NInput,
  NButton
} from 'naive-ui'
import { PersonCircleOutline as UserIcon } from "@vicons/ionicons5"

const username = ref('')
const password = ref('')
const isRegisterForm = ref(true)

const sendForm = async () => {}

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