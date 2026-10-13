<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useConfirm } from 'openvue/useconfirm';
import { useToast } from 'openvue/usetoast';
import IconField from 'openvue/iconfield';
import InputIcon from 'openvue/inputicon';
import { OiChevronLeft, OiSearch } from '@openvue/openicons-vue';
import { GEWISRoosterInternalModelsOrganRole, type User, type UserOrgan } from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import { fieldSelectPt, readableDisabled } from '@/theme/pt';

export type UserRole = (typeof GEWISRoosterInternalModelsOrganRole)[keyof typeof GEWISRoosterInternalModelsOrganRole];

const roleOptions = Object.entries(GEWISRoosterInternalModelsOrganRole).map(([key, value]) => ({
    label: key.replace('Role', ''),
    value,
}));
const roleHelp: Record<UserRole, string> = {
    owner: 'Owners and admins manage rosters, templates and members.',
    admin: 'Owners and admins manage rosters, templates and members.',
    member: 'Members fill in their own answers.',
};

const props = defineProps<{ organId: number }>();

const confirm = useConfirm();
const toast = useToast();

const users = ref<User[]>([]);
const settings = ref<UserOrgan[]>([]);
const selectedId = ref<number | null>(null);
const search = ref('');

const settingFor = (userId: number) => settings.value.find((s) => s.userId === userId);
const roleOf = (userId: number) => settingFor(userId)?.role;
const roleLabel = (role: UserRole | undefined) => roleOptions.find((o) => o.value === role)?.label ?? '';

const filteredUsers = computed(() => {
    const query = search.value.trim().toLowerCase();
    if (!query) return users.value;
    return users.value.filter(
        (user) =>
            user.name.toLowerCase().includes(query) ||
            (settingFor(user.id)?.username ?? '').toLowerCase().includes(query),
    );
});
const selected = computed(() => users.value.find((user) => user.id === selectedId.value) ?? null);

// Form state for the selected member; saved together with "Save changes"
const draftRole = ref<UserRole | null>(null);
const draftName = ref('');
const saving = ref(false);

function select(user: User) {
    selectedId.value = user.id;
    resetDraft();
}
function resetDraft() {
    if (!selected.value) return;
    draftRole.value = roleOf(selected.value.id) ?? null;
    draftName.value = settingFor(selected.value.id)?.username ?? '';
}

const roleChanged = computed(() => !!selected.value && draftRole.value !== roleOf(selected.value.id));
const nameChanged = computed(
    () => !!selected.value && draftName.value !== (settingFor(selected.value.id)?.username ?? ''),
);
const dirty = computed(() => roleChanged.value || nameChanged.value);

async function persist() {
    const user = selected.value;
    if (!user) return;
    saving.value = true;
    try {
        if (roleChanged.value && draftRole.value) {
            const { data } = await ApiService.organ.organIdMemberUserIdRolePatch(props.organId, user.id, {
                role: draftRole.value,
            });
            settings.value = settings.value.map((s) => (s.userId === data.userId ? data : s));
        }
        if (nameChanged.value) {
            await ApiService.organ.organIdMemberUserIdPatch(props.organId, user.id, { username: draftName.value });
            await fetchSettings();
        }
        toast.add({ severity: 'success', summary: 'Changes saved', life: 3000 });
    } catch (error) {
        console.error('Failed to update member:', error);
        toast.add({ severity: 'error', summary: 'Could not save the changes', life: 5000 });
    } finally {
        saving.value = false;
        resetDraft();
    }
}

function save() {
    if (!roleChanged.value) return void persist();
    confirm.require({
        header: 'Change role?',
        message: `${selected.value?.name} becomes ${roleLabel(draftRole.value ?? undefined)} in this organization.`,
        acceptProps: { label: 'Change role' },
        rejectProps: { label: 'Cancel', severity: 'secondary', outlined: true },
        accept: () => void persist(),
    });
}

async function fetchUsers() {
    try {
        users.value = (await ApiService.user.userGet(props.organId)).data;
    } catch (error) {
        console.error(error);
    }
}
async function fetchSettings() {
    try {
        settings.value = (await ApiService.organ.getMembersSettings(props.organId)).data;
    } catch (error) {
        console.error(error);
    }
}
onMounted(() => Promise.all([fetchUsers(), fetchSettings()]));

defineExpose({ count: computed(() => users.value.length) });

const roleTagClass = (role: UserRole | undefined) =>
    role === 'owner'
        ? 'bg-surface-900 text-white dark:bg-surface-0 dark:text-surface-900'
        : 'bg-surface-100 text-surface-700 dark:bg-surface-800 dark:text-surface-100';
const labelClass = 'text-[13px] font-semibold text-surface-700 dark:text-surface-100';
</script>

