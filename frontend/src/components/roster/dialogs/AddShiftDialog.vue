<script setup lang="ts">
import { ref, watch } from 'vue';
import { dialogPt, readableDisabled } from '@/theme/pt';

const visible = defineModel<boolean>('visible', { required: true });
const emit = defineEmits<{ submit: [name: string] }>();

const name = ref('');
watch(visible, (open) => {
    if (open) name.value = '';
});

function submit() {
    const trimmed = name.value.trim();
    if (!trimmed) return;
    emit('submit', trimmed);
    visible.value = false;
}
</script>

<template>
    <Dialog v-model:visible="visible" header="Add shift" modal :pt="dialogPt">
        <form id="add-shift-form" class="flex flex-col gap-1.5" @submit.prevent="submit">
            <label class="text-[13px] font-semibold text-surface-700 dark:text-surface-100" for="add-shift-name">
                Shift name
            </label>
            <InputText
                id="add-shift-name"
                v-model="name"
                autofocus
                class="h-[42px] w-full text-base md:text-sm"
                placeholder="E.g. Bar"
            />
        </form>
        <template #footer>
            <Button class="h-10 text-sm" label="Cancel" outlined severity="secondary" @click="visible = false" />
            <Button
                :class="['h-10 text-sm', readableDisabled]"
                :disabled="!name.trim()"
                form="add-shift-form"
                label="Add shift"
                type="submit"
            />
        </template>
    </Dialog>
</template>
