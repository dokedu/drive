<script lang="ts" setup>
import { HardDrive, Server, Trash2, LogOut } from "lucide-vue-next";
</script>

<template>
  <div class="flex-1 w-full border-r border-gray-100 max-w-64 min-w-64 bg-neutral-50 flex justify-between flex-col">
    <div class="">
      <div class="font-semibold flex items-center h-14 px-4">Drive</div>
      <div class="p-2 flex flex-col gap-1">
        <d-sidebar-link to="/my-files" :active="true">
          <HardDrive :size="18" />
          My files
        </d-sidebar-link>
        <d-sidebar-link to="/shared-drives">
          <Server :size="18" />
          Shared drives
        </d-sidebar-link>
        <d-sidebar-link to="/trash">
          <Trash2 :size="18" />
          Trash
        </d-sidebar-link>
      </div>
    </div>
    <div class="p-2">
      <d-sidebar-link to="/login">
        <LogOut :size="18" />
        Log out
      </d-sidebar-link>
    </div>
  </div>
</template>
