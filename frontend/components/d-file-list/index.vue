<script lang="ts" setup>
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
</script>

<template>
  <d-dropzone :parent-id="parentId">
    <div class="flex flex-col">
      <d-file-list-header />
      <div class="h-full flex-1 overflow-auto">
        <d-file-list-item v-for="file in files" :key="file.id" :file="file" />
      </div>
    </div>
  </d-dropzone>
</template>
