<script setup lang="ts">
import { computed, ref } from 'vue';
import { ShiftGroupCreateRequest } from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import { useOrganStore } from '@/stores/organ.store';

const organStore = useOrganStore();

const props = withDefaults(
  defineProps<{
    open: boolean;
  }>(),
  {
    open: false,
  },
);

const emit = defineEmits(['close']);

const visible = computed({
  get: () => props.open,
  set: () => {
    emit('close');
  },
});

const groupName = ref('');

const createGroup = async () => {
  try {
    const params: ShiftGroupCreateRequest = {
      organId: Number(organStore.organ),
      name: groupName.value,
    };

    await ApiService.shiftGroupApi.createShiftGroup(params);

    groupName.value = '';

    emit('close');
  } catch (error) {
    console.error(error);
  }
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '28rem' }" @hide="emit('close')">
    <template #header>
      <div class="flex flex-col gap-1">
        <h3 class="m-0 text-xl font-semibold text-surface-900">New Shift Group</h3>
        <small class="text-surface-500">Specify the name of the shift group</small>
      </div>
    </template>
    <div class="flex flex-col gap-2 justify-center">
      <InputText v-model="groupName" class="w-full" name="groupName" type="text" />

      <Button icon="pi pi-plus" label="Add" severity="primary" @click="createGroup" />
    </div>
  </Dialog>
</template>

<style scoped></style>
