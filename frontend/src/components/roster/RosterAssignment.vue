<script setup lang="ts">
import { User, SavedShiftUpdateRequest, SavedShift } from '@gewis/grooster-backend-ts';
import { onMounted, reactive, computed, watch } from 'vue';
import { useRosterStore } from '@/stores/roster.store';

const props = defineProps<{
  id: number;
}>();

const rosterStore = useRosterStore();
const savedRoster = computed(() => {
  const data = rosterStore.getSavedRoster(props.id);
  return Array.isArray(data) ? data : (data?.savedShifts ?? []);
});

const savedRosterOrdering = computed(() => {
  const data = rosterStore.getSavedRoster(props.id);
  return Array.isArray(data) ? [] : (data?.savedShiftOrdering ?? []);
});
const shiftAssignedUsers = reactive<Record<number, User[]>>({});

const selectedIds = reactive<Record<number, number | null>>({});

onMounted(async () => {
  await rosterStore.fetchSavedShifts(props.id);
});

watch(
  [savedRoster, savedRosterOrdering],
  ([newSaved, newUsers]) => {
    if (!newSaved || !newUsers) return;

    Object.keys(shiftAssignedUsers).forEach((key) => delete shiftAssignedUsers[parseInt(key)]);
    Object.keys(selectedIds).forEach((key) => delete selectedIds[parseInt(key)]);

    for (const shift of newSaved) {
      shiftAssignedUsers[shift.id] = [...shift.users];
      selectedIds[shift.id] = null;
    }
  },
  { immediate: true, deep: true },
);

const addUser = async (userId: number, shiftId: number) => {
  if (!userId) return;

  const shiftName = savedRoster.value.find((shift) => shift.id === shiftId).rosterShift.name;
  const users = savedRosterOrdering.value.find((ordering) => ordering.shiftName === shiftName).users;
  const user = users.find((user) => user.id === userId);

  if (!user) return;

  const assignedIds = shiftAssignedUsers[shiftId].map((user) => user.id);
  if (assignedIds.includes(user.id)) return;

  shiftAssignedUsers[shiftId].push(user);

  const params: SavedShiftUpdateRequest = {
    users: [...assignedIds, user.id],
  };

  try {
    await rosterStore.updateSavedShift(shiftId, params);
  } catch (error) {
    console.error(error);
    const idx = shiftAssignedUsers[shiftId].findIndex((u) => u.id === user.id);
    if (idx !== -1) shiftAssignedUsers[shiftId].splice(idx, 1);
  }

  selectedIds[shiftId] = null;
};

const removeUserFromShift = async (shiftId: number, index: number) => {
  const removedUser = shiftAssignedUsers[shiftId].splice(index, 1)[0];

  const updatedIds = shiftAssignedUsers[shiftId].map((user) => user.id);
  const params: SavedShiftUpdateRequest = {
    users: [...updatedIds],
  };

  try {
    await rosterStore.updateSavedShift(shiftId, params);
  } catch (error) {
    console.error(error);
    shiftAssignedUsers[shiftId].splice(index, 0, removedUser);
  }
};

const availableUsersForShift = (shift: SavedShift) => {
  const ordering = savedRosterOrdering.value.find((o) => o.shiftName === shift.rosterShift.name);
  if (!ordering) return [];

  return ordering.users.filter((user) => !shiftAssignedUsers[shift.id]?.some((u) => u.id === user.id));
};
</script>

<template>
  <div v-if="savedRoster && savedRoster.length > 0" class="space-y-6 max-w-6xl mx-auto p-4">
    <div class="border-b border-slate-200 pb-4">
      <h2 class="text-xl font-bold text-slate-800">Shift Assignments</h2>
      <p class="text-sm text-slate-500">Assign and manage team members for each active shift.</p>
    </div>

    <div class="grid grid-cols-1 gap-4">
      <div
        v-for="shift in savedRoster"
        :key="shift.id"
        class="bg-white border border-slate-200 rounded-xl shadow-sm hover:shadow-md transition-shadow p-5"
      >
        <div class="flex flex-col md:flex-row md:items-center gap-6">
          <div class="min-w-[160px]">
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-indigo-50 text-indigo-700 border border-indigo-100"
            >
              {{ shift.rosterShift.name }}
            </span>
          </div>

          <div class="flex flex-wrap items-center gap-3 flex-grow">
            <div
              v-for="(user, index) in shiftAssignedUsers[shift.id]"
              :key="user.id + '-' + shift.id"
              class="group flex items-center bg-slate-50 border border-slate-200 rounded-lg pl-3 pr-1 py-1 transition-colors hover:bg-slate-100"
            >
              <span class="text-sm font-medium text-slate-700">{{ user.name }}</span>
              <Button
                class="!p-0 !w-7 !h-7 !text-slate-400 group-hover:!text-red-500"
                icon="pi pi-times"
                rounded
                text
                @click="removeUserFromShift(shift.id, index)"
              />
            </div>

            <div class="relative min-w-[180px]">
              <Select
                :key="'select-' + shift.id"
                v-model="selectedIds[shift.id]"
                class="!w-full !rounded-lg !border-dashed !border-slate-300 !bg-transparent hover:!border-indigo-400 transition-all !text-sm"
                option-label="name"
                option-value="id"
                :options="availableUsersForShift(shift)"
                placeholder="+ Add Team Member"
                @change="(userId) => addUser(userId.value, shift.id)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
