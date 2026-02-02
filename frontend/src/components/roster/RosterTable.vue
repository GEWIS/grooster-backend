<script setup lang="ts">
import { computed, reactive, watch, ref, onMounted } from 'vue';
import {
  RosterAnswer,
  RosterAnswerUpdateRequest,
  Roster,
  RosterAnswerCreateRequest,
  User,
} from '@gewis/grooster-backend-ts';
import { useRoute } from 'vue-router';
import { useRosterStore } from '@/stores/roster.store';
import { getGEWISId } from '@/helpers/TokenHelper';
import ApiService from '@/services/ApiService';

const route = useRoute();

const props = defineProps<{
  id: number;
}>();

const rosterStore = useRosterStore();
const roster = computed<Roster | undefined>(() => rosterStore.getRoster(props.id));

const users = ref<User[]>();
const loadUsers = async () => {
  const response = await ApiService.user.userGet(parseInt(route.params.id as string));
  users.value = response.data;
};

onMounted(loadUsers);

const rosterValues = computed(() => Object.values(roster.value.values || {}));

const shiftAnswers = reactive({});
const visible = ref(false);
const shiftName = ref('');

watch(
  [() => roster.value, () => users.value],
  ([newRoster, newUsers]) => {
    if (!newRoster || !newUsers) return;
    const existingAnswers = newRoster.rosterAnswer || [];

    for (const user of newUsers) {
      if (!shiftAnswers[user.id]) {
        shiftAnswers[user.id] = {};
      }

      if (newRoster.rosterShift) {
        for (const shift of newRoster.rosterShift) {
          const value = existingAnswers.find(
            (answer: RosterAnswer) => answer.rosterShiftId === shift.id && answer.userId === user.id,
          );
          shiftAnswers[user.id][shift.id] = value ?? { id: undefined, value: undefined };
        }
      }
    }
  },
  { immediate: true, deep: true },
);

async function onAnswerChange(user: number, shift: number, newValue: string) {
  if (!shiftAnswers[user][shift].id) {
    const createParams: RosterAnswerCreateRequest = {
      rosterId: props.id,
      rosterShiftId: shift,
      userId: user,
      value: newValue,
    };

    try {
      await rosterStore.createAnswer(createParams);

      shiftAnswers[user][shift] = {
        id: props.id,
        value: newValue,
      };
    } catch (error) {
      console.error(error);
    }
  } else {
    const updateParam: RosterAnswerUpdateRequest = {
      value: newValue,
    };
    const rosterId = shiftAnswers[user][shift].id;

    try {
      await rosterStore.updateAnswer(rosterId, updateParam);

      shiftAnswers[user][shift] = {
        id: props.id,
        value: newValue,
      };
    } catch (error) {
      console.error(error);
    }
  }
}

async function addShift(name: string) {
  visible.value = false;

  try {
    await rosterStore.createShift({ name: name, rosterId: props.id });

    shiftName.value = '';
  } catch (error) {
    console.error('Not able to save shift:', error);
  }
}

async function removeShift(id: number) {
  try {
    await rosterStore.deleteShift(id, props.id);

    for (const userId in shiftAnswers) {
      delete shiftAnswers[userId][id];
    }
  } catch (error) {
    console.error('Could not delete shift:', error);
  }
}

async function saveRoster() {
  try {
    await rosterStore.saveRoster(props.id);
  } catch (error) {
    console.error(error);
  }
}

async function unSaveRoster() {
  await rosterStore.updateRoster(props.id, { saved: false });
}

async function moveShift(shiftId: number, direction: 'left' | 'right') {
  if (!roster.value?.rosterShift) return;

  const sortedShifts = [...roster.value.rosterShift].sort((a, b) => a.order - b.order);

  const currentIndex = sortedShifts.findIndex((s) => s.id === shiftId);
  const targetIndex = direction === 'left' ? currentIndex - 1 : currentIndex + 1;

  if (targetIndex < 0 || targetIndex >= sortedShifts.length) return;

  const currentShift = sortedShifts[currentIndex];
  const targetShift = sortedShifts[targetIndex];

  try {
    await Promise.all([
      rosterStore.updateShift(currentShift.id, { order: targetShift.order }),
      rosterStore.updateShift(targetShift.id, { order: currentShift.order }),
    ]);
  } catch (error) {
    console.error('Failed to reorder shifts:', error);
  }
}

const sortedUsers = computed(() => {
  if (!users.value) return [];

  return [...users.value].sort((a, b) => {
    if (a.id === getGEWISId()) return -1;
    if (b.id === getGEWISId()) return 1;
    return 0;
  });
});

const getStatusColorClass = (value: string) => {
  switch (value) {
    // Emerald/Green -> Soft Mint
    case 'J':
      return 'bg-emerald-300 text-emerald-900 font-semibold';
    // Yellow -> Soft Cream/Amber
    case 'X':
      return 'bg-amber-300 text-amber-900 font-semibold';
    // Red -> Soft Rose
    case 'N':
      return 'bg-rose-300 text-rose-900 font-semibold';
    // Purple -> Soft Lavender
    case 'L':
      return 'bg-purple-300 text-purple-900 font-semibold';
    default:
      return 'bg-transparent text-gray-400';
  }
};
</script>

