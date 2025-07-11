<script setup lang="ts">
import {onMounted, ref} from "vue";
import ApiService from "@/services/ApiService";
import {getGEWISId} from "@/helpers/TokenHelper";
import {Select} from "primevue";
import {Organ} from "../../ts-client";
import {useRouter} from "vue-router";

const router = useRouter();

const user = ref();
const organs = ref<Organ[]>();
const selectedOrgan = ref<number>();

onMounted(() => {
    loadUser();
});

async function loadUser() {
    const gewisId = getGEWISId();
    const response = await ApiService.user.userGet(undefined, gewisId);

    user.value = response.data[0];
    organs.value = user.value.organs;
}

const onSelect = () => {
    router.push({ name: 'roster', params: {id: selectedOrgan.value} });
}

</script>

<template>
    <div v-if="organs && user">
        <Select
            v-model="selectedOrgan"
            :options="organs"
            option-label="Name"
            option-value="id"
            placeholder="Select an organ"
            @change="onSelect"
        />
    </div>
    <div v-else>
        <p>You should not see this!</p>
    </div>
</template>

<style scoped>

</style>