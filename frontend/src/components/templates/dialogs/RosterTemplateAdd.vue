<script setup lang="ts">
import { computed, nextTick, reactive, ref } from 'vue';
import { RosterTemplateCreateRequest } from '@gewis/grooster-backend-ts';
import { useRoute } from 'vue-router';
import { useTemplateStore } from '@/stores/template.store';

const route = useRoute();
const templateStore = useTemplateStore();

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

const initialValues = reactive({
  shiftName: '',
});
const shifts = ref<string[]>([]);
const templateName = ref<string>('');
const shiftNameInput = ref<HTMLInputElement | null>(null);

const resolver = ({ values }) => {
  const errors = { shiftName: [] };

  if (!values.shiftName) {
    errors.shiftName.push({ message: 'A shift name is required.' });
  }

  return {
    values,
    errors,
  };
};

const onFormSubmit = async (e) => {
  if (e.valid) {
    shifts.value.push(e.values.shiftName);
    e.reset();

    await nextTick();
    shiftNameInput.value?.focus();
  }
};

const removeShift = (index: number) => {
  shifts.value.splice(index, 1);
};

const saveTemplate = async () => {
  const params: RosterTemplateCreateRequest = {
    name: templateName.value,
    organId: parseInt(route.params.id as string),
    shifts: shifts.value,
  };
  await templateStore.createTemplate(params);

  shifts.value = [];
  emit('close');
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '28rem' }" @hide="emit('close')">
    <template #header>
      <div class="flex flex-col gap-1">
        <h3 class="m-0 text-xl font-semibold text-surface-900">New Roster Template</h3>
        <small class="text-surface-500">Define the sequence of shifts for this rotation.</small>
      </div>
    </template>

    <div class="flex flex-col gap-6 pt-2">
      <div class="flex flex-col gap-2">
        <label class="font-bold text-sm text-surface-700" for="templateName">Template Name</label>
        <InputText id="templateName" v-model="templateName" class="w-full" placeholder="e.g. Weekly Operations" />
      </div>

      <Divider class="m-0" />

      <div class="flex flex-col gap-2">
        <label class="font-bold text-sm text-surface-700">Add Shift to Sequence</label>
        <Form v-slot="$form" class="flex gap-2" :initial-values :resolver @submit="onFormSubmit">
          <div class="flex-1">
            <InputText
              ref="shiftNameInput"
              class="w-full"
              name="shiftName"
              placeholder="e.g. Morning Shift"
              type="text"
            />
            <Message v-if="$form.shiftName?.invalid" class="mt-1" severity="error" size="small" variant="simple">
              {{ $form.shiftName.error?.message }}
            </Message>
          </div>
          <Button icon="pi pi-plus" label="Add" severity="primary" type="submit" />
        </Form>
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-bold text-sm text-surface-700">Shift Sequence</label>
        <div
          v-if="shifts.length > 0"
          class="grid grid-cols-2 gap-3 bg-surface-50 dark:bg-surface-900 p-3 border-round border-1 border-surface-200 overflow-y-auto"
          style="max-height: 320px"
        >
          <div
            v-for="(shift, index) in shifts"
            :key="index"
            class="flex flex-col bg-white dark:bg-surface-800 border-round shadow-sm border-top-3 border-primary p-3 relative"
          >
            <div class="flex justify-between items-center mb-2">
              <span
                class="flex align-items-center justify-content-center bg-primary-100 text-primary-700 font-bold border-circle text-xs"
                style="min-width: 22px; height: 22px"
              >
                {{ index + 1 }}
              </span>
              <Button
                class="w-2rem h-2rem"
                icon="pi pi-trash"
                rounded
                severity="danger"
                size="small"
                text
                @click="removeShift(index)"
              />
            </div>

            <div class="flex flex-col gap-1">
              <span class="font-medium text-sm truncate" :title="shift">{{ shift }}</span>
            </div>
          </div>

          <div class="col-span-2 text-center mt-2">
            <small class="text-surface-400">Shifts follow the numbered sequence</small>
          </div>
        </div>

        <div v-else class="text-center p-6 border-dashed border-2 border-surface-300 border-round">
          <i class="pi pi-list text-3xl text-surface-300 mb-2"></i>
          <p class="m-0 text-surface-400">Add your first shift to begin the sequence.</p>
        </div>
      </div>

      <Button
        class="w-full py-3 mt-2"
        icon="pi pi-save"
        label="Save Roster Template"
        severity="success"
        @click="saveTemplate"
      />
    </div>
  </Dialog>
</template>

<style scoped></style>
