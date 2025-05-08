<script setup lang="ts">
import AuthForm from "@/components/AuthForm.vue";
import {useUserStore} from "@/stores/user.ts";

import {
  NCard
} from 'naive-ui'

const userStore = useUserStore()
</script>

<template>
  <div class="container flex justify-center items-center">
    <n-card
      v-if="userStore.isAuthorized"
      :title="() => `Welcome back, ${userStore.username}!`"
      class="n-card"
    >
      You are already authorized, so please <RouterLink to="/" class="ref">proceed forward</RouterLink>.
    </n-card>
    <AuthForm v-else/>
  </div>
</template>

<style scoped>
  .n-card {
    text-align: center;
    min-width: 40%;
    width: fit-content;
  }
  .container {
    height: 100%;
  }
  .ref {
    color: rgb(99, 226, 183);
    border-bottom: rgb(99, 226, 183) 0.07em dashed;
  }
</style>