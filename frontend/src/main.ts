/* eslint vue/multi-word-component-names: 0 */
/* eslint vue/no-reserved-component-names: 0 */

import '@fontsource/fira-mono/400.css';
import '@fontsource/fira-mono/500.css';
import '@fontsource/fira-mono/700.css';
import '@fontsource-variable/bricolage-grotesque/opsz.css';
import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';
import OpenVue from 'openvue/config';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';

import Button from 'openvue/button';
import Card from 'openvue/card';
import Select from 'openvue/select';
import Menubar from 'openvue/menubar';
import Dialog from 'openvue/dialog';
import InputText from 'openvue/inputtext';
import ListBox from 'openvue/listbox';
import DatePicker from 'openvue/datepicker';
import { Form } from '@openvue/forms';
import Chip from 'openvue/chip';
import { Divider } from 'openvue';
import Popover from 'openvue/popover';
import Textarea from 'openvue/textarea';
import ToastService from 'openvue/toastservice';
import ConfirmationService from 'openvue/confirmationservice';
import Toast from 'openvue/toast';
import App from './App.vue';
import RosterPreset from '@/theme/openvue-preset';
import router from '@/router/router';
import { useThemeStore } from '@/stores/theme.store';

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

app.use(OpenVue, {
    theme: {
        preset: RosterPreset,
        options: {
            prefix: 'p',
            darkModeSelector: '.app-dark',
            cssLayer: { name: 'primevue', order: 'theme, base, primevue' },
        },
    },
});
app.use(ToastService);
app.use(ConfirmationService);

// Components
app.component('Button', Button);
app.component('Card', Card);
app.component('Select', Select);
app.component('Menubar', Menubar);
app.component('Dialog', Dialog);
app.component('InputText', InputText);
app.component('ListBox', ListBox);
app.component('DatePicker', DatePicker);
app.component('Form', Form);
app.component('Chip', Chip);
app.component('Divider', Divider);
app.component('Popover', Popover);
app.component('Textarea', Textarea);
app.component('Toast', Toast);

app.use(pinia);
useThemeStore().init();
app.use(router);
app.mount('#app');
