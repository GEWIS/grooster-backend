<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { GEWISRoosterInternalModelsOrganRole, User, UserOrgan } from '@gewis/grooster-backend-ts';
import { UpdateMemberSettingsParams } from '@gewis/grooster-backend-ts/dist/api';
import ApiService from '@/services/ApiService';
import { useOrganStore } from '@/stores/organ.store';
import RoleChangeConfirmation from '@/components/organ-settings/dialogs/RoleChangeConfirmation.vue';

export type UserRole = (typeof GEWISRoosterInternalModelsOrganRole)[keyof typeof GEWISRoosterInternalModelsOrganRole];

const roleOptions = Object.entries(GEWISRoosterInternalModelsOrganRole).map(([key, value]) => ({
    label: key.replace('Role', ''),
    value: value,
}));

const organStore = useOrganStore();

const users = ref<User[]>([]);
const userOrganSettings = ref<UserOrgan[]>([]);
const selectedUser = ref<User | null>(null);

const currentRole = computed(() => (selectedUser.value ? getUserRole(selectedUser.value.id) : null));
const pendingRole = ref<UserRole | null>(null);

type Dialog = 'Confirm' | 'None';
const currentDialog = ref<Dialog>('None');

const editName = ref('');
const isSaving = ref(false);

const isChanged = computed(() => {
    if (!selectedUser.value) return false;
    const original = userOrganSettings.value.find((uo) => uo.userId === selectedUser.value.id);
    return original ? editName.value !== original.username : false;
});

const fetchUsers = async () => {
    try {
        const response = await ApiService.user.userGet(organStore.organ);
        users.value = response.data;
    } catch (e) {
        console.error(e);
        users.value = [];
    }
};

const fetchUserOrganSettings = async () => {
    try {
        const response = await ApiService.organ.getMembersSettings(organStore.organ);
        userOrganSettings.value = response.data;
    } catch (e) {
        console.error(e);
    }
};

const getUserRole = (userId: number) => {
    return userOrganSettings.value.find((organ: UserOrgan) => organ.userId == userId)?.role;
};

const handleRoleUpdate = (role: UserRole) => {
    pendingRole.value = role;
    currentDialog.value = 'Confirm';
};

const confirmRoleChange = async (userId: number, role: UserRole) => {
    try {
        const { data } = await ApiService.organ.organIdMemberUserIdRolePatch(organStore.organ, userId, { role });
        userOrganSettings.value = userOrganSettings.value.map((userOrgan) => {
            if (userOrgan.userId === data.userId) {
                return data;
            } else {
                return userOrgan;
            }
        });
    } catch (e) {
        console.error(e);
    } finally {
        pendingRole.value = null;
        currentDialog.value = 'None';
    }
};

const setUser = (user: User) => {
    selectedUser.value = user;
    const settings = userOrganSettings.value.find((uo) => uo.userId === user.id);
    editName.value = settings ? settings.username : '';
};

const resetName = () => {
    if (selectedUser.value) {
        const settings = userOrganSettings.value.find((uo) => uo.userId === selectedUser.value.id);
        editName.value = settings ? settings.username : '';
    }
};

const updateMemberSettings = async (userId: number) => {
    const params: UpdateMemberSettingsParams = {
        username: editName.value,
    };

    isSaving.value = true;
    try {
        await ApiService.organ.organIdMemberUserIdPatch(organStore.organ, userId, params);
        await fetchUserOrganSettings();
    } catch (e) {
        console.error('Failed to update settings:', e);
    } finally {
        isSaving.value = false;
    }
};

