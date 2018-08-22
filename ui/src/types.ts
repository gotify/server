interface IApplication {
    id: number;
    token: string;
    name: string;
    description: string;
    image: string;
}

interface IClient {
    id: number;
    token: string;
    name: string;
}

interface IMessage {
    id: number;
    appid: number;
    message: string;
    title: string;
    priority: number;
    date: string;
    image?: string;
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

interface IAppMessages {
    messages: IMessage[];
    hasMore: boolean;
    nextSince: number;
    id?: number;
}
