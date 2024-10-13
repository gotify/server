import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react-swc';

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), "")

    return {
        build: {
            outDir: 'build',
            emptyOutDir: true,
            sourcemap: true,
        },
        plugins: [ react({ tsDecorators: true }) ],
        define: {
            // Some libraries use the global object, even though it doesn't exist in the browser.
            // Alternatively, we could add `<script>window.global = window;</script>` to index.html.
            // https://github.com/vitejs/vite/discussions/5912
            global: {},
        },
        server: {
            proxy: {
                "/api": {
                    target: 'http://localhost:3000/',
                    changeOrigin: true,
                    secure: false,
                    rewrite: (p) => p.replace(/^\/api/, ""),
                },
            },
            cors: false,
        },
        preview: {
            proxy: {
                "/api": {
                    target: 'http://localhost:3000/',
                    changeOrigin: true,
                    secure: false,
                    rewrite: (p) => p.replace(/^\/api/, ""),
                },
            },
            cors: false,
        },
    };
});
