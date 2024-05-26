<script setup lang="ts">
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

import { useFileStore } from "@/stores/file";

const fileStore = useFileStore();

const checkboxOne = ref(false);
const checkboxTwo = ref(false);
const person = ref("pedro");
</script>

<template>
  <ContextMenuRoot>
    <ContextMenuTrigger as-child>
      <div>
        <slot name="content" />
      </div>
    </ContextMenuTrigger>
    <ContextMenuPortal>
      <ContextMenuContent
        class="min-w-[220px] z-30 border border-gray-100 bg-white outline-none rounded-md p-[5px] shadow will-change-[opacity,transform]"
        :side-offset="5"
      >
        <slot name="menu" />
      </ContextMenuContent>
    </ContextMenuPortal>
  </ContextMenuRoot>
</template>
