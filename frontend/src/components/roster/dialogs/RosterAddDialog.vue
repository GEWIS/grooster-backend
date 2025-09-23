<script setup lang="ts">
import { AddRosterForm } from '@/types/roster.types';
import { computed, ref } from 'vue';
import { RosterCreateRequest } from '@gewis/grooster-backend-ts';
import { useRosterStore } from '@/stores/roster.store';
import { useRoute } from 'vue-router';

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

function addRoster() {
  const createParams: RosterCreateRequest = {
    name: formValues.value.name,
    date: formValues.value.date.toISOString(),
    organId: parseInt(route.params.id as string),
  };

  rosterStore.createRoster(createParams);
  emit('close');
}
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
        <InputText type="text" v-model="formValues.name" class="w-full" autocomplete="off" />
      </div>
      <div class="flex row gap-3">
        <label class="font-semibold">Roster Date</label>
        <DatePicker v-model="formValues.date" :minDate="new Date()" />
      </div>
      <div class="flex justify-center gap-3 pt-4">
        <Button label="Cancel" severity="secondary" @click="emit('close')" class="px-6" />
        <Button label="Add" @click="addRoster()" class="px-6" />
      </div>
    </div>
  </Dialog>
</template>

<style scoped></style>
