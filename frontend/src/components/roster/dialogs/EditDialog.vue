<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import type { Roster } from '@gewis/grooster-backend-ts';
import { useRosterStore } from '@/stores/roster.store';
import { parseRosterDate, toRosterDate } from '@/helpers/date';
import { dialogPt, readableDisabled } from '@/theme/pt';

// "Edit details" from the roster's overflow menu: name and date.
const props = defineProps<{ roster: Roster }>();
const visible = defineModel<boolean>('visible', { required: true });

const rosterStore = useRosterStore();
const name = ref('');
const date = ref<Date | null>(null);
const saving = ref(false);

watch(visible, (open) => {
    if (!open) return;
    name.value = props.roster.name ?? '';
    date.value = parseRosterDate(props.roster.date);
});

const canSave = computed(() => !!name.value.trim() && !!date.value && !saving.value);

async function save() {
    if (!canSave.value || !date.value) return;
    saving.value = true;
    try {
        await rosterStore.updateRoster(props.roster.id, { name: name.value.trim(), date: toRosterDate(date.value) });
        visible.value = false;
    } finally {
        saving.value = false;
    }
}

const labelClass = 'text-[13px] font-semibold text-surface-700 dark:text-surface-100';
</script>

<template>
    <Dialog v-model:visible="visible" header="Edit details" modal :pt="dialogPt">
        <form id="edit-roster-form" class="flex flex-col gap-[18px]" @submit.prevent="save">
            <div class="flex flex-col gap-1.5">
                <label :class="labelClass" for="edit-roster-name">Roster name</label>
                <InputText
                    id="edit-roster-name"
                    v-model="name"
                    autocomplete="off"
                    autofocus
                    class="h-[42px] w-full text-base md:text-sm"
                />
            </div>
            <div class="flex flex-col gap-1.5">
                <label :class="labelClass" for="edit-roster-date">Date</label>
                <DatePicker
                    v-model="date"
                    date-format="D d M yy"
                    input-id="edit-roster-date"
                    :pt="{ pcInputText: { root: { class: 'h-[42px] w-full text-base md:text-sm' } } }"
                />
            </div>
        </form>
        <template #footer>
            <Button class="h-10 text-sm" label="Cancel" outlined severity="secondary" @click="visible = false" />
            <Button
                :class="['h-10 text-sm', readableDisabled]"
                :disabled="!canSave"
                form="edit-roster-form"
                label="Save changes"
                :loading="saving"
                type="submit"
            />
        </template>
    </Dialog>
</template>
