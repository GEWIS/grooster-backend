<script setup lang="ts">
import {
  GEWISRoosterInternalModelsGroupPriority,
  GroupPriorityUpdateParam,
  ShiftGroupPriority,
  User,
} from '@gewis/grooster-backend-ts';
import { computed, onMounted, ref } from 'vue';
import { useOrganStore } from '@/stores/organ.store';
import ApiService from '@/services/ApiService';

export type GroupPriority =
  (typeof GEWISRoosterInternalModelsGroupPriority)[keyof typeof GEWISRoosterInternalModelsGroupPriority];

const priorityOptions = Object.entries(GEWISRoosterInternalModelsGroupPriority).map(([key, value]) => ({
  label: key,
  value: value,
}));

const props = defineProps<{
  groupId: number;
}>();

const organStore = useOrganStore();
const users = ref<User[]>([]);

const userGroupPriorities = ref<ShiftGroupPriority[]>([]);
const selectedUser = ref<User | null>(null);

const currentPriority = computed(() => getUserPriority(selectedUser.value.id));

const fetchUsers = async () => {
  try {
    const response = await ApiService.user.userGet(organStore.organ);
    users.value = response.data;
  } catch (e) {
    console.error(e);
    users.value = [];
  }
};

const fetchUserPriorities = async () => {
  try {
    const response = await ApiService.shiftGroupApi.getShiftGroupPriorities(props.groupId);
    userGroupPriorities.value = response.data ?? [];
  } catch (e) {
    console.error(e);
  }
};

const getUserPriority = (userId: number): GroupPriority => {
  const found = userGroupPriorities.value.find((group: ShiftGroupPriority) => group.userId === userId);

  return found?.priority ?? GEWISRoosterInternalModelsGroupPriority.Default;
};

const handlePriorityUpdate = async (priority: GroupPriority, userId: number) => {
  try {
    const updateParams: GroupPriorityUpdateParam = {
      userId,
      priority,
    };

    await ApiService.shiftGroupApi.updateShiftGroupPriority(props.groupId, updateParams);

    const index = userGroupPriorities.value.findIndex((ugp) => ugp.userId === userId);

    if (index !== -1) {
      userGroupPriorities.value[index].priority = priority;
    } else {
      userGroupPriorities.value.push({
        userId,
        priority,
      });
    }
  } catch (e) {
    console.error(e);
  }
};

onMounted(fetchUsers);
onMounted(fetchUserPriorities);
</script>

<template>
  <div class="flex h-[400px]">
    <div class="w-1/3 border-r border-emerald-50 overflow-y-auto bg-emerald-50/30">
      <div
        v-for="user in users"
        :key="user.id"
        :class="[
          'p-4 cursor-pointer transition-colors border-b border-emerald-50',
          selectedUser?.id === user.id ? 'bg-white font-bold text-emerald-700' : 'hover:bg-emerald-100/50',
        ]"
        @click="selectedUser = user"
      >
        <p class="truncate text-sm">{{ user.name }}</p>
      </div>
    </div>

    <div class="w-2/3 p-6 flex flex-col justify-center items-center text-center">
      <div v-if="selectedUser" class="space-y-4 w-full">
        <div>
          <h3 class="font-bold text-gray-900">{{ selectedUser.name }}</h3>
          <p class="text-xs text-gray-500 uppercase tracking-widest mt-1">
            Current Priority: {{ getUserPriority(selectedUser.id) }}
          </p>
        </div>

        <div class="pt-4">
          <label class="block text-xs font-medium text-gray-400 mb-2 text-left">CHANGE ROLE</label>
          <Select
            class="w-full border-emerald-200 text-sm focus:ring-emerald-500"
            :model-value="currentPriority"
            option-label="label"
            option-value="value"
            :options="priorityOptions"
            placeholder="Select a Role"
            @update:model-value="(value) => handlePriorityUpdate(value, selectedUser.id)"
          />
        </div>
      </div>

      <div v-else class="text-gray-400 italic text-sm">Select a member to manage their role</div>
    </div>
  </div>
</template>

<style scoped></style>
