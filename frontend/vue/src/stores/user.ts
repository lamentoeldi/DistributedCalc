import { ref } from 'vue'
import { defineStore } from 'pinia'
import { treaty } from '@elysiajs/eden'
import type { App } from '../../../bff/src'

export const useUserStore = defineStore('user', () => {
  const isAuthorized = ref(false)
  const userID = ref('')
  const username = ref('')

  const fetchUser = async() => {
    const app = treaty<App>(window.location.origin)

    const { data, error } = await app.api.v1.authorize.get()
    if (!data || error) {
      return
    }

    isAuthorized.value = true
    userID.value = data.user_id
    username.value = data.username
  }

  return {
    fetchUser,
    isAuthorized,
    userID,
    username
  }
})