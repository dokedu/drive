<script lang="ts" setup>
import { useFileStore } from "@/stores/file";

const fileStore = useFileStore();

await useAsyncData("files", async () => {
  fileStore.resetOptions();

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
