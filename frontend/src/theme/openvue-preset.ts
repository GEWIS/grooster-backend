// OpenVue theme preset for the roster app ("Refined ops" design).
// Registered in main.ts as `theme.preset`; dark mode follows the `.app-dark` class (see stores/theme.store.ts).

import { definePreset } from '@openvue/themes';
import Aura from '@openvue/themes/aura';

const RosterPreset = definePreset(Aura, {
    semantic: {
        primary: {
            50: '#eef7f2',
            100: '#d6ede0',
            200: '#b0dbc4',
            300: '#7fc39d',
            400: '#4ea676',
            500: '#2c8a5a',
            600: '#1d7349',
            700: '#175e3b',
            800: '#134b30',
            900: '#0f3d26',
            950: '#0e2a21',
        },
        focusRing: {
            width: '2px',
            style: 'solid',
            color: '{primary.color}',
            offset: '2px',
        },
        formField: {
            borderRadius: '8px',
            focusRing: {
                width: '2px',
                style: 'solid',
                color: '{primary.color}',
                offset: '1px',
            },
        },
        colorScheme: {
            light: {
                // Green-tinted neutrals instead of zinc/slate.
                surface: {
                    0: '#ffffff',
                    50: '#f5f6f4', // page background
                    100: '#eef1ef', // subtle fills, "OK" chip
                    200: '#e2e6e3', // card borders, row dividers
                    300: '#d3d9d5', // input / secondary button borders
                    400: '#b9c3be', // dashed "add" borders
                    500: '#66716c', // muted text (5.1:1 on white)
                    600: '#56635d', // secondary text
                    700: '#3f4a45', // labels
                    800: '#2a3530',
                    900: '#15211c', // body text, "assigned" outline
                    950: '#0c1612',
                },
                primary: {
                    color: '{primary.600}',
                    contrastColor: '#ffffff',
                    hoverColor: '{primary.700}',
                    activeColor: '{primary.800}',
                },
                // HANDOFF §2: body text surface-900, muted surface-500 (Aura defaults to surface-700).
                text: {
                    color: '{surface.900}',
                    hoverColor: '{surface.950}',
                    mutedColor: '{surface.500}',
                    hoverMutedColor: '{surface.600}',
                },
                highlight: {
                    background: '#e6f2eb',
                    focusBackground: '#d6ede0',
                    color: '{primary.900}',
                    focusColor: '{primary.950}',
                },
            },
            dark: {
                // Not designed yet. These are sensible starting values; review them in the app.
                surface: {
                    0: '#ffffff',
                    50: '#e8eeea',
                    100: '#c9d3cd',
                    200: '#a3b0a9',
                    300: '#7e8c85',
                    400: '#5f6c66',
                    500: '#47534d',
                    600: '#333e39',
                    700: '#26302b',
                    800: '#1b2420',
                    900: '#131b17',
                    950: '#0c1210',
                },
                // Aura's dark muted text (surface-400) is below 4.5:1 on these surfaces.
                text: {
                    mutedColor: '{surface.200}',
                    hoverMutedColor: '{surface.100}',
                },
                primary: {
                    color: '{primary.400}',
                    contrastColor: '{surface.950}',
                    hoverColor: '{primary.300}',
                    activeColor: '{primary.200}',
                },
            },
        },
    },
    components: {
        button: {
            root: {
                label: { fontWeight: '600' },
            },
            colorScheme: {
                light: {
                    // Secondary action: white with a surface-300 border and body text (HANDOFF §2)
                    outlined: {
                        secondary: {
                            borderColor: '{surface.300}',
                            color: '{surface.900}',
                            hoverBackground: '{surface.100}',
                            activeBackground: '{surface.200}',
                        },
                    },
                },
                dark: {
                    outlined: {
                        secondary: {
                            borderColor: '{surface.600}',
                            color: '{surface.0}',
                            hoverBackground: '{surface.800}',
                            activeBackground: '{surface.700}',
                        },
                    },
                },
            },
        },
    },
});

export default RosterPreset;
