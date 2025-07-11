import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';

import Button from "primevue/button";
import Card from "primevue/card";
import Select from "primevue/select";
import Menubar from 'primevue/menubar';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import ListBox from 'primevue/listbox';
import DatePicker from 'primevue/datepicker';
import router from "@/router/router";

const app = createApp(App);
const pinia = createPinia()

app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            prefix: 'p',
            darkModeSelector: false,
            cssLayer: false
        }
    }
});

// Components
app.component('Button', Button);
app.component('Card', Card);
app.component('Select', Select);
app.component('Menubar', Menubar);
app.component('Dialog', Dialog);
app.component('InputText', InputText);
app.component('ListBox', ListBox);
app.component('DatePicker', DatePicker);

app.use(pinia);
app.use(router);
app.mount('#app');
