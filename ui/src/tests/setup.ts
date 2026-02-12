import getPort from 'get-port';
import {spawn, ChildProcess} from 'child_process';
import {rimrafSync} from 'rimraf';
import path from 'path';
import puppeteer, {Browser, Page} from 'puppeteer';
import fs from 'fs';
import kill from 'tree-kill';

export interface GotifyTest {
    url: string;
    close: () => Promise<void>;
    browser: Browser;
    page: Page;
}

const windowsPrefix = process.platform === 'win32' ? '.exe' : '';
const appDotGo = path.join(__dirname, '..', '..', '..', 'app.go');
const testBuildPath = path.join(__dirname, 'build');

export const newPluginDir = async (plugins: string[]): Promise<string> => {
    const {dir, generator} = testPluginDir();
    for (const pluginName of plugins) {
        await buildGoPlugin(generator(), pluginName);
    }
    return dir;
};

export const newTest = async (pluginsDir = ''): Promise<GotifyTest> => {
    const port = await getPort();

    const gotifyFile = testFilePath();

    await buildGoExecutable(gotifyFile);

    const gotifyInstance = startGotify(gotifyFile, port, pluginsDir);

    const gotifyURL = 'http://localhost:' + port;
    await waitForGotify(gotifyURL);
    const browser = await puppeteer.launch({
        headless: process.env.CI === 'true',
        args: [`--window-size=1920,1080`, '--no-sandbox'],
    });
    const page = await browser.newPage();
    await page.setViewport({width: 1920, height: 1080});
    await page.goto(gotifyURL);

    return {
        close: async () => {
            await Promise.all([
                browser.close(),
                new Promise((resolve) =>
                    kill(gotifyInstance.pid!, 'SIGKILL', () => resolve(undefined))
                ),
            ]);
            rimrafSync(gotifyFile, {maxRetries: 8});
        },
        url: gotifyURL,
        browser,
        page,
    };
};

const testPluginDir = (): {dir: string; generator: () => string} => {
    const random = Math.random().toString(36).substring(2, 15);
    const dirName = 'gotifyplugin_' + random;
    const dir = path.join(testBuildPath, dirName);
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, {recursive: true, mode: 0o755});
    }
    return {
        dir,
        generator: () => {
            const randomFn = Math.random().toString(36).substring(2, 15);
            return path.join(dir, randomFn + '.so');
        },
    };
};

const testFilePath = (): string => {
    const random = Math.random().toString(36).substring(2, 15);
    const filename = 'gotifytest_' + random + windowsPrefix;
    return path.join(testBuildPath, filename);
};

const waitForGotify = async (url: string): Promise<void> => {
    const deadline = Date.now() + 30000;
    let status = new Error('timeout');
    while (Date.now() < deadline) {
        const abc = new AbortController();
        const timeout = setTimeout(() => {
            abc.abort();
        }, 1000);
        try {
            const res = await fetch(url, {
                signal: abc.signal,
            });
            if (res.status === 200) {
                return;
            }
            status = new Error(`${res.status} ${res.statusText}`);
        } catch (error) {
            if (error instanceof Error) {
                status = error;
            } else {
                status = new Error(String(error));
            }
            await new Promise((resolve) => setTimeout(resolve, 250));
        } finally {
            clearTimeout(timeout);
        }
    }
    throw status;
};

const buildGoPlugin = (filename: string, pluginPath: string): Promise<void> => {
    process.stdout.write(`### Building Plugin ${pluginPath}\n`);
    return new Promise((resolve, err) => {
        const build = spawn('go', ['build', '-o', filename, '-buildmode=plugin', pluginPath], {
            stdio: 'inherit',
        });

        build.on('close', (code) => {
            if (code) {
                err('exit code: ' + err);
            } else {
                resolve();
            }
        });
    });
};

const buildGoExecutable = (filename: string): Promise<void> => {
    const envGotify = process.env.GOTIFY_EXE;
    if (envGotify) {
        if (!fs.existsSync(testBuildPath)) {
            fs.mkdirSync(testBuildPath, {recursive: true});
        }
        fs.copyFileSync(envGotify, filename);
        process.stdout.write(`### Copying ${envGotify} to ${filename}\n`);
        return Promise.resolve();
    } else {
        process.stdout.write(`### Building Gotify ${filename}\n`);
        return new Promise((resolve, err) => {
            const build = spawn(
                'go',
                ['build', '-ldflags=-X main.Mode=prod', '-o', filename, appDotGo],
                {
                    stdio: 'inherit',
                }
            );
            build.on('close', (code) => {
                if (code) {
                    err('exit code: ' + err);
                } else {
                    resolve();
                }
            });
        });
    }
};

const startGotify = (filename: string, port: number, pluginDir: string): ChildProcess => {
    const gotify = spawn(filename, [], {
        env: {
            GOTIFY_SERVER_PORT: '' + port,
            GOTIFY_DATABASE_CONNECTION: 'file::memory:?mode=memory&cache=shared',
            GOTIFY_PLUGINSDIR: pluginDir,
            NODE_ENV: process.env.NODE_ENV,
            PUBLIC_URL: process.env.PUBLIC_URL,
        },
    });
    gotify.stdout.pipe(process.stdout);
    gotify.stderr.pipe(process.stderr);
    return gotify;
};
