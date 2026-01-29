<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { RosterUpdateRequest } from '@gewis/grooster-backend-ts';
import { useRosterStore } from '@/stores/roster.store';

const props = withDefaults(defineProps<{ open: boolean }>(), { open: false });

const emit = defineEmits(['close']);
const rosterStore = useRosterStore();

const newName = ref('');
const isSubmitting = ref(false);

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      const currentRoster = rosterStore.rosters[rosterStore.selectedRosterId];
      newName.value = currentRoster?.name || '';
    }
  },
);

const visible = computed({
  get: () => props.open,
  set: () => {
    emit('close');
  },
});

const updateName = async () => {
  if (!newName.value.trim()) return;

  isSubmitting.value = true;
  try {
    const params: RosterUpdateRequest = { name: newName.value };
    await rosterStore.updateRoster(rosterStore.selectedRosterId, params);
    emit('close');
  } finally {
    isSubmitting.value = false;
  }
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
        <div class="flex items-center justify-center w-11 h-11 rounded-full bg-amber-50 border border-amber-100">
          <i class="pi pi-question-circle text-amber-600 text-xl"></i>
        </div>
        <span class="font-bold text-xl text-gray-800">Change Roster Name</span>
      </div>
    </template>

    <div class="flex flex-col gap-2 pt-4">
      <label class="text-sm font-semibold text-gray-700" for="rosterName">New Display Name</label>
      <div class="relative">
        <InputText
          id="rosterName"
          v-model="newName"
          autofocus
          class="w-full p-3 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
          placeholder="e.g. Summer Shift 2024"
          @keydown.enter="updateName"
        />
      </div>
      <small class="text-gray-500">Give your roster a descriptive name that others will recognize.</small>
    </div>

    <template #footer>
      <Button label="Cancel" severity="secondary" text @click="emit('close')" />
      <Button
        class="px-5 bg-blue-600 border-none hover:bg-blue-700"
        icon="pi pi-check"
        label="Save Changes"
        @click="updateName"
      />
    </template>
  </Dialog>
</template>

<style scoped></style>
