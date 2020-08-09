import path from 'path'
import { Configuration } from '@nuxt/types'

const config: Configuration = {
  mode: 'spa',
  /**
   * Source directory
   */
  srcDir: '.',
  /*
   ** Headers of the page
   */
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      {
        hid: 'description',
        name: 'description',
        content: process.env.npm_package_description || ''
      }
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }]
  },
  /*
   ** Customize the progress-bar color
   */
  loading: { color: '#2ebcfc' },
  /*
   ** Global CSS
   */
  css: [
    'normalize.css',
    '@/assets/css/main.scss',
    'material-design-iconic-font/dist/css/material-design-iconic-font.min.css',
    'vue2-animate/dist/vue2-animate.min.css'
  ],
  styleResources: {
    scss: [
      'sass-mq',
      '@/assets/css/_variables.scss',
      '@/assets/css/_functions.scss',
      '@/assets/css/_mixins.scss'
    ],
  },
  /*
   ** Plugins to load before mounting the App
   */
  plugins: [
    '@/plugins/jQuery.print.min.js',
    '@/plugins/vue-local-storage-decorator.ts'
  ],
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: [
    // Doc: https://github.com/nuxt-community/eslint-module
    '@nuxt/typescript-build',
    '@nuxtjs/eslint-module'
  ],
  /*
   ** Nuxt.js modules
   */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    '@nuxtjs/style-resources'
  ],
  /*
   ** Axios module configuration
   ** See https://axios.nuxtjs.org/options
   */
  axios: {
    baseURL: './'
  },
  /*
   ** Build configuration
   */
  build: {
    /*
     ** You can extend webpack config here
     */
    extend(config, ctx) {}
  }
}

export default config
