<script setup lang="ts">
import { computed } from 'vue';
import { useRosterStore } from '@/stores/roster.store';

const props = withDefaults(
  defineProps<{
    open: boolean;
  }>(),
  {
    open: false,
  },
);

const emit = defineEmits(['close']);
const rosterStore = useRosterStore();

const visible = computed({
  get: () => props.open,
  set: () => {
    emit('close');
  },
});

const deleteRoster = async () => {
  await rosterStore.deleteRoster(rosterStore.selectedRosterId);

  emit('close');
};
</script>

<template>
  <Dialog
    v-model:visible="visible"
    class="mx-4"
    :dismissable-mask="true"
    :draggable="false"
    modal
    :pt="{
      root: 'border-none shadow-xl rounded-xl overflow-hidden',
      header: 'bg-white border-b border-gray-100 p-5',
      content: 'p-0', // We'll handle padding inside
    }"
    :style="{ width: '28rem' }"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="flex items-center justify-center w-10 h-10 rounded-full bg-red-50">
          <i class="pi pi-exclamation-triangle text-red-500 text-lg"></i>
        </div>
        <span class="font-bold text-xl text-gray-800">Delete Roster</span>
      </div>
    </template>

    <div class="p-6">
      <p class="text-gray-600 leading-relaxed">
        Are you sure you want to delete this roster?
        <span class="block mt-1 font-medium text-gray-900">This action cannot be undone.</span>
      </p>
    </div>

    <template #footer>
      <div class="flex gap-3 justify-end w-full p-4 bg-gray-50 border-t border-gray-100">
        <Button
          class="hover:bg-gray-200 transition-colors"
          label="Cancel"
          severity="secondary"
          text
          @click="emit('close')"
        />
        <Button class="px-6 shadow-md shadow-red-100" label="Yes, Delete" severity="danger" @click="deleteRoster" />
      </div>
    </template>
  </Dialog>
</template>
