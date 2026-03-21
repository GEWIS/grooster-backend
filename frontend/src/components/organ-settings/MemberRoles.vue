<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { GEWISRoosterInternalModelsOrganRole, User, UserOrgan } from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import { useOrganStore } from '@/stores/organ.store';
import RoleChangeConfirmation from '@/components/organ-settings/dialogs/RoleChangeConfirmation.vue';

export type UserRole = (typeof GEWISRoosterInternalModelsOrganRole)[keyof typeof GEWISRoosterInternalModelsOrganRole];

const roleOptions = Object.entries(GEWISRoosterInternalModelsOrganRole).map(([key, value]) => ({
  label: key.replace('Role', ''),
  value: value,
}));

const organStore = useOrganStore();

const users = ref<User[]>([]);
const userOrganSettings = ref<UserOrgan[]>([]);
const selectedUser = ref<User | null>(null);

const currentRole = computed(() => getUserRole(selectedUser.value.id));
const pendingRole = ref<UserRole | null>(null);

type Dialog = 'Confirm' | 'None';
const currentDialog = ref<Dialog>('None');

const fetchUsers = async () => {
  try {
    const response = await ApiService.user.userGet(organStore.organ);
    users.value = response.data;
  } catch (e) {
    console.error(e);
    users.value = [];
  }
};

const fetchUserOrganSettings = async () => {
  try {
    const response = await ApiService.organ.getMembersSettings(organStore.organ);
    userOrganSettings.value = response.data;
  } catch (e) {
    console.error(e);
    users.value = [];
  }
};

const getUserRole = (userId: number) => {
  return userOrganSettings.value.find((organ: UserOrgan) => organ.userId == userId).role;
};

const handleRoleUpdate = (role: UserRole) => {
  pendingRole.value = role;
  currentDialog.value = 'Confirm';
};

const confirmRoleChange = async (userId: number, role: UserRole) => {
  try {
    const { data } = await ApiService.organ.organIdMemberUserIdRolePatch(organStore.organ, userId, { role });
    userOrganSettings.value = userOrganSettings.value.map((userOrgan) => {
      if (userOrgan.userId === data.userId) {
        return data;
      } else {
        return userOrgan;
      }
    });
  } catch (e) {
    console.error(e);
  } finally {
    pendingRole.value = null;
    currentDialog.value = 'None';
  }
};

onMounted(fetchUsers);
onMounted(fetchUserOrganSettings);
</script>

<template>
  <div
    v-if="users.length > 0"
    class="bg-white border border-emerald-100 rounded-xl overflow-hidden shadow-sm max-w-2xl"
  >
    <div class="bg-emerald-950 p-6 flex items-center gap-4">
      <div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-900 shadow-sm">
        <i class="pi pi-users text-xl" />
      </div>
      <h2 class="text-white font-semibold text-lg">Organization Members</h2>
    </div>

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
          <div
            class="h-16 w-16 bg-emerald-100 rounded-full flex items-center justify-center mx-auto text-emerald-900 mb-2"
          >
            <span class="text-xl font-bold">{{ selectedUser.name.charAt(0) }}</span>
          </div>
          <div>
            <h3 class="font-bold text-gray-900">{{ selectedUser.name }}</h3>
            <p class="text-xs text-gray-500 uppercase tracking-widest mt-1">
              Current Role: {{ getUserRole(selectedUser.id) }}
            </p>
          </div>

          <div class="pt-4">
            <label class="block text-xs font-medium text-gray-400 mb-2 text-left">CHANGE ROLE</label>
            <Select
              class="w-full border-emerald-200 text-sm focus:ring-emerald-500"
              :model-value="currentRole"
              option-label="label"
              option-value="value"
              :options="roleOptions"
              placeholder="Select a Role"
              @update:model-value="(value) => handleRoleUpdate(value)"
            />
          </div>
        </div>

        <div v-else class="text-gray-400 italic text-sm">Select a member to manage their role</div>
      </div>
    </div>
  </div>
  <RoleChangeConfirmation
    v-if="pendingRole"
    :open="currentDialog === 'Confirm'"
    :role="pendingRole"
    :user-id="selectedUser.id"
    @close="currentDialog = 'None'"
    @confirm="({ userId, role }) => confirmRoleChange(userId, role)"
  />
</template>

<style scoped></style>
