<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { OiPlus } from '@openvue/openicons-vue';
import type { ShiftGroup } from '@gewis/grooster-backend-ts';
import PageHeader from '@/components/common/PageHeader.vue';
import RosterTemplate from '@/components/templates/RosterTemplate.vue';
import RosterTemplateAdd from '@/components/templates/dialogs/RosterTemplateAdd.vue';
import ShiftGroupsSection from '@/components/templates/ShiftGroupsSection.vue';
import NewRosterDialog from '@/components/roster/dialogs/NewRosterDialog.vue';
import { useTemplateStore } from '@/stores/template.store';
import { Role, useAuthStore } from '@/stores/auth.store';

const route = useRoute();
const router = useRouter();
const templateStore = useTemplateStore();
const authStore = useAuthStore();

const organId = computed(() => parseInt(route.params.id as string));
const isAdmin = computed(() => authStore.can([Role.Admin, Role.Owner]));
const templates = computed(() => Object.values(templateStore.templates).filter((t) => t.organId === organId.value));
const shiftGroups = ref<ShiftGroup[]>([]);

// Only one dialog can be open at a time
type DialogType = 'add-template' | 'new-roster';
const activeDialog = ref<DialogType | null>(null);
const closeDialog = () => (activeDialog.value = null);
const useTemplateId = ref<number | null>(null);

function useTemplate(templateId: number) {
    useTemplateId.value = templateId;
    activeDialog.value = 'new-roster';
}

onMounted(() => templateStore.fetchTemplates(organId.value));
</script>

<template>
    <div class="flex flex-col gap-7 px-4 py-6 md:px-10 md:py-8">
        <PageHeader description="Reusable shift sequences. Start a roster from one in a click." title="Templates">
            <template v-if="isAdmin" #actions>
                <Button class="h-10 text-sm" label="New template" @click="activeDialog = 'add-template'">
                    <template #icon><OiPlus aria-hidden="true" :size="16" /></template>
                </Button>
            </template>
        </PageHeader>

        <ul class="grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-3">
            <li v-for="template in templates" :key="template.id" class="flex">
                <RosterTemplate
                    class="grow"
                    :is-admin="isAdmin"
                    :shift-groups="shiftGroups"
                    :template="template"
                    @use="useTemplate(template.id)"
                />
            </li>
            <li v-if="isAdmin" class="flex">
                <button
                    class="flex min-h-[180px] grow flex-col items-center justify-center gap-2.5 rounded-card border-[1.5px] border-dashed border-surface-400 text-primary hover:bg-surface-0 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:border-surface-500 dark:hover:bg-surface-900"
                    type="button"
                    @click="activeDialog = 'add-template'"
                >
                    <span class="flex size-11 items-center justify-center rounded-full bg-highlight">
                        <OiPlus aria-hidden="true" :size="20" />
                    </span>
                    <span class="text-[15px] font-semibold">New template</span>
                    <span class="text-[13px] text-surface-600 dark:text-surface-200">Define a sequence of shifts</span>
                </button>
            </li>
            <li
                v-if="!isAdmin && templates.length === 0"
                class="rounded-card border border-dashed border-surface-400 px-6 py-10 text-center text-sm text-surface-600 dark:border-surface-500 dark:text-surface-200"
            >
                No templates yet.
            </li>
        </ul>

        <ShiftGroupsSection :is-admin="isAdmin" :organ-id="organId" @update:groups="shiftGroups = $event" />
    </div>

    <RosterTemplateAdd :organ-id="organId" :visible="activeDialog === 'add-template'" @update:visible="closeDialog" />
    <NewRosterDialog
        :initial-template-id="useTemplateId"
        :organ-id="organId"
        :visible="activeDialog === 'new-roster'"
        @created="router.push(`/rosters/${organId}`)"
        @update:visible="closeDialog"
    />
</template>
