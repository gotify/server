import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';

try {
    process.loadEnvFile('../gotify-server.env');
} catch {
    // file is optional
}

const GOTIFY_SERVER_PORT = process.env.GOTIFY_SERVER_PORT ?? '80';

export default defineConfig({
    base: './',
    build: {
        outDir: 'build',
        emptyOutDir: true,
        sourcemap: false,
        assetsDir: 'static',
    },
    plugins: [react()],
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
                secure: false,
            },
            '/stream': {
                target: `http://localhost:${GOTIFY_SERVER_PORT}/`,
                ws: true,
            },
        },
        cors: false,
    },
});
