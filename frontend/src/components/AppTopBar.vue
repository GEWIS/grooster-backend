<script setup lang="ts">
import { computed, ref, watch, type Component } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Menu from 'openvue/menu';
import Drawer from 'openvue/drawer';
import {
    OiBars,
    OiCheck,
    OiChevronDown,
    OiClone,
    OiDesktop,
    OiMoon,
    OiSun,
    OiTable,
    OiUsers,
} from '@openvue/openicons-vue';
import { useOrganStore } from '@/stores/organ.store';
import { useUserStore } from '@/stores/user.store';
import { Role, useAuthStore } from '@/stores/auth.store';
import { type ThemeMode, useThemeStore } from '@/stores/theme.store';

interface NavItem {
    label: string;
    to: string;
    icon: Component;
}

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const organStore = useOrganStore();
const userStore = useUserStore();
const themeStore = useThemeStore();

const navItems = computed<NavItem[]>(() => {
    const organ = organStore.organ;
    if (route.path === '/' || !organ) return [];

    const items: NavItem[] = [
        { label: 'Rosters', to: `/rosters/${organ}`, icon: OiTable },
        { label: 'Templates', to: `/templates/${organ}`, icon: OiClone },
    ];
    if (authStore.can([Role.Admin, Role.Owner])) {
        items.push({ label: 'Members', to: `/organ/${organ}/profile`, icon: OiUsers });
    }
    return items;
});

const pageTitle = computed(() => navItems.value.find((item) => item.to === route.path)?.label ?? 'GRooster');

const organs = computed(() => userStore.getUser?.organs ?? []);

const organMenu = ref<InstanceType<typeof Menu>>();
const organMenuItems = computed(() =>
    organs.value.map((organ) => ({
        label: organ.name,
        current: organ.id === organStore.organ,
        command: () => selectOrgan(organ.id),
    })),
);

function selectOrgan(id: number) {
    void router.push({ name: 'rosters', params: { id } });
}

const themeIcons: Record<ThemeMode, Component> = { system: OiDesktop, light: OiSun, dark: OiMoon };
const themeOrder: ThemeMode[] = ['system', 'light', 'dark'];
const nextThemeMode = computed(() => themeOrder[(themeOrder.indexOf(themeStore.mode) + 1) % themeOrder.length]);
const themeLabel = computed(() => `Theme: ${themeStore.mode}. Switch to ${nextThemeMode.value}`);

const isDrawerOpen = ref(false);
watch(
    () => route.fullPath,
    () => (isDrawerOpen.value = false),
);

const focusOnDark = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white';
const focusOnLight = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary';
</script>

