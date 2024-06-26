<template>
  <d-auth-container
    title="Welcome to dokedu drive"
    subtitle="Please enter your login credentials"
  >
    <template #banner>
      <d-banner v-if="error" type="error" :title="error.message"></d-banner>
    </template>
    <template #form>
      <form @submit.prevent="onSubmit" class="flex flex-col gap-5">
        <div class="flex flex-col gap-3">
          <d-input
            size="md"
            label="Email"
            v-model="email"
            type="email"
            name="email"
            id="email"
            required
            autocomplete="email"
            placeholder="Your email"
          />
        </div>
        <d-button submit type="primary"> Log in </d-button>
        <router-link
          class="mx-auto block w-fit rounded-md text-center text-xs font-medium leading-none text-muted hover:text-default focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-neutral-950"
          to="/signup"
        >
          Create new organisation
        </router-link>
        <div v-if="loginLinkSent" class="text-center text-stone-500 text-sm font-normal leading-tight">
          Login link sent to your email. Please check your inbox.
        </div>
        <div v-if="error" class="text-center text-stone-500 text-sm font-normal leading-tight">
          {{ error.message }}
        </div>
      </form>
    </template>
  </d-auth-container>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import DInput from "@/components/d-input/d-input.vue";
import DButton from "@/components/d-button/d-button.vue";
import DAuthContainer from "@/components/_auth/d-auth-container.vue";
import { useAuthStore } from "@/stores/auth";

definePageMeta({
  layout: "auth",
});

const email = ref("");
const error = ref<Error | null>(null);
const loginLinkSent = ref(false);

const authStore = useAuthStore();
const route = useRoute();

if (route.hash) {
  // Get only the part after #token=
  const token = route.hash.split("#token=")[1];

  if (token) {
    const response = await authStore.login(token);
    if (response) {
      navigateTo("/");
    } else {
      navigateTo("/login");
    }
  }
}

async function onSubmit() {
  const response = await authStore.getLoginLink(email.value);
  console.log(response);
  if (response) {
    console.log(response)
    if (response.error) {
      error.value = new Error(response.error);
      return;
    }
    loginLinkSent.value = true;
  } else {
    error.value = new Error("There was an error. Please try again.");
  }
}
</script>
