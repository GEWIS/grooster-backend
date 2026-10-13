<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useToast } from 'openvue/usetoast';
import ToggleButton from 'openvue/togglebutton';
import { useRosterStore } from '@/stores/roster.store';
import { useTemplateStore } from '@/stores/template.store';
import { toRosterDate } from '@/helpers/date';
import { dialogPt, fieldSelectPt, readableDisabled } from '@/theme/pt';

const props = withDefaults(
    defineProps<{
        organId: number;
        /** Pre-selects a template ("Use template" on the Templates page) */
        initialTemplateId?: number | null;
    }>(),
    { initialTemplateId: null },
);
const visible = defineModel<boolean>('visible', { required: true });
const emit = defineEmits<{ created: [rosterId: number] }>();

const rosterStore = useRosterStore();
const templateStore = useTemplateStore();
const toast = useToast();

const name = ref('');
const date = ref<Date>(new Date());
const templateId = ref<number | null>(null);
const selectedShifts = ref<string[]>([]);
const saving = ref(false);

const templates = computed(() => Object.values(templateStore.templates));
const template = computed(() => (templateId.value ? templateStore.getTemplate(templateId.value) : undefined));
const templateShifts = computed(() => template.value?.shifts?.map((shift) => shift.shiftName) ?? []);
const allSelected = computed(() => selectedShifts.value.length === templateShifts.value.length);

watch(visible, (open) => {
    if (!open) return;
    if (!templates.value.length) void templateStore.fetchTemplates(props.organId);
    templateId.value = props.initialTemplateId;
    name.value = template.value?.name ?? '';
    date.value = new Date();
});
watch(templateShifts, (shifts) => (selectedShifts.value = [...shifts]), { immediate: true });

function toggleShift(shift: string, on: boolean) {
    selectedShifts.value = on
        ? templateShifts.value.filter((s) => s === shift || selectedShifts.value.includes(s))
        : selectedShifts.value.filter((s) => s !== shift);
}

const canCreate = computed(
    () => !!name.value.trim() && (!template.value || selectedShifts.value.length > 0) && !saving.value,
);

async function create() {
    if (!canCreate.value) return;
    saving.value = true;
    try {
        const created = await rosterStore.createRoster({
            name: name.value.trim(),
            date: toRosterDate(date.value),
            organId: props.organId,
            ...(template.value ? { shifts: selectedShifts.value, templateId: template.value.id } : {}),
        });
        if (created?.id != null) {
            rosterStore.setSelectedRoster(created.id);
            emit('created', created.id);
        }
        visible.value = false;
    } catch (error) {
        console.error('Failed to create roster:', error);
        toast.add({ severity: 'error', summary: 'Could not create the roster', life: 5000 });
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
                <h2 class="font-display text-[22px] font-semibold tracking-[-0.01em]">New roster</h2>
                <p class="text-sm text-surface-600 dark:text-surface-200">
                    Start from a template, then collect preferences.
                </p>
            </div>
        </template>

        <form id="new-roster-form" class="flex flex-col gap-[18px]" @submit.prevent="create">
            <div class="flex flex-col gap-1.5">
                <label :class="labelClass" for="new-roster-name">Roster name</label>
                <InputText
                    id="new-roster-name"
                    v-model="name"
                    autocomplete="off"
                    autofocus
                    class="h-[42px] w-full text-base md:text-sm"
                />
            </div>
            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div class="flex flex-col gap-1.5">
                    <label :class="labelClass" for="new-roster-date">Date</label>
                    <DatePicker
                        v-model="date"
                        date-format="D d M yy"
                        input-id="new-roster-date"
                        :min-date="new Date()"
                        :pt="{ pcInputText: { root: { class: 'h-[42px] w-full text-base md:text-sm' } } }"
                    />
                </div>
                <div class="flex flex-col gap-1.5">
                    <label :class="labelClass" for="new-roster-template">Template</label>
                    <Select
                        v-model="templateId"
                        class="h-[42px] w-full items-center"
                        input-id="new-roster-template"
                        :option-label="
                            (t) => `${t.name} · ${t.shifts?.length ?? 0} ${t.shifts?.length === 1 ? 'shift' : 'shifts'}`
                        "
                        option-value="id"
                        :options="templates"
                        placeholder="No template"
                        :pt="fieldSelectPt"
                        show-clear
                    />
                </div>
            </div>

            <fieldset v-if="template && templateShifts.length" class="flex flex-col gap-1.5">
                <div class="flex items-baseline">
                    <legend :class="[labelClass, 'grow']">Shifts</legend>
                    <span class="text-xs text-surface-600 dark:text-surface-200">
                        {{ selectedShifts.length }} of {{ templateShifts.length }} selected ·
                        <button
                            class="font-semibold text-primary underline-offset-2 hover:underline focus-visible:outline-2 focus-visible:outline-primary"
                            type="button"
                            @click="selectedShifts = allSelected ? [] : [...templateShifts]"
                        >
                            {{ allSelected ? 'Clear all' : 'Select all' }}
                        </button>
                    </span>
                </div>
                <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
                    <ToggleButton
                        v-for="shift in templateShifts"
                        :key="shift"
                        :model-value="selectedShifts.includes(shift)"
                        :off-label="shift"
                        :on-label="shift"
                        :pt="{
                            root: ({ props: p }) => ({
                                class: [
                                    'h-9 min-w-0 rounded-full border px-3 text-[13px] font-semibold',
                                    p.modelValue
                                        ? 'border-primary bg-highlight'
                                        : 'border-surface-400 bg-surface-0 text-surface-700 dark:border-surface-500 dark:bg-surface-900 dark:text-surface-100',
                                ],
                            }),
                            content: { class: 'min-w-0 justify-start bg-transparent p-0 shadow-none' },
                            label: { class: 'truncate' },
                        }"
                        @update:model-value="toggleShift(shift, $event)"
                    />
                </div>
            </fieldset>
        </form>

        <template #footer>
            <Button class="h-10 text-sm" label="Cancel" outlined severity="secondary" @click="visible = false" />
            <Button
                :class="['h-10 text-sm', readableDisabled]"
                :disabled="!canCreate"
                form="new-roster-form"
                label="Create roster"
                :loading="saving"
                type="submit"
            />
        </template>
    </Dialog>
</template>
