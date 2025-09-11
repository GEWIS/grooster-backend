import {defineStore} from "pinia";
import {RosterTemplate} from "@gewis/grooster-backend-ts";
import ApiService from "@/services/ApiService";

export const useTemplateStore = defineStore('rosterTemplate', {
    state: () => ({
        templates: {} as Record<number, RosterTemplate>
    }),
    getters: {
        getTemplate: (state) => {
            return (templateId: number) => state.templates[templateId];
        }  ,
        getAllTemplates: (state) => {
            return state.templates;
        }
    },
    actions: {
        async fetchTemplates(organId: number) {
            try {
                const response = await ApiService.roster.getRosterTemplates(organId);
                const templates = response.data;

                if (templates && templates.length > 0) {
                    templates.forEach(template => {
                        this.templates[template.id] = template;
                    })
                }
            } catch (error) {
                console.error('Failed to fetch templates', error);
            }
        },
        async createTemplate(params: RosterTemplate) {
            try {
                const response = await ApiService.roster.createRosterTemplate(params);
                const template = response.data as RosterTemplate;

                this.templates[template.id] = template;
            } catch (error) {
                console.error('Failed to create template', error);
            }
        },
        async updateTemplate(params: RosterTemplate) {},
        async deleteTemplate(templateId: number) {
            try {
                await ApiService.roster.deleteRosterTemplate(templateId);

                delete this.templates[templateId];
            } catch (error) {
                console.error('Failed to delete template', error);
            }
        },
    }
})