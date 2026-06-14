import {spawn, execFileSync} from 'child_process';
import getPort from 'get-port';
import fs from 'fs';
import os from 'os';
import path from 'path';
import {rimrafSync} from 'rimraf';
import {stringify} from 'yaml';
// @ts-expect-error no types
import wait from 'wait-on';

// All dex test users share this password. The hash is bcrypt("password").
export const DEX_PASSWORD = 'password';
const DEX_PASSWORD_HASH = '$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W';

const DEX_IMAGE = 'ghcr.io/dexidp/dex:v2.45.1-alpine';

export interface DexUser {
    email: string;
    username: string;
    userID: string;
}

export interface DexInstance {
    issuer: string;
    close: () => void;
}

const dexConfig = (issuerPort: number, redirectURL: string, users: DexUser[]): string =>
    stringify({
        issuer: `http://127.0.0.1:${issuerPort}/dex`,
        storage: {type: 'memory'},
        web: {http: '0.0.0.0:5556'},
        oauth2: {skipApprovalScreen: true},
        staticClients: [
            {
                id: 'gotify',
                name: 'Gotify',
                secret: 'secret',
                redirectURIs: [redirectURL],
            },
        ],
        enablePasswordDB: true,
        staticPasswords: users.map((u) => ({
            email: u.email,
            hash: DEX_PASSWORD_HASH,
            username: u.username,
            preferredUsername: u.username,
            userID: u.userID,
        })),
    });

const waitForDex = (port: number): Promise<void> =>
    new Promise((resolve, reject) => {
        wait(
            {
                resources: [`http-get://127.0.0.1:${port}/dex/.well-known/openid-configuration`],
                timeout: 60000,
            },
            (error: string) => (error ? reject(error) : resolve())
        );
    });

export const startDex = async (redirectURL: string, users: DexUser[]): Promise<DexInstance> => {
    const port = await getPort();
    const configDir = fs.mkdtempSync(path.join(os.tmpdir(), 'gotify-dex-'));
    fs.writeFileSync(path.join(configDir, 'dex.conf'), dexConfig(port, redirectURL, users));

    const userArgs =
        process.getuid && process.getgid
            ? ['--user', `${process.getuid()}:${process.getgid()}`]
            : [];

    const containerName = `gotify-dex-test-${port}`;
    process.stdout.write(`### Starting dex ${containerName}\n`);
    const dex = spawn('docker', [
        'run',
        '--rm',
        ...userArgs,
        '--name',
        containerName,
        '-p',
        `${port}:5556`,
        '-v',
        `${configDir}:/config`,
        DEX_IMAGE,
        'dex',
        'serve',
        '/config/dex.conf',
    ]);
    dex.stdout.pipe(process.stdout);
    dex.stderr.pipe(process.stderr);

    const crashed = new Promise<never>((_, reject) => {
        const abort = (reason: string) => reject(new Error(`dex ${containerName} ${reason}`));
        dex.on('exit', (code, signal) =>
            abort(`exited unexpectedly (code=${code}, signal=${signal})`)
        );
        dex.on('error', (err) => abort(`failed to start: ${err.message}`));
    });

    const cleanup = () => {
        try {
            execFileSync('docker', ['rm', '-f', containerName], {stdio: 'ignore'});
        } catch {
            // container may already be gone (e.g. it crashed with --rm)
        }
        rimrafSync(configDir);
    };

    try {
        await Promise.race([waitForDex(port), crashed]);
    } catch (err) {
        cleanup();
        throw err;
    }

    return {
        issuer: `http://127.0.0.1:${port}/dex`,
        close: cleanup,
    };
};
