const { Neutrino } = require('neutrino');

module.exports = Neutrino({ root: __dirname })
  .use('.neutrinorc.js', {
    eslint: {
      rules: {
        "indent": ["error", "tab"],
        "comma-dangle": ["warn", "always"],
      }
    }
  })
  .call('eslintrc');
