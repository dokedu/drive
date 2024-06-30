import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", () => {
  // User data
  const userToken = ref("");
  const user = ref({});

  // request login link
  async function getLoginLink(email: string) {
    try {

      const response = await $fetch<any>("http://localhost:1323/one_time_login", {
        method: "POST",
        query: {
          email: email
        }
      });

      if (!response) {
        return false;
      }

      return response;
    } catch (error: any) {
      console.log(error)

      if (error.status === 400) {
        return { error: error.data }
      }
    }
  }

  async function login(token: string) {
    const response = await $fetch<any>("http://localhost:1323/login", {
      method: "POST",
      query: {
        token: token
      }
    });

    userToken.value = response.token;
    user.value = response.user;

    console.log(response.user, response.token)

    return response;

  }

  async function register(firstName: string, lastName: string, email: string, organisation: string) {
    try {
      const { response, error } = await $fetch<any>("http://localhost:1323/sign_up", {
        method: "POST",
        query: {
          firstName: firstName,
          lastName: lastName,
          email: email,
          organisation: organisation
        }
      });

      return response;
    } catch (error: any) {
      console.log(error)

      if (error.status === 400) {
        return { error: error.data }
      }
    }
  }

  async function logout() {
    await $fetch<any>("http://localhost:1323/logout", {
      method: "POST",
    });

    userToken.value = "";
    user.value = {};
  }

  return {
    getLoginLink,
    login,
    register,
    logout,
    userToken,
    user
  }
}, {
  persist: true
})
