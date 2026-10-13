<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import PageHeader from '@/components/common/PageHeader.vue';
import MemberRoles from '@/components/organ-settings/MemberRoles.vue';
import { Role, useAuthStore } from '@/stores/auth.store';
import { useOrganStore } from '@/stores/organ.store';

const route = useRoute();
const authStore = useAuthStore();
const organStore = useOrganStore();

const organId = computed(() => parseInt(route.params.id as string));
const memberList = ref<InstanceType<typeof MemberRoles>>();
const count = computed(() => memberList.value?.count);
</script>

<template>
    <div class="flex flex-col gap-6 px-4 py-6 md:px-10 md:py-8">
        <PageHeader
            :count="count"
            :description="`Roles and display names inside ${organStore.organName ?? 'this organization'}.`"
            title="Members"
        />
        <MemberRoles v-if="authStore.can([Role.Admin, Role.Owner])" ref="memberList" :organ-id="organId" />
        <p v-else class="text-sm text-surface-600 dark:text-surface-200">Only admins and owners can manage members.</p>
    </div>
</template>
