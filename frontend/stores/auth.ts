import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", () => {
  // User data
  const userToken = ref("");
  const user = ref({});

  // request login link
  async function getLoginLink(email: string) {
    try {

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

    // Convert json to object
    const res = JSON.parse(response)
    if (!res) {
      return false;
    }

    userToken.value = res.token;
    user.value = res.user;

    console.log(res.user, res.token)

    return response;

  }

  async function register(firstName: string, lastName: string, email: string, organisation: string) {
    try {
      const { response, error } = await $fetch<any>("http://localhost:1323/sign-up", {
        method: "POST",
        query: {
          firstName: firstName,
          lastName: lastName,
          email: email,
          organisation: organisation
        }
      });
      console.log(response, error)

      // Convert json to object
      const res = JSON.parse(response)
      return res;
    } catch (error: any) {
      console.log(error)

      if (error.status === 400) {
        return { error: error.data }
      }
    }
  }


  return {
    getLoginLink,
    login,
    register,
    userToken,
    user
  }
}, {
  persist: true
})
