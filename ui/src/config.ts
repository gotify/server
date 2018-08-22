export interface IConfig {
    url: string;
}

let config: IConfig;

export function set(c: IConfig) {
    config = c;
}

export function get(val: 'url'): string {
    return config[val];
}
