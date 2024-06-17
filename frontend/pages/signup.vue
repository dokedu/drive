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
        <d-button submit type="primary"> Log in </d-button>
        <router-link
          class="mx-auto block w-fit rounded-md text-center text-xs font-medium leading-none text-muted hover:text-default focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-neutral-950"
          to="/forgot-password"
        >
          Already have an account? Log in
        </router-link>
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

async function onSubmit() {
  //
}
</script>
