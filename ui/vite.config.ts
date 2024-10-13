import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    build: {
        outDir: 'build',
        emptyOutDir: true,
        sourcemap: true,
    },
    plugins: [
        react({
            babel: {
                parserOpts: {
                    plugins: [ 'decorators-legacy' ],
                }
            }
        })
    ],
    define: {
        // Some libraries use the global object, even though it doesn't exist in the browser.
        // Alternatively, we could add `<script>window.global = window;</script>` to index.html.
        // https://github.com/vitejs/vite/discussions/5912
        global: {},
    },
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:3000/',
                changeOrigin: true,
                secure: false,
                rewrite: (p) => p.replace(/^\/api/, ''),
            },
            '/api/stream': {
                target: 'ws://localhost:3000/',
                ws: true,
                rewrite: (p) => p.replace(/^\/api/, ''),
                rewriteWsOrigin: true,
            }
        },
        cors: false,
    },
    preview: {
        proxy: {
            '/api': {
                target: 'http://localhost:3000/',
                changeOrigin: true,
                secure: false,
                rewrite: (p) => p.replace(/^\/api/, ''),
            },
            '/api/stream': {
                target: 'ws://localhost:3000/',
                ws: true,
                rewrite: (p) => p.replace(/^\/api/, ''),
                rewriteWsOrigin: true,
            }
        },
        cors: false,
    },

});
