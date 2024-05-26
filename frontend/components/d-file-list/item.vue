<script lang="ts" setup>
import { FileText, FileType, Folder, FileImage, File } from "lucide-vue-next";
import { NuxtLink } from "#components";
import { useFileStore } from "@/stores/file";
import {
  ContextMenuCheckboxItem,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuItemIndicator,
  ContextMenuLabel,
  ContextMenuPortal,
  ContextMenuRadioGroup,
  ContextMenuRadioItem,
  ContextMenuRoot,
  ContextMenuSeparator,
  ContextMenuSub,
  ContextMenuSubContent,
  ContextMenuSubTrigger,
  ContextMenuTrigger,
} from "radix-vue";
const fileStore = useFileStore();

interface File {
  id: string;
  name: string;
  file_size: number;
  mime_type: string;
  created_at: string;
}

interface Props {
  file: File;
}

const { file } = defineProps<Props>();

function formatTime(time: Date) {
  return Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
  }).format(time);
}

function fileIcon(file: File) {
  switch (file.mime_type) {
    case "text/plain":
      return FileText;
    case "image/jpeg":
      return FileImage;
    // pdf
    case "application/pdf":
      return FileType;
    case "directory":
      return Folder;
    default:
      return File;
  }
}

function componentName() {
  return file.mime_type === "directory" ? "router-link" : "div";
}

function prettyBytes(bytes: number) {
  const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  if (Math.abs(bytes) < 1) {
    return bytes + "B";
  }

  const u = Math.min(Math.floor(Math.log10(bytes) / 3), units.length - 1);
  const n = Number((bytes / Math.pow(1000, u)).toFixed(2));
  return `${n} ${units[u]}`;
}

const selected = computed(() => fileStore.fileSelected(file.id));

async function onClick() {
  if (!fileStore.fileSelected(file.id)) {
    fileStore.fileToggleSelected(file.id);
    return;
  }
  if (file.mime_type === "directory") {
    fileStore.clearSelectedFiles();
    await navigateTo(`/folders/${file.id}`);
  }
}

function onContextMenu(event: MouseEvent) {
  event.preventDefault();
  console.log("Right click");
}

async function handleClick() {
  // eslint-disable-next-line no-alert
  await fileStore.createFolder();
}

async function handleDelete() {
  // eslint-disable-next-line no-alert
  await fileStore.deleteFile(file.id);
}

async function handleDownload() {
  // eslint-disable-next-line no-alert
  await fileStore.downloadFile(file);
}

async function rename() {
  await fileStore.updateFileName(file);
}

const editingName = ref(true)
</script>

<template>
  <d-context-menu>
    <template #content>
      <div
        class="grid items-center py-2.5 cursor-default text-sm px-4 text-gray-700 gap-8"
        :class="selected ? 'bg-blue-100 hover:bg-blue-100' : 'hover:bg-gray-100'"
        :style="{ gridTemplateColumns: '6fr 1fr 3fr' }"
        @click="onClick"
        :data-id="file.id"
      >
        <div class="flex items-center gap-2 w-full">
          <component :is="fileIcon(file)" :size="18" :class="file.mime_type === 'directory' ? 'fill-current' : ''" />
          <div v-if="editingName" class="w-full">
            <input
              type="text"
              v-model="file.name"
              class="w-full bg-transparent focus:outline-none"
              @keydown.enter="rename"
            />
          </div>
          <div v-else>
            {{ file.name }}
          </div>
        </div>
        <div class="text-right">{{ prettyBytes(file.file_size) }}</div>
        <div>{{ formatTime(file.createdAt) }}</div>
      </div>
    </template>
    <template #menu>
      <ContextMenuItem
        value="New Folder"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none"
        @click="handleClick"
      >
        New Folder
      </ContextMenuItem>
      <ContextMenuSeparator class="border-t border-neutral-100 my-1" />
      <ContextMenuItem
        value="Rename"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none"
      >
        Rename
      </ContextMenuItem>
      <ContextMenuItem
        value="Delete"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none"
        @click="handleDelete"
      >
        Delete
      </ContextMenuItem>
      <ContextMenuSeparator class="border-t border-neutral-100 my-1" />
      <ContextMenuItem
        value="Download"
        class="text-black text-sm px-2 py-1 rounded-md data-[highlighted]:bg-neutral-100 outline-none"
        @click="handleDownload"
      >
        Download
      </ContextMenuItem>
    </template>
  </d-context-menu>
</template>