onMounted(fetchUsers);
onMounted(fetchUserOrganSettings);
</script>
<template>
    <div
        v-if="users.length > 0"
        class="bg-white border border-emerald-100 rounded-2xl overflow-hidden shadow-xl max-w-3xl mx-auto flex flex-col h-[550px]"
    >
        <div class="bg-emerald-950 p-6 flex items-center justify-between shrink-0">
            <div class="flex items-center gap-4">
                <div
                    class="flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-500/20 text-emerald-400 border border-emerald-500/30"
                >
                    <i class="pi pi-users text-xl" />
                </div>
                <div>
                    <h2 class="text-white font-bold text-lg leading-tight">Organization Members</h2>
                    <p class="text-emerald-400/60 text-xs">Manage roles and local nicknames</p>
                </div>
            </div>
            <div class="text-emerald-500/50 text-sm font-mono">{{ users.length }} Members</div>
        </div>

        <div class="flex flex-1 overflow-hidden">
            <div class="w-1/3 border-r border-emerald-50 overflow-y-auto bg-emerald-50/20">
                <div
                    v-for="user in users"
                    :key="user.id"
                    :class="[
                        'p-4 cursor-pointer transition-all border-b border-emerald-50/50 flex items-center gap-3',
                        selectedUser?.id === user.id
                            ? 'bg-emerald-100/50 border-r-4 border-r-emerald-600'
                            : 'hover:bg-emerald-50',
                    ]"
                    @click="() => setUser(user)"
                >
                    <div
                        class="h-8 w-8 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                        :class="
                            selectedUser?.id === user.id
                                ? 'bg-emerald-600 text-white'
                                : 'bg-emerald-200 text-emerald-800'
                        "
                    >
                        {{ user.name.charAt(0) }}
                    </div>
                    <p
                        :class="[
                            'truncate text-sm',
                            selectedUser?.id === user.id ? 'font-bold text-emerald-900' : 'text-gray-600',
                        ]"
                    >
                        {{ user.name }}
                    </p>
                </div>
            </div>

            <div class="w-2/3 overflow-y-auto bg-white relative">
                <div v-if="selectedUser" class="p-8 h-full flex flex-col">
                    <div class="flex items-center gap-5 mb-8">
                        <div
                            class="h-20 w-20 bg-emerald-100 rounded-2xl flex items-center justify-center text-emerald-900 shadow-inner border border-emerald-200 shrink-0"
                        >
                            <span class="text-3xl font-black">{{ selectedUser.name.charAt(0) }}</span>
                        </div>
                        <div class="overflow-hidden">
                            <h3 class="font-bold text-2xl text-gray-900 truncate">{{ selectedUser.name }}</h3>
                            <div
                                class="inline-flex items-center gap-2 px-2.5 py-0.5 rounded-full bg-emerald-100 text-emerald-800 text-xs font-bold uppercase tracking-wider mt-1"
                            >
                                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                                {{ getUserRole(selectedUser.id) }}
                            </div>
                        </div>
                    </div>

                    <div class="space-y-6 bg-emerald-50/30 p-6 rounded-2xl border border-emerald-100">
                        <div>
                            <label
                                class="block text-[10px] font-bold text-emerald-800 uppercase tracking-widest mb-2 px-1"
                                >Organization Role</label
                            >
                            <Select
                                class="w-full border-emerald-200 text-sm focus:ring-emerald-500 shadow-sm"
                                :model-value="currentRole"
                                option-label="label"
                                option-value="value"
                                :options="roleOptions"
                                placeholder="Select a Role"
                                @update:model-value="(value) => handleRoleUpdate(value)"
                            />
                        </div>

                        <div>
                            <label
                                class="block text-[10px] font-bold text-emerald-800 uppercase tracking-widest mb-2 px-1"
                                >Display Name (Nickname)</label
                            >
                            <div class="relative">
                                <input
                                    id="username"
                                    v-model="editName"
                                    class="w-full px-4 py-2.5 rounded-lg border border-emerald-200 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition-all text-emerald-950 font-medium shadow-sm"
                                    maxlength="25"
                                    placeholder="Their name in this organ..."
                                    type="text"
                                    @keyup.enter="updateMemberSettings(selectedUser.id)"
                                />
                                <i class="pi pi-pencil absolute right-3 top-3 text-emerald-300 pointer-events-none" />
                            </div>
                        </div>
                    </div>

                    <div class="mt-auto pt-6 flex items-center justify-end gap-3 border-t border-emerald-50">
                        <button
                            v-if="isChanged"
                            class="px-5 py-2.5 text-gray-500 hover:text-emerald-700 hover:bg-emerald-50 rounded-xl transition-all text-sm font-semibold"
                            @click="resetName"
                        >
                            Discard
                        </button>

                        <button
                            class="flex items-center gap-2 px-8 py-2.5 bg-emerald-700 text-white rounded-xl font-bold hover:bg-emerald-800 active:scale-[0.98] disabled:opacity-30 disabled:grayscale disabled:cursor-not-allowed transition-all shadow-md shadow-emerald-900/10"
                            :disabled="!isChanged || isSaving"
                            @click="updateMemberSettings(selectedUser.id)"
                        >
                            <i v-if="isSaving" class="pi pi-spin pi-spinner text-xs" />
                            <span>{{ isSaving ? 'Saving...' : 'Update Settings' }}</span>
                        </button>
                    </div>
                </div>

                <div v-else class="h-full flex flex-col items-center justify-center p-12 text-center">
                    <div
                        class="w-20 h-20 bg-gray-50 rounded-full flex items-center justify-center mb-4 border border-dashed border-gray-200"
                    >
                        <i class="pi pi-user-plus text-2xl text-gray-300" />
                    </div>
                    <h3 class="text-gray-900 font-bold">No Member Selected</h3>
                    <p class="text-gray-400 text-sm mt-1 max-w-[200px]">
                        Select a member from the list to manage their profile settings.
                    </p>
                </div>
            </div>
        </div>
    </div>

    <RoleChangeConfirmation
        v-if="pendingRole && selectedUser"
        :open="currentDialog === 'Confirm'"
        :role="pendingRole"
        :user-id="selectedUser.id"
        @close="currentDialog = 'None'"
        @confirm="({ userId, role }) => confirmRoleChange(userId, role)"
    />
</template>
