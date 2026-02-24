<script setup lang="ts">
import {
  RosterTemplate as RosterTemplateModel,
  ShiftGroup,
  TemplateShiftUpdateRequest,
} from '@gewis/grooster-backend-ts';
import { ref } from 'vue';
import RosterTemplateDelete from '@/components/templates/dialogs/RosterTemplateDelete.vue';
import RosterTemplateUseDialog from '@/components/templates/dialogs/RosterTemplateUseDialog.vue';
import RosterTemplatePreferences from '@/components/templates/dialogs/RosterTemplatePreferences.vue';
import ApiService from '@/services/ApiService';

type Dialogs = 'Preferences' | 'Delete' | 'Roster' | 'None';

const props = defineProps<{
  template: RosterTemplateModel;
  shiftGroups: ShiftGroup[];
}>();
const openDialog = ref<Dialogs>('None');

const showDetail = ref<boolean>(false);

const handleGroupChange = async (shiftId: number, groupId: number) => {
  try {
    const params: TemplateShiftUpdateRequest = {
      shiftGroupId: groupId,
    };

    await ApiService.roster.updateRosterTemplateShift(shiftId, params);
  } catch (e) {
    console.error(e);
  }
};

const getGroupName = (id: number) => {
  return props.shiftGroups.find((g) => g.id === id)?.name || '';
};
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
                class="flex items-center justify-between p-2 bg-gray-50 rounded-md border border-gray-100 gap-3"
              >
                <div class="flex items-center min-w-0 flex-1 gap-1">
                  <i class="pi pi-clock mr-2 text-blue-500 text-sm"></i>
                  <span class="text-sm font-medium text-gray-700 truncate">
                    {{ shift.shiftName }}
                  </span>
                </div>

                <div class="flex-shrink-0 w-40 h-8">
                  <Select
                    v-model="shift.shiftGroupId"
                    class="w-full h-8 text-xs flex items-center border-gray-200"
                    option-label="name"
                    option-value="id"
                    :options="shiftGroups"
                    placeholder="Assign..."
                    show-clear
                    @change="(e) => handleGroupChange(shift.id, e.value)"
                  >
                    <template #value="slotProps">
                      <div v-if="slotProps.value" class="flex items-center gap-2 text-xs truncate">
                        <span class="truncate">{{ getGroupName(slotProps.value) }}</span>
                      </div>
                      <span v-else class="text-xs text-gray-400">Assign...</span>
                    </template>

                    <template #option="slotProps">
                      <div class="flex items-center gap-2 text-xs">
                        <div
                          class="w-2 h-2 rounded-full flex-shrink-0"
                          :style="{ backgroundColor: slotProps.option.color || '#cbd5e1' }"
                        ></div>
                        <span class="truncate">{{ slotProps.option.name }}</span>
                      </div>
                    </template>
                  </Select>
                </div>
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
