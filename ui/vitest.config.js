import {defineConfig} from 'vitest/config';

const timeout = process.env.CI === 'true' ? 60000 : 30000;

export default defineConfig({
    test: {
        testTimeout: timeout,
        hookTimeout: timeout,
    },
});
