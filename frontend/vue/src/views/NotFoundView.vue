<script setup lang="ts">
import {
  NCard,
} from 'naive-ui'
</script>

<template>
  <div class="container flex justify-center items-center">
    <n-card
        title="Ooops... Page not found"
    >
      Looks like the page you are looking for does not exist. Redirect to the <RouterLink to="/" class="ref">home page</RouterLink>?
    </n-card>
  </div>
</template>

<style scoped>
  .container {
    height: 100%;
    width: 100%;
  }
  .n-card {
    width: 40%;
    text-align: center;
  }
  .ref {
    color: rgb(99, 226, 183);
    border-bottom: rgb(99, 226, 183) 0.07em dashed;
  }
</style>