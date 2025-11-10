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
  <div class="flex flex-col gap-2">
    <Card>
      <template #content>
        <div class="flex flex-row gap-2 justify-center">
          <Button label="Create Template" @click="addDialog = true" />
        </div>
      </template>
    </Card>
    <div class="grid grid-cols-3 gap-2">
      <div v-if="templates.length > 0">
        <div v-for="template in templates" :key="template.id">
          <RosterTemplate :template="template" />
        </div>
      </div>
    </div>
  </div>
  <RosterTemplateAdd :open="addDialog" @close="addDialog = false" />
</template>

<style scoped></style>
4
