const glob = require('glob')
const path = require("path");
const webpack = require("webpack");

let files = glob.sync('./src/pages/**/index.jsx').reduce((acc, path) => {
  const entry = path.replace('./src/', '').replace('/index.jsx', '')
  acc[entry] = path
  return acc
}, {});

console.log(files)

module.exports = {
  entry: files,

output: {
    filename: './build/[name].js',
    path: path.resolve(__dirname)
},
optimization: {
  minimize: true,
  usedExports: true,
},
  mode: "production",
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /(node_modules|bower_components)/,
        loader: "babel-loader",
        options: { presets: ["@babel/env"] }
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"]
      }
    ]
  },
  resolve: { extensions: ["*", ".js", ".jsx"] },
};

