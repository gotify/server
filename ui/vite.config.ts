import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { writeFileSync } from 'fs';
import { resolve } from 'path';
import type { OutputBundle, OutputChunk } from 'rollup';

export default defineConfig({
    build: {
        outDir: 'build',
        assetsDir: 'static',
        emptyOutDir: true,
        sourcemap: 'inline',
        rollupOptions: {
            plugins: [
                {
                    name: 'generate-asset-manifest',
                    generateBundle(_, bundle: OutputBundle) {
                        const files: Record<string, string> = {};
                        const entrypoints: string[] = [];
                        let entrypointCss: string = "";

                        for (const [ , value ] of Object.entries(bundle)) {
                            const fileName = value.fileName.split('/').pop() || '';
                            const [ baseFileName, ...extParts ] = fileName.split('.');
                            const extension = extParts.pop();
                            const cleanBaseFileName = baseFileName.replace(/-[\w-]{8}$/, '');

                            if ((value as OutputChunk).isEntry) {
                                entrypoints.push(value.fileName);
                                entrypointCss = `${cleanBaseFileName}.css`;
                            }

                            files[`${cleanBaseFileName}.${extension}`] = value.fileName;
                        }

                        if (entrypointCss in files) {
                            entrypoints.push(files[entrypointCss]);
                        }

                        const manifest = {
                            files,
                            entrypoints,
                        };

                        writeFileSync(resolve(__dirname, 'build/asset-manifest.json'), JSON.stringify(manifest, null, 2));
                    }
                }
            ]
        }
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
        //global: {},
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
