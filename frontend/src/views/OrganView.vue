<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Select } from 'primevue';
import type { Organ } from '@gewis/grooster-backend-ts';
import { useRouter } from 'vue-router';
import { getGEWISId } from '@/helpers/TokenHelper';
import ApiService from '@/services/ApiService';

const router = useRouter();

const user = ref();
const organs = ref<Organ[]>();
const selectedOrgan = ref<number>();

onMounted(async () => {
  await loadUser();
});

async function loadUser() {
  const gewisId = getGEWISId();
  const response = await ApiService.user.userGet(undefined, gewisId);

  user.value = response.data[0];
  organs.value = user.value.organs;

  if (organs.value.length === 1) {
    await onSelect(organs.value[0].id);
  }
}

const onSelect = async (organId: number) => {
  await router.push({ name: 'rosters', params: { id: organId } });
};
</script>

<template>
  <div v-if="organs && user">
    <Select
      v-model="selectedOrgan"
      option-label="name"
      option-value="id"
      :options="organs"
      placeholder="Select an organ"
      @change="onSelect(selectedOrgan)"
    />
  </div>
  <div v-else>
    <p>You should not see this!</p>
  </div>
</template>

<style scoped></style>
