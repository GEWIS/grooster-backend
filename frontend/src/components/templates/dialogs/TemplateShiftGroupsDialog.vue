<script setup lang="ts">
import type { RosterTemplate, ShiftGroup } from '@gewis/grooster-backend-ts';
import { useToast } from 'openvue/usetoast';
import ApiService from '@/services/ApiService';
import { dialogPt, fieldSelectPt } from '@/theme/pt';

defineProps<{
    template: RosterTemplate;
    shiftGroups: ShiftGroup[];
}>();
const visible = defineModel<boolean>('visible', { required: true });
const toast = useToast();

async function setGroup(shiftId: number, shiftGroupId: number | null) {
    try {
        await ApiService.roster.updateRosterTemplateShift(shiftId, { shiftGroupId });
    } catch (error) {
        console.error(error);
        toast.add({ severity: 'error', summary: 'Could not change the shift group', life: 5000 });
    }
}
</script>

<template>
    <Dialog v-model:visible="visible" modal :pt="dialogPt">
        <template #header>
            <div class="flex flex-col gap-1">
                <h2 class="font-display text-[22px] font-semibold tracking-[-0.01em]">Shift groups</h2>
                <p class="text-sm text-surface-600 dark:text-surface-200">
                    Who can be scheduled on each shift of {{ template.name }}.
                </p>
            </div>
        </template>
        <ul class="flex flex-col divide-y divide-surface-100 dark:divide-surface-800">
            <li v-for="shift in template.shifts" :key="shift.id" class="flex items-center gap-3 py-2">
                <label class="min-w-0 grow truncate text-sm font-semibold" :for="`shift-group-${shift.id}`">
                    {{ shift.shiftName }}
                </label>
                <Select
                    v-model="shift.shiftGroupId"
                    class="h-10 w-48 items-center"
                    :input-id="`shift-group-${shift.id}`"
                    option-label="name"
                    option-value="id"
                    :options="shiftGroups"
                    placeholder="No group"
                    :pt="fieldSelectPt"
                    show-clear
                    @update:model-value="setGroup(shift.id, $event)"
                />
            </li>
        </ul>
        <template #footer>
            <Button class="h-10 text-sm" label="Done" @click="visible = false" />
        </template>
    </Dialog>
</template>
