<script setup lang="ts">
import { computed, ref } from 'vue';
import { RosterCreateRequest } from '@gewis/grooster-backend-ts';
import { useRoute } from 'vue-router';
import { useRosterStore } from '@/stores/roster.store';
import { AddRosterForm } from '@/types/roster.types';

const route = useRoute();

const props = withDefaults(
  defineProps<{
    open: boolean;
  }>(),
  {
    open: false,
  },
);

const emit = defineEmits(['close']);

const rosterStore = useRosterStore();
const visible = computed({
  get: () => props.open,
  set: () => {
    emit('close');
  },
});

const formValues = ref<AddRosterForm>({
  name: '',
  date: new Date(),
});

const addRoster = async () => {
  const d = formValues.value.date;
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const dateString = `${year}-${month}-${day}T00:00:00Z`;

  const createParams: RosterCreateRequest = {
    name: formValues.value.name,
    date: dateString,
    organId: parseInt(route.params.id as string),
  };

  const created = await rosterStore.createRoster(createParams);
  if (created?.id != null) {
    rosterStore.setSelectedRoster(created.id);
  }

  formValues.value = { name: '', date: new Date() };
  emit('close');
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '25rem' }" @hide="emit('close')">
    <template #header>
      <div class="w-full text-center">
        <span class="font- text-lg">Add Roster</span>
      </div>
    </template>
    <div class="flex flex-col gap-4 p-4">
      <div class="flex row gap-3">
        <label class="font-semibold">Roster Name</label>
        <InputText v-model="formValues.name" autocomplete="off" class="w-full" type="text" />
      </div>
      <div class="flex row gap-3">
        <label class="font-semibold">Roster Date</label>
        <DatePicker v-model="formValues.date" :min-date="new Date()" />
      </div>
      <div class="flex justify-center gap-3 pt-4">
        <Button class="px-6" label="Cancel" severity="secondary" @click="emit('close')" />
        <Button class="px-6" label="Add" @click="addRoster()" />
      </div>
    </div>
  </Dialog>
</template>

<style scoped></style>
