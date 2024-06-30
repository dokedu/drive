<script lang="ts" setup>
import { onKeyDown, onClickOutside } from "@vueuse/core";
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

onKeyDown("Escape", () => {
  store.clearSelectedFiles();
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

async function handleDelete() {
  // eslint-disable-next-line no-alert
  if (!selectedFile.value) return
  await store.deleteFile(selectedFile.value.id);
}

async function handleDownload() {
  // eslint-disable-next-line no-alert
  if (!selectedFile.value) return
  await store.downloadFile(selectedFile.value);
}

const fileList = ref<HTMLElement | null>(null);
onClickOutside(fileList, () => {
  store.clearSelectedFiles();
});
</script>

<template>
  <d-file-preview :open="previewOpen" v-if="previewOpen && selectedFile" @close="previewOpen = false"
    :file="selectedFile"></d-file-preview>

  <d-context-menu>
    <template #content>
      <d-dropzone :parent-id="parentId">
        <div class="flex flex-col">
          <d-file-list-header />
          <div class="h-full flex-1 overflow-auto" ref="fileList">
            <d-file-list-item v-for="file in files" :key="file.id" :file="file" @dblclick="togglePreview" @contextmenu="onContextMenu" />
          </div>
        </div>
      </d-dropzone>
    </template>
    <template #menu>
      <ContextMenuItem value="New Folder"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none">
        New Folder
      </ContextMenuItem>
      <div v-if="selectedFile">

      <ContextMenuSeparator class="border-t border-neutral-100 my-1" />
      <ContextMenuItem value="Rename"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none">
        Rename
      </ContextMenuItem>
      <ContextMenuItem value="Delete"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none"
        @click="handleDelete">
        Delete
      </ContextMenuItem>
      <ContextMenuSeparator class="border-t border-neutral-100 my-1" />
      <ContextMenuItem value="Download"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none"
        @click="handleDownload">
        Download
      </ContextMenuItem>

      </div>
    </template>
  </d-context-menu>
</template>
