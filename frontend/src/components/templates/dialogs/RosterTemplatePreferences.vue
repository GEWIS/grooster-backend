<script setup lang="ts">
import { ref, watch } from 'vue';
import type { RosterTemplate, RosterTemplateShiftPreference } from '@gewis/grooster-backend-ts';
import { useToast } from 'openvue/usetoast';
import PreferenceSelect from '@/components/common/PreferenceSelect.vue';
import ApiService from '@/services/ApiService';
import { useUserStore } from '@/stores/user.store';
import { dialogPt } from '@/theme/pt';

const props = defineProps<{ template: RosterTemplate }>();
const visible = defineModel<boolean>('visible', { required: true });

const OPTIONS = ['J', 'N', 'X', 'L'];
const userStore = useUserStore();
const toast = useToast();
const preferences = ref<Record<number, RosterTemplateShiftPreference>>({});

watch(visible, async (open) => {
    if (!open || !userStore.user) return;
    try {
        const { data } = await ApiService.roster.getRosterTemplateShiftPreferences(
            userStore.user.id,
            props.template.id,
        );
        preferences.value = Object.fromEntries(data.map((pref) => [pref.rosterTemplateShiftID, pref]));
    } catch (error) {
        console.error(error);
    }
});

async function setPreference(templateShiftId: number, value: string) {
    const existing = preferences.value[templateShiftId];
    try {
        const { data } = existing?.id
            ? await ApiService.roster.updateRosterTemplateShiftPreference(existing.id, { preference: value })
            : await ApiService.roster.createRosterTemplateShiftPreference({
                  preference: value,
                  rosterTemplateShiftID: templateShiftId,
                  userId: userStore.user.id,
              });
        preferences.value = { ...preferences.value, [templateShiftId]: data };
    } catch (error) {
        console.error(error);
        toast.add({ severity: 'error', summary: 'Could not save your preference', life: 5000 });
    }
}
</script>

<template>
    <Dialog v-model:visible="visible" modal :pt="dialogPt">
        <template #header>
            <div class="flex flex-col gap-1">
                <h2 class="font-display text-[22px] font-semibold tracking-[-0.01em]">My preferences</h2>
                <p class="text-sm text-surface-600 dark:text-surface-200">
                    Your usual answer for each shift in {{ template.name }}. Saved automatically.
                </p>
            </div>
        </template>
        <ul class="flex flex-col divide-y divide-surface-100 dark:divide-surface-800">
            <li v-for="shift in template.shifts" :key="shift.id" class="flex min-h-[60px] items-center gap-3 py-2">
                <span class="min-w-0 grow truncate text-base font-semibold">{{ shift.shiftName }}</span>
                <PreferenceSelect
                    :label="`${shift.shiftName}: your usual answer`"
                    :model-value="preferences[shift.id]?.value ?? null"
                    :options="OPTIONS"
                    size="row"
                    @update:model-value="setPreference(shift.id, $event)"
                />
            </li>
        </ul>
        <template #footer>
            <Button class="h-10 text-sm" label="Done" @click="visible = false" />
        </template>
    </Dialog>
</template>
