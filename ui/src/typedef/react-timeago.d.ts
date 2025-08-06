declare module 'react-timeago' {
    import React from 'react';

    export type FormatterOptions = {
        style?: 'long' | 'short' | 'narrow';
    };
    export type Formatter = (options: FormatterOptions) => React.ReactNode;

    export interface ITimeAgoProps {
        date: string;
        formatter?: Formatter;
    }

    export default class TimeAgo extends React.Component<ITimeAgoProps, unknown> {}
}

declare module 'react-timeago/defaultFormatter' {
    declare function makeIntlFormatter(options: FormatterOptions): Formatter;
}
