declare module 'react-timeago' {
    import React from 'react';

    export interface ITimeAgoProps {
        date: string;
    }

    export default class TimeAgo extends React.Component<ITimeAgoProps, unknown> {}
}
