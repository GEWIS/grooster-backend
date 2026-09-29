<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import {
    RosterTemplateShiftPreference,
    RosterTemplate,
    TemplateShiftPreferenceCreateRequest,
    TemplateShiftPreferenceUpdateRequest,
} from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import { useUserStore } from '@/stores/user.store';

const props = withDefaults(
    defineProps<{
        open: boolean;
        template: RosterTemplate;
    }>(),
    {
        open: false,
    },
);

const emit = defineEmits(['close']);

const userStore = useUserStore();
const shiftPreferences = ref<Record<number, RosterTemplateShiftPreference>>({});
const options = ['J', 'N', 'X', 'L'];

onMounted(async () => {
    try {
        const { data: preferences } = await ApiService.roster.getRosterTemplateShiftPreferences(
            userStore.user.id,
            props.template.id,
        );

        const prefMap: Record<number, RosterTemplateShiftPreference> = {};

        props.template.shifts.forEach((shift) => {
            const existingPref = preferences.find((p) => p.rosterTemplateShiftID === shift.id);

            if (existingPref) {
                prefMap[shift.id] = existingPref;
            }
        });

        shiftPreferences.value = prefMap;
    } catch (e) {
        console.error(e);
    }
});

const visible = computed({
    get: () => props.open,
    set: () => {
        emit('close');
    },
});

const displayShifts = computed(() => {
    return props.template.shifts.map((shift) => {
        const pref = Object.values(shiftPreferences.value).find((p) => p.rosterTemplateShiftID === shift.id) || {
            id: null,
            rosterTemplateShiftID: shift.id,
            value: null,
        };

        return {
            ...pref,
            displayName: shift.shiftName,
        };
    });
});

const onPreferenceChange = async (shiftId: number, templateShiftId: number, value: string) => {
    if (!shiftId) {
        try {
            const params: TemplateShiftPreferenceCreateRequest = {
                preference: value,
                rosterTemplateShiftID: templateShiftId,
                userId: userStore.user.id,
            };

            const { data } = await ApiService.roster.createRosterTemplateShiftPreference(params);

            shiftPreferences.value[templateShiftId] = data;
        } catch (e) {
            console.error(e);
        }
    } else {
        try {
            const params: TemplateShiftPreferenceUpdateRequest = {
                preference: value,
            };

            const { data } = await ApiService.roster.updateRosterTemplateShiftPreference(shiftId, params);

            shiftPreferences.value[templateShiftId] = data;
        } catch (e) {
            console.error(e);
        }
    }
};
</script>

<template>
    <Dialog v-model:visible="visible" modal :style="{ width: '28rem' }" @hide="emit('close')">
        <template #header>
            <div class="flex flex-col gap-1">
                <h3 class="m-0 text-xl font-semibold text-color">Roster Template Preferences</h3>
                <small class="text-muted-color">Give your preferences for each shift in this template.</small>
            </div>
        </template>
        <div
            class="max-w-2xl mx-auto bg-surface-0 dark:bg-surface-900 rounded-xl shadow-sm border border-surface overflow-hidden"
        >
            <div class="px-6 py-4 border-b border-surface bg-surface-50/50 dark:bg-surface-800/50">
                <h3 class="text-sm font-semibold text-muted-color uppercase tracking-wider">Shift Preferences</h3>
            </div>

            <div class="divide-y divide-surface-100 dark:divide-surface-700">
                <div
                    v-for="item in displayShifts"
                    :key="item.rosterTemplateShiftID"
                    class="grid grid-cols-2 items-center px-6 py-3 hover:bg-emphasis transition-colors duration-200"
                >
                    <div class="flex items-center gap-3">
                        <span class="text-sm font-medium text-muted-color">{{ item.displayName }}</span>
                    </div>

                    <div class="flex justify-end">
                        <Select
                            v-model="item.value"
                            class="w-full max-w-[200px]"
                            :options="options"
                            placeholder="Select preference..."
                            @update:model-value="
                                (value) => onPreferenceChange(item.id, item.rosterTemplateShiftID, value)
                            "
                        />
                    </div>
                </div>
            </div>
        </div>
    </Dialog>
</template>

<style scoped></style>
