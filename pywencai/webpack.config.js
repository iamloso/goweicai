const path = require('path');

module.exports = {
    target: ['node', 'es5'],
    entry: './hexin-v.js',
    output: {
      filename: 'hexin-v.bundle.js',
      path: __dirname,
      environment: {
        arrowFunction: false,
        const: false,
        destructuring: false,
        forOf: false,
      }
    },
    externals: {
      'canvas': 'commonjs canvas'
    },
    optimization: {
      minimize: false
    },
    node: {
      __dirname: false,
      __filename: false,
    },
    module: {
      rules: [
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: [
                ['@babel/preset-env', {
                  targets: {
                    ie: '9'
                  },
                  useBuiltIns: false,
                  modules: false
                }]
              ]
            }
          }
        }
      ]
    }
  };
