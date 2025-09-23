<script setup lang="ts">
import { User, SavedShiftUpdateRequest } from '@gewis/grooster-backend-ts';
import { onMounted, reactive, computed, watch, ref } from 'vue';
import { useRosterStore } from '@/stores/roster.store';
import ApiService from '@/services/ApiService';
import { useRoute } from 'vue-router';

const route = useRoute();

const props = defineProps<{
  id: number;
}>();

const rosterStore = useRosterStore();
const savedRoster = computed(() => rosterStore.getSavedRoster(props.id));

const users = ref<User[]>([]);
const loadUsers = async () => {
  const response = await ApiService.user.userGet(parseInt(route.params.id as string));
  users.value = response.data;
};

const selectedId = ref(null);

const shiftAssignedUsers = reactive<Record<number, User[]>>({});

const selectedIds = reactive<Record<number, number | null>>({});

onMounted(async () => {
  await rosterStore.fetchSavedShifts(props.id);
  await loadUsers();
});

watch(
  [savedRoster, users],
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

const addUser = async (userId, shiftId: number) => {
  if (!userId) return;

  const user = users.value.find((user) => user.id === userId);
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
</script>

<template>
  <div v-if="savedRoster && savedRoster.length > 0">
    <div class="grid grid-cols-1 gap-4">
      <div v-for="shift in savedRoster" :key="shift.id">
        <div class="flex gap-4">
          <strong>{{ shift.rosterShift.name }}</strong>
          <!-- Render select for each selected user -->
          <div
            v-for="(user, index) in shiftAssignedUsers[shift.id]"
            :key="user.id + '-' + shift.id"
            class="flex items-center gap-2"
          >
            <Select
              v-model="user.id"
              :options="users"
              option-label="name"
              option-value="id"
              class="w-40"
              :disabled="true"
            />
            <Button
              icon="pi pi-times"
              class="p-button-rounded p-button-text p-button-danger"
              @click="removeUserFromShift(shift.id, index)"
              :aria-label="'Remove ' + user.name"
            />
          </div>

          <Select
            :key="'select-' + shift.id"
            v-model="selectedIds[shift.id]"
            :options="users.filter((user) => !shiftAssignedUsers[shift.id].some((u) => u.id === user.id))"
            optionLabel="name"
            optionValue="id"
            placeholder="Add a user"
            class="w-40"
            @change="(userId) => addUser(userId.value, shift.id)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mb-4 {
  margin-bottom: 1rem;
}

.mb-2 {
  margin-bottom: 0.5rem;
}
</style>
