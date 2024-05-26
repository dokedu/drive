<script lang="ts" setup>
import { nanoid } from "nanoid";

const files = ref([
  {
    id: 1,
    name: "Pictures",
    size: 0,
    type: "directory",
    createdAt: new Date(),
  },
  {
    id: 2,
    name: "file1.txt",
    size: 1024,
    type: "text/plain",
    createdAt: new Date(),
  },
  {
    id: 3,
    name: "file2.txt",
    size: 2048,
    type: "text/plain",
    createdAt: new Date(),
  },
]);
</script>

<template>
  <div class="w-full">
    <d-header>My files</d-header>
    <div>preview of file</div>
  </div>
</template>
