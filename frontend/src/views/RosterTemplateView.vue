<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import RosterTemplate from '@/components/templates/RosterTemplate.vue';
import { useOrganStore } from '@/stores/organ.store';
import RosterTemplateAdd from '@/components/templates/dialogs/RosterTemplateAdd.vue';
import { useTemplateStore } from '@/stores/template.store';

const organStore = useOrganStore();
const templateStore = useTemplateStore();

const templates = computed(() => Object.values(templateStore.templates));
const addDialog = ref(false);

onMounted(async () => {
  await templateStore.fetchTemplates(organStore.organ);
});
</script>

<template>
  <div class="p-4 max-w-7xl mx-auto">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Roster Templates</h1>
        <p class="text-gray-500">Manage and deploy your recurring shift schedules.</p>
      </div>
      <Button class="p-button-raised shadow-md" icon="pi pi-plus" label="Create Template" @click="addDialog = true" />
    </div>

    <div v-if="templates.length === 0" class="text-center py-20 border-2 border-dashed border-gray-200 rounded-xl">
      <i class="pi pi-calendar-plus text-5xl text-gray-300 mb-4"></i>
      <p class="text-gray-500 text-lg">No templates found. Create your first one to get started!</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <RosterTemplate v-for="template in templates" :key="template.id" :template="template" />
    </div>
  </div>

  <RosterTemplateAdd :open="addDialog" @close="addDialog = false" />
</template>

<style scoped></style>
4
