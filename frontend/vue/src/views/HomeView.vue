<script>
  import { defineComponent } from "vue";
  import {
    LogoGithub
  } from "@vicons/ionicons5";
  import { darkTheme } from 'naive-ui'

  export default defineComponent({
    components: {
      LogoGithub
    },
    setup() {
      return {
        darkTheme,
        LogoGithub
      }
    }
  })
</script>

<template>
  <main class="container grid grid-cols-12 grid-rows-12 mc">
    <a
        href="https://github.com/lamentoeldi/DistributedCalc"
        target="_blank"
        rel="noopener noreferrer"
        class="col-start-12 flex items-center justify-center gi"
    >
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-icon size="25">
            <LogoGithub />
          </n-icon>
        </template>
        Visit our GitHub
      </n-tooltip>
    </a>
    <h1 class="row-start-1 row-end-1 col-start-6 text-center text-xl">
      DistributedCalc
    </h1>
    <div class="row-start-2 col-start-3 col-end-11">
      <n-card title="What is DistributedCalc?" :theme="darkTheme">
        DistributedCalc is top-class online calculation platform
      </n-card>
      <n-card title="Why DistributedCalc?" :theme="darkTheme">
        DistributedCalc provides smooth and delightful experience via extensive and performant API
      </n-card>
      <n-card title="How do I start using DistributedCalc?" :theme="darkTheme">
        Proceed to <RouterLink to="auth" class="ref">registration</RouterLink> and enjoy your stay
      </n-card>
    </div>
  </main>
</template>

<style scoped>
  .ref {
    color: rgb(99, 226, 183);
    border-bottom: rgb(99, 226, 183) 0.07em dashed;
  }
  .n-card {
    width: 100%;
    border-radius: 1em;
    margin: 1em;
  }
  h1 {
    font-size: 1.7em;
  }
  .mc {
    margin: 0.5em 0.5em 0.5em 0.5em;
  }

  .gi {
    background-color: #19191D;
    border: 0.13em solid #2A2A2E;
    border-radius: 1em;
    width: 4em;
    height: 2.5em;
  }
</style>