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
  <Dialog v-model:visible="visible" modal :style="{ width: '25rem' }" @hide="emit('close')">
    <template #header>
      <div class="w-full text-center">
        <span class="font- text-lg">Add Roster Template</span>
      </div>
    </template>
    <div>
      <Form v-slot="$form" :initial-values :resolver @submit="onFormSubmit">
        <div class="flex flex-row gap-1">
          <InputText
            ref="shiftNameInput"
            name="shiftName"
            placeholder="Shift name"
            :pt="{
              root: {
                ref: shiftNameInput,
              },
            }"
            type="text"
          />
          <Message v-if="$form.shiftName?.invalid" severity="error" size="small" variant="simple">{{
            $form.shiftName.error?.message
          }}</Message>
          <Button label="Add Shift" severity="secondary" type="submit" />
        </div>
      </Form>
    </div>
    <div v-for="(shift, index) in shifts" :key="index">
      {{ shift }}
      <Button icon="pi pi-times" style="font-size: 1rem" @click="removeShift(index)" />
    </div>

    <div class="flex flex-col justify-center gap-1">
      <InputText v-model="templateName" placeholder="Template name" type="text" />
      <Button label="Save Template" severity="success" @click="saveTemplate" />
    </div>
  </Dialog>
</template>

<style scoped></style>
