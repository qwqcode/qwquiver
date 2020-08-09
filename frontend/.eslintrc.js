module.exports = {
  root: true,
  env: {
    browser: true,
    node: true
  },
  extends: [
    '@nuxtjs/eslint-config-typescript',
    'prettier',
    'prettier/vue',
    'plugin:nuxt/recommended'
  ],
  plugins: ['prettier', 'import'],
  // add your custom rules here
  rules: {
    '@typescript-eslint/no-unused-vars': 'off',
    'nuxt/no-cjs-in-config': 'off',
    'no-unused-vars': 'off',
    'require-await': 'off',
    'lines-between-class-members': 'off'
  },
  settings: {
    'import/core-modules': ['nuxt-property-decorator'],
    'import/resolver': {
      // use <root>/tsconfig.json
      'typescript': {},
    }
  }
}
