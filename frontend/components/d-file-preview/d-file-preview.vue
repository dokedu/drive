<script lang="ts" setup>
import { Download, Folder, X, FileIcon } from "lucide-vue-next";
import { type File } from "@/components/d-file-list/item.vue";
import { useFileStore } from "#imports";

// @ts-ignore
import * as pdfjs from "pdfjs-dist/build/pdf";
import pdfjsWorker from "pdfjs-dist/build/pdf.worker?worker";

interface Props {
  open: boolean;
  file: File;
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
});
const emit = defineEmits(["close"]);
const store = useFileStore();
const previewUrl = ref(await store.getPreviewUrl(props.file.id));

const fileType = computed(() => {
  switch (fileRef.value.mime_type) {
    case "image/jpeg":
    case "image/png":
    case "image/webp":
    case "image/gif":
      return "image";
    case "video/mp4":
    case "video/ogg":
    case "video/webm":
      return "video";
    case "audio/mpeg":
    case "audio/ogg":
    case "audio/wav":
      return "audio";
    case "application/pdf":
      return "pdf";
    case "directory":
      return "directory";
    default:
      return "unknown";
  }
});

const fileRef = toRef(props, "file");
const mimeType = computed(() => fileRef.value.mime_type);

onMounted(async () => {
  if (fileType.value === "pdf") {
    renderPDF(previewUrl.value);
  }
});

watch(fileRef, async () => {
  previewUrl.value = await store.getPreviewUrl(fileRef.value.id);
  if (fileType.value === "pdf") {
    renderPDF(previewUrl.value);
  }
});

const onClose = () => {
  emit("close");
};

const canvas = ref<HTMLCanvasElement | null>(null);

function init(): void {
  try {
    if (typeof window === "undefined" || !("Worker" in window)) {
      throw new Error("Web Workers not supported in this environment.");
    }

    // @ts-ignore
    window.pdfjsWorker = pdfjsWorker;
    pdfjs.GlobalWorkerOptions.workerSrc = `/pdfjs-4.2.67-dist/build/pdf.worker.mjs`;
  } catch (error) {
    throw new Error("PDF.js failed to load. ");
  }
}

async function renderPDF(url: string) {
  try {
    init();
  } catch (error) {
    console.error(error);
  }

  const loadingTask = pdfjs.getDocument(url);

  loadingTask.onPassword = (callback: Function, reason: number) => {
    if (reason == 1) {
      const enteredPassword = prompt("Enter password");
      if (enteredPassword !== null) {
        callback(enteredPassword);
      } else {
        emit("close");
      }
    } else {
      const enteredPassword = prompt("Password incorrect, please try again");
      if (enteredPassword !== null) {
        callback(enteredPassword);
      } else {
        emit("close");
      }
    }
  };

  // @ts-expect-error
  await loadingTask.promise.then(async (pdf) => {
    const page = await pdf.getPage(1);

    const viewport = page.getViewport({ scale: 1.5 });

    const context = canvas.value?.getContext("2d");

    if (context) {
      canvas.value!.height = viewport.height;
      canvas.value!.width = viewport.width;

      const renderContext = {
        canvasContext: context,
        viewport: viewport,
      };

      await page.render(renderContext);
    }
  });
}
</script>

<template>
  <DialogRoot :open="open">
    <DialogPortal>
      <DialogOverlay
        class="bg-gray-800/20 data-[state=open]:animate-overlayShow fixed inset-0 z-30"
      />
      <DialogContent
        class="data-[state=open]:animate-contentShow fixed top-[50%] left-[50%] translate-x-[-50%] translate-y-[-50%] focus:outline-none z-[100] bg-gray-100 rounded-xl shadow-xl"
        @keydown.esc="onClose"
      >
        <div
          class="flex justify-between items-center gap-4 py-2 px-2 text-gray-700"
        >
          <DialogTitle class="ml-2 text-gray-600 text-sm font-semibold">
            {{ fileRef.name }}
          </DialogTitle>
          <div class="flex gap-1">
            <DIconButton
              :icon="Download"
              @click="store.downloadFile(fileRef)"
            ></DIconButton>
            <DialogClose @click="onClose" class="outline-none">
              <DIconButton :icon="X"></DIconButton>
            </DialogClose>
          </div>
        </div>

        <div
          class="max-h-[85vh] max-w-4xl flex justify-center mx-1 mb-1 bg-gray-50 rounded-lg overflow-hidden shadow-sm"
          :key="previewUrl"
        >
          <img
            v-if="fileType === 'image'"
            :src="previewUrl"
            :alt="fileRef.name"
            class="object-contain"
          />

          <video
            v-else-if="fileType === 'video'"
            controls="true"
            class="object-contain"
          >
            <source :src="previewUrl" :type="mimeType" />
          </video>

          <audio
            v-else-if="fileType === 'audio'"
            :controls="true"
            class="w-full"
          >
            <source :src="previewUrl" :type="mimeType" />
          </audio>

          <canvas
            @click.stop
            v-if="fileType === 'pdf'"
            ref="canvas"
            class="mx-auto block h-fit max-h-full w-fit object-contain"
          ></canvas>
          <div
            v-else-if="fileType === 'directory' || fileType === 'unknown'"
            class="flex items-center gap-4 p-10"
          >
            <div
              class="size-14 bg-white shadow-md grid place-items-center rounded-2xl"
            >
              <Folder
                v-if="fileType === 'directory'"
                class="size-7 fill-blue-600 stroke-blue-600"
              />
              <FileIcon v-else class="size-7 stroke-blue-600" />
            </div>
            <div>
              <div class="text-gray-800 text-2xl">
                {{ fileRef.name }}
              </div>
              <div class="text-gray-500 text-sm">
                {{ formatTime(new Date(fileRef.created_at)) }}
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
