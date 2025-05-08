<script setup lang="ts">
import {RouterLink, RouterView} from 'vue-router'
import type { Component } from "vue";
import { h, ref, onMounted } from 'vue'
import {
  NIcon,
  NMenu,
  NConfigProvider,
  NLayout,
  NLayoutSider,
  NMessageProvider,
  darkTheme
} from 'naive-ui'
import {
  HomeOutline as HomeIcon,
  PersonOutline as AuthIcon,
  AddCircleOutline as CalculateIcon,
  SearchOutline as FindIcon,
  ListOutline as ExpressionsIcon
} from '@vicons/ionicons5'
import {useUserStore} from "@/stores/user.ts";

const renderIcon = (icon: Component) => {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = [
  {
    label: () =>
      h(
          RouterLink,
          {
            to: {
              name: 'home'
            }
          },
          {
            default: () => 'Home'
          }
      ),
    key: 'home',
    icon: renderIcon(HomeIcon)
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'auth'
              }
            },
            {
              default: () => 'Auth'
            }
        ),
    key: 'auth',
    icon: renderIcon(AuthIcon)
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'expressions'
              }
            },
            {
              default: () => 'Expressions'
            }
        ),
    key: 'expressions',
    icon: renderIcon(ExpressionsIcon)
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'calculate'
              }
            },
            {
              default: () => 'Calculate'
            }
        ),
    key: 'calculate',
    icon: renderIcon(CalculateIcon)
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'find'
              }
            },
            {
              default: () => 'Find'
            }
        ),
    key: 'find',
    icon: renderIcon(FindIcon)
  }
]

const collapsed = ref(true)

const userStore = useUserStore()
onMounted(userStore.fetchUser)
</script>

<template>
  <n-config-provider
    :theme="darkTheme"
  >
    <n-message-provider>
      <n-layout has-sider
      >
        <n-layout-sider
            bordered
            collapse-mode="width"
            :collapsed-width="64"
            :width="160"
            :collapsed="collapsed"
            show-trigger
            @collapse="collapsed = true"
            @expand="collapsed = false"
        >
          <n-menu
              :collapsed="collapsed"
              :collapsed-width="64"
              :collapsed-icon-size="22"
              :options="menuOptions"
          />
        </n-layout-sider>
        <RouterView/>
      </n-layout>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
  .n-layout, .ec, .n-config-provider, .n-message-provider {
    height: 100%;
  }
  .bd {
    border: 0.1em solid;
    border-radius: 3em;
    width: 70%;
    height: 70%;
  }
</style>