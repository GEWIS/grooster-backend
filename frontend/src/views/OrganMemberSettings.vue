<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue';
import { UserOrgan } from '@gewis/grooster-backend-ts';
import { UpdateMemberSettingsParams } from '@gewis/grooster-backend-ts/dist/api';
import ApiService from '@/services/ApiService';
import { useUserStore } from '@/stores/user.store';
import { useOrganStore } from '@/stores/organ.store';

const userStore = useUserStore();
const organStore = useOrganStore();

const memberSettings = ref<UserOrgan | null>(null);
const editName = ref('');
const isSaving = ref(false);

const isChanged = computed(() => {
  return editName.value !== memberSettings.value?.username && editName.value.trim() !== '';
});

const fetchMemberData = async () => {
  if (!userStore.user?.id || !organStore.organ) return;

  try {
    const { data } = await ApiService.organ.organIdMemberUserIdGet(userStore.user.id, organStore.organ);
    memberSettings.value = data;
    editName.value = data.username || '';
  } catch (e) {
    console.error('Failed to fetch member settings:', e);
  }
};

watch(
  () => userStore.user?.id,
  async (newId) => {
    if (newId) await fetchMemberData();
  },
  { immediate: true },
);

onMounted(async () => {
  if (userStore.user?.id) {
    await fetchMemberData();
  }
});

const updateMemberSettings = async () => {
  if (!isChanged.value) return;

  isSaving.value = true;
  const params: UpdateMemberSettingsParams = {
    username: editName.value,
  };

  try {
    const { data } = await ApiService.organ.organIdMemberUserIdPatch(userStore.user.id, organStore.organ, params);
    memberSettings.value = data;
    editName.value = data.username;
  } catch (e) {
    console.error('Failed to update settings:', e);
  } finally {
    isSaving.value = false;
  }
};

const resetName = () => {
  if (memberSettings.value) {
    editName.value = memberSettings.value.username;
  }
};
</script>

<template>
  <div class="max-w-2xl mx-auto mt-10 px-4">
    <div v-if="memberSettings" class="bg-emerald-950/5 border border-emerald-100 rounded-xl overflow-hidden shadow-sm">
      <div class="bg-emerald-950 p-6 flex items-center gap-4">
        <div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-900 shadow-sm">
          <i class="pi pi-user text-xl" />
        </div>
        <div>
          <h2 class="text-white font-semibold text-lg leading-tight">Member Profile</h2>
          <p class="text-emerald-200/70 text-xs">Organ ID: {{ organStore.organ }}</p>
        </div>
      </div>

      <div class="p-8 space-y-6">
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold uppercase tracking-wider text-emerald-900/60 ml-1" for="username">
            Display Name
          </label>

          <div class="flex flex-col sm:flex-row gap-3">
            <div class="relative flex-1">
              <input
                id="username"
                v-model="editName"
                class="w-full px-4 py-2.5 rounded-lg border border-emerald-200 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition-all text-emerald-950 font-medium"
                placeholder="Your name in this organ..."
                type="text"
                @keyup.enter="updateMemberSettings"
              />
            </div>

            <div class="flex gap-2">
              <button
                v-if="isChanged"
                class="px-4 py-2 text-emerald-700 hover:bg-emerald-100 rounded-lg transition-colors text-sm font-medium"
                @click="resetName"
              >
                Cancel
              </button>

              <button
                class="flex items-center gap-2 px-6 py-2.5 bg-emerald-800 text-white rounded-lg font-medium hover:bg-emerald-900 disabled:opacity-30 disabled:cursor-not-allowed transition-all shadow-sm"
                :disabled="!isChanged || isSaving"
                @click="updateMemberSettings"
              >
                <i v-if="isSaving" class="pi pi-spin pi-spinner text-xs" />
                <span>{{ isSaving ? 'Saving' : 'Save Changes' }}</span>
              </button>
            </div>
          </div>

          <p class="text-[11px] text-emerald-600 italic ml-1 mt-1">
            * This name is specific to this organ and will appear on all generated rosters.
          </p>
        </div>
      </div>
    </div>

    <div
      v-else
      class="animate-pulse bg-emerald-50/50 rounded-xl p-8 border border-emerald-100 h-64 flex flex-col justify-center items-center gap-4"
    >
      <div class="h-12 w-12 bg-emerald-200 rounded-full"></div>
      <div class="h-4 w-48 bg-emerald-200 rounded"></div>
    </div>
  </div>
</template>
