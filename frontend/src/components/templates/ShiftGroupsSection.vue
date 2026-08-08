<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ShiftGroup } from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import { useOrganStore } from '@/stores/organ.store';
import { Role, useAuthStore } from '@/stores/auth.store';
import ShiftGroupAdd from '@/components/templates/dialogs/ShiftGroupAdd.vue';
import ShiftGroupPriority from '@/components/templates/ShiftGroupPriority.vue';

const organStore = useOrganStore();
const authStore = useAuthStore();

const shiftGroups = ref<ShiftGroup[]>([]);

type Dialogs = 'AddGroup' | 'None';
const openDialog = ref<Dialogs>('None');

const emit = defineEmits(['update:groups']);

const expandedGroupId = ref<number>(null);

const fetchGroups = async () => {
    try {
        const response = await ApiService.shiftGroupApi.getShiftGroups(organStore.organ);
        const newGroups = response.data ?? [];

        shiftGroups.value = newGroups;
        emit('update:groups', newGroups);
    } catch (error) {
        console.error('Failed to fetch groups:', error);
    }
};

const toggleGroup = (groupId: number) => {
    if (groupId === expandedGroupId.value) {
        expandedGroupId.value = null;
    } else {
        expandedGroupId.value = groupId;
    }
};

const handleGroupClose = async () => {
    openDialog.value = 'None';
    await fetchGroups();
};

onMounted(async () => {
    await fetchGroups();
});
</script>

<template>
    <section class="flex flex-col gap-2">
        <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
            <i class="pi pi-users text-primary"></i> Shift Groups
            <Button
                v-if="authStore.can([Role.Admin, Role.Owner])"
                class="p-button-outlined"
                icon="pi pi-users"
                label="Add shift group"
                @click="openDialog = 'AddGroup'"
            />
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 items-start">
            <template v-for="group in shiftGroups" :key="group.id">
                <div
                    :class="[
                        'p-4 bg-white border rounded-lg shadow-sm transition-all cursor-pointer hover:border-primary',
                        expandedGroupId === group.id ? 'border-primary ring-2 ring-primary' : 'border-gray-200',
                    ]"
                    @click="toggleGroup(group.id)"
                >
                    <div class="flex items-center justify-between">
                        <p class="font-bold text-gray-800">{{ group.name }}</p>
                        <i class="pi pi-chevron-right text-gray-400 text-sm"></i>
                    </div>
                </div>
            </template>
        </div>

        <div v-if="expandedGroupId">
            <div class="fixed inset-0 bg-black/40 backdrop-blur-sm z-[60]" @click="expandedGroupId = null"></div>

            <div
                class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-[70] bg-white border border-gray-200 rounded-2xl shadow-2xl p-8 w-[90vw] max-w-[800px] min-h-[500px] overflow-y-auto"
            >
                <div class="flex justify-between items-center mb-8">
                    <div>
                        <h3 class="text-2xl font-bold text-gray-800">
                            Assigning priority for: {{ shiftGroups.find((g) => g.id === expandedGroupId)?.name }}
                        </h3>
                    </div>
                    <Button
                        class="p-button-rounded p-button-text p-button-lg text-gray-400 hover:text-red-500"
                        icon="pi pi-times"
                        @click="expandedGroupId = null"
                    />
                </div>

                <ShiftGroupPriority :group-id="expandedGroupId" />
            </div>
        </div>
    </section>
    <ShiftGroupAdd :open="openDialog === 'AddGroup'" @close="handleGroupClose" />
</template>

<style scoped></style>
