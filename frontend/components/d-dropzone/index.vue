<script setup lang="ts">
import { useDropZone } from "@vueuse/core";
import { useFileStore } from "@/stores/file";

interface Props {
  parentId?: string;
}

const fileStore = useFileStore();

const { parentId } = defineProps<Props>();

const dropZoneRef = ref<HTMLDivElement>();

async function onDrop(files: File[] | null) {
  if (!files) return;

  await fileStore.uploadFiles(files, parentId);
}

const { isOverDropZone } = useDropZone(dropZoneRef, {
  onDrop,
  // specify the types of data to be received.
  //dataTypes: ["image/jpeg"],
});
</script>

<template>
  <div
    ref="dropZoneRef"
    class="h-full flex-1 overflow-hidden"
    style="height: calc(100vh - 57px)"
    :class="{ 'bg-red-200': isOverDropZone }"
  >
    <slot />
  </div>
</template>
