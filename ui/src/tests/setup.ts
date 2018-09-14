import getPort from 'get-port';
import {spawn, exec, ChildProcess} from 'child_process';
import rimraf from 'rimraf';
import path from 'path';
import puppeteer, {Browser, Page} from 'puppeteer';
// @ts-ignore
import wait from 'wait-on';
import kill from 'tree-kill';

export interface GotifyTest {
    url: string;
    close: () => Promise<void>;
    browser: Browser;
    page: Page;
}

const windowsPrefix = process.platform === 'win32' ? '.exe' : '';
const appDotGo = path.join(__dirname, '..', '..', '..', 'app.go');

export const newTest = async (): Promise<GotifyTest> => {
    const port = await getPort();
    const gotifyFile = testFilePath();

    await buildGoExecutable(gotifyFile);

    const gotifyInstance = startGotify(gotifyFile, port);

    const gotifyURL = 'http://localhost:' + port;
    await waitForGotify('http-get://localhost:' + port);
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
                new Promise((resolve) => kill(gotifyInstance.pid, 'SIGKILL', () => resolve())),
            ]);
            rimraf.sync(gotifyFile, {maxBusyTries: 8});
        },
        url: gotifyURL,
        browser,
        page,
    };
};

const testFilePath = (): string => {
    const random = Math.random()
        .toString(36)
        .substring(2, 15);
    const filename = 'gotifytest_' + random + windowsPrefix;
    return path.join(__dirname, 'build', filename);
};

const waitForGotify = (url: string): Promise<void> => {
    return new Promise((resolve, err) => {
        wait({resources: [url], timeout: 40000}, (error: string) => {
            if (error) {
                console.log(error);
                err(error);
            } else {
                resolve();
            }
        });
    });
};

const buildGoExecutable = (filename: string): Promise<void> => {
    return new Promise((resolve) => exec(`go build  -o ${filename} ${appDotGo}`, () => resolve()));
};

const startGotify = (filename: string, port: number): ChildProcess => {
    const gotify = spawn(filename, [], {
        env: {
            GOTIFY_SERVER_PORT: '' + port,
            GOTIFY_DATABASE_CONNECTION: 'file::memory:?mode=memory&cache=shared',
        },
    });
    gotify.stdout.pipe(process.stdout);
    gotify.stderr.pipe(process.stderr);
    return gotify;
};
