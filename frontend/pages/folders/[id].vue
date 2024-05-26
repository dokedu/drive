<script lang="ts" setup>
import { useFileStore } from "@/stores/file";

const route = useRoute();

const fileStore = useFileStore();

await useAsyncData(`file-${route.params.id}`, async () => {
  fileStore.resetOptions();
  fileStore.parentId = route.params.id;

  await fileStore.fetchFiles();

  return true;
});
</script>

<template>
  <div class="w-full">
    <d-header>My files</d-header>
    <d-file-list :files="fileStore.files" />
  </div>
</template>
