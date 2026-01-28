<script setup lang="ts">
import { RosterTemplate as RosterTemplateModel } from '@gewis/grooster-backend-ts';
import { ref } from 'vue';
import RosterTemplateDelete from '@/components/templates/dialogs/RosterTemplateDelete.vue';
import RosterTemplateUseDialog from '@/components/templates/dialogs/RosterTemplateUseDialog.vue';

const props = defineProps<{
  template: RosterTemplateModel;
}>();
const openDeleteDialog = ref(false);
const openRosterDialog = ref(false);

const showDetail = ref<boolean>(false);
</script>

<template>
  <div>
    <Card>
      <template #content>
        <div class="grid grid-cols-3 gap-4">
          <div class="col-span-2 flex items-center font-bold">{{ props.template.name }}</div>
          <div class="flex flex-row gap-2 justify-end">
            <!--                    <Button icon="pi pi-pencil" @click="openEditDialog=true" aria-label="Edit"/>-->
            <Button icon="pi pi-plus" @click="openRosterDialog = true" />
            <Button aria-label="Delete" icon="pi pi-trash" @click="openDeleteDialog = true" />
          </div>
          <div class="col-span-3 border-t border-gray-300 mx-auto my-4"></div>
          <div></div>
          <div>
            <Button class="px-5" label="View details" variant="text" @click="showDetail = !showDetail" />
          </div>
        </div>

        <!-- Detail section of card -->
        <div v-if="showDetail">
          <ul>
            <li v-for="shift in props.template.shifts" :key="shift.id" class="flex justify-between items-center">
              <span>{{ shift.shiftName }}</span>
            </li>
          </ul>
        </div>
      </template>
    </Card>
  </div>

  <RosterTemplateDelete :open="openDeleteDialog" :template-id="props.template.id" @close="openDeleteDialog = false" />
  <RosterTemplateUseDialog
    :name="props.template.name"
    :open="openRosterDialog"
    :shifts="props.template.shifts.map((s) => s.shiftName)"
    @close="openRosterDialog = false"
  />
  <!--    <RosterTemplateEdit :open="openEditDialog" :template="props.template" @close="openEditDialog=false"/>-->
</template>

<style scoped></style>
