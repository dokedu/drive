<script lang="ts" setup>
import { onKeyDown } from "@vueuse/core";
import { useFileStore } from "#imports";

interface File {
  id: string;
  name: string;
  file_size: number;
  mime_type: string;
  created_at: string;
}

interface Props {
  files: File[];
  parentId?: string;
}

const { files, parentId } = defineProps<Props>();
const store = useFileStore();

onKeyDown("ArrowDown", () => {
  store.selectNextFile();
});

onKeyDown("ArrowUp", () => {
  store.selectPreviousFile();
});

const selectedFile = computed(() =>
  files.find((file) => file.id === store.selectedFiles[0]),
);
const previewOpen = ref(false);
function togglePreview() {
  previewOpen.value = !previewOpen.value;
}

const onSpace = (event: KeyboardEvent) => {
  togglePreview();
};

onKeyDown(" ", onSpace);
</script>

<template>
  <d-file-preview
    :open="previewOpen"
    v-if="previewOpen && selectedFile"
    @close="previewOpen = false"
    :file="selectedFile"
  ></d-file-preview>

  <d-dropzone :parent-id="parentId">
    <div class="flex flex-col">
      <d-file-list-header />
      <div class="h-full flex-1 overflow-auto">
        <d-file-list-item
          v-for="file in files"
          :key="file.id"
          :file="file"
          @dblclick="togglePreview"
        />
      </div>
    </div>
  </d-dropzone>
</template>
