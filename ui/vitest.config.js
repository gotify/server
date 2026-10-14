import {defineConfig} from 'vitest/config';
import babel from '@rolldown/plugin-babel';

const timeout = globalThis.process.env.CI === 'true' ? 60000 : 30000;

// The UI source uses TC39 stage-3 decorators (`@observable`/`@action` and the
// `accessor` keyword).  The production bundle compiles them via the babel
// plugin configured in vite.config.ts, but vitest's transform pipeline needs
// the same plugin in order to import those modules from tests.
export default defineConfig({
    plugins: [
        babel({
            presets: [
                {
                    preset: () => ({
                        plugins: [['@babel/plugin-proposal-decorators', {version: '2022-03'}]],
                    }),
                    rolldown: {filter: {code: '@'}},
                },
            ],
        }),
    ],
    test: {
        setupFiles: ['./src/tests/jsdom-shim.ts'],
        testTimeout: timeout,
        hookTimeout: timeout,
    },
});
