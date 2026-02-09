<script setup lang="ts">
import { RosterTemplate as RosterTemplateModel } from '@gewis/grooster-backend-ts';
import { ref } from 'vue';
import RosterTemplateDelete from '@/components/templates/dialogs/RosterTemplateDelete.vue';
import RosterTemplateUseDialog from '@/components/templates/dialogs/RosterTemplateUseDialog.vue';
import RosterTemplatePreferences from '@/components/templates/dialogs/RosterTemplatePreferences.vue';

type Dialogs = 'Preferences' | 'Delete' | 'Roster' | 'None';

const props = defineProps<{
  template: RosterTemplateModel;
}>();
const openDialog = ref<Dialogs>('None');

const showDetail = ref<boolean>(false);
</script>

<template>
  <div class="p-2">
    <Card class="overflow-hidden border border-gray-200 shadow-sm hover:shadow-md transition-shadow duration-300">
      <template #content>
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800 truncate">
            {{ props.template.name }}
          </h3>
          <div class="flex gap-2">
            <Button icon="pi pi-heart" rounded @click="openDialog = 'Preferences'" />
            <Button icon="pi pi-plus" rounded severity="secondary" text @click="openDialog = 'Roster'" />
            <Button icon="pi pi-trash" rounded severity="danger" text @click="openDialog = 'Delete'" />
          </div>
        </div>

        <Divider class="my-0" />

        <div class="flex justify-between items-center mt-3">
          <div class="flex items-center text-sm text-gray-500">
            <i class="pi pi-calendar mr-2"></i>
            {{ props.template.shifts?.length || 0 }} Shifts
          </div>
          <Button
            class="p-button-sm"
            :icon="showDetail ? 'pi pi-chevron-up' : 'pi pi-chevron-down'"
            icon-pos="right"
            :label="showDetail ? 'Hide details' : 'View details'"
            variant="text"
            @click="showDetail = !showDetail"
          />
        </div>

        <transition name="p-toggleable-content">
          <div v-if="showDetail" class="mt-4 pt-4 border-t border-dashed border-gray-200">
            <ul class="space-y-2">
              <li
                v-for="shift in props.template.shifts"
                :key="shift.id"
                class="flex items-center p-2 bg-gray-50 rounded-md border border-gray-100"
              >
                <i class="pi pi-clock mr-3 text-blue-500"></i>
                <span class="text-sm font-medium text-gray-700">{{ shift.shiftName }}</span>
              </li>
            </ul>
          </div>
        </transition>
      </template>
    </Card>
  </div>

  <RosterTemplateDelete :open="openDialog === 'Delete'" :template-id="props.template.id" @close="openDialog = 'None'" />
  <RosterTemplateUseDialog
    :name="props.template.name"
    :open="openDialog === 'Roster'"
    :shifts="props.template.shifts.map((s) => s.shiftName)"
    :template-id="props.template.id"
    @close="openDialog = 'None'"
  />
  <RosterTemplatePreferences
    :open="openDialog === 'Preferences'"
    :template="props.template"
    @close="openDialog = 'None'"
  />
</template>

<style scoped></style>
