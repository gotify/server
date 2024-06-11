import {ElementHandle, JSHandle, Page} from 'puppeteer';

export const innerText = async (page: ElementHandle | Page, selector: string): Promise<string> => {
    const element = await page.$(selector);
    const handle = await element!.getProperty('innerText');
    const value = await handle.jsonValue();
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return (value as any).toString().trim();
};

export const clickByText = async (page: Page, selector: string, text: string): Promise<void> => {
    await waitForExists(page, selector, text);
    text = text.toLowerCase();
    await page.evaluate(
        (_selector, _text) => {
            (
                Array.from(document.querySelectorAll(_selector)).filter(
                    (element) => element.textContent?.toLowerCase().trim() === _text
                )[0] as HTMLButtonElement
            ).click();
        },
        selector,
        text
    );
};

export const count = async (page: Page, selector: string): Promise<number> =>
    page.$$(selector).then((elements) => elements.length);

export const waitToDisappear = async (page: Page, selector: string): Promise<JSHandle> =>
    page.waitForFunction((_selector: string) => !document.querySelector(_selector), {}, selector);

export const waitForCount = async (
    page: Page,
    selector: string,
    amount: number
): Promise<JSHandle> =>
    page.waitForFunction(
        (_selector: string, _amount: number) =>
            document.querySelectorAll(_selector).length === _amount,
        {},
        selector,
        amount
    );

export const waitForExists = async (page: Page, selector: string, text: string): Promise<void> => {
    text = text.toLowerCase();
    await page.waitForFunction(
        (_selector: string, _text: string) =>
            Array.from(document.querySelectorAll(_selector)).filter(
                (element) => element.textContent!.toLowerCase().trim() === _text
            ).length > 0,
        {},
        selector,
        text
    );
};

export const clearField = async (element: ElementHandle | Page, selector: string) => {
    const elementHandle = await element.$(selector);
    if (!elementHandle) {
        fail();
    }
    await elementHandle.click();
    await elementHandle.focus();
    // click three times to select all
    await elementHandle.click({clickCount: 3});
    await elementHandle.press('Backspace');
};
