let config;

export function set(c) {
    config = c;
}

export function get(val): string {
    return config[val];
}