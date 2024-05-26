<script lang="ts" setup>
import { useFileStore } from "@/stores/file";

const route = useRoute();

const fileStore = useFileStore();

await useAsyncData("files", async () => {
  fileStore.resetOptions();
  fileStore.sharedDrive = "yes";

  await fileStore.fetchFiles();

  return true;
});
</script>

<template>
  <div class="w-full">
    <d-header>Shared drives</d-header>
    <d-file-list :files="fileStore.files" />
  </div>
</template>