<template>
    <header class="flex h-14 shrink-0 items-center gap-2 bg-topbar px-2 text-white md:h-16 md:gap-6 md:pl-5 md:pr-6">
        <!-- Phone: menu button + page title -->
        <button
            :aria-expanded="isDrawerOpen"
            aria-haspopup="dialog"
            aria-label="Open menu"
            :class="[
                'flex size-11 items-center justify-center rounded-lg hover:bg-topbar-hover md:hidden',
                focusOnDark,
            ]"
            type="button"
            @click="isDrawerOpen = true"
        >
            <OiBars aria-hidden="true" :size="22" />
        </button>
        <span class="grow truncate font-display text-[17px] font-semibold md:hidden">{{ pageTitle }}</span>

        <!-- Desktop: logo, nav, organization switcher -->
        <RouterLink :class="['hidden items-center gap-2.5 rounded-lg md:flex', focusOnDark]" to="/">
            <span class="flex size-9 items-center justify-center overflow-hidden rounded-[9px] bg-white p-0.5">
                <img alt="" class="size-full object-contain" src="@/assets/grooster-logo.png" />
            </span>
            <span class="font-display text-lg font-semibold tracking-[-0.01em]">GRooster</span>
        </RouterLink>

        <nav v-if="navItems.length" aria-label="Main" class="hidden gap-1 md:flex">
            <RouterLink
                v-for="item in navItems"
                :key="item.to"
                v-slot="{ href, navigate, isExactActive }"
                custom
                :to="item.to"
            >
                <a
                    :aria-current="isExactActive ? 'page' : undefined"
                    :class="[
                        'inline-flex h-10 items-center gap-2 rounded-lg px-3.5 text-sm font-medium',
                        isExactActive
                            ? 'bg-topbar-active text-white'
                            : 'text-topbar-fg hover:bg-topbar-hover hover:text-white',
                        focusOnDark,
                    ]"
                    :href="href"
                    @click="navigate"
                >
                    <component :is="item.icon" aria-hidden="true" :size="18" />
                    {{ item.label }}
                </a>
            </RouterLink>
        </nav>

        <div class="hidden grow md:block" />

        <button
            v-if="organs.length"
            aria-controls="organization-menu"
            aria-haspopup="menu"
            :class="[
                'hidden h-10 items-center gap-2.5 rounded-lg border border-topbar-border bg-topbar-field px-3 hover:bg-topbar-hover md:flex',
                focusOnDark,
            ]"
            type="button"
            @click="organMenu?.toggle($event)"
        >
            <span class="text-xs text-topbar-muted">Organization</span>
            <span class="text-sm font-semibold">{{ organStore.organName ?? 'Choose' }}</span>
            <OiChevronDown aria-hidden="true" :size="16" />
        </button>
        <Menu id="organization-menu" ref="organMenu" :model="organMenuItems" popup>
            <template #item="{ item, props }">
                <a v-bind="props.action" :aria-current="item.current ? 'true' : undefined" class="gap-2">
                    <span class="grow">{{ item.label }}</span>
                    <OiCheck v-if="item.current" aria-hidden="true" class="text-primary" :size="16" />
                </a>
            </template>
        </Menu>

        <button
            :aria-label="themeLabel"
            :class="[
                'hidden size-10 items-center justify-center rounded-lg text-topbar-icon hover:bg-topbar-hover md:flex',
                focusOnDark,
            ]"
            :title="themeLabel"
            type="button"
            @click="themeStore.cycleMode()"
        >
            <component :is="themeIcons[themeStore.mode]" aria-hidden="true" :size="20" />
        </button>

        <!-- Pages teleport their phone actions (e.g. "⋯") here -->
        <div id="topbar-actions" class="flex items-center md:hidden" />
    </header>

    <Drawer
        v-model:visible="isDrawerOpen"
        header="GRooster"
        position="left"
        :pt="{
            root: { class: 'w-[min(20rem,85vw)]' },
            title: { class: 'font-display text-lg font-semibold' },
            pcCloseButton: { root: { class: 'size-11' } },
        }"
    >
        <nav v-if="navItems.length" aria-label="Main" class="flex flex-col gap-1">
            <RouterLink
                v-for="item in navItems"
                :key="item.to"
                v-slot="{ href, navigate, isExactActive }"
                custom
                :to="item.to"
            >
                <a
                    :aria-current="isExactActive ? 'page' : undefined"
                    :class="[
                        'flex h-12 items-center gap-3 rounded-lg px-3 text-[15px] font-medium',
                        isExactActive ? 'bg-highlight' : 'text-color hover:bg-emphasis',
                        focusOnLight,
                    ]"
                    :href="href"
                    @click="navigate"
                >
                    <component :is="item.icon" aria-hidden="true" :size="20" />
                    {{ item.label }}
                </a>
            </RouterLink>
        </nav>

        <section v-if="organs.length" aria-labelledby="drawer-organizations" class="mt-6 flex flex-col gap-1">
            <h2 id="drawer-organizations" class="px-3 pb-1 text-xs font-semibold uppercase text-muted-color">
                Organization
            </h2>
            <button
                v-for="item in organMenuItems"
                :key="item.label"
                :aria-current="item.current ? 'true' : undefined"
                :class="[
                    'flex h-12 items-center gap-3 rounded-lg px-3 text-left text-[15px] font-medium',
                    item.current ? 'bg-highlight' : 'text-color hover:bg-emphasis',
                    focusOnLight,
                ]"
                type="button"
                @click="item.command()"
            >
                <span class="grow">{{ item.label }}</span>
                <OiCheck v-if="item.current" aria-hidden="true" :size="18" />
            </button>
        </section>

        <button
            :aria-label="themeLabel"
            :class="[
                'mt-6 flex h-12 w-full items-center gap-3 rounded-lg px-3 text-[15px] font-medium text-color hover:bg-emphasis',
                focusOnLight,
            ]"
            type="button"
            @click="themeStore.cycleMode()"
        >
            <component :is="themeIcons[themeStore.mode]" aria-hidden="true" :size="20" />
            Theme: {{ themeStore.mode }}
        </button>
    </Drawer>
</template>
