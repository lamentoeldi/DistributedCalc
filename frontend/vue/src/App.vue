<script lang="ts">
import {RouterLink, RouterView} from 'vue-router'
import type MenuOption from 'naive-ui'
import type { Component } from "vue";
import { defineComponent, h, ref } from 'vue'
import { NIcon, darkTheme } from 'naive-ui'
import {
  HomeOutline as HomeIcon,
  PersonOutline as AuthIcon,
  AddCircleOutline as CalculateIcon,
  SearchOutline as FindIcon,
  ListOutline as ExpressionsIcon,
} from '@vicons/ionicons5'

const renderIcon = (icon: Component) => {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
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

export default defineComponent({
  setup() {
    return {
      darkTheme,
      collapsed: ref(true),
      menuOptions
    }
  }
})
</script>

<template>
  <n-layout has-sider
        :theme="darkTheme"
  >
    <n-layout-sider
        bordered
        collapse-mode="width"
        :theme="darkTheme"
        :collapsed-width="64"
        :width="160"
        :collapsed="collapsed"
        show-trigger
        @collapse="collapsed = true"
        @expand="collapsed = false"
    >
      <n-menu
          :theme="darkTheme"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="menuOptions"
      />
    </n-layout-sider>
    <RouterView/>
  </n-layout>
</template>

<style scoped>
  .n-layout {
    height: 100%;
  }
  .ec {
    height: 100%;
  }
  .bd {
    border: 0.1em solid;
    border-radius: 3em;
    width: 70%;
    height: 70%;
  }
</style>