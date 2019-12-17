module.exports = {
  root: true,
  parserOptions: {
    ecmaVersion: 2019,
    sourceType: 'module',
  },
  env: {
    es6: true,
    browser: true,
  },
  plugins: [
    'svelte3',
  ],
  extends: [
    'airbnb-base',
  ],
  overrides: [
    {
      files: ['**/*.svelte'],
      processor: 'svelte3/svelte3',
      rules: {
        // These rules don't evaluate correctly for svelte files
        // https://github.com/sveltejs/eslint-plugin-svelte3/blob/master/OTHER_PLUGINS.md#eslint-plugin-import
        'import/first': 0,
        'import/no-duplicates': 0,
        'import/no-mutable-exports': 0,
      },
    },
  ],
  rules: {
    'arrow-parens': ['error', 'as-needed', { requireForBlockBody: true }],
    'global-require': 0,
    'import/no-extraneous-dependencies': 0, // Doesn't like @sapper imports
    'import/prefer-default-export': 0,
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-param-reassign': [
      'error',
      {
        props: true,
        ignorePropertyModificationsFor: [
          'error',
          'acc',
        ],
      },
    ],
    'no-underscore-dangle': 0,
    'prefer-arrow-callback': 0,
    'space-before-function-paren': ['error', 'always'],
  },
  settings: {
    'import/resolver': {
      'eslint-import-resolver-custom-alias': {
        extensions: ['.mjs', '.js', '.json', '.svelte', '.html'],
        alias: {
          svelte: './node_modules/svelte',
          src: './src',
        },
      },
    },
  },
};
