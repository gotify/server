import * as os from 'os';
import {Page} from 'puppeteer';
import axios from 'axios';

import * as auth from './authentication';
import * as selector from './selector';
import {GotifyTest, newTest, newPluginDir} from './setup';
import {innerText, waitForCount, waitForExists} from './utils';

const pluginSupported = ['linux', 'darwin'].indexOf(os.platform()) !== -1;

let page: Page;
let gotify: GotifyTest;

beforeAll(async () => {
    const gotifyPluginDir = pluginSupported
        ? await newPluginDir(['github.com/gotify/server/v2/plugin/example/echo'])
        : '';
    gotify = await newTest(gotifyPluginDir);
    page = gotify.page;
});

afterAll(async () => await gotify.close());

enum Col {
    ID = 1,
    SetEnabled = 2,
    Name = 3,
    Token = 4,
    Details = 5,
}

const hiddenToken = '•••••••••••••••';

const $table = selector.table('#plugin-table');

const switchSelctor = (id: number) => $table.cell(id, Col.SetEnabled, '[data-enabled]');

const enabledState = async (id: number) =>
    (await page.$eval(switchSelctor(id), (el) => el.getAttribute('data-enabled'))) === 'true';

const toggleEnabled = async (id: number) => {
    const origEnabled = (await enabledState(id)).toString();
    await page.click(switchSelctor(id));
    await page.waitForFunction(
        `document.querySelector("${switchSelctor(
            id
        )}").getAttribute("data-enabled") !== "${origEnabled}"`
    );
};

const pluginInfo = async (className: string) =>
    await innerText(page, `.plugin-info .${className} > span`);

const getDisplayer = async () => await innerText(page, '.displayer');

const hasReceivedMessage = async (title: RegExp, content: RegExp) => {
    await page.click('#message-navigation a');
    await waitForExists(page, selector.heading(), 'All Messages');

    expect(await innerText(page, '.title')).toMatch(title);
    expect(await innerText(page, '.content')).toMatch(content);

    await page.click('#navigate-plugins');
    await waitForExists(page, selector.heading(), 'Plugins');
};

const inDetailPage = async (id: number, callback: () => Promise<void>) => {
    const name = await innerText(page, $table.cell(id, Col.Name));
    await page.click($table.cell(id, Col.Details, 'button'));
    await waitForExists(page, '.plugin-info .name > span', name);
    await callback();
    await page.click('#navigate-plugins');
    await waitForExists(page, selector.heading(), 'Plugins');
    await page.waitForSelector($table.selector());
};

describe('plugin', () => {
    describe('navigation', () => {
        it('does login', async () => await auth.login(page));
        it('navigates to plugins', async () => {
            await page.click('#navigate-plugins');
            await waitForExists(page, selector.heading(), 'Plugins');
        });
    });
    if (!pluginSupported) {
        return;
    }
    describe('functionality test', () => {
        describe('initial status', () => {
            it('has echo plugin', async () => {
                await waitForCount(page, $table.rows(), 1);
                expect(await innerText(page, $table.cell(1, Col.Name))).toEqual('test plugin');
                expect(await innerText(page, $table.cell(1, Col.Token))).toBe(hiddenToken);
                expect(parseInt(await innerText(page, $table.cell(1, Col.ID)), 10)).toBeGreaterThan(
                    0
                );
            });
            it('is disabled by default', async () => {
                expect(await enabledState(1)).toBe(false);
            });
        });
        describe('enable and disable plugin', () => {
            it('enable', async () => {
                await toggleEnabled(1);
                expect(await enabledState(1)).toBe(true);
            });

            it('disable', async () => {
                await toggleEnabled(1);
                expect(await enabledState(1)).toBe(false);
            });
        });
        describe('details page', () => {
            it('has plugin info', async () => {
                await inDetailPage(1, async () => {
                    expect(await pluginInfo('module-path')).toBe(
                        'github.com/gotify/server/v2/plugin/example/echo'
                    );
                });
            });
            it('has displayer', async () => {
                await inDetailPage(1, async () => {
                    expect(await getDisplayer()).toBeTruthy();
                });
            });
            it('has configurer', async () => {
                await inDetailPage(1, async () => {
                    expect(await page.$('.configurer')).toBeTruthy();
                });
            });
            it('updates configurer', async () => {
                await inDetailPage(1, async () => {
                    expect(
                        await (
                            await (await page.$('.config-save'))!.getProperty('disabled')
                        ).jsonValue()
                    ).toBe(true);
                    await page.waitForSelector('.CodeMirror .CodeMirror-code');
                    await page.waitForFunction(
                        'document.querySelector(".CodeMirror .CodeMirror-code").innerText.toLowerCase().indexOf("loading")<0'
                    );
                    await page.click('.CodeMirror .CodeMirror-code > div');
                    await page.keyboard.press('x');
                    await page.waitForFunction(
                        'document.querySelector(".config-save") && !document.querySelector(".config-save").disabled'
                    );
                    await page.click('.config-save');
                    await page.waitForFunction('document.querySelector(".config-save").disabled');
                });
            });
            it('configurer updated', async () => {
                await inDetailPage(1, async () => {
                    expect(
                        await (
                            await (await page.$('.config-save'))!.getProperty('disabled')
                        ).jsonValue()
                    ).toBe(true);
                    await page.waitForSelector('.CodeMirror .CodeMirror-code > div');
                    await page.waitForFunction(
                        'document.querySelector(".CodeMirror .CodeMirror-code > div").innerText.toLowerCase().indexOf("loading")<0'
                    );
                    expect(await innerText(page, '.CodeMirror .CodeMirror-code > div')).toMatch(
                        /x$/
                    );
                });
            });
            it('sends messages', async () => {
                if (!(await enabledState(1))) {
                    await toggleEnabled(1);
                }
                await inDetailPage(1, async () => {
                    await page.waitForSelector('.displayer a');
                    const hook = await page.$eval('.displayer a', (el) => el.getAttribute('href'));
                    await axios.get(hook);
                });
            });
            it('has received message', async () => {
                await hasReceivedMessage(
                    /^.+received$/,
                    /^echo server received a hello message \d+ times$/
                );
            });
        });
    });
});
