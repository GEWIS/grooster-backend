<script setup lang="ts">
import { computed, ref } from 'vue';
import Menu from 'openvue/menu';
import { OiEllipsisH, OiHeart, OiTrash, OiUsers } from '@openvue/openicons-vue';
import type { RosterTemplate as RosterTemplateModel, ShiftGroup } from '@gewis/grooster-backend-ts';
import RosterTemplatePreferences from '@/components/templates/dialogs/RosterTemplatePreferences.vue';
import TemplateShiftGroupsDialog from '@/components/templates/dialogs/TemplateShiftGroupsDialog.vue';
import { useTemplateStore } from '@/stores/template.store';
import { useConfirmDelete } from '@/composables/useConfirmDelete';

const props = defineProps<{
    template: RosterTemplateModel;
    shiftGroups: ShiftGroup[];
    isAdmin: boolean;
}>();
const emit = defineEmits<{ use: [] }>();

const templateStore = useTemplateStore();
const confirmDelete = useConfirmDelete();

// Only one dialog can be open at a time
type DialogType = 'preferences' | 'shift-groups';
const activeDialog = ref<DialogType | null>(null);
const closeDialog = () => (activeDialog.value = null);
const menu = ref<InstanceType<typeof Menu>>();

const shifts = computed(() => props.template.shifts ?? []);

const menuItems = computed(() => [
    { label: 'My preferences', iconComponent: OiHeart, command: () => (activeDialog.value = 'preferences') },
    ...(props.isAdmin
        ? [
              { label: 'Shift groups', iconComponent: OiUsers, command: () => (activeDialog.value = 'shift-groups') },
              { separator: true },
              {
                  label: 'Delete template',
                  iconComponent: OiTrash,
                  danger: true,
                  command: () =>
                      confirmDelete({
                          header: 'Delete template?',
                          message: `"${props.template.name}" is deleted. Rosters made from it stay.`,
                          acceptLabel: 'Delete template',
                          onAccept: () => templateStore.deleteTemplate(props.template.id),
                      }),
              },
          ]
        : []),
]);
</script>

<template>
    <article
        :aria-labelledby="`template-${template.id}`"
        class="flex flex-col gap-4 rounded-card border border-surface bg-surface-0 p-5 dark:bg-surface-900"
    >
        <div class="flex items-center gap-2">
            <h2 :id="`template-${template.id}`" class="min-w-0 grow truncate text-lg font-semibold">
                {{ template.name }}
            </h2>
            <button
                :aria-controls="`template-menu-${template.id}`"
                aria-haspopup="menu"
                :aria-label="`Actions for ${template.name}`"
                class="flex size-9 items-center justify-center rounded-control text-surface-700 hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:text-surface-100 dark:hover:bg-surface-800"
                type="button"
                @click="menu?.toggle($event)"
            >
                <OiEllipsisH aria-hidden="true" :size="18" />
            </button>
        </div>

        <ol :aria-label="`Shifts in ${template.name}`" class="flex flex-wrap items-center gap-1.5">
            <li v-for="(shift, index) in shifts" :key="shift.id" class="inline-flex items-center gap-1.5">
                <span
                    class="inline-flex h-7 items-center rounded-md bg-surface-100 px-2.5 text-xs font-semibold text-surface-700 dark:bg-surface-800 dark:text-surface-100"
                >
                    {{ shift.shiftName }}
                </span>
                <!-- Separator after the chip, so a wrapped line ends with "·" instead of starting with it -->
                <span v-if="index < shifts.length - 1" aria-hidden="true" class="text-surface-500 dark:text-surface-300"
                    >·</span
                >
            </li>
        </ol>

        <div
            class="mt-auto flex justify-center gap-3 border-t border-surface-100 pt-3.5 text-[13px] text-surface-600 dark:border-surface-800 dark:text-surface-200"
        >
            <Button
                v-if="isAdmin"
                class="h-9 bg-surface-0 text-sm dark:bg-surface-900"
                label="Use template"
                outlined
                severity="secondary"
                @click="emit('use')"
            />
        </div>

        <Menu :id="`template-menu-${template.id}`" ref="menu" :model="menuItems" popup>
            <template #item="{ item, props: itemProps }">
                <a v-bind="itemProps.action" :class="['gap-2', item.danger ? 'text-red-700 dark:text-red-300' : '']">
                    <span>{{ item.label }}</span>
                </a>
            </template>
        </Menu>
        <RosterTemplatePreferences
            :template="template"
            :visible="activeDialog === 'preferences'"
            @update:visible="closeDialog"
        />
        <TemplateShiftGroupsDialog
            :shift-groups="shiftGroups"
            :template="template"
            :visible="activeDialog === 'shift-groups'"
            @update:visible="closeDialog"
        />
    </article>
</template>
