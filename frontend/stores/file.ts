import { defineStore } from "pinia";

export const useFileStore = defineStore("file", () => {
  const selectedFiles = ref<string[]>([]);
  const authStore = useAuthStore();

  const files = ref([]);

  const sharedDrive = ref(null);
  const parentId = ref<string | null>(null);
  const deleted = ref(false);

  const options = computed(() => {
    return {
      sharedDrive: sharedDrive.value,
      parentId: parentId.value,
      deleted: deleted.value,
    };
  });

  watch(options, async () => {
    await fetchFiles();
  });

  function resetOptions() {
    sharedDrive.value = null;
    parentId.value = null;
    deleted.value = false;
  }

  async function fetchFiles() {
    const response = await $fetch<any>("http://localhost:1323/files/", {
      method: "GET",
      query: {
        parent_id: parentId.value,
        shared_drive: sharedDrive.value,
        deleted: deleted.value,
      },
      headers: {
        Authorization: `${authStore.userToken}`,
      },
    });

    // @ts-ignore
    files.value = response.data || [];

    return true;
  }

  async function createFolder(parentId = null) {
    const formData = new FormData();
    formData.append("name", "Untitled Folder");
    formData.append("is_folder", "true");

    const response = await $fetch<any>("http://localhost:1323/files/", {
      method: "POST",
      body: formData,
      headers: {
        Authorization: `${authStore.userToken}`,
      },
    });

    await fetchFiles();
  }

  async function deleteFile(id: string) {
    // already remove from files
    files.value = files.value.filter((file) => file.id !== id);

    await $fetch(`http://localhost:1323/files/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `${authStore.userToken}`,
      },
    });

    await fetchFiles();
  }

  async function updateFileName(file: File) {
    await $fetch(`http://localhost:1323/files/${file.id}`, {
      method: "PATCH",
      body: JSON.stringify({
        id: file.id,
        name: file.name,
      }),
      headers: {
        Authorization: `${authStore.userToken}`,
      },
    });

    await fetchFiles();
  }

  async function uploadFiles(files: File[]) {
    for (const file of files) {
      const form = new FormData();
      form.append("file", file);

      if (parentId.value) {
        form.append("parent_id", parentId.value);
      }

      await $fetch("http://localhost:1323/files/", {
        method: "POST",
        body: form,
        headers: {
          Authorization: `${authStore.userToken}`,
        },
      });

      await fetchFiles();
    }
  }

  function fileSelected(id: string) {
    return selectedFiles.value.includes(id);
  }

  function fileToggleSelected(id: string) {
    if (fileSelected(id)) {
      selectedFiles.value = selectedFiles.value.filter(
        (fileId) => fileId !== id,
      );
    } else {
      selectedFiles.value = [id];
    }
  }

  function clearSelectedFiles() {
    selectedFiles.value = [];
  }

  async function downloadFile(file: File) {
    const response = await fetch(
      `http://localhost:1323/files/${file.id}/download`,
      {
        method: "GET",
        headers: {
          Authorization: `${authStore.userToken}`,
        },
      },
    );

    const blob = await response.blob();

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = file.name;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
  }

  async function getPreviewUrl(id: string) {
    const response = await fetch(`http://localhost:1323/files/${id}/preview`, {
      method: "GET",
      headers: {
        Authorization: `${authStore.userToken}`,
      },
    });
    const json = await response.json();

    return json.url;
  }

  function selectNextFile() {
    const index = files.value.findIndex(
      (file) => file.id === selectedFiles.value[0],
    );
    if (index === -1) {
      return;
    }

    if (index + 1 < files.value.length) {
      selectedFiles.value = [files.value[index + 1].id];
    }
  }

  function selectPreviousFile() {
    const index = files.value.findIndex(
      (file) => file.id === selectedFiles.value[0],
    );
    if (index === -1) {
      return;
    }

    if (index - 1 >= 0) {
      selectedFiles.value = [files.value[index - 1].id];
    }
  }

  return {
    files,
    parentId,
    sharedDrive,
    fetchFiles,
    deleteFile,
    createFolder,
    resetOptions,
    uploadFiles,
    selectedFiles,
    fileSelected,
    fileToggleSelected,
    clearSelectedFiles,
    downloadFile,
    updateFileName,
    getPreviewUrl,
    selectNextFile,
    selectPreviousFile,
  };
});
