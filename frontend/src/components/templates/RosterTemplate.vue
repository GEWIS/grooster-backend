<script setup lang="ts">
import {RosterTemplate as RosterTemplateModel} from "@gewis/grooster-backend-ts";
import {ref} from "vue";
import RosterTemplateDelete from "@/components/templates/dialogs/RosterTemplateDelete.vue";
import RosterTemplateEdit from "@/components/templates/dialogs/RosterTemplateEdit.vue";

const props = defineProps<{
    template: RosterTemplateModel,
}>();
const openEditDialog = ref(false);
const openDeleteDialog = ref(false);

const showDetail = ref<boolean>(false);

</script>

<template>
    <div>
    <Card class="w-1/3">
        <template #content>
            <div class="grid grid-cols-3 gap-4">
                <div class="col-span-2 flex items-center font-bold">{{props.template.name}}</div>
                <div class="flex flex-row gap-2">
                    <Button icon="pi pi-pencil" @click="openEditDialog=true" aria-label="Edit"/>
                    <Button icon="pi pi-trash"  @click="openDeleteDialog=true" aria-label="Delete"/>
                </div>
                <div class="col-span-3 border-t border-gray-300  mx-auto my-4"></div>
                <div></div>
                <div>
                    <Button label="View details" @click="showDetail=!showDetail" variant="text" class="px-5"/>
                </div>
            </div>

            <!-- Detail section of card -->
            <div v-if="showDetail">
                <ul>
                    <li v-for="shift in props.template.shifts" :key="shift" class="flex justify-between items-center">
                        <span>{{ shift }}</span>
                    </li>
                </ul>
            </div>
        </template>
    </Card>
    </div>

    <RosterTemplateDelete :open="openDeleteDialog" :template-id="props.template.id" @close="openDeleteDialog=false"/>
    <RosterTemplateEdit :open="openEditDialog" :template="props.template" @close="openEditDialog=false"/>
</template>

<style scoped>

</style>