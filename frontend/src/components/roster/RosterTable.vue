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
</script>

<template>
  <div v-if="roster && Object.keys(shiftAnswers).length" class="flex flex-row gap-2 w-full">
    <div class="max-w-5xl overflow-x-auto">
      <table class="table-auto w-full min-w-full text-xs sm:text-sm lg:text-base">
        <thead>
          <tr>
            <th class="user-column text-[10px] sm:text-sm lg:text-base px-1 py-0.5">Users</th>
            <template v-if="roster.rosterShift">
              <th
                v-for="shift in roster.rosterShift"
                :key="shift.id"
                class="shift-column text-[10px] sm:text-sm lg:text-base px-1 py-0.5"
              >
                <div class="flex flex-row gap-1">
                  <span class="block mb-1 sm:mb-2">{{ shift.name }}</span>
                  <div v-if="!roster.saved">
                    <Button class="sm:scale-50 lg:scale-60" @click="removeShift(shift.id)">
                      <i class="pi pi-times ml-1 sm:ml-2 lg:ml-4 text-[10px] sm:text-sm"></i>
                    </Button>
                  </div>
                </div>
              </th>
            </template>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td class="user-cell text-[10px] sm:text-sm lg:text-base px-1 py-0.5">
              {{ user.name }}
            </td>
            <template v-if="roster.rosterShift">
              <td
                v-for="shift in roster.rosterShift"
                :key="user.id + '-' + shift.id"
                class="shift-cell px-1 sm:px-2 lg:px-4 py-1 sm:py-2"
              >
                <Select
                  v-model="shiftAnswers[user.id][shift.id].value"
                  class="w-full text-[10px] scale-90"
                  :disabled="roster.saved || user.gewis_id != getGEWISId()"
                  :options="rosterValues"
                  placeholder="Answer"
                  @update:model-value="(value) => onAnswerChange(user.id, shift.id, value)"
                />
              </td>
            </template>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="flex flex-col gap-2 justify-center">
      <Button
        class="text-xs sm:text-sm lg:text-base px-2 sm:px-4 py-1 sm:py-2"
        :disabled="roster.saved"
        label="Add Shift"
        @click="visible = true"
      >
        Add Shift
      </Button>
      <Button
        class="text-xs sm:text-sm lg:text-base px-2 sm:px-4 py-1 sm:py-2"
        :disabled="roster.saved"
        label="Save Roster"
        @click="saveRoster"
      >
        Save Roster
      </Button>
    </div>
    <Dialog v-model:visible="visible" class="p-4" :style="{ width: '25rem' }">
      <template #header>
        <div class="w-full text-center">
          <span class="font- text-lg">Add Shift</span>
        </div>
      </template>
      <div class="flex flex-col gap-4 p-4">
        <div class="flex row gap-3">
          <label class="font-semibold">Shift Name</label>
          <InputText v-model="shiftName" autocomplete="off" class="w-full" type="text" />
        </div>
        <div class="flex justify-center gap-3 pt-4">
          <Button class="px-6" label="Cancel" severity="secondary" @click="visible = false" />
          <Button class="px-6" label="Add" @click="addShift(shiftName)" />
        </div>
      </div>
    </Dialog>
  </div>
</template>

<style scoped>
th,
td {
  border: 1px solid #ddd;
  text-align: center;
}

thead {
  background-color: #f2f2f2;
}

.shift-column {
  min-width: 120px;
  text-align: center;
  font-weight: bold;
}

.user-cell {
  font-weight: 500;
}

.shift-cell {
  text-align: center;
}
</style>