<template>
  <div
    v-if="roster && Object.keys(shiftAnswers).length"
    class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden"
  >
    <div class="overflow-x-auto custom-scrollbar">
      <table class="w-full border-collapse">
        <thead>
          <tr class="bg-gray-50 border-b border-gray-200">
            <th class="text-left px-4 py-3 text-xs font-bold uppercase tracking-wider text-gray-500 min-w-[150px]">
              Users
            </th>

            <template v-if="roster.rosterShift">
              <th
                v-for="(shift, index) in roster.rosterShift.sort((a, b) => a.order - b.order)"
                :key="shift.id"
                class="px-3 py-3 text-center min-w-[140px] border-l border-gray-100"
              >
                <div class="flex flex-col items-center gap-1.5">
                  <span class="text-sm font-semibold text-gray-700 leading-tight">{{ shift.name }}</span>

                  <div v-if="!roster.saved" class="flex items-center gap-1">
                    <Button
                      class="!p-0 !w-6 !h-6"
                      :disabled="index === 0"
                      icon="pi pi-angle-left"
                      rounded
                      size="small"
                      text
                      @click="moveShift(shift.id, 'left')"
                    />

                    <Button
                      class="!p-0 !w-6 !h-6"
                      icon="pi pi-trash"
                      rounded
                      severity="danger"
                      text
                      @click="removeShift(shift.id)"
                    />

                    <Button
                      class="!p-0 !w-6 !h-6"
                      :disabled="index === roster.rosterShift.length - 1"
                      icon="pi pi-angle-right"
                      rounded
                      size="small"
                      text
                      @click="moveShift(shift.id, 'right')"
                    />
                  </div>
                </div>
              </th>
            </template>
          </tr>
        </thead>

        <tbody class="divide-y divide-gray-100">
          <tr v-for="user in sortedUsers" :key="user.id" class="hover:bg-blue-50/30 transition-colors">
            <td class="px-4 py-3 text-sm font-medium text-gray-900 border-r border-gray-50">
              {{ user.name }}
            </td>

            <template v-if="roster.rosterShift">
              <td
                v-for="shift in roster.rosterShift"
                :key="user.id + '-' + shift.id"
                class="p-1 transition-colors duration-200"
                :class="getStatusColorClass(shiftAnswers[user.id][shift.id].value)"
              >
                <Select
                  v-model="shiftAnswers[user.id][shift.id].value"
                  class="w-full !text-xs !border-none !ring-0 !shadow-none !bg-transparent"
                  :disabled="roster.saved || user.gewis_id != getGEWISId()"
                  :options="rosterValues"
                  placeholder="Select..."
                  @update:model-value="(value) => onAnswerChange(user.id, shift.id, value)"
                />
              </td>
            </template>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      class="bg-gray-50 px-4 py-4 border-t border-gray-200 flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <Button
        class="w-full sm:w-auto !justify-center"
        :disabled="roster.saved"
        icon="pi pi-plus"
        label="Add New Shift"
        severity="secondary"
        size="small"
        text
        @click="visible = true"
      />

      <div class="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
        <Button
          v-if="roster.saved"
          class="w-full sm:w-auto !justify-center"
          icon="pi pi-lock-open"
          label="Unlock Roster"
          outlined
          severity="warning"
          size="small"
          @click="unSaveRoster"
        />

        <Button
          v-else
          class="w-full sm:w-auto !justify-center"
          icon="pi pi-save"
          label="Lock Roster"
          severity="success"
          size="small"
          @click="saveRoster"
        />
      </div>
    </div>

    <Dialog v-model:visible="visible" class="p-fluid" header="Add New Shift" modal :style="{ width: '24rem' }">
      <div class="flex flex-col gap-4 py-2">
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-gray-700" for="shiftName">Shift Name</label>
          <InputText id="shiftName" v-model="shiftName" autofocus placeholder="e.g. Evening Standby" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <Button label="Cancel" severity="secondary" text @click="visible = false" />
          <Button icon="pi pi-check" label="Add Shift" @click="addShift(shiftName)" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
/* Optional: Make the horizontal scrollbar look a bit cleaner */
.custom-scrollbar::-webkit-scrollbar {
  height: 8px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #f1f1f1;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.overflow-x-auto {
  display: block;
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch; /* Momentum scrolling for iOS */
}

@media (max-width: 640px) {
  /* 1. Make the Select components easier to tap on mobile */
  :deep(.p-select) {
    height: 44px; /* Apple/Google recommended touch target size */
  }

  /* 2. Prevent text wrapping in the user column */
  td:first-child {
    white-space: nowrap;
    padding-right: 2rem;
  }

  /* 3. Make the shift buttons (arrows/trash) larger or spaced out */
  .flex.items-center.gap-1 {
    gap: 0.75rem;
    padding: 0.5rem 0;
  }
}

/* Ensure Select components don't look too bulky in the table */
:deep(.p-select) {
  background: transparent;
}
:deep(.p-select-label) {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
}

:deep(.bg-green-500 .p-select-label),
:deep(.bg-red-500 .p-select-label),
:deep(.bg-purple-500 .p-select-label) {
  color: white !important;
}

/* Ensure the dropdown icon also adapts */
:deep(.text-white .p-select-dropdown) {
  color: white !important;
}

/* Remove default PrimeVue background-color on the inner input */
:deep(.p-select) {
  background-color: transparent !important;
  border: none !important;
}
</style>
