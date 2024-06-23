<script lang="ts" setup>
import { useFileStore } from "@/stores/file";

const fileStore = useFileStore();

await useAsyncData("files", async () => {
  fileStore.resetOptions();
  await fileStore.fetchFiles();

  return true;
});

function newFolder() {
  fileStore.createFolder();
}
</script>

<template>
  <div class="w-full">
    <d-header class="justify-between">
      <div>My files</div>
      <div>
        <DButton @click="newFolder()">New Folder</DButton>
      </div>
    </d-header>
    <d-file-list :files="fileStore.files" />
  </div>
</template>
