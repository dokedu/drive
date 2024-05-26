<script lang="ts" setup>
import { useRouteQuery } from "@vueuse/router";

const search = useRouteQuery("search");
</script>

<template>
  <div class="py-1.5 min-h-0 px-4 border-gray-100 border-b w-full text-sm h-14 flex items-center">
    <slot />
  </div>
</template>
