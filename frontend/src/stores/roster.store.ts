import { defineStore } from 'pinia'
import {
    RosterAnswerUpdateRequest,
    RosterShiftCreateRequest,
    SavedShift,
    SavedShiftUpdateRequest, RosterCreateRequest, Roster, RosterAnswerCreateRequest
} from "@gewis/grooster-backend-ts";
import ApiService from "@/services/ApiService";

export const useRosterStore = defineStore('roster', {
    state: () => ({
        rosters: {} as Record<number, Roster>,
        savedRoster: {} as Record<number, SavedShift[]>,
    }),
    getters: {
        getRoster: (state) => {
            return (rosterId: number) => state.rosters[rosterId];
        },
        getAllRosters: (state) => {
            return state.rosters;
        },
        getSavedRoster: (state) => {
            return (rosterId: number) => state.savedRoster[rosterId];
        }
    },
    actions: {
        async fetchRosters(organId: number) {
            await ApiService.roster.getRosters(undefined, organId).then((res) => {
                res.data.forEach((roster: Roster) => {
                    this.rosters[roster.id] = roster;
                });
                return this.rosters;
            });
        },
        async createRoster(params: RosterCreateRequest) {
            await ApiService.roster.createRoster(params).then((res) => {
                const roster = res.data;
                this.rosters[roster.id] = roster;
            });
        },
        async createAnswer(params: RosterAnswerCreateRequest) {
            await ApiService.rosterAnswer.createRosterAnswer(params).then((res) => {
                const roster = this.rosters[res.data.rosterId];

                const existingAnswers = roster.rosterAnswer ?? [];

                const updatedRoster: Roster = {
                    ...roster,
                    rosterAnswer: [...existingAnswers, res.data],
                }
                this.rosters = {
                    ...this.rosters,
                    [roster.id]: updatedRoster,
                };
            });
        },
        async updateAnswer(id: number, params: RosterAnswerUpdateRequest) {
            await ApiService.rosterAnswer.updateRosterAnswer(id, params).then((res) => {
                console.log(res);
                const roster = this.rosters[res.data.rosterId];

                const answers = roster.rosterAnswer ?? [];
                const answerIndex = answers.findIndex(answer => answer.id === id);

                if (answerIndex !== -1) {
                    const updatedAnswers = [
                        ...answers.slice(0, answerIndex),
                        res.data,
                        ...answers.slice(answerIndex + 1),
                    ];

                    const updatedRoster = {
                        ...roster,
                        rosterAnswer: updatedAnswers,
                    };

                    this.rosters = {
                        ...this.rosters,
                        [roster.id]: updatedRoster,
                    };
                }
            });
        },
        async createShift(params: RosterShiftCreateRequest) {
            await ApiService.rosterShift.createRosterShift(params).then((res) => {
                const roster = this.rosters[res.data.rosterId];

                const existingShifts = roster.rosterShift ?? [];

                const updatedRoster = {
                    ...roster,
                    rosterShift: [...existingShifts, res.data],
                }

                this.rosters = {
                    ...this.rosters,
                    [roster.id]: updatedRoster,
                };
            });
        },
        async deleteShift(shiftId: number, rosterId: number) {
            await ApiService.rosterShift.deleteRosterShift(shiftId).then(() => {
                const roster = this.rosters[rosterId];
                if (!roster) return;

                const updatedRosterShift = roster.rosterShift.filter(shift => shift.id !== shiftId);
                const updatedRosterAnswer = roster.rosterAnswer.filter(answer => answer.rosterShiftId !== shiftId);

                const updatedRoster = {
                    ...roster,
                    rosterShift: updatedRosterShift,
                    rosterAnswer: updatedRosterAnswer,
                };

                this.rosters = {
                    ...this.rosters,
                    [rosterId]: updatedRoster,
                };
            });
        },
        async saveRoster(id: number) {
            await ApiService.savedShiftApi.rosterSave(id).then(async () => {
                const roster: Roster = this.rosters[id];
                const updatedRoster = {
                    ...roster,
                    saved: true
                }

                this.rosters = {
                    ...this.rosters,
                    [id]: updatedRoster,
                };
            });
        },
        async fetchSavedShifts(id: number) {
            await ApiService.savedShiftApi.getSavedRoster(id).then((res) => {
                this.savedRoster[id] = res.data;
            })
        },
        async updateSavedShift(shiftId: number, params: SavedShiftUpdateRequest) {
            await ApiService.savedShiftApi.updateSavedShift(shiftId, params).then((res) => {
                const updatedShift = res.data;

                const shifts = this.savedRoster[updatedShift.rosterId];
                if (!shifts) return;

                const index = shifts.findIndex((s) => s.id === shiftId);
                if (index !== -1) {
                    this.savedRoster[updatedShift.rosterId].splice(index, 1, updatedShift);
                }
            });
        }
    }
})