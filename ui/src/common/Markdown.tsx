import React from 'react';
import ReactMarkdown, {defaultUrlTransform} from 'react-markdown';
import type {UrlTransform} from 'react-markdown';
import gfm from 'remark-gfm';

// Copy from mlflow/server/js/src/shared/web-shared/genai-markdown-renderer/GenAIMarkdownRenderer.tsx
// Related PR: https://github.com/mlflow/mlflow/pull/16761
const urlTransform: UrlTransform = (value) => {
    if (value.startsWith('data:image/png;') || value.startsWith('data:image/jpeg;') || value.startsWith('data:image/gif;')) {
        return value;
    }
    return defaultUrlTransform(value);
};

export const Markdown = ({
    children,
    onImageLoaded = () => {},
}: {
    children: string;
    onImageLoaded?: () => void;
}) => (
    <ReactMarkdown
        components={{img: ({...props}) => <img onLoad={onImageLoaded} {...props} />}}
        remarkPlugins={[gfm]}
        urlTransform={urlTransform}>
        {children}
    </ReactMarkdown>
);
