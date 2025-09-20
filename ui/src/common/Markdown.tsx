import React from 'react';
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';

export const Markdown = ({
    children,
    onImageLoaded = () => {},
}: {
    children: string;
    onImageLoaded?: () => void;
}) => (
    <ReactMarkdown
        components={{img: ({...props}) => <img onLoad={onImageLoaded} {...props} />}}
        remarkPlugins={[gfm]}>
        {children}
    </ReactMarkdown>
);
