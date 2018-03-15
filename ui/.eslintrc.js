module.exports = {
    'extends': ['google', 'plugin:react/recommended'],
    'parser': 'babel-eslint',
    'env': {
        'browser': true,
    },
    'globals': {
        'React': true,
    },
    'plugins': ['react', 'babel'],
    'rules': {
        'linebreak-style': 0,
        'max-len': ['error', 140],
        'require-jsdoc': 0,
        'babel/no-invalid-this': 1,
        'no-invalid-this': 0,
    },
};