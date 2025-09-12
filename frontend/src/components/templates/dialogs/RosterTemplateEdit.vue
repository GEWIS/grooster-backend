<script setup lang="ts">
import {computed} from "vue";
import {useTemplateStore} from "@/stores/template.store";
import {RosterTemplate} from "@gewis/grooster-backend-ts";

const templateStore = useTemplateStore();

const props = withDefaults(
    defineProps<{
        open: boolean;
        template: RosterTemplate;
    }>(),
    {
        open: false,
    }
)

const emit = defineEmits(["close"]);

const visible = computed({
    get: () => props.open,
    set: () => {
        emit('close');
    },
});

</script>

<template>
    <Dialog
        v-model:visible="visible"
        modal
        :style="{ width: '25rem'}"
        @hide="emit('close')"
    >
        <template #header>
            <div class="w-full text-center">
                <span class="font- text-lg">Edit a template</span>
            </div>
        </template>
        <div class="flex flex-col gap-2 justify-center">
            <div v-for="shift in props.template.shifts">
                {{shift}}
                <Button icon="pi pi-pencil" aria-label="Edit"/>
                <Button icon="pi pi-trash" aria-label="Delete"/>
            </div>
            <div class="flex flex-row gap-2 justify-center">
                <Button label="Cancel"/>
                <Button label="Save"/>
            </div>
        </div>

    </Dialog>
</template>

<style scoped>

</style>