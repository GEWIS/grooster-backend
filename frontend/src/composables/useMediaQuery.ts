import { onBeforeUnmount, ref } from 'vue';

/** Reactive `matchMedia`, e.g. `useMediaQuery('(min-width: 1024px)')` for Tailwind's `lg`. */
export function useMediaQuery(query: string) {
    const media = window.matchMedia(query);
    const matches = ref(media.matches);
    const update = (event: MediaQueryListEvent) => (matches.value = event.matches);

    media.addEventListener('change', update);
    onBeforeUnmount(() => media.removeEventListener('change', update));

    return matches;
}
