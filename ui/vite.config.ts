import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';
import babel from '@rolldown/plugin-babel';

const GOTIFY_SERVER_PORT = process.env.GOTIFY_SERVER_PORT ?? '80';

function decoratorPreset(options: Record<string, unknown>) {
    return {
        preset: () => ({
            plugins: [['@babel/plugin-proposal-decorators', options]],
        }),
        rolldown: {
            filter: {code: '@'},
        },
    };
}

export default defineConfig({
    base: './',
    build: {
        outDir: 'build',
        emptyOutDir: true,
        sourcemap: false,
        assetsDir: 'static',
    },
    plugins: [
        react(),
        babel({
            presets: [decoratorPreset({version: '2022-03'})],
        }),
    ],
    define: {
        // Some libraries use the global object, even though it doesn't exist in the browser.
        // Alternatively, we could add `<script>window.global = window;</script>` to index.html.
        // https://github.com/vitejs/vite/discussions/5912
        global: {},
    },
    server: {
        host: '0.0.0.0',
        proxy: {
            '^/(application|message|client|current|user|plugin|version|image|auth)': {
                target: `http://localhost:${GOTIFY_SERVER_PORT}/`,
                changeOrigin: true,
                secure: false,
            },
            '/stream': {
                target: `ws://localhost:${GOTIFY_SERVER_PORT}/`,
                ws: true,
                rewriteWsOrigin: true,
            },
        },
        cors: false,
    },
});
