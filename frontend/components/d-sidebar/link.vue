<script lang="ts" setup>
const route = useRoute();

interface Props {
  to: string;
}

const { to } = defineProps<Props>();

const active = computed(() => {
  if (to === "/") {
    return route.path === to;
  }
  return route.path.startsWith(to);
});
</script>

<template>
  <nuxt-link
    :to="to"
    class="px-3 py-2 flex gap-2.5 items-center text-sm rounded-md cursor-default"
    :class="active ? 'bg-neutral-200 text-neutral-950 hover:bg-neutral-200' : 'hover:bg-neutral-100 text-neutral-600 '"
  >
    <slot />
  </nuxt-link>
</template>
