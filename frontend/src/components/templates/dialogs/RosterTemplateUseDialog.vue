<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { RosterCreateRequest } from '@gewis/grooster-backend-ts';
import { useRoute } from 'vue-router';
import { useRosterStore } from '@/stores/roster.store';

const route = useRoute();
const rosterStore = useRosterStore();

const props = withDefaults(
  defineProps<{
    open: boolean;
    name: string;
    shifts: string[];
    templateId: number;
  }>(),
  {
    open: false,
  },
);
const emit = defineEmits(['close']);

const rosterName = ref(props.name);

watch(
  () => props.name,
  (newVal) => {
    rosterName.value = newVal;
  },
);

const visible = computed({
  get: () => props.open,
  set: () => {
    emit('close');
  },
});

const date = ref(new Date());

const addRoster = async () => {
  const d = date.value;
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const dateString = `${year}-${month}-${day}T00:00:00Z`;

  const createParams: RosterCreateRequest = {
    name: rosterName.value,
    date: dateString,
    organId: parseInt(route.params.id as string),
    shifts: props.shifts,
    templateId: props.templateId,
  };

  const created = await rosterStore.createRoster(createParams);
  if (created?.id != null) {
    rosterStore.setSelectedRoster(created.id);
  }
  rosterName.value = '';
  date.value = new Date();
  emit('close');
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '25rem' }" @hide="emit('close')">
    <template #header>
      <div class="w-full text-center">
        <span class="font-bold text-lg">Add Roster</span>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <label class="font-semibold" for="rosterName">Roster Name</label>
        <InputText id="rosterName" v-model="rosterName" class="w-full" />
      </div>

      <div class="flex flex-col gap-2">
        <label class="font-semibold">Roster Date</label>
        <DatePicker v-model="date" class="w-full" :min-date="new Date()" />
      </div>
    </div>

    <div class="flex justify-center gap-3 pt-6">
      <Button class="px-6" label="Cancel" severity="secondary" @click="emit('close')" />
      <Button class="px-6" label="Add" @click="addRoster()" />
    </div>
  </Dialog>
</template>

<style scoped></style>
