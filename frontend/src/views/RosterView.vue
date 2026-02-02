<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import type { Roster } from '@gewis/grooster-backend-ts';
import { useRoute } from 'vue-router';
import RosterTable from '@/components/roster/RosterTable.vue';
import RosterAssignment from '@/components/roster/RosterAssignment.vue';
import { useRosterStore } from '@/stores/roster.store.js';
import AddDialog from '@/components/roster/dialogs/AddDialog.vue';
import ApiService from '@/services/ApiService';
import DeleteDialog from '@/components/roster/dialogs/DeleteDialog.vue';
import EditDialog from '@/components/roster/dialogs/EditDialog.vue';

type DialogType = 'add' | 'edit' | 'delete';

const route = useRoute();

const rosterStore = useRosterStore();
const rosters = computed(() => Object.values(rosterStore.rosters));
const users = ref();
const selectedRosterId = computed(() => rosterStore.selectedRosterId);

const selectedRoster = ref<Roster | null>(null);

const activeDialog = ref<DialogType>(null);

onMounted(async () => {
  await fetchRosters();
});

watch(
  () => rosters.value,
  (newRosters: Roster[]) => {
    if (!newRosters || newRosters.length === 0) {
      selectedRoster.value = null;
      rosterStore.setSelectedRoster(null);
      return;
    }

    const storeSelectedId = rosterStore.selectedRosterId;
    if (storeSelectedId != null) {
      const updated = newRosters.find((response) => response.id === storeSelectedId);
      if (updated) {
        selectedRoster.value = updated;
        return;
      }
    }

    if (!selectedRoster.value) {
      const first = newRosters[0];
      if (first?.id != null) {
        selectedRoster.value = first;
        rosterStore.setSelectedRoster(first.id);
      }
    }
  },
  { immediate: true, deep: true },
);

watch(
  () => selectedRoster.value?.id ?? null,
  (newId) => {
    if (newId !== rosterStore.selectedRosterId) {
      rosterStore.setSelectedRoster(newId);
    }
  },
);

watch(
  () => selectedRosterId.value,
  (newId) => {
    if (newId == null) {
      selectedRoster.value = null;
      return;
    }
    const updated = rosters.value.find((response) => response.id === newId);
    if (updated) {
      selectedRoster.value = updated;
    }
  },
);

async function fetchRosters() {
  try {
    await rosterStore.fetchRosters(parseInt(route.params.id as string));

    const response = await ApiService.user.userGet(parseInt(route.params.id as string));
    users.value = response.data;
  } catch (e) {
    console.error(e);
  }
}

const openDialog = (action: DialogType) => {
  activeDialog.value = action;
};

const closeDialog = () => {
  activeDialog.value = null;
};
</script>

<template>
  <Card v-if="rosters" class="h-full w-full">
    <template #content>
      <div class="flex flex-col gap-5 h-full">
        <div v-if="selectedRoster?.id" class="flex flex-col gap-5 flex-1">
          <RosterTable :id="selectedRoster.id" />
          <RosterAssignment v-if="selectedRoster.saved" :id="selectedRoster.id" />
        </div>
        <div class="flex items-center justify-center gap-2">
          <Select
            v-model="selectedRoster"
            class="w-full md:w-56 mr-4"
            :option-value="(roster) => roster"
            :options="rosters"
            placeholder="Select a roster"
          >
            <template #value="slotProps">
              <div v-if="slotProps.value">{{ slotProps.value.name }} - {{ slotProps.value.date.split('T')[0] }}</div>
              <span v-else>{{ slotProps.placeholder }}</span>
            </template>
            <template #option="slotProps">
              {{ slotProps.option.name }} ({{ slotProps.option.date.split('T')[0] }})
            </template>
          </Select>
          <div class="flex gap-2 p-4">
            <Button class="p-button-success" icon="pi pi-plus" label="Add Roster" @click="openDialog('add')" />

            <Button
              v-if="selectedRoster"
              class="p-button-outlined"
              icon="pi pi-pencil"
              label="Edit"
              @click="openDialog('edit')"
            />

            <Button
              v-if="selectedRoster"
              class="p-button-danger p-button-outlined"
              icon="pi pi-trash"
              label="Delete"
              @click="openDialog('delete')"
            />
          </div>
        </div>
      </div>
    </template>
  </Card>
  <AddDialog :open="activeDialog === 'add'" @close="closeDialog" />
  <DeleteDialog :open="activeDialog === 'delete'" @close="closeDialog" />
  <EditDialog :open="activeDialog === 'edit'" @close="closeDialog" />
</template>

<style scoped></style>
