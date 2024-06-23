import { get } from "@vueuse/core";
import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", () => {
  // User data
  const userToken = ref("");
  const user = ref<any>(null);

  // request login link
  async function getLoginLink(email: string) {
    const response = await $fetch<any>("http://localhost:1323/one-time-login", {
      method: "POST",
      query: {
        email: email
      }
    });

    if (!response) {
      return false;
    }

    return response;
  }

  async function login(token: string) {
    const response = await $fetch<any>("http://localhost:1323/login", {
      method: "POST",
      query: {
        token: token
      }
    });

    // Convert json to object
    const res = JSON.parse(response)
    console.log(res.token)

    if (!res) {
      return false;
    }

    userToken.value = res.token;
    user.value = res.user;

    console.log(userToken.value, user.value)

    return response;
  }


  return {
    getLoginLink,
    login,
    userToken,
    user
  }
})
