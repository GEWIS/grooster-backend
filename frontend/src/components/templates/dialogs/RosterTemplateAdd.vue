<script setup lang="ts">
import {computed, reactive, ref} from "vue";
import ApiService from "@/services/ApiService";
import {RosterTemplateCreateRequest} from "@gewis/grooster-backend-ts";
import {useRoute} from "vue-router";
import {useTemplateStore} from "@/stores/template.store";

const route = useRoute();
const  templateStore = useTemplateStore();

const props = withDefaults(
    defineProps<{
        open: boolean;
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

const initialValues = reactive({
    shiftName: '',
});
const shifts = ref<string[]>([]);

const resolver = ({ values }) => {
    const errors = { shiftName: [] };

    if (!values.shiftName) {
        errors.shiftName.push({ message: 'A shift name is required.' });
    }

    return {
        values,
        errors
    };
};

const onFormSubmit = (e) => {
    if (e.valid) {
        shifts.value.push(e.values.shiftName);
    }
}

const removeShift = (index: number) => {
    shifts.value.splice(index, 1);
}

const saveTemplate = async () => {
    const params: RosterTemplateCreateRequest = {
        name: 'Test',
        organId: parseInt(route.params.id as string),
        shifts: shifts.value,
    }
    await templateStore.createTemplate(params);

    shifts.value = [];
    emit('close');
}

</script>

<template>
    <Dialog
        v-model:visible="visible"
        modal
        :style="{ width: '25rem' }"
        @hide="emit('close')"
    >
        <template #header>
            <div class="w-full text-center">
                <span class="font- text-lg">Add Roster Template</span>
            </div>
        </template>
        <div>
            <Form v-slot="$form" :resolver :initial-values @submit="onFormSubmit">
                <div class="flex flex-col gap-1">
                    <InputText name="shiftName" type="text" placeholder="ShiftName" />
                    <Message v-if="$form.shiftName?.invalid" severity="error" size="small" variant="simple">{{ $form.shiftName.error?.message }}</Message>
                </div>
                <Button type="submit" severity="secondary" label="Add roster" />
                <Button @click="saveTemplate" severity="success" label="Save Template"/>
            </Form>
        </div>
        <div v-for="(shift, index) in shifts">
            {{index}}: {{shift}}
            <Button icon="pi pi-times" @click="removeShift(index)" style="font-size: 1rem"/>
        </div>
    </Dialog>
</template>

<style scoped>

</style>