<template>
    <div class="flex overflow-hidden rounded-card border border-surface bg-surface-0 md:h-[660px] dark:bg-surface-900">
        <!-- List -->
        <div
            :class="[
                'w-full shrink-0 flex-col border-surface md:flex md:w-[400px] md:border-r',
                selected ? 'hidden' : 'flex',
            ]"
        >
            <div class="border-b border-surface px-4 py-3.5">
                <IconField>
                    <InputIcon><OiSearch aria-hidden="true" :size="16" /></InputIcon>
                    <InputText
                        v-model="search"
                        aria-label="Search members"
                        class="h-10 w-full text-base md:text-sm"
                        placeholder="Search by name"
                    />
                </IconField>
            </div>
            <ul aria-label="Members" class="grow overflow-y-auto">
                <li v-for="user in filteredUsers" :key="user.id">
                    <button
                        :aria-current="selectedId === user.id ? 'true' : undefined"
                        :class="[
                            'flex h-[60px] w-full items-center gap-3 border-l-[3px] px-4 text-left focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary',
                            selectedId === user.id
                                ? 'border-l-primary bg-highlight'
                                : 'border-l-transparent hover:bg-surface-50 dark:hover:bg-surface-800',
                        ]"
                        type="button"
                        @click="select(user)"
                    >
                        <span class="flex min-w-0 grow flex-col">
                            <span class="truncate text-sm font-semibold">{{ user.name }}</span>
                            <span
                                v-if="settingFor(user.id)?.username"
                                class="truncate text-xs text-surface-600 dark:text-surface-200"
                            >
                                Shown as {{ settingFor(user.id)?.username }}
                            </span>
                        </span>
                        <span
                            :class="['rounded-full px-2 py-0.5 text-xs font-semibold', roleTagClass(roleOf(user.id))]"
                        >
                            {{ roleLabel(roleOf(user.id)) }}
                        </span>
                    </button>
                </li>
                <li v-if="!filteredUsers.length" class="px-4 py-6 text-sm text-surface-600 dark:text-surface-200">
                    No members match "{{ search }}".
                </li>
            </ul>
        </div>

        <!-- Detail -->
        <div :class="['min-w-0 grow flex-col md:flex', selected ? 'flex' : 'hidden']">
            <template v-if="selected">
                <button
                    class="flex h-11 items-center gap-1 px-3 text-sm font-semibold text-primary md:hidden"
                    type="button"
                    @click="selectedId = null"
                >
                    <OiChevronLeft aria-hidden="true" :size="18" /> All members
                </button>
                <div class="flex flex-col gap-1.5 border-b border-surface px-5 pb-6 pt-2 md:px-8 md:py-7">
                    <div class="flex flex-wrap items-center gap-2.5">
                        <h2 class="font-display text-[26px] font-semibold">{{ selected.name }}</h2>
                        <span
                            :class="[
                                'rounded-full px-2 py-0.5 text-xs font-semibold',
                                roleTagClass(roleOf(selected.id)),
                            ]"
                        >
                            {{ roleLabel(roleOf(selected.id)) }}
                        </span>
                    </div>
                    <span class="text-sm text-surface-600 dark:text-surface-200">
                        {{
                            settingFor(selected.id)?.username
                                ? `Shown as ${settingFor(selected.id)?.username} on rosters`
                                : 'No display name set; rosters show the full name'
                        }}
                    </span>
                </div>
                <form
                    id="member-form"
                    class="grid max-w-[720px] grid-cols-1 gap-5 px-5 py-6 md:grid-cols-2 md:px-8 md:py-7"
                    @submit.prevent="save"
                >
                    <div class="flex flex-col gap-1.5">
                        <label :class="labelClass" for="member-role">Organization role</label>
                        <Select
                            v-model="draftRole"
                            aria-describedby="member-role-help"
                            class="h-[42px] w-full items-center"
                            input-id="member-role"
                            option-label="label"
                            option-value="value"
                            :options="roleOptions"
                            :pt="fieldSelectPt"
                        />
                        <span id="member-role-help" class="text-xs text-surface-600 dark:text-surface-200">
                            {{ draftRole ? roleHelp[draftRole] : '' }}
                        </span>
                    </div>
                    <div class="flex flex-col gap-1.5">
                        <label :class="labelClass" for="member-name">Display name</label>
                        <InputText
                            id="member-name"
                            v-model="draftName"
                            aria-describedby="member-name-help"
                            class="h-[42px] w-full text-base md:text-sm"
                            maxlength="25"
                            :placeholder="selected.name"
                        />
                        <span id="member-name-help" class="text-xs text-surface-600 dark:text-surface-200">
                            Shown on rosters and exports in this organization.
                        </span>
                    </div>
                </form>
                <span class="grow" />
                <div
                    class="flex flex-wrap items-center justify-end gap-2.5 border-t border-surface bg-surface-50 px-5 py-4 md:px-8 dark:bg-surface-800"
                >
                    <span
                        v-if="dirty"
                        aria-live="polite"
                        class="mr-auto text-[13px] text-surface-600 dark:text-surface-200"
                    >
                        Unsaved changes
                    </span>
                    <Button
                        class="h-10 bg-surface-0 text-sm dark:bg-surface-900"
                        :disabled="!dirty"
                        label="Discard"
                        outlined
                        severity="secondary"
                        @click="resetDraft"
                    />
                    <Button
                        :class="['h-10 text-sm', readableDisabled]"
                        :disabled="!dirty"
                        form="member-form"
                        label="Save changes"
                        :loading="saving"
                        type="submit"
                    />
                </div>
            </template>
            <p v-else class="m-auto px-6 text-center text-sm text-surface-600 dark:text-surface-200">
                Pick a member to change their role or display name.
            </p>
        </div>
    </div>
</template>
