<script setup lang="ts">
import { computed } from 'vue';
import { useTemplateStore } from '@/stores/template.store';

const templateStore = useTemplateStore();

const props = withDefaults(
  defineProps<{
    open: boolean;
    templateId: number;
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

const deleteTemplate = async () => {
  await templateStore.deleteTemplate(props.templateId);

  emit('close');
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '25rem' }" @hide="emit('close')">
    <template #header>
      <div class="w-full text-center">
        <span class="font- text-lg">Are you sure you want to delete this template</span>
      </div>
    </template>
    <div class="flex flex-row gap-2 justify-center">
      <Button label="Cancel" @click="emit('close')" />
      <Button label="Delete" @click="deleteTemplate" />
    </div>
  </Dialog>
</template>

<style scoped></style>
