module.exports = (grunt) ->
  grunt.initConfig
    typescript:
      main:
        src: 'src/main/typescript/Ignite.ts',
        dest: 'src/main/webapp/scripts/'
        options:
          target: 'es5'
          base_path: 'src/main/typescript'
          sourcemap: false,
          declaration_file: false

  grunt.loadNpmTasks 'grunt-typescript'

  grunt.registerTask 'default', [
    'typescript'
  ]