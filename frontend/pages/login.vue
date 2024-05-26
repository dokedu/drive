<script lang="ts" setup>
import { ref } from "vue";

const email = ref("");
const password = ref("");

async function onSubmit() {
  console.log(email.value, password.value);
}
</script>

<template>
  <div>
    <form @submit.prevent="onSubmit">
      <input type="email" placeholder="Email" />
      <input type="password" placeholder="Password" />
      <button type="submit">Login</button>
    </form>
  </div>
</template>
