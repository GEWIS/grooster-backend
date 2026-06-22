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
  { open: false },
);
const emit = defineEmits(['close']);

const rosterName = ref(props.name);
const date = ref(new Date());
const selectedShifts = ref<string[]>([...props.shifts]);

const visible = computed({
  get: () => props.open,
  set: () => emit('close'),
});

watch(
  () => props.name,
  (val) => (rosterName.value = val),
);
watch(
  () => props.shifts,
  (val) => (selectedShifts.value = [...val]),
);

const toggleShift = (shift: string) => {
  const idx = selectedShifts.value.indexOf(shift);
  if (idx === -1) {
    selectedShifts.value.push(shift);
  } else {
    selectedShifts.value.splice(idx, 1);
  }
};

const addRoster = async () => {
  const d = date.value;
  const dateString = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}T00:00:00Z`;

  const created = await rosterStore.createRoster({
    name: rosterName.value,
    date: dateString,
    organId: parseInt(route.params.id as string),
    shifts: selectedShifts.value,
    templateId: props.templateId,
  } satisfies RosterCreateRequest);

  if (created?.id != null) {
    rosterStore.setSelectedRoster(created.id);
  }
  rosterName.value = '';
  date.value = new Date();
  emit('close');
};
</script>

<template>
  <Dialog v-model:visible="visible" modal :style="{ width: '26rem' }" @hide="emit('close')">
    <template #header>
      <div class="w-full text-center">
        <span class="font-bold text-lg">Add Roster</span>
      </div>
    </template>

    <div class="flex flex-col gap-5">
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-semibold" for="rosterName">Roster Name</label>
        <InputText id="rosterName" v-model="rosterName" class="w-full" />
      </div>

      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-semibold">Roster Date</label>
        <DatePicker v-model="date" class="w-full" :min-date="new Date()" />
      </div>

      <div v-if="props.shifts.length > 0" class="flex flex-col gap-2">
        <div class="flex items-center justify-between">
          <label class="text-sm font-semibold">Shifts</label>
          <span class="text-xs text-gray-400">{{ selectedShifts.length }} / {{ props.shifts.length }} selected</span>
        </div>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="shift in props.shifts"
            :key="shift"
            class="px-3 py-1 rounded-full text-sm font-medium border transition-all duration-150 cursor-pointer select-none"
            :class="
              selectedShifts.includes(shift)
                ? 'bg-[--p-primary-color] border-[--p-primary-color] text-[--p-primary-contrast-color]'
                : 'bg-transparent border-gray-300 text-gray-400 hover:border-[--p-primary-color] hover:text-[--p-primary-color]'
            "
            type="button"
            @click="toggleShift(shift)"
          >
            {{ shift }}
          </button>
        </div>
      </div>
    </div>

    <Divider />

    <div class="flex justify-end gap-2">
      <Button label="Cancel" severity="secondary" @click="emit('close')" />
      <Button :disabled="props.shifts.length > 0 && selectedShifts.length === 0" label="Add" @click="addRoster()" />
    </div>
  </Dialog>
</template>

<style scoped></style>
