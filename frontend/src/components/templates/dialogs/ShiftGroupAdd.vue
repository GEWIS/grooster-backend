<script setup lang="ts">
import { ref, watch } from 'vue';
import { useToast } from 'openvue/usetoast';
import ApiService from '@/services/ApiService';
import { dialogPt, readableDisabled } from '@/theme/pt';

const props = defineProps<{ organId: number }>();
const visible = defineModel<boolean>('visible', { required: true });
const emit = defineEmits<{ created: [] }>();

const toast = useToast();
const name = ref('');
const saving = ref(false);
watch(visible, (open) => open && (name.value = ''));

async function create() {
    if (!name.value.trim()) return;
    saving.value = true;
    try {
        await ApiService.shiftGroupApi.createShiftGroup({ organId: props.organId, name: name.value.trim() });
        emit('created');
        visible.value = false;
    } catch (error) {
        console.error(error);
        toast.add({ severity: 'error', summary: 'Could not create the group', life: 5000 });
    } finally {
        saving.value = false;
    }
}
</script>

<template>
    <Dialog v-model:visible="visible" header="New shift group" modal :pt="dialogPt">
        <form id="shift-group-form" class="flex flex-col gap-1.5" @submit.prevent="create">
            <label class="text-[13px] font-semibold text-surface-700 dark:text-surface-100" for="shift-group-name">
                Group name
            </label>
            <InputText
                id="shift-group-name"
                v-model="name"
                autocomplete="off"
                autofocus
                class="h-[42px] w-full text-base md:text-sm"
            />
        </form>
        <template #footer>
            <Button class="h-10 text-sm" label="Cancel" outlined severity="secondary" @click="visible = false" />
            <Button
                :class="['h-10 text-sm', readableDisabled]"
                :disabled="!name.trim()"
                form="shift-group-form"
                label="Create group"
                :loading="saving"
                type="submit"
            />
        </template>
    </Dialog>
</template>
