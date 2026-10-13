<script setup lang="ts">
import { ref } from 'vue';
import Popover from 'openvue/popover';
import type { RosterGrid, RosterMember } from '@/composables/useRosterGrid';

const props = defineProps<{ grid: RosterGrid }>();

const MAX_LENGTH = 250;
const popover = ref<InstanceType<typeof Popover>>();
const member = ref<RosterMember | null>(null);
const draft = ref('');
const saving = ref(false);

function open(event: Event, target: RosterMember) {
    member.value = target;
    draft.value = props.grid.commentFor(target.id);
    popover.value?.show(event);
}

async function save() {
    if (!member.value) return;
    saving.value = true;
    try {
        await props.grid.saveComment(member.value.id, draft.value);
        popover.value?.hide();
    } finally {
        saving.value = false;
    }
}

defineExpose({ open });
</script>

<template>
    <Popover ref="popover" :pt="{ content: { class: 'w-72 p-4' } }">
        <div v-if="member" class="flex flex-col gap-2">
            <span class="text-sm font-semibold">{{
                member.isMe ? 'Your comment' : `Comment by ${member.displayName}`
            }}</span>
            <template v-if="member.isMe">
                <Textarea
                    v-model="draft"
                    aria-label="Your comment"
                    auto-resize
                    class="w-full text-sm"
                    :maxlength="MAX_LENGTH"
                    placeholder="E.g. only available after 18:00"
                    rows="3"
                />
                <div class="flex items-center gap-2">
                    <span aria-live="polite" class="grow text-xs text-surface-600 dark:text-surface-200">
                        {{ draft.length }} / {{ MAX_LENGTH }}
                    </span>
                    <Button class="h-9 text-sm" label="Save" :loading="saving" @click="save" />
                </div>
            </template>
            <p v-else class="whitespace-pre-wrap text-sm text-surface-600 dark:text-surface-200">
                {{ grid.commentFor(member.id) || 'No comment yet.' }}
            </p>
        </div>
    </Popover>
</template>
