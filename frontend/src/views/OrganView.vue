<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { OiChevronRight } from '@openvue/openicons-vue';
import type { Organ } from '@gewis/grooster-backend-ts';
import { useRouter } from 'vue-router';
import { getGEWISId } from '@/helpers/TokenHelper';
import ApiService from '@/services/ApiService';
import PageHeader from '@/components/common/PageHeader.vue';

const router = useRouter();

const user = ref();
const organs = ref<Organ[]>();

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
    <div class="mx-auto flex w-full max-w-[560px] flex-col gap-6 px-4 py-6 md:py-10">
        <PageHeader description="Pick the organization you want to plan for." title="Organizations" />
        <ul
            v-if="organs?.length"
            class="overflow-hidden rounded-card border border-surface bg-surface-0 dark:bg-surface-900"
        >
            <li v-for="organ in organs" :key="organ.id" class="border-t border-surface first:border-t-0">
                <button
                    class="flex h-14 w-full items-center gap-3 px-5 text-left text-[15px] font-semibold hover:bg-surface-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary dark:hover:bg-surface-800"
                    type="button"
                    @click="onSelect(organ.id)"
                >
                    <span class="grow">{{ organ.name }}</span>
                    <OiChevronRight aria-hidden="true" class="text-surface-600 dark:text-surface-200" :size="18" />
                </button>
            </li>
        </ul>
        <p v-else-if="user" class="text-sm text-surface-600 dark:text-surface-200">
            You are not a member of any organization yet.
        </p>
    </div>
</template>
