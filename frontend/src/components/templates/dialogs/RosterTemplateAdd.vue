<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { OiPlus, OiTimes } from '@openvue/openicons-vue';
import { useTemplateStore } from '@/stores/template.store';
import { dialogPt, readableDisabled } from '@/theme/pt';

const props = defineProps<{ organId: number }>();
const visible = defineModel<boolean>('visible', { required: true });

const templateStore = useTemplateStore();
const name = ref('');
const shiftName = ref('');
const shifts = ref<string[]>([]);
const saving = ref(false);

watch(visible, (open) => {
    if (!open) return;
    name.value = '';
    shiftName.value = '';
    shifts.value = [];
});

function addShift() {
    const trimmed = shiftName.value.trim();
    if (!trimmed) return;
    shifts.value.push(trimmed);
    shiftName.value = '';
}

const canSave = computed(() => !!name.value.trim() && shifts.value.length > 0 && !saving.value);

async function save() {
    if (!canSave.value) return;
    saving.value = true;
    try {
        await templateStore.createTemplate({ name: name.value.trim(), organId: props.organId, shifts: shifts.value });
        visible.value = false;
    } finally {
        saving.value = false;
    }
}

const labelClass = 'text-[13px] font-semibold text-surface-700 dark:text-surface-100';
</script>

<template>
    <Dialog v-model:visible="visible" modal :pt="dialogPt">
        <template #header>
            <div class="flex flex-col gap-1">
                <h2 class="font-display text-[22px] font-semibold tracking-[-0.01em]">New template</h2>
                <p class="text-sm text-surface-600 dark:text-surface-200">Define the sequence of shifts.</p>
            </div>
        </template>

        <div class="flex flex-col gap-1.5">
            <label :class="labelClass" for="template-name">Template name</label>
            <InputText
                id="template-name"
                v-model="name"
                autocomplete="off"
                autofocus
                class="h-[42px] w-full text-base md:text-sm"
            />
        </div>

        <form class="flex flex-col gap-1.5" @submit.prevent="addShift">
            <label :class="labelClass" for="template-shift">Add a shift</label>
            <div class="flex gap-2">
                <InputText
                    id="template-shift"
                    v-model="shiftName"
                    autocomplete="off"
                    class="h-[42px] grow text-base md:text-sm"
                    placeholder="E.g. Bar"
                />
                <Button
                    class="h-[42px] bg-surface-0 text-sm dark:bg-surface-900"
                    :disabled="!shiftName.trim()"
                    label="Add"
                    outlined
                    severity="secondary"
                    type="submit"
                >
                    <template #icon><OiPlus aria-hidden="true" :size="16" /></template>
                </Button>
            </div>
        </form>

        <div class="flex flex-col gap-1.5">
            <span :class="labelClass">Shift sequence</span>
            <ol v-if="shifts.length" class="flex flex-wrap gap-1.5">
                <li
                    v-for="(shift, index) in shifts"
                    :key="`${index}-${shift}`"
                    class="inline-flex h-8 items-center gap-1.5 rounded-md bg-surface-100 pl-2.5 pr-1 text-[13px] font-semibold text-surface-700 dark:bg-surface-800 dark:text-surface-100"
                >
                    <b class="font-bold text-primary">{{ index + 1 }}</b
                    >{{ shift }}
                    <button
                        :aria-label="`Remove ${shift}`"
                        class="flex size-6 items-center justify-center rounded hover:bg-surface-200 focus-visible:outline-2 focus-visible:outline-primary dark:hover:bg-surface-700"
                        type="button"
                        @click="shifts.splice(index, 1)"
                    >
                        <OiTimes aria-hidden="true" :size="16" />
                    </button>
                </li>
            </ol>
            <p v-else class="text-sm text-surface-600 dark:text-surface-200">
                Add your first shift to start the sequence.
            </p>
        </div>

        <template #footer>
            <Button class="h-10 text-sm" label="Cancel" outlined severity="secondary" @click="visible = false" />
            <Button
                :class="['h-10 text-sm', readableDisabled]"
                :disabled="!canSave"
                label="Create template"
                :loading="saving"
                @click="save"
            />
        </template>
    </Dialog>
</template>
