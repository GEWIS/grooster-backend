<script setup lang="ts">
import {computed, ref} from "vue";
import {RosterCreateRequest} from "@gewis/grooster-backend-ts";
import {useRoute} from "vue-router";
import ApiService from "@/services/ApiService";

const route = useRoute();

const props = withDefaults(
    defineProps<{
        open: boolean;
        name: string;
        shifts: string[];
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

const date = ref(new Date());

const addRoster = async () => {
    const createParams: RosterCreateRequest = {
        name: props.name,
        date: date.value.toISOString(),
        organId: parseInt(route.params.id as string),
        shifts: props.shifts,
    }

    await ApiService.roster.createRoster(createParams);
    emit("close");}
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
                <span class="font- text-lg">Add Roster</span>
            </div>
        </template>
        <div class="flex row gap-3">
            <label class="font-semibold">Roster Date</label>
            <DatePicker v-model="date" :minDate="new Date()" />
        </div>
        <div class="flex justify-center gap-3 pt-4">
            <Button label="Cancel" severity="secondary" @click="emit('close')" class="px-6"/>
            <Button label="Add" @click="addRoster()" class="px-6"/>
        </div>
    </Dialog>
</template>

<style scoped>

</style>