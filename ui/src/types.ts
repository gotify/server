interface IApplication {
    id: number;
    token: string;
    name: string;
    description: string;
    image: string;
    internal: boolean;
}

interface IClient {
    id: number;
    token: string;
    name: string;
}

interface IPlugin {
    id: number;
    token: string;
    name: string;
    modulePath: string;
    enabled: boolean;
    author?: string;
    website?: string;
    license?: string;
    capabilities: Array<'webhooker' | 'displayer' | 'configurer' | 'messenger' | 'storager'>;
}

interface IMessage {
    id: number;
    appid: number;
    message: string;
    title: string;
    priority: number;
    date: string;
    image?: string;
    extras?: IMessageExtras;
}

interface IMessageExtras {
    [key: string]: any; // tslint:disable-line no-any
}

interface IPagedMessages {
    paging: IPaging;
    messages: IMessage[];
}

interface IPaging {
    next?: string;
    since?: number;
    size: number;
    limit: number;
}

interface IUser {
    id: number;
    name: string;
    admin: boolean;
}

interface IVersion {
    version: string;
    commit: string;
    buildDate: string;
}
