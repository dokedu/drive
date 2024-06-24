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
            label="First Name"
            v-model="firstName"
            type="text"
            name="firstname"
            id="firstname"
            required
            placeholder="Your first name"
          />
          <d-input
            size="md"
            label="Last Name"
            v-model="lastName"
            type="text"
            name="lastname"
            id="lastname"
            required
            placeholder="Your last name"
          />
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
          <d-input
            size="md"
            label="Organisation name"
            v-model="organisation"
            type="text"
            name="organisation"
            id="organisation"
            required
            placeholder="Your organisation name"
          />
        </div>
        <d-button submit type="primary"> Sign up </d-button>
        <router-link
          class="mx-auto block w-fit rounded-md text-center text-xs font-medium leading-none text-muted hover:text-default focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-neutral-950"
          to="/login"
        >
          Already have an account? Log in
        </router-link>
        <div v-if="loginLinkSent" class="text-center text-stone-500 text-sm font-normal leading-tight">
          Signed up successfully. Login link sent to your email. Please check your inbox.
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

definePageMeta({
  layout: "auth",
});

const firstName = ref("");
const lastName = ref("");
const email = ref("");
const organisation = ref("");

const error = ref<Error | null>(null);
const loginLinkSent = ref(false);
const authStore = useAuthStore();

async function onSubmit() {
  if (!firstName.value || !lastName.value || !email.value || !organisation.value) {
    error.value = new Error("Please fill in all fields");
    return;
  }
  const response = await authStore.register(firstName.value, lastName.value, email.value, organisation.value);
  console.log(response);
  if (response) {
    loginLinkSent.value = true;
  } else {
    error.value = new Error("There was an error. Please try again.");
  }
}
</script>
