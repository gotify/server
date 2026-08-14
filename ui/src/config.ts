import {IVersion} from './types';

export interface IConfig {
    url: string;
    register: boolean;
    version: IVersion;
    oidc: boolean;
    localAuth: boolean;
}

declare global {
    interface Window {
        config?: Partial<IConfig>;
    }
}

const config: IConfig = {
    url: 'unset',
    register: false,
    version: {commit: 'unknown', buildDate: 'unknown', version: 'unknown'},
    oidc: false,
    localAuth: true,
    ...window.config,
};

export function set<Key extends keyof IConfig>(key: Key, value: IConfig[Key]): void {
    config[key] = value;
}

export function get<K extends keyof IConfig>(key: K): IConfig[K] {
    return config[key];
}
