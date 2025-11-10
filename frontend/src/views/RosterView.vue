<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import type { Roster } from '@gewis/grooster-backend-ts';
import { useRoute } from 'vue-router';
import RosterTable from '@/components/roster/RosterTable.vue';
import RosterAssignment from '@/components/roster/RosterAssignment.vue';
import { useRosterStore } from '@/stores/roster.store.js';
import RosterAddDialog from '@/components/roster/dialogs/RosterAddDialog.vue';
import ApiService from '@/services/ApiService';

const route = useRoute();

const rosterStore = useRosterStore();
const rosters = computed(() => Object.values(rosterStore.rosters));
const users = ref();

const selectedRoster = ref<Roster | null>(null);
const addDialog = ref(false);

onMounted(async () => {
  await fetchRosters();
});

watch(
  () => rosters.value,
  (newRosters: Roster[]) => {
    if (newRosters) {
      if (!selectedRoster.value) return;

      const updated = newRosters.find((response) => response.id === selectedRoster.value?.id);

      // Check if the object deeply changed or not
      if (updated && JSON.stringify(updated) !== JSON.stringify(selectedRoster.value)) {
        selectedRoster.value = updated;
      }
    }
  },
  { immediate: true, deep: true },
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
</script>

<template>
  <Card v-if="rosters" class="h-full w-full">
    <template #content>
      <div class="flex flex-col gap-5 h-full">
        <div v-if="selectedRoster?.id" class="flex flex-col gap-5 flex-1">
          <RosterTable :id="selectedRoster.id" />
          <RosterAssignment v-if="selectedRoster.saved" :id="selectedRoster.id" />
        </div>
        <div class="flex items-center justify-center">
          <Select
            v-model="selectedRoster"
            class="w-full md:w-56 mr-4"
            option-label="name"
            :option-value="(roster) => roster"
            :options="rosters"
            placeholder="Select a roster"
          />
          <Button label="Add Roster" @click="addDialog = true" />
        </div>
      </div>
    </template>
  </Card>
  <RosterAddDialog :open="addDialog" @close="addDialog = false" />
</template>

<style scoped></style>
