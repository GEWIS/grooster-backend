import { computed, reactive, ref, watch, type Ref } from 'vue';
import { useToast } from 'openvue/usetoast';
import type { RosterAnswer, RosterShift, User } from '@gewis/grooster-backend-ts';
import { useRosterStore } from '@/stores/roster.store';
import { Role, useAuthStore } from '@/stores/auth.store';
import { getGEWISId } from '@/helpers/TokenHelper';
import ApiService from '@/services/ApiService';

export interface RosterMember {
    id: number;
    /** Full name from the user account */
    name: string;
    /** Per-organ nickname if set, otherwise the full name */
    displayName: string;
    isMe: boolean;
}

/**
 * Shared state for one roster page: members, shifts, answers, comments, shift management,
 * lock/fill and (once locked) assignments. Used by the desktop grid, the phone list,
 * the assign panel and "My preferences", so they never duplicate logic.
 */
export function useRosterGrid(rosterId: Ref<number>, organId: Ref<number>) {
    const rosterStore = useRosterStore();
    const authStore = useAuthStore();
    const toast = useToast();

    const roster = computed(() => rosterStore.getRoster(rosterId.value));
    const isAdmin = computed(() => authStore.can([Role.Admin, Role.Owner]));
    const isLocked = computed(() => !!roster.value?.saved);
    const values = computed<string[]>(() => Object.values(roster.value?.values ?? {}));
    const shifts = computed<RosterShift[]>(() =>
        [...(roster.value?.rosterShift ?? [])].sort((a, b) => a.order - b.order || a.id - b.id),
    );

    function notifyError(summary: string, error: unknown) {
        console.error(summary, error);
        toast.add({ severity: 'error', summary, detail: 'Please try again.', life: 5000 });
    }

    // --- Members -----------------------------------------------------------------------------
    const members = ref<RosterMember[]>([]);
    const membersById = computed(() => new Map(members.value.map((member) => [member.id, member])));

    async function loadMembers() {
        try {
            const [userResponse, settingsResponse] = await Promise.all([
                ApiService.user.userGet(organId.value),
                ApiService.organ.getMembersSettings(organId.value),
            ]);
            const nicknames = new Map(settingsResponse.data.map((setting) => [setting.userId, setting.username]));
            const myId = getGEWISId();

            const list = userResponse.data.map((user) => ({
                id: user.id,
                name: user.name,
                displayName: nicknames.get(user.id) || user.name,
                isMe: user.gewis_id === myId,
            }));
            // Current user first, otherwise keep the backend order
            members.value = [...list.filter((m) => m.isMe), ...list.filter((m) => !m.isMe)];
        } catch (error) {
            notifyError('Failed to load members', error);
        }
    }

    const me = computed(() => members.value.find((member) => member.isMe));
    const toMember = (user: User): RosterMember =>
        membersById.value.get(user.id) ?? { id: user.id, name: user.name, displayName: user.name, isMe: false };

    // --- Answers -----------------------------------------------------------------------------
    const key = (userId: number, shiftId: number) => `${userId}:${shiftId}`;
    const answers = computed(() => {
        const map = new Map<string, RosterAnswer>();
        for (const answer of roster.value?.rosterAnswer ?? [])
            map.set(key(answer.userId, answer.rosterShiftId), answer);
        return map;
    });
    // Optimistic values while a request is in flight, so the label follows the select immediately
    const pendingAnswers = reactive(new Map<string, string>());

    const answerValue = (userId: number, shiftId: number): string | null =>
        pendingAnswers.get(key(userId, shiftId)) ?? answers.value.get(key(userId, shiftId))?.value ?? null;

    const canEdit = (member: RosterMember) => !isLocked.value && member.isMe;

    async function setAnswer(userId: number, shiftId: number, value: string) {
        const k = key(userId, shiftId);
        const existing = answers.value.get(k);
        pendingAnswers.set(k, value);
        try {
            if (existing?.id) {
                await rosterStore.updateAnswer(existing.id, { value });
            } else {
                await rosterStore.createAnswer({ rosterId: rosterId.value, rosterShiftId: shiftId, userId, value });
            }
        } catch (error) {
            notifyError('Could not save your answer', error);
        } finally {
            pendingAnswers.delete(k);
        }
    }

    // --- Comments ----------------------------------------------------------------------------
    const commentFor = (userId: number) =>
        rosterStore.getComments(rosterId.value).find((comment) => comment.userId === userId)?.comment ?? '';

    async function saveComment(userId: number, comment: string) {
        await rosterStore.updateComment({ rosterId: rosterId.value, userId, comment });
    }

    // --- Shifts (admin, while open) ----------------------------------------------------------
    async function addShift(name: string) {
        try {
            await rosterStore.createShift({ name, rosterId: rosterId.value });
        } catch (error) {
            notifyError('Could not add the shift', error);
        }
    }

    async function removeShift(shiftId: number) {
        try {
            await rosterStore.deleteShift(shiftId, rosterId.value);
        } catch (error) {
            notifyError('Could not delete the shift', error);
        }
    }

    async function moveShift(shiftId: number, direction: 'left' | 'right') {
        const list = shifts.value;
        const index = list.findIndex((shift) => shift.id === shiftId);
        const target = list[direction === 'left' ? index - 1 : index + 1];
        const current = list[index];
        if (!current || !target) return;

        await Promise.all([
            rosterStore.updateShift(current.id, { order: target.order }),
            rosterStore.updateShift(target.id, { order: current.order }),
        ]);
    }

    // --- Lock / fill -------------------------------------------------------------------------
    async function lock() {
        try {
            await rosterStore.saveRoster(rosterId.value);
        } catch (error) {
            notifyError('Could not lock the roster', error);
        }
    }

    async function unlock() {
        await rosterStore.updateRoster(rosterId.value, { saved: false });
    }

    async function autoFill() {
        await rosterStore.fillRoster(rosterId.value);
    }

    // --- Assignments (only exist once the roster is locked) ----------------------------------
    const savedRoster = computed(() => rosterStore.getSavedRoster(rosterId.value));
    const savedShiftFor = (shiftId: number) =>
        savedRoster.value?.savedShifts?.find((saved) => saved.rosterShiftId === shiftId);

    const assignedMembers = (shiftId: number): RosterMember[] => (savedShiftFor(shiftId)?.users ?? []).map(toMember);
    const isAssigned = (userId: number, shiftId: number) =>
        !!savedShiftFor(shiftId)?.users?.some((user) => user.id === userId);

    /** Everyone not yet assigned, in the backend's ranked order (best match first). */
    function candidates(shiftId: number): RosterMember[] {
        const shift = shifts.value.find((s) => s.id === shiftId);
        const ranked = savedRoster.value?.savedShiftOrdering?.find((o) => o.shiftName === shift?.name)?.users;
        const pool = ranked ? ranked.map(toMember) : members.value;
        return pool.filter((member) => !isAssigned(member.id, shiftId));
    }

    const busyAssignments = reactive(new Set<string>());
    const isBusy = (userId: number, shiftId: number) => busyAssignments.has(key(userId, shiftId));

    async function setAssigned(userId: number, shiftId: number, assigned: boolean) {
        const saved = savedShiftFor(shiftId);
        if (!saved || isAssigned(userId, shiftId) === assigned) return;

        const ids = (saved.users ?? []).map((user) => user.id);
        const users = assigned ? [...ids, userId] : ids.filter((id) => id !== userId);
        busyAssignments.add(key(userId, shiftId));
        try {
            await rosterStore.updateSavedShift(saved.id, { users });
        } catch (error) {
            notifyError(assigned ? 'Could not assign' : 'Could not remove the assignment', error);
        } finally {
            busyAssignments.delete(key(userId, shiftId));
        }
    }

    // --- Loading -----------------------------------------------------------------------------
    watch(
        organId,
        (id) => {
            if (id) void loadMembers();
        },
        { immediate: true },
    );
    watch(
        rosterId,
        (id) => {
            if (id) void rosterStore.fetchComments(id);
        },
        { immediate: true },
    );
    watch(
        [rosterId, isLocked],
        ([id, locked]) => {
            if (id && locked) void rosterStore.fetchSavedShifts(id);
        },
        { immediate: true },
    );

    return {
        roster,
        isAdmin,
        isLocked,
        values,
        shifts,
        members,
        me,
        answerValue,
        canEdit,
        setAnswer,
        commentFor,
        saveComment,
        addShift,
        removeShift,
        moveShift,
        lock,
        unlock,
        autoFill,
        assignedMembers,
        candidates,
        isAssigned,
        isBusy,
        setAssigned,
    };
}

export type RosterGrid = ReturnType<typeof useRosterGrid>;
