<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ShiftGroup } from '@gewis/grooster-backend-ts';
import RosterTemplate from '@/components/templates/RosterTemplate.vue';
import { useOrganStore } from '@/stores/organ.store';
import RosterTemplateAdd from '@/components/templates/dialogs/RosterTemplateAdd.vue';
import { useTemplateStore } from '@/stores/template.store';
import ShiftGroupAdd from '@/components/templates/dialogs/ShiftGroupAdd.vue';
import ApiService from '@/services/ApiService';

type Dialogs = 'AddTemplate' | 'AddGroup' | 'None';

const organStore = useOrganStore();
const templateStore = useTemplateStore();

const templates = computed(() => Object.values(templateStore.templates));
const openDialog = ref<Dialogs>('None');

const shiftGroups = ref<ShiftGroup[]>([]);

const fetchGroups = async () => {
  try {
    const response = await ApiService.shiftGroupApi.getShiftGroups(organStore.organ);
    shiftGroups.value = response.data ?? [];
  } catch (error) {
    console.error('Failed to fetch groups:', error);
  }
};

const handleGroupClose = async () => {
  openDialog.value = 'None';
  await fetchGroups();
};

onMounted(async () => {
  await templateStore.fetchTemplates(organStore.organ);
  await fetchGroups();
});
</script>

<template>
  <div class="p-4 max-w-7xl mx-auto">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Roster Management</h1>
        <p class="text-gray-500">Manage and deploy your recurring shift schedules and groups.</p>
      </div>
      <div class="flex gap-2">
        <Button class="p-button-outlined" icon="pi pi-users" label="Add shift group" @click="openDialog = 'AddGroup'" />
        <Button
          class="p-button-raised shadow-md"
          icon="pi pi-plus"
          label="Create Template"
          @click="openDialog = 'AddTemplate'"
        />
      </div>
    </div>

    <section class="mb-12">
      <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
        <i class="pi pi-copy text-primary"></i> Roster Templates
      </h2>

      <div v-if="templates.length === 0" class="text-center py-12 border-2 border-dashed border-gray-200 rounded-xl">
        <i class="pi pi-calendar-plus text-5xl text-gray-300 mb-4"></i>
        <p class="text-gray-500 text-lg">No templates found.</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <RosterTemplate
          v-for="template in templates"
          :key="template.id"
          :shift-groups="shiftGroups"
          :template="template"
        />
      </div>
    </section>

    <hr class="border-gray-100 mb-12" />

    <section>
      <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center gap-2">
        <i class="pi pi-users text-primary"></i> Shift Groups
      </h2>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div
          v-for="group in shiftGroups"
          :key="group.id"
          class="p-4 bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer flex items-center justify-between"
        >
          <div class="flex items-center gap-3">
            <div>
              <p class="font-bold text-gray-800">{{ group.name }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>

  <RosterTemplateAdd :open="openDialog === 'AddTemplate'" @close="openDialog = 'None'" />
  <ShiftGroupAdd :open="openDialog === 'AddGroup'" @close="handleGroupClose" />
</template>

<style scoped></style>
4
