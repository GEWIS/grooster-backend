import { defineStore } from 'pinia';
import type {
  RosterAnswerUpdateRequest,
  RosterShiftCreateRequest,
  SavedShiftUpdateRequest,
  RosterCreateRequest,
  Roster,
  RosterAnswerCreateRequest,
  SavedShiftResponse,
  RosterUpdateRequest,
  RosterShiftUpdateRequest,
} from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';

export const useRosterStore = defineStore('roster', {
  state: () => ({
    rosters: {} as Record<number, Roster>,
    savedRoster: {} as Record<number, SavedShiftResponse>,
    selectedRosterId: null as number | null,
  }),
  persist: true,
  getters: {
    getRoster: (state) => {
      return (rosterId: number) => state.rosters[rosterId];
    },
    getAllRosters: (state) => {
      return state.rosters;
    },
    getSavedRoster: (state) => {
      return (rosterId: number) => state.savedRoster[rosterId];
    },
    selectedRoster: (state) => {
      return state.selectedRosterId ? state.rosters[state.selectedRosterId] : null;
    },
  },
  actions: {
    setSelectedRoster(id: number | null) {
      this.selectedRosterId = id;
    },
    clearRosters() {
      this.rosters = {};
      this.savedRoster = {};
      this.selectedRosterId = null;
    },
    async fetchRosters(organId: number) {
      await ApiService.roster.getRosters(undefined, organId).then((res) => {
        res.data.forEach((roster: Roster) => {
          if (roster.id != null) {
            this.rosters[roster.id] = roster;
          }
        });
        return this.rosters;
      });
    },
    async createRoster(params: RosterCreateRequest) {
      const res = await ApiService.roster.createRoster(params);
      const roster = res.data;
      if (roster.id != null) {
        this.rosters[roster.id] = roster;
      }
      return roster;
    },
    async updateRoster(id: number, params: RosterUpdateRequest) {
      try {
        await ApiService.roster.updateRoster(id, params);
        const currentRoster = this.rosters[id];

        // TODO Make sure backend also returns preloaded relations
        if (currentRoster) {
          this.rosters = {
            ...this.rosters,
            [id]: {
              ...currentRoster,
              ...params,
            },
          };
        }
      } catch (error) {
        console.error('Failed to update an roster:', error);
      }
    },
    async deleteRoster(id: number) {
      try {
        await ApiService.roster.deleteRoster(id);

        delete this.rosters[id];

        this.selectedRosterId = null;
      } catch (error) {
        console.error('Failed to delete an roster:', error);
      }
    },
    async createAnswer(params: RosterAnswerCreateRequest) {
      try {
        const response = await ApiService.rosterAnswer.createRosterAnswer(params);
        const rosterId = response.data?.rosterId;
        const roster = rosterId && this.rosters[rosterId];

        if (roster) {
          this.rosters[rosterId] = {
            ...roster,
            rosterAnswer: [...(roster.rosterAnswer ?? []), response.data],
          };
        }

        return response.data;
      } catch (error) {
        console.error('Failed to create roster answer:', error);
      }
    },
    async updateAnswer(id: number, params: RosterAnswerUpdateRequest) {
      try {
        const response = await ApiService.rosterAnswer.updateRosterAnswer(id, params);
        const rosterId = response.data?.rosterId;

        if (!rosterId || !this.rosters[rosterId]) {
          console.warn('Invalid roster ID or roster not found');
          return;
        }

        const roster = this.rosters[rosterId];
        const answers = roster.rosterAnswer ?? [];
        const answerIndex = answers.findIndex((answer) => answer.id === id);

        if (answerIndex !== -1) {
          const updatedAnswers = [...answers.slice(0, answerIndex), response.data, ...answers.slice(answerIndex + 1)];

          this.rosters = {
            ...this.rosters,
            [rosterId]: {
              ...roster,
              rosterAnswer: updatedAnswers,
            },
          };
        }
      } catch (error) {
        console.error('Failed to update roster answer:', error);
      }
    },
    async createShift(params: RosterShiftCreateRequest) {
      await ApiService.rosterShift.createRosterShift(params).then((res) => {
        const roster = this.rosters[res.data.rosterId];

        const existingShifts = roster.rosterShift ?? [];

        const updatedRoster = {
          ...roster,
          rosterShift: [...existingShifts, res.data],
        };

        this.rosters = {
          ...this.rosters,
          [roster.id]: updatedRoster,
        };
      });
    },
    async updateShift(shiftId: number, params: RosterShiftUpdateRequest) {
      try {
        await ApiService.rosterShift.updateRosterShift(shiftId, params);

        const rosters: Roster[] = Object.values(this.rosters);
        const roster = rosters.find((r) => r.rosterShift?.some((s) => s.id === shiftId));

        if (roster && roster.rosterShift) {
          const shiftIndex = roster.rosterShift.findIndex((s) => s.id === shiftId);
          if (shiftIndex !== -1) {
            roster.rosterShift[shiftIndex] = {
              ...roster.rosterShift[shiftIndex],
              ...params,
            };

            roster.rosterShift.sort((a, b) => {
              if (a.order !== b.order) return a.order - b.order;
              return a.id - b.id;
            });
          }
        }
      } catch (error) {
        console.error('Failed to update shift in store:', error);
      }
    },
    async deleteShift(shiftId: number, rosterId: number) {
      await ApiService.rosterShift.deleteRosterShift(shiftId).then(() => {
        const roster = this.rosters[rosterId];
        if (!roster) return;

        const updatedRosterShift = roster.rosterShift.filter((shift) => shift.id !== shiftId);
        const updatedRosterAnswer = roster.rosterAnswer.filter((answer) => answer.rosterShiftId !== shiftId);

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
      await ApiService.savedShiftApi.rosterSave(id).then(() => {
        const roster: Roster = this.rosters[id];
        const updatedRoster = {
          ...roster,
          saved: true,
        };

        this.rosters = {
          ...this.rosters,
          [id]: updatedRoster,
        };
      });
    },
    async fetchSavedShifts(id: number) {
      try {
        const res = await ApiService.savedShiftApi.getSavedRoster(id);
        this.savedRoster = {
          ...this.savedRoster,
          [id]: res.data,
        };
      } catch (error) {
        console.error('Fetch failed', error);
      }
    },
    async updateSavedShift(shiftId: number, params: SavedShiftUpdateRequest) {
      await ApiService.savedShiftApi.updateSavedShift(shiftId, params).then((res) => {
        const updatedShift = res.data;

        const shiftResponses = this.savedRoster[updatedShift.rosterId];
        if (!shiftResponses) return;

        const shifts = shiftResponses.savedShifts;

        const index = shifts.findIndex((s) => s.id === shiftId);
        if (index !== -1) {
          this.savedRoster[updatedShift.rosterId].savedShifts.splice(index, 1, updatedShift);
        }
      });
    },
  },
});
