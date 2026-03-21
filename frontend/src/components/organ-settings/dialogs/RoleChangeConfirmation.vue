<script setup lang="ts">
import { computed } from 'vue';
import { UserRole } from '@/components/organ-settings/MemberRoles.vue';

const props = withDefaults(
  defineProps<{
    open: boolean;
    userId: number;
    role: UserRole;
  }>(),
  {
    open: false,
  },
);

const visible = computed({
  get: () => props.open,
  set: () => {
    emit('close');
  },
});

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm', payload: { userId: number; role: UserRole }): void;
}>();

const handleConfirm = () => {
  emit('confirm', {
    userId: props.userId,
    role: props.role,
  });
  // Usually, you want to close the dialog after confirming
  emit('close');
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '25rem' }" @hide="emit('close')">
    <template #header>
      <div class="w-full text-center">
        <span class="text-lg font-bold">Role Change Confirmation</span>
      </div>
    </template>
    <div class="flex flex-col gap-4 p-4">
      <p class="text-center">Are you sure you want to change this role?</p>
      <div class="flex justify-center gap-3 pt-4">
        <Button class="px-6" label="No" severity="secondary" @click="emit('close')" />
        <Button class="px-6" label="Yes" @click="handleConfirm" />
      </div>
    </div>
  </Dialog>
</template>

<style scoped></style>
