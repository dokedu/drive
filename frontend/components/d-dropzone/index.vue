<script setup lang="ts">
import { useDropZone } from "@vueuse/core";
import { useFileStore } from "@/stores/file";
import { Upload } from "lucide-vue-next";

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
    class="h-full flex-1 overflow-hidden relative"
    style="height: calc(100vh - 57px)"
  >
    <transition>
      <div
        class="absolute top-0 left-0 w-full h-full grid place-items-center p-1 pointer-events-none transition"
        :class="!isOverDropZone ? 'scale-95' : 'scale-100'"
      >
        <div
          class="w-full h-full rounded-lg flex items-end justify-center transition"
          :class="
            !isOverDropZone
              ? 'bg-transparent'
              : 'bg-blue-700/20 broder-blue-500'
          "
        >
          <div
            class="text-white text-lg p-1 rounded-xl flex items-center gap-3 mb-10 shadow-xl bg-blue-600 transition-all"
            :class="
              !isOverDropZone ? 'translate-y-12 opacity-0' : 'translate-y-0'
            "
          >
            <div class="size-8 grid place-items-center bg-white rounded-lg">
              <Upload class="size-4 stroke-blue-800 stroke" />
            </div>
            <div class="mr-2">Drop files to upload them</div>
          </div>
        </div>
      </div>
    </transition>
    <slot />
  </div>
</template>
